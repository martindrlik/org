package main

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/pkcs12"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: $ fakeapns cert.p12 password")
		os.Exit(2)
	}
	cfg, err := config(os.Args[1], os.Args[2])
	fatal(err)
	ln, err := tls.Listen("tcp", ":8084", cfg)
	fatal(err)
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("fakeapns: %v\n", err)
			continue
		}
		fmt.Printf("fakeapns: accept: %v\n", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func fatal(err error) {
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, "fakeapns: %v\n", err)
	os.Exit(1)
}

func config(certFile, password string) (*tls.Config, error) {
	p12, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	blocks, err := pkcs12.ToPEM(p12, password)
	if err != nil {
		return nil, err
	}
	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}
	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		return nil, err
	}
	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return cfg, nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	b := make([]byte, 500)
	n, err := conn.Read(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fakeapns: %v: %v\n", conn.RemoteAddr(), err)
		return
	}
	os.Stdout.Write(b[:n])
}
