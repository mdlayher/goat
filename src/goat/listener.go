package goat

import (
	"fmt"
	"net"
)

type Listener interface {
	Listen(port string)
}
type HttpListener struct {
}

func (h HttpListener) Listen(port string) {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := l.Accept()
		handle := new(HttpConnHandler)
		if err != nil {
			fmt.Println(err)
		}
		go handle.Handle(conn)
	}
}

type UdpListener struct {
}

func (u UdpListener) Listen(port string) {

}
