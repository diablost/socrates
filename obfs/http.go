package obfs

import (
	"bytes"
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	mathRand "math/rand"
)


// HTTPRequest obfs http header
//func HTTPRequest(w net.Conn, data []byte, localConn net.Conn) error {
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
		"Content-Length: %d\r\n\r\n"

	httpRequestTemplate = fmt.Sprintf(httpRequestTemplate,
		"/", "www.bing.com", rand.Int()%51, rand.Int()%2, base64.URLEncoding.EncodeToString(c),
		int64(len(data)))
	var buf bytes.Buffer
	buf.WriteString(httpRequestTemplate)
	buf.Write(data)

	_, e := w.Write(buf.Bytes())

	return e
}
