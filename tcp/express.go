package tcp

import (
	"fmt"
	"log"
	"net"

	"github.com/xtaci/smux"
)

var (
	sessions = make(chan *smux.Session, 512)
)

func ExpressDial(protocl, dst string, usesmux bool, first ...[]byte) (con net.Conn, reply []byte, err error) {

	if usesmux {
		var sess *smux.Session
		sess, err = WithASession(protocl, dst, sessions)
		if err != nil {
			return nil, []byte{}, err
		}
		// fmt.Println("Get a session")
		con, err = sess.OpenStream()
		if err != nil {
			sess.Close()
			log.Println("connect forward error:", err)
			return
		}
		if first != nil {
			con.Write([]byte(first[0]))
		}

	} else {
		switch protocl {
		case "tls":
			tlsConf := UseDefaultTlsConfig(dst)
			con, reply, err = initiTlsConnection(tlsConf, first...)
		case "wss":
			con, reply, err = ConnectWssAndFirstBuf(dst, first...)
		case "ws":
			con, reply, err = ConnectWsAndFirstBuf(dst, first...)
		}
	}
	return
}

type ExpressListener struct {
	protocol string
	usesmux  bool
	nowClose bool
	lastErr  error
	Sessions chan *smux.Session
	rawLst   net.Listener
	subLst   net.Listener
	acceptCh chan net.Conn
}

func (expressLst *ExpressListener) Accept() (con net.Conn, err error) {
	if expressLst.usesmux {
		con = <-expressLst.acceptCh
		return
	}
	con, err = expressLst.rawLst.Accept()
	return
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (expressLst *ExpressListener) Close() error {
	expressLst.nowClose = true
	return nil
}

// Addr returns the listener's network address.
func (expressLst *ExpressListener) Addr() net.Addr {
	return expressLst.rawLst.Addr()
}

func (lst *ExpressListener) Backend() {
	defer fmt.Println("END")

	go Average(lst.Sessions, 600)
	for {
		if lst.nowClose {
			log.Println("Close listener")
			break
		}

		// fmt.Println("con connected b")
		con, err := lst.rawLst.Accept()
		lst.lastErr = err
		if err != nil {
			log.Println("accept err :", err)
			break
		}
		smuxListner, err := smux.Server(con, nil)
		if err != nil {
			log.Println("smux server wrap err:", err)
			break
		}
		for {
			if lst.nowClose {

				log.Println("Close session listener")
				break
			}

			// fmt.Println("session connected b")
			scon, err := smuxListner.AcceptStream()
			// fmt.Println("session connected")
			lst.lastErr = err
			if err != nil {
				log.Println("this tunel accep one sesion err:", err)
				break
			}

			// go func(chs chan net.Conn) {
			lst.acceptCh <- scon
			// }(lst.acceptCh)
		}

	}

}

// func (lst *ExpressListener) Handle(do func(acceptCon net.Conn)) {
// 	defer fmt.Println("END")
// 	chs := make(chan *smux.Session, 100)
// 	go Average(chs, 600)
// 	for {
// 		var con net.Conn
// 		var err error
// 		if lst.usesmux {
// 			con, err = lst.rawLst.Accept()
// 			// fmt.Println("ss")
// 		} else {
// 			con, err = lst.Accept()
// 		}
// 		if err != nil {
// 			log.Println("accept err:", err)
// 			break
// 		}
// 		if lst.usesmux {

// 			smuxListner, err := smux.Server(con, nil)
// 			// fmt.Println("ss 2")
// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			for {
// 				scon, err := smuxListner.AcceptStream()
// 				if err != nil {
// 					log.Println("sess accept err:", err)
// 				}
// 				// fmt.Println("ss 3")

// 				go do(scon)
// 			}
// 		} else {
// 			go do(con)
// 		}
// 	}
// }

func ExpressListenWith(protocl, addr string, useSmux bool) (lst *ExpressListener, err error) {
	lst = new(ExpressListener)
	lst.protocol = protocl
	lst.usesmux = useSmux
	switch protocl {
	case "tls":
		// log.Println("parse tls before:")
		tlsConf := UseDefaultTlsConfig(addr)
		lst.rawLst, err = tlsConf.WithTlsListener()
		// log.Println("parse tls err:", err)
	case "wss":
		lst.rawLst, err = UseWebSocketListener(addr, true)
	case "ws":
		lst.rawLst, err = UseWebSocketListener(addr, false)
	}
	if lst.usesmux {
		lst.acceptCh = make(chan net.Conn, 600)
		lst.Sessions = make(chan *smux.Session, 256)
		go lst.Backend()
	}

	return
}

func Socks5ConnectionForward(fr net.Conn, protocol, server string, socksheader []byte, chs chan *smux.Session) {
	// dstCon, err := UseDefaultTlsConfig(server).WithConn()
	session, err := WithASession(protocol, server, chs)
	dstCon, err := session.OpenStream()
	if err != nil {
		log.Println("connect forward error:", err)
		return
	}
	_, err = dstCon.Write(socksheader)
	if err != nil {
		log.Println("forawrd err:", err)
		fr.Close()
		return
	}
	Pipe(fr, dstCon)
}

func ExpressConnectTo(fr net.Conn, protocol, remoteAddr string, raw []byte, multi bool) (err error) {
	dstcon, reply, err := ExpressDial(protocol, remoteAddr, multi, raw)
	// Socks5ConnectionForward(con, protocol, remoteAddr, raw, ifmulti)
	if err != nil {
		log.Println("connect reply err :", err)
		return err
	}
	fr.Write(reply)
	Pipe(fr, dstcon)
	return nil
}

func ExpressSocks5ToListener(protocol, listenAddr string, remoteAddr string, ifmulti ...bool) {
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	useMulti := false
	if ifmulti != nil && ifmulti[0] {
		useMulti = true
	}
	for {
		con, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			if err := Socks5HandShake(con); err != nil {
				log.Println("socks5 handle failed!", err)
				return
			}
			raw, _, err := ParseSocks5Header(con)
			if err != nil {
				log.Println("no socks5 addr header")
				return
			}
			ExpressConnectTo(con, protocol, remoteAddr, raw, useMulti)

		}()

	}
}
