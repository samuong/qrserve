package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/mdp/qrterminal"
)

func findAddr(addrs []net.Addr, err error) (net.IP, error) {
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, err
		} else if !ip.IsLoopback() && !ip.IsLinkLocalUnicast() {
			return ip, nil
		}
	}
	return nil, fmt.Errorf("couldn't find address of local host")
}

func main() {
	port := flag.Uint("port", 8080, "port number to listen on")
	flag.Parse()

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: http.FileServer(http.Dir(pwd)),
	}

	ip, err := findAddr(net.InterfaceAddrs())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	u := &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(ip.String(), fmt.Sprint(*port)),
	}
	qrterminal.Generate(u.String(), qrterminal.L, os.Stdout)
	fmt.Printf("Serving %s on %s\n", pwd, u.String())

	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
