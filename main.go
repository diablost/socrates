package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/sic-project/socrates/core"
	"github.com/sic-project/socrates/socks"
	sct "github.com/sic-project/socrates/socrates"
)

var config struct {
	Verbose    bool
	UDPTimeout time.Duration
}

func main() {

	var flags struct {
		Client     string
		Server     string
		Cipher     string
		Key        string
		Password   string
		Keygen     int
		Socks      string
		RedirTCP   string
		RedirTCP6  string
		TCPTun     string
		UDPTun     string
		UDPSocks   bool
		AccessList string
		Discovery  bool
		Obfs       string
		HTTPProxy  string
	}

	flag.BoolVar(&sct.Verbose, "verbose", false, "verbose mode")
	flag.BoolVar(&flags.Discovery, "discovery", false, "(client-only) Proxy address discovery mode")
	flag.StringVar(&flags.Cipher, "cipher", "AEAD_CHACHA20_POLY1305", "available ciphers: "+strings.Join(core.ListCipher(), " "))
	flag.StringVar(&flags.Key, "key", "", "base64url-encoded key (derive from password if empty)")
	flag.IntVar(&flags.Keygen, "keygen", 0, "generate a base64url-encoded random key of given length in byte")
	flag.StringVar(&flags.Password, "password", "", "password")
	flag.StringVar(&flags.Server, "s", "", "server listen address or url")
	flag.StringVar(&flags.Client, "c", "", "client connect address or url")
	flag.StringVar(&flags.Socks, "socks", "", "(client-only) SOCKS listen address")
	flag.BoolVar(&flags.UDPSocks, "u", false, "(client-only) Enable UDP support for SOCKS")
	flag.StringVar(&flags.RedirTCP, "redir", "", "(client-only) redirect TCP from this address")
	flag.StringVar(&flags.RedirTCP6, "redir6", "", "(client-only) redirect TCP IPv6 from this address")
	flag.StringVar(&flags.TCPTun, "tcptun", "", "(client-only) TCP tunnel (laddr1=raddr1,laddr2=raddr2,...)")
	flag.StringVar(&flags.UDPTun, "udptun", "", "(client-only) UDP tunnel (laddr1=raddr1,laddr2=raddr2,...)")
	flag.DurationVar(&sct.UDPTimeout, "udptimeout", 5*time.Minute, "UDP tunnel timeout")
	flag.StringVar(&flags.AccessList, "accesslist", "", "(server-only) Remote access whitelist")
	flag.StringVar(&flags.Obfs, "obfs", "http", "Obfuscating by http/tls")
	flag.StringVar(&flags.HTTPProxy, "httpproxy", "", "(client-only) HTTP listen address")
	flag.Parse()

	if flags.Keygen > 0 {
		key := make([]byte, flags.Keygen)
		io.ReadFull(rand.Reader, key)
		fmt.Println(base64.URLEncoding.EncodeToString(key))
		return
	}

	if flags.Client == "" && flags.Server == "" {
		flag.Usage()
		return
	}

	var key []byte
	if flags.Key != "" {
		k, err := base64.URLEncoding.DecodeString(flags.Key)
		if err != nil {
			log.Fatal(err)
		}
		key = k
	}

	var obfs string
	if flags.Obfs == "tls" {
		obfs = "tls"
	} else {
		obfs = "http"
	}
	//sct.logf("use obfs-%v", obfs)

	var ac sct.AccessControl
	if flags.AccessList != "" {
		mapAccess := make(map[string]regexp.Regexp)
		ar := strings.Split(flags.AccessList, ",")
		for i := range ar {
			regStr := fmt.Sprintf(`%s`, ar[i])
			reg, _ := regexp.Compile(regStr)
			mapAccess[ar[i]] = *reg
		}
		ac.WhiteList = mapAccess
	}

	//err := StartParity()
	//if err != nil {
	//	log.Fatal(err)
	//}

	if flags.Client != "" { // client mode
		addr := flags.Client
		cipher := flags.Cipher
		password := flags.Password
		var err error

		if strings.HasPrefix(addr, "ss://") {
			addr, cipher, password, err = parseURL(addr)
			if err != nil {
				log.Fatal(err)
			}
		}

		udpAddr := addr

		ciph, err := core.PickCipher(cipher, key, password)
		if err != nil {
			log.Fatal(err)
		}

		if flags.UDPTun != "" {
			for _, tun := range strings.Split(flags.UDPTun, ",") {
				p := strings.Split(tun, "=")
				go sct.UdpLocal(p[0], udpAddr, p[1], ciph.PacketConn)
			}
		}

		if flags.TCPTun != "" {
			for _, tun := range strings.Split(flags.TCPTun, ",") {
				p := strings.Split(tun, "=")
				go sct.TcpTun(p[0], addr, p[1], ciph.StreamConn, obfs)
			}
		}

		// use http local proxy
		if flags.HTTPProxy != "" {
			go sct.HTTP2socksProxy(flags.HTTPProxy)
		}

		if flags.Socks != "" {
			socks.UDPEnabled = flags.UDPSocks
			go sct.SocksLocal(flags.Socks, addr, ciph.StreamConn, obfs)
			if flags.UDPSocks {
				go sct.UdpSocksLocal(flags.Socks, udpAddr, ciph.PacketConn)
			}
		}

		if flags.RedirTCP != "" {
			go sct.RedirLocal(flags.RedirTCP, addr, ciph.StreamConn)
		}

		if flags.RedirTCP6 != "" {
			go sct.Redir6Local(flags.RedirTCP6, addr, ciph.StreamConn)
		}
	}

	if flags.Server != "" { // server mode
		addr := flags.Server
		cipher := flags.Cipher
		password := flags.Password
		var err error

		if strings.HasPrefix(addr, "ss://") {
			addr, cipher, password, err = parseURL(addr)
			if err != nil {
				log.Fatal(err)
			}
		}

		udpAddr := addr

		ciph, err := core.PickCipher(cipher, key, password)
		if err != nil {
			log.Fatal(err)
		}

		go sct.UdpRemote(udpAddr, ciph.PacketConn)
		go sct.TcpRemote(addr, ciph.StreamConn, ac, obfs)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func parseURL(s string) (addr, cipher, password string, err error) {
	u, err := url.Parse(s)
	if err != nil {
		return
	}

	addr = u.Host
	if u.User != nil {
		cipher = u.User.Username()
		password, _ = u.User.Password()
	}
	return
}
