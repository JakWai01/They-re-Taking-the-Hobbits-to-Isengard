package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strings"
)

// This signaling protocol is heavily inspired by the weron project created by @pojntfx
// Take a look at the specification by clicking the following link: https://github.com/pojntfx/weron/blob/main/docs/signaling-protocol.txt#L12

type Opcode string

const (
	application  Opcode = "application"
	acceptance   Opcode = "acceptance"
	rejection    Opcode = "rejection"
	ready        Opcode = "ready"
	introduction Opcode = "introduction"
	offer        Opcode = "offer"
	answer       Opcode = "answer"
	candidate    Opcode = "candidate"
	exited       Opcode = "exited"
	resignation  Opcode = "resignation"
)

type Application struct {
	Opcode    string `json:"opcode"`
	Community string `json:"community"`
	Mac       string `json:"mac"`
}

type Acceptance struct{}

type Rejection struct{}

type Ready struct{}

type Introduction struct {
	Mac string `json:"mac"`
}

type Offer struct {
	Mac     string `json:"mac"`
	Payload string `json:"payload"`
}

type Answer struct {
	Mac     string `json:"mac"`
	Payload string `json:"payload"`
}

type Candidate struct {
	Mac     string `json:"mac"`
	Payload string `json:"payload"`
}

type Exited struct{}

type Resignation struct {
	Mac string `json:"mac"`
}

func handleConnection(c net.Conn) {
	for {
		message, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			panic(err)
		}

		values := make(map[string]json.RawMessage)

		err = json.Unmarshal([]byte(message), &values)
		if err != nil {
			panic(err)
		}

		switch Opcode(strings.ReplaceAll(string(values["opcode"]), "\"", "")) {
		case application:
			fmt.Println("application")
			var opcode Application

			err := json.Unmarshal([]byte(message), &opcode)
			if err != nil {
				panic(err)
			}

			// byteArray, err := json.MarshalIndent(opcode, "", "  ")
			// if err != nil {
			// 	panic(err)
			// }

			// fmt.Println(string(byteArray))
		case ready:
			fmt.Println("ready")
			var opcode Ready

			err := json.Unmarshal([]byte(message), &opcode)
			if err != nil {
				panic(err)
			}
		case offer:
			fmt.Println("offer")
			var opcode Offer

			err := json.Unmarshal([]byte(message), &opcode)
			if err != nil {
				panic(err)
			}
		case answer:
			fmt.Println("answer")
			var opcode Answer

			err := json.Unmarshal([]byte(message), &opcode)
			if err != nil {
				panic(err)
			}
		case candidate:
			fmt.Println("candidate")
			var opcode Candidate

			err := json.Unmarshal([]byte(message), &opcode)
			if err != nil {
				panic(err)
			}
		case exited:
			fmt.Println("exited")
			var opcode Exited

			err := json.Unmarshal([]byte(message), &opcode)
			if err != nil {
				panic(err)
			}
		default:
			panic("Invalid message. Please use a valid opcode.")
		}
	}
}

func main() {
	var laddr = flag.String("laddr", "localhost:8080", "listen address")
	flag.Parse()

	fmt.Println(*laddr)

	l, err := net.Listen("tcp4", *laddr)
	if err != nil {
		panic(err)
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(c)
	}
}
