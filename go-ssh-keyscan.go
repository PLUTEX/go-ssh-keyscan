package main

import (
	"bufio"
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
var supportedHostKeyAlgos = []string{
	ssh.KeyAlgoRSA,
	ssh.KeyAlgoDSA,
	ssh.KeyAlgoECDSA256,
	ssh.KeyAlgoECDSA384,
	ssh.KeyAlgoECDSA521,
	ssh.KeyAlgoED25519,
	// not yet supported in crypto@0.3.0: ssh.KeyAlgoRSASHA256,
	// not yet supported in crypto@0.3.0: ssh.KeyAlgoRSASHA512,
}


func GetKeyScanCallback(alias string, wg *sync.WaitGroup) func(string, net.Addr, ssh.PublicKey) error {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		wg.Add(1)
		Ch <- fmt.Sprintf("%s %s", alias, string(ssh.MarshalAuthorizedKey(key)))
		return nil
	}
}

func dial(server string, alias string, hostkeyalgo string, wg *sync.WaitGroup) {
	config := &ssh.ClientConfig{
		User:              Username,
		Auth:              []ssh.AuthMethod{},
		HostKeyAlgorithms: []string{hostkeyalgo},
		HostKeyCallback:   GetKeyScanCallback(alias, wg),
		Timeout:           1e9,
	}

	ssh.Dial("tcp", fmt.Sprintf("%s:%d", server, DefaultPort), config)
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
		wg.Add(len(supportedHostKeyAlgos))
		for _, hostkeyalgo := range supportedHostKeyAlgos {
			go dial(server, alias, hostkeyalgo, &wg)
		}
	}
	wg.Wait()
}
