package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/Qingluan/Tunnel/tcp"
)

var (
	P = "tls"
)

func main() {
	ifmulti := false

	flag.StringVar(&P, "p", "tls", "set protocol")
	flag.BoolVar(&ifmulti, "m", false, "ifmulti ")

	flag.Parse()
	fmt.Println(ifmulti)

	lst, err := tcp.ExpressListenWith(P, "127.0.0.1:12345", ifmulti)
	if err != nil {
		log.Fatal(err)
	}
	// do := func(con net.Conn) {
	// 	if con == nil {
	// 		fmt.Println("null")
	// 		return
	// 	}
	// 	buf := make([]byte, 1024)
	// 	if n, err := con.Read(buf); err == nil {
	// 		log.Println("R:", string(buf[:n]))
	// 	}

	// 	// log.Println("Accept:", con.RemoteAddr().String())

	// 	con.Write([]byte("hello\n"))
	// 	con.Close()
	// 	fmt.Printf("%s\n", "[o]")
	// }

	for {
		con, _ := lst.Accept()
		go tcp.Socks5ForwardServer(P, con, ifmulti)
	}

	// var lst net.Listener
	// var err error
	// addr := "127.0.0.1:12345"
	// switch P {
	// case "tls":
	// 	// log.Println("parse tls before:")
	// 	tlsConf := tcp.UseDefaultTlsConfig(addr)
	// 	lst, err = tlsConf.WithTlsListener()
	// 	log.Println("parse tls err:", err)
	// case "wss":
	// 	lst, err = tcp.UseWebSocketListener(addr, true)
	// case "ws":
	// 	lst, err = tcp.UseWebSocketListener(addr, false)
	// }
	// sess := make(chan *smux.Session, 600)
	// for {
	// 	con, _ := lst.Accept()
	// 	if ifmulti {
	// 		ss, _ := smux.Server(con, nil)
	// 		for {
	// 			con, err := ss.AcceptStream()
	// 			if err != nil {
	// 				log.Fatal(err)
	// 			}
	// 			go tcp.Socks5Forward(P, con, sess)
	// 		}
	// 	}
	// 	// go tcp.Socks5Forward(P, con, lst.Sessions)
	// }

}
