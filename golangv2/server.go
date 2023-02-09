package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
)

type Server struct {
	// server listener address
	listenAddr string

	// quit trigger for golang channel
	quitChan chan struct{}

	// listener
	listener net.Listener
}

// function to create a new connection
func NewServer(addr string) *Server {
	return &Server{
		listenAddr: addr,
		quitChan:   make(chan struct{}),
	}
}

func (s *Server) StartServer() error {
	lst, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	defer lst.Close()
	s.listener = lst

	go s.acceptConnection()

	<-s.quitChan

	return nil
}

func (s *Server) acceptConnection() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("error when accepting connection: %s\n", err.Error())
			continue
		}

		fmt.Printf("accepting connection from %s\n", conn.RemoteAddr())

		go s.readConnection(conn)
	}
}

func (s *Server) readConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	var buff bytes.Buffer

	for {
		b, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("read error: %s\n", err)
			continue
		}

		buff.Write(b)
		if !isPrefix {
			break
		}
	}

	write(conn, buff.String())
}

func write(conn net.Conn, message string) {
	parsedIncoming := strings.Split(message, " ")
	fmt.Printf("\n\nMethod\t\t:%s\nPath\t\t:%s\nProtocol\t:%s\n\n", parsedIncoming[0], parsedIncoming[1], parsedIncoming[2])
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	fmt.Fprint(conn, "<p>anda mengakses ", parsedIncoming[1], "</p>")
}
