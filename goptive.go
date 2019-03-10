package main

import (
    "os"
    "os/signal"
    "fmt"
    "strings"
    "strconv"
    "net"
    "bufio"
    "time"
    "sync"
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

func dispatcher(listener net.Listener, connections chan net.Conn) {
    for {
        log("Dispatcher waiting for next connection...")
        conn, err := listener.Accept()
        if err != nil {
            log("Dispatcher:", err.Error())
            return
        } else {
            log("Dispatching connection")
            connections <- conn
        }
    }
}

func worker(id int, connections chan net.Conn, wg *sync.WaitGroup) {
    for {
        log("Worker", strconv.Itoa(id), "waiting for more jobs...")
        conn, open := <-connections
        if open {

            log("*** Worker", strconv.Itoa(id), "established connection ***")

            scanner := bufio.NewScanner(bufio.NewReader(conn))
            for scanner.Scan() {
                line := scanner.Text()
                fmt.Printf("> %s\n", line)
                if line == "" {
                    break
                }
            }

            if err := scanner.Err(); err != nil {
                log("Worker", strconv.Itoa(id), "encountered a problem:", err.Error())
                conn.Close()
                continue
            }

            time.Sleep(10*time.Second)
            conn.Write([]byte(ResponseOK))

            log("Worker", strconv.Itoa(id), "closing connection")

            conn.Close()

        } else {
            log("Worker", strconv.Itoa(id), "exiting")
            break
        }
    }
    wg.Done()
}

func log(message... string) {
    fmt.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), strings.Join(message, " "))
}

func main() {

    log("Engines starting up...")

    c := make(chan os.Signal, 1)
    connections := make(chan net.Conn, 10)

    listener, err := net.Listen("tcp", ":1337")

    if err != nil {
        log("Failed to open listener:", err.Error())
        os.Exit(1)
    }

    var wg sync.WaitGroup
    for i := 0; i < 4; i++ {
        wg.Add(1)
        go worker(i, connections, &wg)
    }

    go dispatcher(listener, connections)

    log("=== Goptive server started ===")

    signal.Notify(c, os.Interrupt)
    <-c

    log("*** Received interrupt signal ***")

    listener.Close()
    close(connections)

    log("Waiting for all jobs to finish...")
    wg.Wait()

    log("=== Goptive server shut down ===")
}
