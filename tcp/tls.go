package tcp

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

var (
	PASWD              = "dGxzOi0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQpNSUlFckRDQ0E1U2dBd0lCQWdLQ0FRQUgrYjcvSy93Rm0yL243bGd5YmJMZWFzS1hpZndNTGc4U1BabXhqUXRjCnYrYzh6clErOVRiaThteFVwUWpIK0JjWnoweXN1UTF0NVpFb3F1WnpWRFBPWGRxaFJnaFlDdE54Z2xzNXBWaTIKdkQ5SkZJZEU3Q2ZxQlR4RXVZWlhNY0JCT2dtKzcvN0REODN3ak9QK1V5ZmxieWRsZW1Jc0Ixa0tHQWVrRThOaApRU3d3QTlSR2xBSGI3Nkp0RHJpdHEyK2dSTmxpM3NiUjRMWDR5NlJSTTJRVGwvUTRtdTgzYzc0bDB6TkZTbXN1CmN3U1lPWXRtZkxKanJQVk4xVnNtbDRpYkl4TFdpSXF4cHZPNTgrdE91RWlSNCtsTXd6RUpaSWxNQzE4c1NnNE8KVVc5TmhRZkVEYXVrK3h6TnNFOURuMnFwVW5hNG9PRnlCSUl6K2ozMlFSbkRNQTBHQ1NxR1NJYjNEUUVCQ3dVQQpNRjh4Q3pBSkJnTlZCQVlUQWtwUU1RNHdEQVlEVlFRSUV3VlViMnQ1YnpFT01Bd0dBMVVFQnhNRlZHOXJlVzh4CkVEQU9CZ05WQkFrVEIxUmhjMk5wYTI4eEVUQVBCZ05WQkJFVENERXdMVEl3TUMwME1Rc3dDUVlEVlFRS0V3SkQKYnpBZUZ3MHlNVEEzTVRJeE5qVXlNREJhRncweU1qQTNNVEl4TmpVeU1EQmFNRjh4Q3pBSkJnTlZCQVlUQWtwUQpNUTR3REFZRFZRUUlFd1ZVYjJ0NWJ6RU9NQXdHQTFVRUJ4TUZWRzlyZVc4eEVEQU9CZ05WQkFrVEIxUmhjMk5wCmEyOHhFVEFQQmdOVkJCRVRDREV3TFRJd01DMDBNUXN3Q1FZRFZRUUtFd0pEYnpDQ0FTSXdEUVlKS29aSWh2Y04KQVFFQkJRQURnZ0VQQURDQ0FRb0NnZ0VCQU0vWFp0eDhxM25SWmxCaGViRlVNNnN0WDdSTWRwTXh3SUFsVXRWQQo5YVBmL3FETUdaUDJndmpvMXpnRjZRWld6Y2MxV21KR0VEMWhRbHpQSlVhTFYvcjdMUHpiK2dpR0pJM0hIOXMyClFvTDlreFRLK1I0WTEyK2Q1bjVhd2svanhqYVpwbWhzWEwvanJBYzNhbS96eDdEK1BRUnI5bVdNNktVU1dQUVEKcTlqL1Y4R2VaM2xrZVFNRzZuR3FqdFowRVFrUmRDUWc5NUl6VHpFb2VoVTMxUmVZbUFjOFpQaTFwRVdpVEM3ago0UWpaZFBXNzdneUdNWTVEY083cmlmaDJMRnk0ZkZkajcraXU0VGJuRGFlRUx1M1dDb0kweTJwdllZWXhrU003CmFtV2NSVXRFdVlWYkErbXluQ2c2SitXaTM1VFk1aWg2K084R2FWcEJhQlVtUlBzQ0F3RUFBYU55TUhBd0RnWUQKVlIwUEFRSC9CQVFEQWdLa01CMEdBMVVkSlFRV01CUUdDQ3NHQVFVRkJ3TUJCZ2dyQmdFRkJRY0RBakFQQmdOVgpIUk1CQWY4RUJUQURBUUgvTUIwR0ExVWREZ1FXQkJRaGE0S3FUOTlDZzRORjVVeExTYktwNzRXTStEQVBCZ05WCkhSRUVDREFHaHdSL0FBQUJNQTBHQ1NxR1NJYjNEUUVCQ3dVQUE0SUJBUUF1a0FWaEdlNXZhb3l0VVFidXBPN1IKanhXMVBFTXVydHd4OXJWV3BIaVdVYXFQWUFldnlyZVdxYzlWb2V5V3BxdGhIYnF0WUtKZ0pVeXFNZzJGODNVcQptaGttckFna2ZQZTBFVUJkRHhZTExuMTE2Vy9QVHhUMkI0U0UwbVZ4bXBqcHQ0Wktib2VSTHRjZDNqSW1WR2k0CnRDcmZZU2FTdmgyaVBmRFRpelE2Vy95bzlSaVVqSkExdkFhaXNRbnVuL2lJWFd5bTNjZ00xMGtod2t2Q2ZHUGMKTUpuWnM5ODcvalhGcDJtL1p6ZFdPMjdpSlNuQkZzOGVza3BGMFoxRWJ4aE1uY0lRLzFoVDZKU081RkZRR0tPNAp5c0tLeEFINnFrWEd6bnV3Z0VxS0VRUUFhTWVzYWI0d0RQS1JEMkFJbWNrZitwQ0IwYnlIMjVzSmltZHBLSkRNCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KPFNFUD4tLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS0KTUlJRXZnSUJBREFOQmdrcWhraUc5dzBCQVFFRkFBU0NCS2d3Z2dTa0FnRUFBb0lCQVFEUDEyYmNmS3Q1MFdaUQpZWG14VkRPckxWKzBUSGFUTWNDQUpWTFZRUFdqMy82Z3pCbVQ5b0w0Nk5jNEJla0dWczNITlZwaVJoQTlZVUpjCnp5VkdpMWY2K3l6ODIvb0loaVNOeHgvYk5rS0MvWk1VeXZrZUdOZHZuZVorV3NKUDQ4WTJtYVpvYkZ5LzQ2d0gKTjJwdjg4ZXcvajBFYS9abGpPaWxFbGowRUt2WS8xZkJubWQ1WkhrREJ1cHhxbzdXZEJFSkVYUWtJUGVTTTA4eApLSG9WTjlVWG1KZ0hQR1Q0dGFSRm9rd3U0K0VJMlhUMXUrNE1oakdPUTNEdTY0bjRkaXhjdUh4WFkrL29ydUUyCjV3Mm5oQzd0MWdxQ05NdHFiMkdHTVpFak8ycGxuRVZMUkxtRld3UHBzcHdvT2lmbG90K1UyT1lvZXZqdkJtbGEKUVdnVkprVDdBZ01CQUFFQ2dnRUFHTFlUbE1OOTM4MzF6RGpLcnpyRHFlYUxLblNzNUhOVjMrYVFJcXJHaWVWbgo4TlpUZERRcS9rMHd5WWhxWEVhVjBNbVhKWGdQRmowcUZyN0tQOGp4ZlBYQU01QXoraTk0eVNLVHR3eU1lK2FLClJqNmJ2K2hVTlFFaEZvVFQ5WXV2Vlc2UFNuc1h5L2pWeTBYYmNkUWFPZ0VBWURZMTVYMExiQlR5WHBKYzNEL0kKMStTNlZQMXFoNFhmdjl2RE9YMnpuajFqNXNqd04vb3I0MHA0NmRENFc3eEVOdktrQUF4VG9lL0JkdXp4MlduRQoza2FWcTZKd21xU3J2NW1kNEg5ZTBHdU81N0QyV2I0VEZqZVBxVnNTOHlSOFVPNE5vSW9yVHVKUjg5TWh6NWlQCjRnRmJjKzc1QUhxbWI1dTJkRVZwbUhUSXlhV2xpSkh4U0ZSMmkvVXlnUUtCZ1FEY3dvblZHV1ZMSENtZmJ4Qk0KYlMxdUlXQ2lVSjBPanBnOGhnTmYxc3czWFllb1kxTndoVnB6ZktGZXJ4akVpcTdKRTZYUDhZTnFWQkR4Q0Z0QQpRU3NLS25hanF5MTZQODY0Y2ZkYlZ1aVJha1IyVUE5a08xZVNNUytXaDdJNFFEQmF0RTgyU2hNY2hjN0tBaCt0CkJiSUFtNDA5Y1BXRWZudy9sK1VsWjl0dFZRS0JnUUR4QlBGSGFUUkptWlA4ZjFWTHMxMVZ3K21GZDZDWGdXeXYKN1N1SkJNeEVlTHRCQTFzdzBHTjY5aG1GdVNTMlRmTUIxQXJ5N1lsNC82ekZKZjQ4UGZScnkrL2RUQVVjZUY1cwpjSHJPSlNhRFNnYlNDL1VXelphaDA4Vi83TmpwZzdOc202WTRhYmtGNVp3YVdwbCtCNEJ5d2FER3UwQ2VqVFN1CmFtS2ErYWRwRHdLQmdDKzVmc0txOHRNQ0M1MkVEU240a0ZHMnRZVmhtQktGa1RQRVpRWmJFVnNHeFVVSjlBZm4KVEY2Ykd1ZTFRczE2amI2Nm1LWVR1QzFNYnowc2FVM0N5T014dXVRM0pXWEFWWmhhU3ZkY0duc1ZWaWVkdHpZYgovaHFFdDV6N1NQUVFMR1B2aHhCZGd6UXZXVDBpcUdablRNd0UxTXZybDZoRkQrZFBseUkzQ2FEaEFvR0JBSjNkCmJPUW5SblpHdjZXTkNhTlIwbnFRbmU0cjB1WXBnYlVobFhoanloKzNMSCtDNFVIeDNKYjNodlFOdHJ3cUJsZFcKT2Uyd0pIRTZsa0Z6aHppU0kveFBVY2NUT2UzUjhaYjVmWlowd3VzUG03UU0zUFdZcDJkTHZIcnorWWZLR2NXegpsWVl4eVZ6UmJoUGM3OWlhdFVsMUJnZmxNb2duV1RBOHdtVmJ6SW1GQW9HQkFLLzJhUWtlWnRWSzNZSjhCajQxCjR4S296RWYwdlJyY2FCekNqNjVia3JtSE5Qdmw4ZzVKMi9PN0pSbEdOTVlrODBFWXZ1T2NxR1V4YTlJUVB6WGUKb25zeHNLN29GLzFMUkl6bEFZK1lJVCtibzV4OFdRRndZbmQ5a0Fsa3ZiKzBLT0FBMlJNT3ZIYk4yekRDOUNNWQpOSnVUQ2ZESnZyM1QwOTlDb2VXa3FjYXQKLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQpAMTI3LjAuMC4xOjEyMzQ1"
	CENTRY_SERVER_ADDR = ""
)

type Config struct {
	Server       interface{} `json:"server"`
	ServerPort   int         `json:"server_port"`
	LocalPort    int         `json:"local_port"`
	LocalAddress string      `json:"local_address"`
	Password     string      `json:"password"`
	Method       string      `json:"method"` // encryption method

	// following options are only used by server
	PortPassword map[string]string `json:"port_password"`
	Timeout      int               `json:"timeout"`

	// following options are only used by client

	// The order of servers in the client config is significant, so use array
	// instead of map to preserve the order.
	ServerPassword string `json:"server_password"`
}

type TlsConfig struct {
	Ca     x509.CertPool
	Cert   tls.Certificate
	Server string
	priKey rsa.PrivateKey
}

func (tlsConfig *TlsConfig) GenerateConfig() (config tls.Config) {
	// tlsConfig.Ca.AppendCertsFromPEM(tlsConfig.Cert.)
	config = tls.Config{
		Certificates: []tls.Certificate{tlsConfig.Cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    &tlsConfig.Ca,
	}
	config.Rand = rand.Reader
	return
}

func SetPaswd(pas string) {
	PASWD = pas
}

// func (tlsConfig *TlsConfig) WithConn() (conn *tls.Conn, err error) {
// 	config := tls.Config{
// 		Certificates:       []tls.Certificate{tlsConfig.Cert},
// 		InsecureSkipVerify: true,
// 	}
// 	serverAddress := tlsConfig.Server
// 	conn, err = tls.Dial("tcp", serverAddress, &config)
// 	if err != nil {
// 		log.Println("tls connect:", serverAddress)
// 		return
// 	}
// 	state := conn.ConnectionState()
// 	// for _, v := range state.PeerCertificates {
// 	// 	log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
// 	// }
// 	if !state.HandshakeComplete {
// 		return nil, errors.New("Not TLS Handleshare finished!!")
// 	}

// 	return
// }

func (tlsConfig *TlsConfig) WithTlsListener() (listenr net.Listener, err error) {
	config := tlsConfig.GenerateConfig()
	listenr, err = tls.Listen("tcp", tlsConfig.Server, &config)
	return
}

func (config *Config) ToTlsConfig() (tlsConfig *TlsConfig, err error) {
	if config.Method != "tls" {
		return
	}
	tlsConfig = new(TlsConfig)
	tlsConfig.Server = fmt.Sprintf("%s:%d", config.Server.(string), config.ServerPort)

	// ColorL("raw:", config.Password)
	pems := strings.SplitN(config.Password, "<SEP>", 2)

	pemBlock := []byte(strings.TrimSpace(pems[0]))
	keyBlock := []byte(strings.TrimSpace(pems[1]))

	// preName := ".tmp." + strconv.Itoa(random.Int())

	// ioutil.WriteFile(preName+".pem", pemBlock, os.ModePerm)
	// ioutil.WriteFile(preName+".key", keyBlock, os.ModePerm)
	// defer os.Remove(preName + ".pem")
	// defer os.Remove(preName + ".key")
	// crtx, err2 := x509.ParseCertificate(pemBlock.Bytes)
	// crt, err2 := tls.LoadX509KeyPair(preName+".pem", preName+".key")
	crt, err2 := tls.X509KeyPair(pemBlock, keyBlock)
	if err2 != nil {
		return nil, err2
	}

	tlsConfig.Cert = crt
	// tlsConfig.priKey = *key
	tlsConfig.Ca = *x509.NewCertPool()
	tlsConfig.Ca.AppendCertsFromPEM(pemBlock)
	return
}

func (config *Config) ToString() string {
	return fmt.Sprintf("%s:%d", config.Server.(string), config.ServerPort)
}

func (config *Config) ToFile(dst string) (err error) {
	if f, err := json.Marshal(config); err == nil {
		if err := ioutil.WriteFile(dst, f, 0644); err != nil {
			return err
		}
	} else {
		return err
	}
	return
}

func (config *Config) ToJson() string {
	if f, err := json.Marshal(config); err == nil {
		return string(f)
	}
	return ""
}

func (config *Config) ToUri() string {
	base := fmt.Sprintf("%s:%s@%s:%d", config.Method, config.Password, config.Server.(string), config.ServerPort)
	encoder := base64.StdEncoding.EncodeToString([]byte(base))
	return fmt.Sprintf("ss://%s", encoder)
}

// GetServerArray get server
func (config *Config) GetServerArray() []string {
	// Specifying multiple servers in the "server" options is deprecated.
	// But for backward compatibility, keep this.
	if config.Server == nil {
		return nil
	}
	single, ok := config.Server.(string)
	if ok {
		return []string{single}
	}
	arr, ok := config.Server.([]interface{})
	if ok {
		serverArr := make([]string, len(arr), len(arr))
		for i, s := range arr {
			serverArr[i], ok = s.(string)
			if !ok {
				goto typeError
			}
		}
		return serverArr
	}
typeError:
	panic(fmt.Sprintf("Config.Server type error %v", reflect.TypeOf(config.Server)))
}

// ParseConfig parse path to json
func ParseConfig(path string) (config *Config, err error) {
	file, err := os.Open(path) // For read access.
	if err != nil {
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	config = &Config{}
	if err = json.Unmarshal(data, config); err != nil {
		return nil, err
	}
	config.Timeout = 5000
	return
}

// UpdateConfig : Useful for command line to override options specified in config file  Debug is not updated.
func UpdateConfig(old, new *Config) {
	// Using reflection here is not necessary, but it's a good exercise.
	// For more information on reflections in Go, read "The Laws of Reflection"
	// http://golang.org/doc/articles/laws_of_reflection.html
	newVal := reflect.ValueOf(new).Elem()
	oldVal := reflect.ValueOf(old).Elem()

	// typeOfT := newVal.Type()
	for i := 0; i < newVal.NumField(); i++ {
		newField := newVal.Field(i)
		oldField := oldVal.Field(i)
		// log.Printf("%d: %s %s = %v\n", i,
		// typeOfT.Field(i).Name, newField.Type(), newField.Interface())
		switch newField.Kind() {
		case reflect.Interface:
			if fmt.Sprintf("%v", newField.Interface()) != "" {
				oldField.Set(newField)
			}
		case reflect.String:
			s := newField.String()
			if s != "" {
				oldField.SetString(s)
			}
		case reflect.Int:
			i := newField.Int()
			if i != 0 {
				oldField.SetInt(i)
			}
		}
	}

	old.Timeout = new.Timeout

}

func GetMainDomain(urlOrHost string) string {
	host := urlOrHost
	if strings.HasPrefix(urlOrHost, "http") {
		u, _ := url.Parse(urlOrHost)
		host = u.Host
	}
	dotCount := strings.Count(host, ".")
	if dotCount > 1 {
		return strings.Join(strings.Split(host, ".")[dotCount-1:], ".")
	} else {
		return host
	}
}

func ParseURI(u string) (config *Config) {
	config = new(Config)
	_, err := parseURI(u, config)
	if err != nil {
		log.Fatal("parse config:", err)
	}
	return
}

func parseURI(u string, cfg *Config) (string, error) {
	if u == "" {
		return "", nil
	}
	invalidURI := errors.New("invalid URI")
	// ss://base64(method:password)@host:port
	// ss://base64(method:password@host:port)
	u = strings.TrimLeft(u, "ss://")
	i := strings.IndexRune(u, '@')
	var headParts, tailParts [][]byte
	if i == -1 {
		dat, err := base64.StdEncoding.DecodeString(u)
		if err != nil {
			return "", err
		}
		parts := bytes.Split(dat, []byte("@"))
		if len(parts) < 2 {
			parts = append(parts, []byte("127.0.0.1:12345"))
			// return "", invalidURI
		}
		headParts = bytes.SplitN(parts[0], []byte(":"), 2)
		tailParts = bytes.SplitN(parts[1], []byte(":"), 2)

	} else {
		if i+1 >= len(u) {
			return "", invalidURI
		}
		tailParts = bytes.SplitN([]byte(u[i+1:]), []byte(":"), 2)
		dat, err := base64.StdEncoding.DecodeString(u[:i])
		if err != nil {
			return "", err
		}
		headParts = bytes.SplitN(dat, []byte(":"), 2)
	}
	if len(headParts) != 2 {
		return "", invalidURI
	}

	if len(tailParts) != 2 {
		return "", invalidURI
	}
	cfg.Method = string(headParts[0])

	cfg.Password = string(headParts[1])
	p, e := strconv.Atoi(string(tailParts[1]))
	if e != nil {
		return "", e
	}
	cfg.Server = string(tailParts[0])
	cfg.ServerPort = p
	return string(tailParts[0]), nil

}

func (tlsConfig *TlsConfig) WithConn() (conn *tls.Conn, err error) {
	config := tls.Config{
		Certificates:       []tls.Certificate{tlsConfig.Cert},
		InsecureSkipVerify: true,
	}
	serverAddress := tlsConfig.Server
	conn, err = tls.Dial("tcp", tlsConfig.Server, &config)
	if err != nil {
		log.Println("tls connect:", serverAddress)
		return
	}
	state := conn.ConnectionState()
	// for _, v := range state.PeerCertificates {
	// 	log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
	// }
	if !state.HandshakeComplete {
		return nil, errors.New("Not TLS Handleshare finished!!")
	}
	return
}

func UseDefaultTlsConfig(addr string) (tlsConfig *TlsConfig) {
	config := ParseURI(PASWD)
	// log.Println("parse tls 2:", PASWD)

	ts := strings.SplitN(addr, ":", 2)
	config.Server = ts[0]
	if len(ts) == 1 {
		config.ServerPort = 80
	} else {
		config.ServerPort, _ = strconv.Atoi(ts[1])
	}
	tlsConfig, err := config.ToTlsConfig()

	// fmt.Println(tlsConfig)
	if err != nil {
		log.Fatal("Create tls config failed: ", err)
	}
	return
}

func SplitAddr(addr string) (h string, port int) {
	ts := strings.SplitN(addr, ":", 2)
	h = ts[0]
	port, _ = strconv.Atoi(ts[1])
	return
}

/*
will cache to locl file syste,with  ${tmp}/name.crt | ${tmp}/name.key
*/
func (config *Config) ToCertAndKey(name string) (crt_key_paths []string) {
	tmp := os.TempDir()
	preName := filepath.Join(tmp, name)
	crt_key_paths = append(crt_key_paths, preName+".crt")
	crt_key_paths = append(crt_key_paths, preName+".key")
	if _, err := os.Stat(preName + ".crt"); err != nil {
		pems := strings.SplitN(config.Password, "<SEP>", 2)
		ioutil.WriteFile(preName+".crt", []byte(pems[0]), os.ModePerm)
		ioutil.WriteFile(preName+".key", []byte(pems[1]), os.ModePerm)
	}
	return
}

// func initiTlsConnection(tlsConfig *TlsConfig, firstData ...[]byte) (con net.Conn, roleReply []byte, err error) {
// 	con, err = tlsConfig.WithConn()
// 	if err != nil {
// 		log.Println("Err Conn remote tls server:", err)
// 		return
// 	}
// 	if firstData != nil {
// 		con.Write(firstData[0])
// 		buf := make([]byte, 20)
// 		n, err := con.Read(buf)
// 		if err != nil {
// 			con.Close()
// 			log.Println("Err in role confirm reply:", err)
// 			return nil, []byte{}, err
// 		}
// 		roleReply = buf[:n]
// 	}
// 	return
// }
