package main

import (
  "os"
  "fmt"
  "net"
  "bufio"
  "strings"
  "time"
)

func main() {

    listener, err := net.Listen("tcp", ":1337")

    if err != nil {
        log(err.Error())
        os.Exit(1)
    }

	log("=== Goptive server started ===")

    for {
		log("Waiting for connection...")
        conn, err := listener.Accept()
        if err != nil {
            log(err.Error())
        } else {
			log("Handling connection")
            go handleConnection(conn)
        }
    }
}

func handleConnection(conn net.Conn) {
	
	<-time.NewTimer(2 * time.Second).C

    message, err := bufio.NewReader(conn).ReadString('\n')

    if err != nil {
        log(err.Error())
		conn.Close()
        return
    }

    log("Message received: " + strings.TrimRight(message, "\n"))

    if strings.TrimRight(message, "\n") == "ping" {
        conn.Write([]byte("pong\n"))
    } else {
        conn.Write([]byte("does not compute\n"))
    }

	conn.Close()
}

func log(message string) {
	fmt.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), message)
}

