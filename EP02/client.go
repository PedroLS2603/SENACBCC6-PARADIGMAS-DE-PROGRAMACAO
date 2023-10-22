package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
		return
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	fmt.Println("Connected!")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("conexão encerrada")
		done <- true
	}()

	mustCopy(conn, os.Stdin)
	conn.Close()
	<-done
}
