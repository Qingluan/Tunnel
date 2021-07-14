package tcp

import (
	"fmt"
	"log"
	"net"

	"github.com/xtaci/smux"
)

type ExpressListener struct {
	protocol string
	usesmux  bool
	nowClose bool
	lastErr  error
	Sessions chan *smux.Session
	rawLst   net.Listener
	// subLst   net.Listener
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
		lst.Sessions <- smuxListner
		if err != nil {
			log.Println("smux server wrap err:", err)
			break
		}
	TUNNEL:
		for {
			if lst.nowClose {

				log.Println("Close session listener")
				break
			}
			if smuxListner.IsClosed() {
				log.Println("session closed wait new")

				break
			}
			// fmt.Println("session connected b")
			scon, err := smuxListner.AcceptStream()
			// fmt.Println("session connected")
			lst.lastErr = err
			if err != nil {
				log.Println("this tunel accep one sesion err:", err)
				break TUNNEL
			}

			// go func(chs chan net.Conn) {
			lst.acceptCh <- scon
			// }(lst.acceptCh)
		}

	}

}
