// +build !linux

package socrates

import "net"

func RedirLocal(addr, server string, shadow func(net.Conn, string) net.Conn) {
	logf("TCP redirect not supported")
}

func Redir6Local(addr, server string, shadow func(net.Conn, string) net.Conn) {
	logf("TCP6 redirect not supported")
}
