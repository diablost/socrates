package socrates

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// Create a HTTP server listening on addr and proxy to server.
// http proxy listen 18080
// redirect to socrates socks5 proxy(1080)
func HTTP2socksProxy(addr string) {
	l, ok := net.Listen("tcp", addr)
	if ok != nil {
		logf("Http proxy listen %v failed.", addr)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			logf("Http proxy accept %v failed.", addr)
			continue
		}
		go httpProxy(c)
	}
}

func httpProxy(conn net.Conn) {
	defer conn.Close()

	var cr io.Reader = bufio.NewReader(conn)
	buff, err := cr.(*bufio.Reader).Peek(3)
	if err != nil {
		logf("Read buffer from conn failed.")
		return
	}
	var peer net.Conn

	if bytes.Equal(buff, []byte("CON")) {
		logf("accept request conn http request.")
		peer, err = buildHttpConnProxy(cr, conn)
	} else if buff[0] >= 'A' && buff[0] <= 'Z' {
		logf("accept request http request.")
		buildHTTPProxy(cr, conn)
	} else if buff[0] == 5 {
		// sock5走sock5配置代理
		// 0x81模式的附带bit not加密
		// 0x81模式是自己goproxy连接goproxy的模式
		logf("accept request socks5 request.")
		//cr, c, peer, direct, err = buildSocks5Proxy(cr, c)
	} else {
		logf("unknown protocol:%v", string(buff))
		return
	}

	if err != nil {
		log.Println("build proxy failed:", err)
		return
	}
	if peer == nil {
		return
	}

	defer peer.Close()

	go func() {
		defer peer.Close()
		defer conn.Close()
		io.Copy(conn, peer)
	}()
	io.Copy(peer, cr)
}

// http method connect
func buildHttpConnProxy(cr io.Reader, c net.Conn) (peer net.Conn, err error) {
	buff, err := cr.(*bufio.Reader).ReadSlice(' ')
	if err != nil {
		logf("Conn.Read failed: %v", err)
		return
	}
	if !bytes.Equal(buff, []byte("CONNECT ")) {
		logf("Protocol error: %v", string(buff))
		return
	}
	buff, err = cr.(*bufio.Reader).ReadSlice(':')
	if err != nil {
		logf("Conn.Read failed: %v", err)
		return
	}
	if len(buff) <= 1 {
		logf("CONNECT protocol error: host not found")
		return
	}
	domain := string(buff[:len(buff)-1])
	buff, err = cr.(*bufio.Reader).ReadSlice(' ')
	if err != nil {
		logf("Conn.Read failed: %v", err)
		return
	}
	if len(buff) <= 1 {
		logf("CONNECT protocol error: port not found")
		return
	}
	_port := string(buff[:len(buff)-1])
	port, err := strconv.Atoi(_port)
	if err != nil {
		logf("CONNECT protocol error: port format error: %v %v", err, _port)
		return
	}
	for {
		if buff, _, err = cr.(*bufio.Reader).ReadLine(); err != nil {
			logf("Conn.Read failed: %s", err)
			return
		} else if len(buff) == 0 {
			break
		}
	}
	logf("http connect proxy connect socks5 domain:%s, port:%d", domain, port)
	peer, err = connectLocalSocks5(domain, uint16(port))
	if err != nil {
		logf("connect socks5 failed: %v", err)
		return
	}
	if err != nil {
		logf("connect failed:%v", err)
		return
	}
	_, err = c.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
	if err != nil {
		logf("write resp failed:%v", err)
	}
	return
}

func buildHTTPProxy(cr io.Reader, conn net.Conn) (peer net.Conn, err error) {
	buff, err := cr.(*bufio.Reader).ReadBytes('\n')
	if err != nil {
		logf("Conn.Read failed: %v", err)
		return
	}

	n := len(buff)
	p1 := bytes.Index(buff[:n], []byte("http://"))
	if p1 == -1 {
		logf("http proxy format error, host not found")
		return
	}
	p2 := bytes.Index(buff[p1+7:n], []byte("/"))
	if p2 == -1 {
		logf("http proxy format error, host not finish")
		return
	}
	url := string(buff[p1+7 : p1+7+p2])
	buff = append(buff[:p1], buff[p1+7+p2:]...)
	n -= (7 + p2)
	p3 := strings.IndexByte(url, ':')
	port := 80
	_port := "80"
	domain := url
	if p3 == -1 {
		url += ":80"
	} else {
		domain = url[:p3]
		_port = string(url[p3+1:])
		port, err = strconv.Atoi(_port)
		if err != nil {
			logf("http port format error: %v", _port)
			return
		}
	}
	logf("http proxy connect socks5 domain:%v, port:%v", domain, port)

	peer, err = connectLocalSocks5(domain, uint16(port))
	if err != nil {
		logf("connect socks5 failed: %v", err)
		return
	}
	_, err = peer.Write(buff[:n])
	if err != nil {
		peer.Close()
		peer = nil
		logf("Conn.Write failed: %v", err)
		return
	}
	return peer, err
}

func connectLocalSocks5(domain string, port uint16) (net.Conn, error) {
	socks5 := "127.0.0.1:1088"
	c2, err := net.Dial("tcp", socks5)
	if err != nil {
		logf("Conn.Dial failed:%v, %v", err, socks5)
		return nil, err
	}
	c2.SetDeadline(time.Now().Add(10 * time.Second))
	// 此处要求0,0x81两种模式
	// goproxy会返回0x81模式
	// 普通socks5则只会返回0
	c2.Write([]byte{5, 2, 0, 0x81})
	resp := make([]byte, 2)
	n, err := c2.Read(resp)
	if err != nil {
		logf("Conn.Read failed:%v", err)
		return nil, err
	}
	if n != 2 {
		logf("socks5 response error:%v", resp)
		return nil, errors.New("socks5_error")
	}
	method := resp[1]
	if method != 0 && method != 0x81 {
		logf("socks5 not support 'NO AUTHENTICATION REQUIRED'")
		return nil, errors.New("socks5_error")
	}
	send := make([]byte, 0, 512)
	send = append(send, []byte{5, 1, 0, 3, byte(len(domain))}...)
	if method == 0 {
		send = append(send, []byte(domain)...)
	} else {
		edomain := []byte(domain)
		for i, c := range edomain {
			edomain[i] = ^c
		}
		send = append(send, edomain...)
	}
	send = append(send, byte(port>>8))
	send = append(send, byte(port&0xff))
	_, err = c2.Write(send)
	if err != nil {
		logf("Conn.Write failed:%v", err)
		return nil, err
	}
	n, err = c2.Read(send[0:10])
	if err != nil {
		logf("Conn.Read failed:%v", err)
		return nil, err
	}
	if send[1] != 0 {
		switch send[1] {
		case 1:
			logf("socks5 general SOCKS server failure")
		case 2:
			logf("socks5 connection not allowed by ruleset")
		case 3:
			logf("socks5 Network unreachable")
		case 4:
			logf("socks5 Host unreachable")
		case 5:
			logf("socks5 Connection refused")
		case 6:
			logf("socks5 TTL expired")
		case 7:
			logf("socks5 Command not supported")
		case 8:
			logf("socks5 Address type not supported")
		default:
			logf("socks5 Unknown eerror:%v", send[1])
		}
		return nil, errors.New("socks5_error")
	}
	c2.SetDeadline(time.Time{})
	if method == 0 {
		return c2, nil
	} else {
		return c2, nil
	}
}
