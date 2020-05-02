package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type addr struct {
	network, str string
}

func (a addr) Network() string {
	return a.network
}

func (a addr) String() string {
	return a.str
}

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestFindAddr(t *testing.T) {
	addrs := []net.Addr{
		addr{"tcp", "127.0.0.1/8"},                  // loopback IPv4
		addr{"tcp", "::1/128"},                      // loopback IPv6
		addr{"tcp", "169.254.0.196/24"},             // link-local IPv4
		addr{"tcp", "fe80::21f7:779c:6328:1b85/64"}, // link-local IPv6
		addr{"tcp", "192.168.0.196/24"},             // public IPv4
	}
	ip, err := findAddr(addrs, nil)
	require.NoError(t, err)
	assert.Equal(t, "192.168.0.196", ip.String())
}

func TestFindAddrWithErrorFromInput(t *testing.T) {
	_, err := findAddr(nil, errors.New("error from net.InterfaceAddrs()"))
	assert.Error(t, err, "error from net.InterfaceAddrs()")
}

func TestFindAddrWithNoAddresses(t *testing.T) {
	_, err := findAddr([]net.Addr{}, nil)
	assert.Error(t, err)
}

func TestIPv6AddrIsEnclosedInSquareBrackets(t *testing.T) {
}
