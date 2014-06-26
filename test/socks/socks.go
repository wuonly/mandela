package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	// sample1()
	// sample2()
	sample3()
}

func sample1() {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	ln, _ := net.ListenTCP("tcp", addr)
	for {
		conn, _ := ln.AcceptTCP()
		defer conn.Close()
		var ver, nMethods byte
		binary.Read(conn, binary.BigEndian, &ver)
		log.Println(ver)
		time.Sleep(10 * time.Millisecond)
		binary.Read(conn, binary.BigEndian, &nMethods)
		log.Println(nMethods)
	}
}

func sample2() {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	ln, _ := net.ListenTCP("tcp", addr)
	for {
		conn, _ := ln.AcceptTCP()
		defer conn.Close()
		// var ver, nMethods byte
		buf := make([]byte, 1024)
		n, _ := conn.Read(buf)
		log.Println(buf[:n])
	}
}

func sample3() {
	addr := &net.TCPAddr{IP: []byte{0, 0, 0, 0}, Port: 8080}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", 8080, err)
	}

	for {
		client, err := listener.AcceptTCP()
		if err != nil {
			log.Fatalf("Failed to accept new client connection: %v", err)
			e := err.(net.Error)
			if !e.Temporary() {
				os.Exit(1)
			}
		}
		raddr := client.RemoteAddr()
		defer client.Close()

		var versionMethod [1]byte
		_, err = io.ReadFull(client, versionMethod[:])
		if err != nil {
			log.Printf("%v: Failed to read the version and methods number: %v", raddr, err)
			return
		}
		log.Println(versionMethod)

		_, err = io.ReadFull(client, versionMethod[:])
		if err != nil {
			log.Printf("%v: Failed to read the version and methods number: %v", raddr, err)
			return
		}
		log.Println(versionMethod)

	}
}
