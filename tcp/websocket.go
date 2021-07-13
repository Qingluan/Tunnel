package tcp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Qingluan/Tunnel/tcp/websocket"

	"sync/atomic"
)

var (
	ErrWebsocketListenerClosed = errors.New("websocket listener closed")
)

var (
	FrpWebsocketPath2 = "/image123124x123214_download"
	FrpWebsocketPath  = "/benchmark_1000MB.bin"
	PreDomain         = ""
)

type WebsocketListener struct {
	ln       net.Listener
	acceptCh chan net.Conn

	server    *http.Server
	httpMutex *http.ServeMux
}

type CloseNotifyConn struct {
	net.Conn

	// 1 means closed
	closeFlag int32

	closeFn func()
}

func (cc *CloseNotifyConn) Close() (err error) {
	pflag := atomic.SwapInt32(&cc.closeFlag, 1)
	if pflag == 0 {
		err = cc.Close()
		if cc.closeFn != nil {
			cc.closeFn()
		}
	}
	return
}

// closeFn will be only called once
func WrapCloseNotifyConn(c net.Conn, closeFn func()) net.Conn {
	return &CloseNotifyConn{
		Conn:    c,
		closeFn: closeFn,
	}
}

// NewWebsocketListener to handle websocket connections
// ln: tcp listener for websocket connections
func NewWebsocketListener(ln net.Listener, certFile_keyFile ...string) (wl *WebsocketListener) {
	wl = &WebsocketListener{
		acceptCh: make(chan net.Conn),
	}

	muxer := http.NewServeMux()
	muxer.HandleFunc("/test", func(resp http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(resp, "Test Start ok:%s", r.Host)
	})
	muxer.Handle(FrpWebsocketPath, websocket.Handler(func(c *websocket.Conn) {
		// host := c.Request().Header.Get("Host")
		// log.Println("Recv Host:", host)
		notifyCh := make(chan struct{})
		conn := WrapCloseNotifyConn(c, func() {
			close(notifyCh)
		})
		wl.acceptCh <- conn
		<-notifyCh
	}))
	muxer.Handle(FrpWebsocketPath2, websocket.Handler(func(c *websocket.Conn) {
		notifyCh := make(chan struct{})
		conn := WrapCloseNotifyConn(c, func() {
			close(notifyCh)
		})
		wl.acceptCh <- conn
		<-notifyCh
	}))

	if certFile_keyFile != nil && len(certFile_keyFile) == 2 {

		wl.server = &http.Server{
			Addr:    ln.Addr().String(),
			Handler: muxer,
		}

		go wl.server.ServeTLS(ln, certFile_keyFile[0], certFile_keyFile[1])
		// go wl.server.Serve(ln)
	} else {

		wl.server = &http.Server{
			Addr:    ln.Addr().String(),
			Handler: muxer,
		}
		go wl.server.Serve(ln)
	}

	return
}

func UseWebSocketListener(listenAddr string, useTls bool) (*WebsocketListener, error) {
	config := ParseURI(PASWD)

	tcpLn, err := net.Listen("tcp", listenAddr)
	if err != nil {

		return nil, err
	}
	if useTls {
		files := config.ToCertAndKey("wss-" + strings.Replace(listenAddr, ":", "-", -1))
		l := NewWebsocketListener(tcpLn, files[0], files[1])
		return l, nil
	} else {
		l := NewWebsocketListener(tcpLn)
		return l, nil

	}
}

func (p *WebsocketListener) Accept() (net.Conn, error) {
	c, ok := <-p.acceptCh
	if !ok {
		return nil, ErrWebsocketListenerClosed
	}
	return c, nil
}

func (p *WebsocketListener) Close() error {
	return p.server.Close()
}

func (p *WebsocketListener) Addr() net.Addr {
	return p.ln.Addr()
}

// addr: domain:port
func connectWebsocketServer(addr string, isSecure bool) (net.Conn, error) {
	// var addr string
	if isSecure {
		ho := strings.Split(addr, ":")
		ip, err := net.ResolveIPAddr("ip", ho[0])
		ip_addr := ip.String() + ":" + ho[1]
		if err != nil {
			return nil, err
		}
		addr = "wss://" + ip_addr + FrpWebsocketPath
	} else {
		addr = "ws://" + addr + FrpWebsocketPath
	}
	uri, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	var origin string
	if isSecure {
		ho := strings.Split(uri.Host, ":")
		ip, err := net.ResolveIPAddr("ip", ho[0])
		ip_addr := ip.String() + ":" + ho[1]
		if err != nil {
			return nil, err
		}
		origin = "https://" + ip_addr
	} else {
		origin = "http://" + uri.Host
	}

	cfg, err := websocket.NewConfig(addr, origin)
	cfg.TlsConfig = &tls.Config{InsecureSkipVerify: true}
	// cfg.TlsConfig = UseDefaultTlsConfig("").GenerateConfig()
	if err != nil {
		return nil, err
	}
	cfg.Dialer = &net.Dialer{
		Timeout: 10 * time.Second,
	}
	if PreDomain != "" {
		cfg.Header.Set("Host", PreDomain)
	}

	conn, err := websocket.DialConfig(cfg)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// func ConnectWsAndFirstBuf(addr string, first ...[]byte) (con net.Conn, firstReply []byte, err error) {
// 	con, err = ConnectWebsocketServer(addr, false)
// 	if err != nil {
// 		return
// 	}
// 	if first != nil {
// 		// var n int
// 		_, err = con.Write(first[0])
// 		if err != nil {
// 			con.Close()
// 			return
// 		}
// 		buf := make([]byte, 20)
// 		if nr, eerr := con.Read(buf); err != nil {
// 			con.Close()
// 			return con, firstReply, eerr
// 		} else {
// 			firstReply = buf[:nr]
// 			return
// 		}
// 	}
// 	return
// }

// func ConnectWssAndFirstBuf(addr string, first ...[]byte) (con net.Conn, firstReply []byte, err error) {
// 	con, err = ConnectWebsocketServer(addr, true)
// 	if err != nil {
// 		return
// 	}
// 	if first != nil {
// 		// var n int
// 		_, err = con.Write(first[0])
// 		if err != nil {
// 			con.Close()
// 			return
// 		}
// 		buf := make([]byte, 20)
// 		if nr, eerr := con.Read(buf); err != nil {
// 			con.Close()
// 			return con, firstReply, eerr
// 		} else {
// 			firstReply = buf[:nr]
// 			return
// 		}
// 	}
// 	return
// }

// func dialWithDialer(dialer *net.Dialer, config *Config) (conn net.Conn, err error) {
// 	switch config.Location.Scheme {
// 	case "ws":
// 		conn, err = dialer.Dial("tcp", parseAuthority(config.Location))

// 	case "wss":
// 		config.TlsConfig = &tls.Config{
// 			InsecureSkipVerify: true,
// 		}
// 		conn, err = tls.DialWithDialer(dialer, "tcp", parseAuthority(config.Location), config.TlsConfig)

// 	default:
// 		err = ErrBadScheme
// 	}
// 	return
// }
