package main

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

type MyArgs struct {
}

type Raft struct {
}

var vote float64 = 0
var timer int
var term int = 0
var a int = 0
var app bool = false
var jumlahanggota float64 = 1

func append(peerAddress string, state int) {
	a += 1
	client, err := rpc.DialHTTP("tcp", peerAddress)
	handleConnError(err, peerAddress, state)

	if err == nil {

		clientargs := term
		var localResult bool
		err = client.Call("Raft.AppendEntries", clientargs, &localResult)
		handleConnError(err, peerAddress, state)
		result := localResult
		if result == true {

			time.Sleep(1 * time.Second)

		}

	}

}

func leader(peerAddress map[int]string) {
	state := 2
	fmt.Printf("jadi leader. term: %d\n", term)
	for true {

		for i := 0; i < len(os.Args)-2; i++ {
			go append(peerAddress[i], state)
		}
		if a >= len(peerAddress) {
			time.Sleep(10 * time.Second)
			a = 0
		}

	}
}

func voting(peerAddress string, state int) {
	client, err := rpc.DialHTTP("tcp", peerAddress)
	handleConnError(err, peerAddress, state)
	if err == nil {
		clientargs := ""
		var localResult bool
		err = client.Call("Raft.RequestVote", clientargs, &localResult)
		handleConnError(err, peerAddress, state)
		result := localResult
		if err != nil {
			jumlahanggota += 1
			result = false
		}
		if result == true && err == nil {
			jumlahanggota += 1
			vote += 1
		}
	}

}

func candidate(peerAddress map[int]string) {
	state := 1
	for true {
		fmt.Println("jadi candidate")
		term += 1
		vote += 1
		for i := 0; i < len(os.Args)-2; i++ {
			voting(peerAddress[i], state)
		}

		if vote > 0.5*jumlahanggota {
			leader(peerAddress)
		} else {
			follower(peerAddress)
		}
	}

}

func (t *Raft) RequestVote(empty string, result *bool) error {
	*result = true
	return nil
}

func (t *Raft) AppendEntries(termm int, result *bool) error {
	term = termm
	if term == termm {

		app = true
	}
	*result = true
	return nil
}


func follower(peerAddress map[int]string) {
	rand.Seed(time.Now().UnixNano())
	timer = rand.Intn(20) + 30
	for true {

		if app == true {
			fmt.Println("menerima heartbeat")
			timer = rand.Intn(20) + 30
			app = false
		}
		fmt.Printf("%d term: %d\n", timer, term)
		time.Sleep(1 * time.Second)
		timer -= 1

		if timer == 0 {
			candidate(peerAddress)
		}
	}
}

func main() {
	portNumber := os.Args[1]
	var peersAddress = make(map[int]string)
	for i := 0; i < len(os.Args)-2; i++ {
		peersAddress[i] = os.Args[i+2]
	}

	go follower(peersAddress)
	// Inisiasi struct raft
	raft := &Raft{}
	// Registrasikan struct dan method ke RPC
	rpc.Register(raft)
	// Deklarasikan bahwa kita menggunakan protokol HTTP sebagai msekanisme pengiriman pesan
	rpc.HandleHTTP()
	// Deklarasikan listerner HTTP dengan layer transport TCP dan Port 
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

func handleConnError(err error, peersAdrress string, state int) {
	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())

		time.Sleep(5 * time.Second)
		if state == 1 {
			voting(peersAdrress, state)

		} else if state == 2 {
			append(peersAdrress, state)
		}

	}

}
