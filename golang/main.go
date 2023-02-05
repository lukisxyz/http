package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	listen, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Printf("listener error: %s\n", err)
		os.Exit(1)
	}
	defer listen.Close()

	log.Println("listening to port 3000")
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Printf("connection error: %s\n", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			fmt.Println("read the incoming request")

			// read the incoming request
			reader := bufio.NewReader(conn)
			var buff bytes.Buffer
			for {
				b, isPrefix, err := reader.ReadLine()
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Printf("read error: %s\n", err)
					continue
				}
				buff.Write(b)
				if !isPrefix {
					break
				}
			}

			// parse the string
			incoming := string(buff.String())
			parsedIncoming := strings.Split(incoming, " ")
			fmt.Printf("\n\nMethod\t\t:%s\nPath\t\t:%s\nProtocol\t:%s\n\n", parsedIncoming[0], parsedIncoming[1], parsedIncoming[2])

			// write response
			fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
			fmt.Fprint(conn, "Content-Type: text/html\r\n")
			fmt.Fprint(conn, "\r\n")
			fmt.Fprint(conn, "<p>anda mengakses ", parsedIncoming[1], "</p>")
		}(conn)
	}
}
