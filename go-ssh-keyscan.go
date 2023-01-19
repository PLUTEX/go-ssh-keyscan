package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"log"
	"os"
	"sync"
	"strings"

	"golang.org/x/crypto/ssh"
)

const (
	Username    = "dummy"
	DefaultPort = 22
)

var Ch chan string = make(chan string)
var IgnoreError = errors.New("Ignore this error")


func GetKeyScanCallback(alias string) func(string, net.Addr, ssh.PublicKey) error {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		Ch <- fmt.Sprintf("%s %s", alias, string(ssh.MarshalAuthorizedKey(key)))
		return IgnoreError
	}
}

func dial(server string, alias string, wg *sync.WaitGroup) {
	config := &ssh.ClientConfig{
		User:            Username,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: GetKeyScanCallback(alias),
		Timeout:         1e9,
	}

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
	var wg sync.WaitGroup
	var alias string
	var server string
	go out(&wg)
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		line = line[:len(line)-1] // chomp
		fields := strings.Fields(line)
		if len(fields) == 1 {
			server = fields[0]
			alias = fields[0]
		} else if len(fields) == 2 {
			server = fields[0]
			alias = fields[1]
		} else {
			log.Fatalln("Too many whitespaces in input line:", line)
		}
		wg.Add(2)                       // dial and print
		go dial(server, alias, &wg)
	}
	wg.Wait()
}
