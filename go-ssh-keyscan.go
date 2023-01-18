package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"

	"golang.org/x/crypto/ssh"
)

const (
	Username    = "dummy"
	DefaultPort = 22
)

var Ch chan string = make(chan string)
var IgnoreError = errors.New("Ignore this error")

func KeyScanCallback(hostname string, remote net.Addr, key ssh.PublicKey) error {
	Ch <- fmt.Sprintf("%s %s", hostname[:len(hostname)-3], string(ssh.MarshalAuthorizedKey(key)))
	return IgnoreError
}

func dial(server string, config *ssh.ClientConfig, wg *sync.WaitGroup) {
	_, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", server, DefaultPort), config)
	// For errors.Is() to work, x/crypto/ssh/client.go needs to be patched
	// to used %w instead of %v
	if err != nil && !errors.Is(err, IgnoreError) {
		// Don't expect a key from out()
		wg.Done()
	}
	wg.Done()

}

func out(wg *sync.WaitGroup) {
	for s := range Ch {
		fmt.Printf("%s", s)
		wg.Done()
	}
}

func main() {
	auths := []ssh.AuthMethod{}

	config := &ssh.ClientConfig{
		User:            Username,
		Auth:            auths,
		HostKeyCallback: KeyScanCallback,
	}

	var wg sync.WaitGroup
	go out(&wg)
	reader := bufio.NewReader(os.Stdin)
	for {
		server, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		server = server[:len(server)-1] // chomp
		wg.Add(2)                       // dial and print
		go dial(server, config, &wg)
	}
	wg.Wait()
}
