// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.

package main

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"flag"
	"fmt"

	"github.com/Qingluan/Tunnel/config"
)

var (
	host = flag.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
	isCA = flag.Bool("ca", false, "whether this cert should be its own Certificate Authority")
	key  = flag.String("key", "", "set some key for generate certificate")
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func main() {
	flag.Parse()

	// config.SetKey(*key)
	key, _ := config.CreateCertificate("127.0.0.1:12345", true)

	// e.Write([]byte(pem))
	fmt.Println(key)
	// fmt.Printf("pem: %x\n", md5.New().Sum([]byte(pem))[:16])
	// fmt.Printf("key: %x\n", md5.New().Sum([]byte(key))[:16])
}
