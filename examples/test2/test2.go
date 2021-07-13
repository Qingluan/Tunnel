package main

import (
	"flag"
	"fmt"

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

	tcp.ExpressSocks5ListenerTo("127.0.0.1:1080", "127.0.0.1:12345", P, ifmulti)
	// one := func(f string) {
	// 	con, _, err := tcp.ExpressDial(P, "127.0.0.1:12345", ifmulti)
	// 	fmt.Println("->")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	con.Write([]byte("ok................" + f))
	// 	b, err := ioutil.ReadAll(con)

	// 	fmt.Println("R:", string(b), err)
	// 	con.Close()
	// }

	// go one("1")
	// go one("2")
	// go one("3")
	// go one("4")
	// one("end")

}
