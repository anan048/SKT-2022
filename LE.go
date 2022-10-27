package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"strings"
	"time"
)

type MyArgs struct {
}

type Arith struct {
}

var vote float64 = 0
var a int
var term int = 0
var terml int
var leade bool = false
var app bool = false
var jumlahanggota float64 = 1

func leader() {
	leade = true
	fmt.Printf("jadi leader term: %d", term)
	for true {
	}
}

func candidate(peerAddress string) {
	for true {
		fmt.Println("jadi candidate")
		term += 1
		vote += 1
		client, err := rpc.DialHTTP("tcp", peerAddress)
		handleConnError(err, peerAddress)
		clientargs := ""
		var localResult bool
		err = client.Call("Arith.RequestVote", clientargs, &localResult)
		result := localResult
		if result == true {
			jumlahanggota += 1
			vote += 1
			if vote > 0.5*jumlahanggota {
				leader()
			}
		}
	}

}

func (t *Arith) RequestVote(empty string, result *bool) error {
	*result = true
	return nil
}

func (t *Arith) AppendEntries(termm *int, result *string) error {

	var b int

	if leade == true {
		b = term
	}
	s1 := strconv.Itoa(b)
	s2 := strconv.FormatBool(leade)
	*result = fmt.Sprint(s1 + "," + s2)
	return nil
}

func append(peerAddress string) {
	client, err := rpc.DialHTTP("tcp", peerAddress)
	handleConnError(err, peerAddress)
	var split []string
	for true {

		clientargs := terml
		var localResult string
		err = client.Call("Arith.AppendEntries", clientargs, &localResult)
		handleConnError(err, peerAddress)
		split = strings.Split(localResult, ",")

		leade, err = strconv.ParseBool(split[1])
		if leade == true {
			app = true
			term, err = strconv.Atoi(split[0])
		}
		if app == true {
			time.Sleep(10 * time.Second)
		}

	}
}

func follower(peerAddress string) {
	go append(peerAddress)
	for true {

		if app == true {
			fmt.Println("menerima heartbeat")
			a = rand.Intn(20) + 30
			app = false
		}
		fmt.Printf("%d term: %d\n", a, term)
		time.Sleep(1 * time.Second)
		a -= 1

		if a == 0 {
			candidate(peerAddress)
		}
	}

}

func main() {
	portNumber := os.Args[1]
	peersAddress := os.Args[2]

	rand.Seed(time.Now().UnixNano())
	a = rand.Intn(20) + 30
	go follower(peersAddress)
	// Inisiasi struct arith
	arith := &Arith{}
	// Registrasikan struct dan method ke RPC
	rpc.Register(arith)
	// Deklarasikan bahwa kita menggunakan protokol HTTP sebagai msekanisme pengiriman pesan
	rpc.HandleHTTP()
	// Deklarasikan listerner HTTP dengan layer transport TCP dan Port 1234
	listener, err := net.Listen("tcp", ":"+string(portNumber))
	handleError(err)
	// Jalankan server HTTP
	http.Serve(listener, nil)

}

func handleError(err error) {

	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())

	}
}

func handleConnError(err error, peersAdrress string) {
	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())
		time.Sleep(1 * time.Second)
		follower(peersAdrress)
	}

}
