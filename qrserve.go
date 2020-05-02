package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/mdp/qrterminal"
)

func findAddr() (net.IP, error) {
	var ipv6 net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				return nil, err
			} else if ip.To4() != nil {
				return ip, nil
			} else {
				ipv6 = ip
			}
		}
	}
	return ipv6, nil
}

func handler(w http.ResponseWriter, req *http.Request) {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, err := f.Stat()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, req, f.Name(), info.ModTime(), f)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Errorf("usage: %s file", os.Args[0])
		os.Exit(1)
	}

	ip, err := findAddr()
	if err != nil {
		log.Fatal(err)
	}
	path := "/" + uuid.New().String()
	http.HandleFunc(path, handler)

	url := "http://" + ip.String() + ":8080" + path
	log.Print(url)
	qrterminal.Generate(url, qrterminal.L, os.Stdout)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
