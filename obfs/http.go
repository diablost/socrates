package obfs

import (
	//"bufio"
	"bytes"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	//"net"
	"fmt"
	mathRand "math/rand"
	"io"
	//"time"

)

func Test(){
	fmt.Println("1111111111111")
}

// HTTPRequest obfs http header 
func HTTPRequest(w io.Writer, data []byte) error {
	b := make([]byte, 8)
	cryptoRand.Read(b)
	rand := mathRand.New(mathRand.NewSource(int64(binary.BigEndian.Uint64(b))))

	c := make([]byte, 16)
	rand.Read(c)

	httpRequestTemplate := "GET %s HTTP/1.1\r\n" +
    	"Host: %s\r\n" +
    	"User-Agent: curl/7.%d.%d\r\n" +
    	"Upgrade: websocket\r\n" +
    	"Connection: Upgrade\r\n" +
    	"Sec-WebSocket-Key: %s\r\n" +
    	"Content-Length: %v\r\n\r\n"

	httpRequestTemplate = fmt.Sprintf(httpRequestTemplate, 
		"", "www.bing.com", rand.Int()%51, rand.Int()%2, base64.URLEncoding.EncodeToString(c),
		int64(len(data)))
    fmt.Println(httpRequestTemplate)
	var buf bytes.Buffer
	buf.WriteString(httpRequestTemplate)
	buf.Write(data)
	//fmt.Println(buf)

	m, e := w.Write(buf.Bytes())
	fmt.Println("---------------write size,data:", m, data)


	//n := len(data)
	//req, _ := http.NewRequest("GET", fmt.Sprintf("http://%s/", "bing.com"), bytes.NewBuffer(data))
	//req.Header.Set("User-Agent", fmt.Sprintf("curl/7.%d.%d", rand.Int()%51, rand.Int()%2))
	//req.Header.Set("Upgrade", "websocket")
	//req.Header.Set("Connection", "Upgrade")
	//req.Header.Set("Sec-WebSocket-Key", base64.URLEncoding.EncodeToString(c))
	//req.ContentLength = int64(n)
	////req.Host = "bing.com"
	////fmt.Println(req)
	//err := req.Write(w)

	return e
}