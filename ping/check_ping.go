package main

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type messageType int

const (
	icmpV4TypeEchoRequest = 8
	icmpV4TypeEchoReply   = 0
	icmpV6TypeEchoRequest = 128
	icmpV6TypeEchoReply   = 129
)

type messageBody struct {
	ID   int
	Seq  int
	Data []byte
}

func (body *messageBody) Marshal() []byte {
	b := make([]byte, 8+len(body.Data))
	b[0] = byte(body.ID >> 8)
	b[1] = byte(body.ID & 0xFF)
	b[2] = byte(body.Seq >> 8)
	b[3] = byte(body.Seq & 0xFF)

	copy(b[4:], body.Data)
	return b
}

type message struct {
	Type     messageType
	Code     int
	CheckSum int
	Body     *messageBody
}

func (msg *message) Marshal() []byte {
	b := []byte{byte(msg.Type), byte(msg.Code), 0, 0}

	bodyByte := msg.Body.Marshal()
	b = append(b, bodyByte...)

	if msg.Type == icmpV6TypeEchoRequest {
		return b
	}

	sum := checksum(b)
	b[2] ^= byte(sum & 0xFF)
	b[3] ^= byte(sum >> 8)
	return b
}

func checksum(b []byte) uint32 {
	lens := len(b) - 1
	sum := uint32(0)
	for i := 0; i < lens; i += 2 {
		sum += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if sum&1 == 0 {
		sum += uint32(b[lens])
	}

	sum = sum>>16 + sum&0xFFFF
	sum = sum + sum>>16
	return ^sum
}

func combineData(packetSize, seq int, ifIPV6 bool) []byte {
	if packetSize <= 8 {
		packetSize = 9
	}
	bodySize := packetSize - 8
	body := messageBody{
		ID:   os.Getpid() & 0xFFFF,
		Seq:  seq,
		Data: bytes.Repeat([]byte("c"), bodySize),
	}

	var msg message
	if ifIPV6 {
		msg.Type = icmpV6TypeEchoRequest
	} else {
		msg.Type = icmpV4TypeEchoRequest
	}
	msg.Body = &body
	msgBody := msg.Marshal()
	return msgBody
}

func stripIPv4Header(b []byte, lens int) []byte {
	if len(b) < 20 {
		return b
	}

	l := int(b[0]&0x0F) << 2
	if 20 > l || l > len(b) {
		return b
	}

	if b[0]>>4 != 4 {
		return b
	}

	return b[l:]
}

func parseICMPMessage(b []byte) (id, seqnum int) {
	id = int(b[4])<<8 | int(b[5])
	seqnum = int(b[6])<<8 | int(b[7])
	return
}

func transfer(conn net.Conn, ifIPV6 bool, msgBody []byte) (err error) {
	if _, err = conn.Write(msgBody); err != nil {
		return
	}
	var lens int
	reply := make([]byte, len(msgBody)+128)
	for {
		lens, err = conn.Read(reply)
		if err != nil {
			return
		}

		if ifIPV6 {
			if reply[0] != icmpV6TypeEchoReply {
				continue
			}
		} else {
			reply = stripIPv4Header(reply, lens)
			if reply[0] != icmpV4TypeEchoReply {
				continue
			}
		}

		requestID, requestNum := parseICMPMessage(msgBody)
		responseID, responseNum := parseICMPMessage(reply)
		if requestID != responseID || requestNum != responseNum {
			err = fmt.Errorf("request %d %d response %d %d", requestID, requestNum, responseID, responseNum)
			return
		}
		break
	}
	return
}

func ping(address string, seq, packetSize, timeout int, ifIPV6 bool) (err error) {

	var conn net.Conn
	if ifIPV6 {
		conn, err = net.Dial("ip6:ipv6-icmp", address)
	} else {
		conn, err = net.Dial("ip4:icmp", address)
	}
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))

	msgBody := combineData(packetSize, seq, ifIPV6)
	before := time.Now().UnixNano()
	if err = transfer(conn, ifIPV6, msgBody); err != nil {
		if strings.Contains(err.Error(), "i/o timeout") {
			fmt.Printf("Ping %s icmp_seq:%d timeout\n", address, seq)
			err = nil
		}
		return
	}
	duration := time.Now().UnixNano() - before
	fmt.Printf("Ping %s icmp_seq=%d time=%dms\n", address, seq, duration/1000000)
	return
}

func run(address string, count, packetSize, timeout int, ifIPV6 bool) (err error) {

	for i := 0; i < count; i++ {
		if err = ping(address, i+1, packetSize, timeout, ifIPV6); err != nil {
			return
		}
	}

	return
}

var host = flag.String("h", "www.baidu.com", "地址")
var count = flag.Int("c", 10, "请求次数")
var packetSize = flag.Int("s", 10, "请求大小，单位byte")
var timeout = flag.Int("t", 1000, "超时时间，单位：毫秒")
var ipv6 = flag.Bool("ipv6", false, "是否请求ipv6")

func main() {
	flag.Parse()
	err := run(*host, *count, *packetSize, *timeout, *ipv6)
	if err != nil {
		fmt.Println("Find Error ", err)
	}
}
