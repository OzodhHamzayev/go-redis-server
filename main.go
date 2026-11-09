package main

import (
	"fmt"
	"net"
	"strings"
)


func main() {
    l, err := net.Listen("tcp", ":6379")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer l.Close()

    db := make(map[string]string)

    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        handleConn(conn, db)
    }
}

func handleConn(conn net.Conn, db map[string]string) {
    defer conn.Close()
    resp := NewResp(conn)
    writer := NewWriter(conn)

    for {
        command, err := resp.Read()
        if err != nil {
            return
        }

        fmt.Println("dadasdasd")

        reply := handleCommand(command, db)
        if err := writer.Write(reply); err != nil {
            fmt.Println(err)
            return
        }
    }
}

func handleCommand(v Value, db map[string]string) Value {
    if v.typ != "array" || len(v.array) == 0 {
        return Value{typ: "error", str: "noto‘g‘ri buyruq"}
    }

    command := strings.ToUpper(v.array[0].bulk)

    switch command {
    case "PING":
        return Value{typ: "string", str: "PONG"}

    case "SET":
        if len(v.array) != 3 {
            return Value{typ: "error", str: "noto‘g‘ri SET buyrug‘i"}
        }
        key := v.array[1].bulk
        value := v.array[2].bulk
        db[key] = value
        return Value{typ: "string", str: "OK"}

    case "GET":
        if len(v.array) != 2 {
            return Value{typ: "error", str: "noto‘g‘ri GET buyrug‘i"}
        }
        value, exists := db[v.array[1].bulk]
        if !exists {
            return Value{typ: "null"}
        }
        return Value{typ: "bulk", bulk: value}

    default:
        return Value{typ: "error", str: "noma’lum buyruq"}
    }


}