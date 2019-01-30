package tsl

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
)

func StartTCLServer() {

	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println(err)
		return
	}
	if err != nil {
		log.Println(err)
		return
	}
	defer ln.Close()
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		println(msg)
		switch msg {
		case "ehlo\n":
			n, err := conn.Write([]byte("250 OK\n"))
			if err != nil {
				log.Println(n, err)
				return
			}
		default:
			n, err := conn.Write([]byte("world\n"))
			if err != nil {
				log.Println(n, err)
				return
			}
		}

	}
}

func StartTLSServer() {
	fmt.Println("start https server port 443")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Write([]byte("This is an example server.\n"))
	})
	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		//		CipherSuites: []uint16{
		//			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		//			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		//			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		//			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		//		},
	}
	srv := &http.Server{
		Addr:         ":443",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}
