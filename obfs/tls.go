package obfs

import (
	"fmt"
	"crypto/tls"
	"crypto/x509"
	//"io"

	"log"

)

func Connect(tgt string) *tls.Conn {
	cert, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
    if err != nil {
		log.Printf("server: loadkeys: %s", err)
    }
    config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", tgt, &config)
	if err != nil {
        //log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
    //log.Println("client: connected to: ", conn.RemoteAddr())

	state := conn.ConnectionState()
    for _, v := range state.PeerCertificates {
        fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
        fmt.Println(v.Subject)
    }
    //log.Println("client: handshake: ", state.HandshakeComplete)
	//log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
	
	return conn

    //message := "Hello\n"
    //n, err := io.WriteString(conn, message)
    //if err != nil {
    //    log.Fatalf("client: write: %s", err)
    //}
    //log.Printf("client: wrote %q (%d bytes)", message, n)

    //reply := make([]byte, 256)
    //n, err = conn.Read(reply)
    //log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
    //log.Print("client: exiting")

}


//func TlsServerListen() {
//	cert, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
//    if err != nil {
//        log.Fatalf("server: loadkeys: %s", err)
//    }
//    config := tls.Config{Certificates: []tls.Certificate{cert}}
//    config.Rand = rand.Reader
//    service := "0.0.0.0:8000"
//    listener, err := tls.Listen("tcp", service, &config)
//    if err != nil {
//        log.Fatalf("server: listen: %s", err)
//    }
//    log.Print("server: listening")
//    for {
//        conn, err := listener.Accept()
//        if err != nil {
//            log.Printf("server: accept: %s", err)
//            break
//        }
//        defer conn.Close()
//        log.Printf("server: accepted from %s", conn.RemoteAddr())
//        tlscon, ok := conn.(*tls.Conn)
//        if ok {
//            log.Print("ok=true")
//            state := tlscon.ConnectionState()
//            for _, v := range state.PeerCertificates {
//                log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
//            }
//        }
//        go handleClient(conn)
//    }
//}
//
//func handler() {
//
//}