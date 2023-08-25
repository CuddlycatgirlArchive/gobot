package attacks

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"syscall"
)

func localIPPort(dstip net.IP) (net.IP, int) {
	serverAddr, err := net.ResolveUDPAddr("udp", dstip.String()+":12345")
	if err != nil {
		log.Fatal(err)
	}

	// We don't actually connect to anything, but we can determine
	// based on our destination ip what source ip we should use.
	if con, err := net.DialUDP("udp", nil, serverAddr); err == nil {
		if udpaddr, ok := con.LocalAddr().(*net.UDPAddr); ok {
			return udpaddr.IP, udpaddr.Port
		}
	}
	if err != nil {
		log.Fatal("could not get local ip: " + err.Error())
	}
	return nil, -1
}
func sendRawPacket(dst net.IP, rawPacket []byte) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, 0xff)
	file := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
		file = nil
	}(file)
	conn, err := net.FilePacketConn(file)
	defer func(conn net.PacketConn) {
		err := conn.Close()
		if err != nil {

		}
		conn = nil
	}(conn)
	if err != nil {
		println(err.Error())
		return
	}
	_, err = conn.WriteTo(rawPacket, &net.IPAddr{IP: dst})
	if err != nil {
		println(err.Error())
		return
	}
}
func sendProtoPacket(dst net.IP, protocol string, rawPacket []byte) {
	conn, err := net.ListenPacket("ip4:"+protocol, "0.0.0.0")
	if err != nil {
		log.Fatal(err)
		return
	}
	if _, err := conn.WriteTo(rawPacket, &net.IPAddr{IP: dst}); err != nil {
		log.Fatal(err)
	}
	conn.Close()
	conn = nil
}
func getProtoConn(protocol string) net.PacketConn {
	conn, _ := net.ListenPacket("ip4:"+protocol, "0.0.0.0")
	return conn
}
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
