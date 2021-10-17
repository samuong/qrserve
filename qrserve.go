package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/mdp/qrterminal"
)

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func handler(path string) (http.Handler, string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, "", err
	}
	if info.IsDir() {
		return http.FileServer(http.Dir(path)), "", nil
	}
	base := filepath.Base(path)
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/"+base,
		func(w http.ResponseWriter, req *http.Request) {
			http.ServeFile(w, req, path)
		},
	)
	return mux, base, nil
}

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
	path := flag.Arg(0)
	if path == "" {
		var err error
		path, err = os.Getwd()
		check(err)
	}
	h, suffix, err := handler(path)
	check(err)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: &logMiddleware{Handler: h},
	}
	ip, err := findAddr(net.InterfaceAddrs())
	check(err)
	u := &url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(ip.String(), fmt.Sprint(*port)),
		Path:   suffix,
	}
	qrterminal.Generate(u.String(), qrterminal.L, os.Stdout)
	fmt.Printf("Serving %s on %s\n", path, u.String())
	check(server.ListenAndServe())
}
