package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	"io"
	"net"
	"riemann"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:5555")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	fmt.Println("connected")
	state := "ok"
	host := "localhost"
	service := "test"
	metric := int64(100)
	ttl := float32(100000.0)

	evt := riemann.Event{
		State:        &state,
		Host:         &host,
		Service:      &service,
		MetricSint64: &metric,
		Ttl:          &ttl,
	}

	msg := riemann.Msg{}
	msg.Events = append(msg.Events, &evt)

	data, err := proto.Marshal(&msg)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, uint32(len(data)))
	if err != nil {
		panic(err)
	}
	fmt.Printf("sending header (%d bytes) of %d\n", len(buf.Bytes()), len(data))
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}

	fmt.Println("sending body")
	_, err = conn.Write(data)
	if err != nil {
		panic(err)
	}

	fmt.Println("receiving")
	var header uint32
	err = binary.Read(conn, binary.BigEndian, &header)
	if err != nil {
		panic(err)
	}

	fmt.Printf("receiving %d\n", header)
	response := make([]byte, header)
	read, err := io.ReadFull(conn, response)
	if err != nil {
		panic(err)
	}
	fmt.Printf("received %d of %d\n", read, len(response))

	msg.Reset()
	err = proto.Unmarshal(response, &msg)
	if err != nil {
		panic(err)
	}
	if !msg.GetOk() {
		panic(msg.GetError())
	}

	fmt.Println("done")
}
