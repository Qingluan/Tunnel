package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	typeRedirect         = 9
	CHAIN                = "<=CHAIN=>"
	AddrMask        byte = 0xf
	socksVer5            = 5
	socksCmdConnect      = 1
	socksCmdUdp          = 3
	idVer                = 0
	idCmd                = 1
	idType               = 3 // address type index
	idIP0                = 4 // ip address start index
	idDmLen              = 4 // domain address length index
	idDm0                = 5 // domain address start index

	typeIPv4   = 1 // type is ipv4 address
	typeDm     = 3 // type is domain address
	typeIPv6   = 4 // type is ipv6 address
	typeChange = 5 // type is ss change config

	lenIPv4   = 3 + 1 + net.IPv4len + 2 // 3(ver+cmd+rsv) + 1addrType + ipv4 + 2port
	lenIPv6   = 3 + 1 + net.IPv6len + 2 // 3(ver+cmd+rsv) + 1addrType + ipv6 + 2port
	lenDmBase = 3 + 1 + 1 + 2           // 3 + 1addrType + 1addrLen + 2port, plus addrLen
	// lenHmacSha1 = 10
)

var (
	errAddrType      = errors.New("socks addr type not supported")
	errVer           = errors.New("socks version not supported")
	errMethod        = errors.New("socks only support 1 method now")
	errAuthExtraData = errors.New("socks authentication get extra data")
	errReqExtraData  = errors.New("socks request get extra data")
	errCmd           = errors.New("socks command not supported")
	SOCKS5_REPY      = []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43}
	GLOBALTIMEOUT    = 1000
	// GLOBALPROTOCOL   = "tls"
	// smuxConfig = smux.DefaultConfig()
	ServeHandle func(raw []byte, host string, con net.Conn)
)

// func SetProtocol(i string) {
// 	locker.Lock()
// 	defer locker.Unlock()
// 	// GLOBALPROTOCOL = i
// }

func SetGlobalTimeout(i int) {
	locker.Lock()
	defer locker.Unlock()

	GLOBALTIMEOUT = i
}

func SetReadTimeout(c net.Conn) {
	(c).SetReadDeadline(time.Now().Add(time.Duration(GLOBALTIMEOUT) * time.Second))
}

func SetWriteTimeout(c *net.Conn) {
	(*c).SetWriteDeadline(time.Now().Add(time.Duration(GLOBALTIMEOUT) * time.Second))
}

func ParseSocks5Header(conn net.Conn) (rawaddr []byte, host string, err error) {

	// refer to getRequest in server.go for why set buffer size to 263
	buf := make([]byte, 263)
	var n int
	SetReadTimeout(conn)
	// read till we get possible domain length field
	if n, err = io.ReadAtLeast(conn, buf, idDmLen+1); err != nil {
		return
	}
	// ColorL("->", buf[:10])
	// check version and cmd
	if buf[idVer] != socksVer5 {

		err = errors.New("Sock5 error: " + string(buf[idVer]))
		return
	}

	reqLen := -1
	switch buf[idType] {
	case typeIPv4:
		reqLen = lenIPv4
	case typeIPv6:
		reqLen = lenIPv6
	case typeDm:
		reqLen = int(buf[idDmLen]) + lenDmBase
		host = string(buf[idDm0 : idDm0+buf[idDmLen]])

	case typeChange:
		reqLen = int(buf[idDmLen]) + lenDmBase - 2
		host = string(buf[idDm0 : idDm0+buf[idDmLen]])
		// ColorL("hh", host)
	default:
		err = errAddrType
		return
	}
	// ColorL("hq", buf[:10])

	if n == reqLen {
		// common case, do nothing
	} else if n < reqLen { // rare case
		if _, err = io.ReadFull(conn, buf[n:reqLen]); err != nil {
			return
		}
	} else {
		fmt.Println(n, reqLen, buf)
		err = errReqExtraData
		return
	}

	// rawaddr = buf[idType:reqLen]
	rawaddr = buf[:reqLen]
	switch buf[idType] {
	case typeIPv4:
		host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
	case typeIPv6:
		host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
	case typeDm:
		host = string(buf[idDm0 : idDm0+buf[idDmLen]])
	case typeChange:
		host = string(buf[idDm0 : idDm0+buf[idDmLen]])
		// ColorL("hm", host)
		return
	}
	port := binary.BigEndian.Uint16(buf[reqLen-2 : reqLen])
	// host = net.JoinHostPort(host, strconv.Itoa(int(port)))

	host = fmt.Sprintf("%s:%d", host, int(port))
	log.Println("host:", host)
	return
}

func Socks5HandShake(conn net.Conn) (err error) {
	const (
		idVer     = 0
		idNmethod = 1
	)
	// version identification and method selection message in theory can have
	// at most 256 methods, plus version and nmethod field in total 258 bytes
	// the current rfc defines only 3 authentication methods (plus 2 reserved),
	// so it won't be such long in practice
	// SetReadTimeout(conn)
	buf := make([]byte, 258)
	var n int
	if n, err = io.ReadAtLeast(conn, buf, idNmethod+1); err != nil {
		return
	}
	if buf[idVer] != socksVer5 {
		log.Println(buf)
		return errVer
	}
	nmethod := int(buf[idNmethod])
	msgLen := nmethod + 2
	if n == msgLen { // handshake done, common case
		// do nothing, jump directly to send confirmation
	} else if n < msgLen { // has more methods to read, rare case
		if _, err = io.ReadFull(conn, buf[n:msgLen]); err != nil {
			return
		}
	} else { // error, should not get extra data
		log.Println(buf)
		return errAuthExtraData
	}
	// send confirmation: version 5, no authentication required
	if _, err = (conn).Write([]byte{socksVer5, 0}); err != nil {
		return err
	}
	return
}

func Socks5ConnectedReply(p1 net.Conn) (err error) {
	_, err = p1.Write(SOCKS5_REPY)
	return
}

func SetSocks5Handle(f func(raw []byte, host string, con net.Conn)) {
	locker.Lock()
	defer locker.Unlock()
	ServeHandle = f
}

func Socks5Serve(lc net.Conn, configs ...interface{}) (err error) {

	raw, host, err := ParseSocks5Header(lc)
	if err != nil {
		log.Println("err:", err)
		return
	}

	if ServeHandle != nil {
		ServeHandle(raw, host, lc)
	} else {
		TcpEnd(host, lc, SOCKS5_REPY)
	}
	// if strings.Contains(host, "config://menu") {
	// 	log.Println("go config>>>")
	// } else if strings.Contains(host, "config://alive") {
	// 	_, err = lc.Write([]byte("ok"))
	// } else if strings.HasPrefix(host, "proxys://") {
	// 	fields := strings.Split(host, CHAIN)
	// 	nexts := strings.Split(fields[0], "proxys://")
	// 	c := rand.Int() % len(nexts)
	// 	next := nexts[c]
	// 	Special := strings.Join(fields[1:], CHAIN)
	// 	log.Println("--->", next, "||", Special)
	// 	configs = append(configs, config.Config(Socks5Padding(Special)))
	// 	ExpressPipeTo(lc, next, configs...)
	// } else if strings.HasPrefix(host, "proxy://") {
	// 	fields := strings.Split(host, CHAIN)
	// 	next := strings.SplitN(fields[0], "proxy://", 2)[1]
	// 	Special := strings.Join(fields[1:], CHAIN)
	// 	log.Println("--->", next, "||", Special)
	// 	configs = append(configs, config.Config(Socks5Padding(Special)))
	// 	ExpressPipeTo(lc, next, configs...)
	// } else {
	// 	TcpEnd(host, lc, SOCKS5_REPY)
	// }
	return
}

func Socks5Padding(payload string, ports ...int) []byte {

	buf := bytes.NewBuffer([]byte{0x5, 0x1, 0x0, 0x3})
	l := len(payload)
	var port int
	fmt.Println(payload, "!!")
	if ports != nil {
		port = ports[0]
	} else {
		if strings.Contains(payload, ":") && payload != "config://menu" {
			fs := strings.Split(payload, ":")
			payload = strings.Join(fs[:len(fs)-1], ":")
			l = len(payload)
			port, _ = strconv.Atoi(fs[len(fs)-1])
		}
	}
	lbuf := make([]byte, 2)
	binary.BigEndian.PutUint16(lbuf, uint16(l))
	buf.WriteByte(lbuf[1])
	buf.Write([]byte(payload))
	if port != 0 {
		portbuf := make([]byte, 2)
		binary.BigEndian.PutUint16(portbuf, uint16(port))
		buf.Write(portbuf)
	} else {
		port = 80
		portbuf := make([]byte, 2)
		binary.BigEndian.PutUint16(portbuf, uint16(port))
		buf.Write(portbuf)
	}
	return buf.Bytes()
}

func Socks5Init(con net.Conn, host string) (ok bool) {
	_, err := con.Write([]byte{0x5, 0x1, 0x0})
	if err != nil {
		return
	}
	buf := make([]byte, 3)
	_, err = con.Read(buf)
	if err != nil {
		return
	}
	if bytes.Equal(buf[:2], []byte{0x5, 0x0}) {
		ok = true
	}
	hostBUf := Socks5Padding(host)
	_, err = con.Write(hostBUf)
	if err != nil {
		ok = false
		return
	}
	return
}
