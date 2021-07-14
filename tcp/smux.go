package tcp

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/xtaci/smux"
)

var (
	allcon               = make(map[string]*smux.Session)
	GlobalSessions       = make(chan *smux.Session, 512)
	GlobalAverageRunning = false
	locker               = sync.RWMutex{}
)

type scavengeSession struct {
	session *smux.Session
	ts      time.Time
}

func Average(ch chan *smux.Session, ttl int) net.Conn {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	var sessionList []scavengeSession
	for {
		select {
		case sess := <-ch:
			sessionList = append(sessionList, scavengeSession{sess, time.Now()})
			log.Println("session marked as expired", sess.RemoteAddr())
		case <-tick.C:
			var newList []scavengeSession
			for k := range sessionList {
				s := sessionList[k]
				if s.session.NumStreams() == 0 || s.session.IsClosed() {
					log.Println("session normally closed", s.session.RemoteAddr())
					s.session.Close()
				} else if ttl >= 0 && time.Since(s.ts) >= time.Duration(ttl)*time.Second {
					log.Println("session reached scavenge ttl", s.session.RemoteAddr())
					s.session.Close()
				} else {
					newList = append(newList, sessionList[k])
				}
			}
			sessionList = newList
		}
	}
}

func UpdateSession(proxy string, con *smux.Session) {
	locker.Lock()
	defer locker.Unlock()
	allcon[proxy] = con
}

func createSession(protocol, proxy string) (*smux.Session, error) {
	for {
		con, _, err := ExpressDial(proxy, protocol, false)
		if err != nil {
			log.Println(fmt.Sprintf(protocol+" connect error \"%s\" , wait 1s to reconnect!!!", proxy))
			time.Sleep(1 * time.Second)
			continue
		}
		for i := 0; i < 3; i++ {
			session, err := smux.Client(con, nil)
			if err != nil {
				log.Println(fmt.Sprintf("smux session connect error \"%s\" , wait 1s to reconnect!!!", proxy))
				time.Sleep(1 * time.Second)
			} else {
				return session, nil
			}
		}
		con.Close()
	}

}

func WithASession(protocol, proxy string, chs chan *smux.Session) (*smux.Session, error) {
	if !GlobalAverageRunning {
		go Average(GlobalSessions, 600)
	}
	if sess, ok := allcon[proxy]; ok {
		if !sess.IsClosed() {
			return sess, nil
		} else {
			sess, err := createSession(protocol, proxy)
			chs <- sess
			if err != nil {
				return nil, err
			}
			UpdateSession(proxy, sess)
			return sess, nil
		}
	} else {
		sess, _ := createSession(protocol, proxy)
		chs <- sess
		UpdateSession(proxy, sess)
		return sess, nil
	}
}
