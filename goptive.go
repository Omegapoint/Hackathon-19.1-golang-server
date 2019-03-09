package main

import (
  "os"
  "fmt"
  "net"
  "bufio"
  "time"
)

const ResponseOK string = `HTTP/1.1 200 OK
Server: Goptive
Content-Type: text/html; charset=utf-8

<html>
<head>
<title>Goptive Non-blocking server</title>
</head>
<body>
<h1>Seems OK!</h1>
</html>
`

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
            go handleHttpConnection(conn)
        }
    }
}

func handleHttpConnection(conn net.Conn) {

    <-time.NewTimer(1 * time.Second).C

    log("*** Connection established ***")

    scanner := bufio.NewScanner(bufio.NewReader(conn))
    for scanner.Scan() {
        line := scanner.Text()
        fmt.Printf("> %s\n", line)
        if line == "" {
            break
        }
    }

    if err := scanner.Err(); err != nil {
        log(err.Error())
        conn.Close()
        return
    }

    conn.Write([]byte(ResponseOK))

    log("Closing connection")

    conn.Close()
}

func log(message string) {
    fmt.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), message)
}
