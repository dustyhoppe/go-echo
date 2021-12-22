package proxy

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dustyhoppe/go-echo/config"
	"github.com/dustyhoppe/go-echo/routing"
)

type ProxyServer interface {
	Start() error
}

type proxyServer struct {
	port     int
	server   *http.Server
	protocol string
	pemFile  string
	keyFile  string
}

func NewProxyServer(appConfig *config.AppConfig, routes *routing.Route) ProxyServer {

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", appConfig.Port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler(w, r)
		}),
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	return &proxyServer{
		server:   server,
		port:     appConfig.Port,
		protocol: appConfig.Protocol,
		pemFile:  appConfig.PemFile,
		keyFile:  appConfig.KeyFile,
	}
}

func (p *proxyServer) Start() error {

	if strings.EqualFold(p.protocol, "http") {
		log.Fatal(p.server.ListenAndServe())
	} else {
		log.Fatal(p.server.ListenAndServeTLS(p.pemFile, p.keyFile))
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Host %q", r.URL.Host)
	if r.URL.Host == "www.draftkings.com:443" {
		log.Fatal("DK sucks")
		return
	}

	if r.Method == http.MethodConnect {
		handleTunneling(w, r)
	} else {
		handleHTTP(w, r)
	}
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}

	//log.Printf("%s %s", r.Method, r.URL.String())

	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	//log.Printf("%s %s", req.Method, req.URL.String())

	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(destination, source http.Header) {
	for k, vv := range source {
		for _, v := range vv {
			destination.Add(k, v)
		}
	}
}
