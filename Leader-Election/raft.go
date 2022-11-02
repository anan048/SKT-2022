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
var app bool = false
var jumlahanggota float64 = 1

func append(peerAddress string, state int) bool {
	client, err := rpc.DialHTTP("tcp", peerAddress)
	handleConnError(err, peerAddress, state)
	clientargs := term
	var localResult bool
	err = client.Call("Raft.AppendEntries", clientargs, &localResult)
	handleConnError(err, peerAddress, state)
	return localResult
}

func leader(peerAddress map[int]string) {
	var result bool
	state := 2
	fmt.Printf("jadi leader. term: %d", term)
	for true {
		for i := 0; i < len(os.Args)-2; i++ {
			result = append(peerAddress[i], state)
		}
		if result == true {
			time.Sleep(10 * time.Second)
		} else {
			leader(peerAddress)
		}

	}
}

func voting(peerAddress string, state int) bool {
	client, err := rpc.DialHTTP("tcp", peerAddress)
	handleConnError(err, peerAddress, state)
	clientargs := ""
	var localResult bool
	err = client.Call("Raft.RequestVote", clientargs, &localResult)
	handleConnError(err, peerAddress, state)
	return localResult
}

func candidate(peerAddress map[int]string) {
	var result bool
	state := 1
	for true {
		fmt.Println("jadi candidate")
		term += 1
		vote += 1
		for i := 0; i < len(os.Args)-2; i++ {
			result = voting(peerAddress[i], state)
			if result == true {
				jumlahanggota += 1
				vote += 1
			}
		}

		if vote > 0.5*jumlahanggota {
			leader(peerAddress)
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
	var a bool
	if err != nil {
		fmt.Println("Terdapat error : ", err.Error())

		time.Sleep(5 * time.Second)
		if state == 1 {
			a = voting(peersAdrress, state)
			if a == true {
				fmt.Println("")
			}
		} else {
			a = append(peersAdrress, state)
			if a == true {
				fmt.Println("")
			}
		}

	}

}
