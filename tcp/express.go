package tcp

import (
	"log"
	"net"

	"github.com/Qingluan/Tunnel/config"
	"github.com/xtaci/smux"
)

func ExpressDial(dst string, configs ...interface{}) (con net.Conn, reply []byte, err error) {
	defaultconfig := config.ParseConfigs(configs...)

	T("Pro:", defaultconfig.Protocol, "use mutl:", defaultconfig.Multi, "first:", defaultconfig.First)

	if defaultconfig.Multi {
		var sess *smux.Session
		sess, err = WithASession(defaultconfig.Protocol, dst)
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

	} else {

		con, err = ConnectTo(defaultconfig.Protocol, dst)
		if err != nil {
			return
		}

	}
	if defaultconfig.First != nil {
		_, err = con.Write(defaultconfig.First)
		if err != nil {
			return
		}
		// buf := make([])
	}
	return
}

func ExpressListenWith(addr string, configs ...interface{}) (lst *ExpressListener, err error) {
	lst = new(ExpressListener)
	defaultconfig := config.ParseConfigs(configs...)

	lst.protocol = defaultconfig.Protocol
	lst.usesmux = defaultconfig.Multi
	switch lst.protocol {
	case "tls":
		// log.Println("parse tls before:")
		tlsConf := UseDefaultTlsConfig(addr)
		lst.rawLst, err = tlsConf.WithTlsListener()
		// log.Println("parse tls err:", err)
	case "wss":
		lst.rawLst, err = UseWebSocketListener(addr, true)
	case "ws":
		lst.rawLst, err = UseWebSocketListener(addr, false)
	case "tcp":
		lst.rawLst, err = net.Listen("tcp", addr)
	case "unix":
		lst.rawLst, err = net.Listen("unix", addr)
	}
	log.Printf("Run Server At : %s://%s (Start Multi Channel:%v)\n", lst.protocol, addr, lst.usesmux)
	if err != nil {
		return
	}
	if lst.usesmux {
		lst.acceptCh = make(chan net.Conn, 600)
		lst.Sessions = make(chan *smux.Session, 256)
		go lst.Backend()
	}

	return
}

func ExpressPipeTo(localConn net.Conn, remoteAddr string, configs ...interface{}) error {
	dstcon, _, err := ExpressDial(remoteAddr, configs...)
	T(2)

	if err != nil {
		log.Println("connect reply err :", err)
		return err
	}

	Pipe(localConn, dstcon)
	return nil
}

func ExpressSocks5ListenerTo(listenAddr string, remoteAddr string, configs ...interface{}) {

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	d := config.ParseConfigs(configs...)
	log.Printf("Run Local Socks5 Server At tcp://%s <---> %s://%s (Start Multi Channel:%v)\n", listenAddr, d.Protocol, remoteAddr, d.Multi)

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

			configs = append(configs, raw)
			T(1)
			ExpressPipeTo(con, remoteAddr, configs...)

		}()

	}
}

func ExpressSocks5InitConnection(con net.Conn, realHost string) (ok bool) {
	ok = Socks5Init(con, realHost)
	return
}
