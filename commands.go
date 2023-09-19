package main

import (
	"fmt"
	"net"
)

func PING(conn net.Conn){
	conn.Write([]byte("+PONG\r\n"))
}

func ECHO(conn net.Conn,c []interface{}) {

	errorMsg := "-ERR wrong number of arguments for command\r\n";

	if len(c) != 2{
		conn.Write([]byte(errorMsg))
		return;
	}

	if str, ok := c[1].(string); ok{
		conn.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(str), str)))
	}else{
		conn.Write([]byte(errorMsg))
	}
}