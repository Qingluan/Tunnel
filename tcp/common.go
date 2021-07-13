package tcp

import (
	"io"
	"net"
)

// Base -=-------------------------------------------------------------------------------------
func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}
	// fallback to standard io.CopyBuffer
	buf := make([]byte, 8192)
	return io.CopyBuffer(dst, src, buf)
}

func Pipe(p1, p2 net.Conn) {
	// start tunnel & wait for tunnel termination
	streamCopy := func(dst io.Writer, src io.ReadCloser, fr, to net.Addr) {
		// startAt := time.Now()
		Copy(dst, src)
		p1.Close()
		p2.Close()
	}
	go streamCopy(p1, p2, p2.RemoteAddr(), p1.RemoteAddr())
	streamCopy(p2, p1, p1.RemoteAddr(), p2.RemoteAddr())
}

func TcpEnd(toDst string, conFrom net.Conn, midledata []byte) error {
	defer conFrom.Close()
	remoteConn, err := ConnectTo("tcp", toDst)
	if err != nil {
		return err
	}
	if midledata != nil {
		_, err = conFrom.Write(midledata)
	}
	if err != nil {
		return err
	}
	Pipe(conFrom, remoteConn)
	return nil
}

func ConnectTo(protocl, dst string) (con net.Conn, err error) {
	switch protocl {
	case "tcp":
		con, err = net.Dial("tcp", dst)
	case "unix":
		con, err = net.Dial("unix", dst)
	case "tls":
		con, err = UseDefaultTlsConfig(dst).WithConn()
	case "wss":
		con, err = connectWebsocketServer(dst, true)
	case "ws":
		con, err = connectWebsocketServer(dst, false)
	}
	return
}

func AddConfigs(old []interface{}, new interface{}) (out []interface{}) {
	out = append(old, new)
	return
}

func T(i ...interface{}) {
	// fmt.Println("ok :", fmt.Sprint(i...))
}
