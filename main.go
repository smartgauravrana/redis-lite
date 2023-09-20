package main

import (
	"fmt"
	"net"
	"reflect"
	"strings"
)


func handleClient(conn net.Conn) {
    defer conn.Close()
    fmt.Println("Client connected:", conn.RemoteAddr())

    for {
        buf := make([]byte, 1024)
        _, err := conn.Read(buf)
        if err != nil {
            fmt.Println("Error reading from client:", err)
            return
        }

        cmd, err := DeserializeRESP(buf)
        if err != nil {
            fmt.Println("Error parsing RESP:", err)
            return
        }

		fmt.Println("cmd: ", cmd)
	
		fmt.Println("isSLice: ", IsSlice(cmd))


		cmdSlice := cmd.([]interface{});

		str,ok := cmdSlice[0].(string)
		var command string

		if(ok){
			command = strings.ToUpper(str)
		} else{
			fmt.Println("not string!", cmd)
			return;
		}

        switch command {
        case "PING":
            PING(conn)
        case "ECHO":
            ECHO(conn, cmdSlice)
		case "SET":
			SET(conn, cmdSlice)
		case "GET":
			GET(conn, cmdSlice)
        default:
            fmt.Println("Unknown command:", cmd)
			conn.Write([]byte("*0\r\n"))
        }
    }
}

func IsSlice(v interface{}) bool {
    return reflect.TypeOf(v).Kind() == reflect.Slice
}

func startServer() {
    listenAddr := ":6379"

    listener, err := net.Listen("tcp", listenAddr)
    if err != nil {
        fmt.Println("Error starting server:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Server is listening on", listenAddr)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }
        go handleClient(conn)
    }
}

func main() {
    startServer()
}
