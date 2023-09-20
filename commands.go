package main

import (
	"fmt"
	"net"
)

var mapStore map[string]string

const okReply = "+OK\r\n"

const nullBulkStringReply = "$-1\r\n"

func getBulkStringReply(v string) string{
	return fmt.Sprintf("$%d\r\n%s\r\n", len(v), v)
}

func SET(conn net.Conn,c []interface{}){

	if len(mapStore) == 0{
		mapStore = make(map[string]string)
	}

	errorMsg := "-ERROR Syntax error\r\n";

	cmds := make([]string, len(c))
	for i, v := range c {
		cmds[i] = fmt.Sprint(v)
	}
	
	cmdLen := len(cmds)
	if cmdLen > 3{
		conn.Write([]byte(errorMsg))
		return;
	}
	var v string

	if cmdLen == 2 {
		v = ""
	}else{
		v = cmds[2]
	}

	key := cmds[1]
	mapStore[key] = v
	conn.Write([]byte(okReply))
}

func GET(conn net.Conn,c []interface{}){

	if len(mapStore) == 0{
		mapStore = make(map[string]string)
	}

	errorMsg := "-ERR Syntax errorr\n";

	cmds := make([]string, len(c))
	for i, v := range c {
		cmds[i] = fmt.Sprint(v)
	}

	if len(cmds) > 2{
		conn.Write([]byte(errorMsg))
		return;
	}

	key := cmds[1]
	value := mapStore[key];

	if value == ""{
		conn.Write([]byte(nullBulkStringReply));
		return;
	}
	conn.Write([]byte(getBulkStringReply(value)))
}

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