package main

import (
	"github.com/docker/go-p9p"
	"github.com/porterjamesj/gist9p"
	"golang.org/x/net/context"
	"log"
	"net"
	"os"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("done goofed")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("done goofed")
		}
		log.Println("serving connection")
		githubToken := os.Args[1]
		session := gist9p.NewGistSession(githubToken)
		go p9p.ServeConn(context.Background(), conn, p9p.Dispatch(session))
	}
}
