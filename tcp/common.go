package tcp

import (
	"io"
	"log"
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

func TcpConnected(toDst string, conFrom net.Conn, midledata []byte) {
	defer conFrom.Close()
	remoteConn, err := net.Dial("tcp", toDst)
	if err != nil {
		log.Println(err)
		return
	}
	if midledata != nil {
		_, err = conFrom.Write(midledata)
	}
	Pipe(conFrom, remoteConn)
}
