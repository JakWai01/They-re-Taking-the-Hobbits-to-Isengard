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

var communities = map[string][]string{}
var macs = []string{}

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

type Acceptance struct {
	Opcode string `json:"opcode"`
}

type Rejection struct {
	Opcode string `json:"opcode"`
}

type Ready struct {
	Opcode string `json:"opcode"`
}

type Introduction struct {
	Opcode string `json:"opcode"`
	Mac    string `json:"mac"`
}

type Offer struct {
	Opcode  string `json:"opcode"`
	Mac     string `json:"mac"`
	Payload string `json:"payload"`
}

type Answer struct {
	Opcode  string `json:"opcode"`
	Mac     string `json:"mac"`
	Payload string `json:"payload"`
}

type Candidate struct {
	Opcode  string `json:"opcode"`
	Mac     string `json:"mac"`
	Payload string `json:"payload"`
}

type Exited struct {
	Opcode string `json:"opcode"`
}

type Resignation struct {
	Opcode string `json:"opcode"`
	Mac    string `json:"mac"`
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
			// we get community and mac. Check if community exists. If not create it. Only allow unused macs.

			// Community maps string to tuple. Macs is an array and must be unique.
			fmt.Println("application")
			var opcode Application

			err := json.Unmarshal([]byte(message), &opcode)
			if err != nil {
				panic(err)
			}

			// check if community exists and if there are less than 2 members inside
			if val, ok := communities[opcode.Community]; ok {
				// check if length smaller than 2
				if len(val) >= 2 {
					// send Rejection. This community is full
					byteArray, err := json.Marshal(Rejection{Opcode: string(rejection)})
					if err != nil {
						panic(err)
					}

					fmt.Println(string(byteArray))

					_, err = c.Write(byteArray)
					if err != nil {
						panic(err)
					}
					return
				}
			} else {
				// Community does not exist. Create community and insert mac
				communities[opcode.Community] = append(communities[opcode.Community], opcode.Mac)
				fmt.Println(communities)

				macs = append(macs, opcode.Mac)
				fmt.Println(macs)

				// send Acceptance
				byteArray, err := json.Marshal(Acceptance{Opcode: string(acceptance)})
				if err != nil {
					panic(err)
				}

				fmt.Println(string(byteArray))

				_, err = c.Write(byteArray)
				if err != nil {
					panic(err)
				}
				return
			}

			for i := 0; i < len(macs); i++ {
				if macs[i] == opcode.Mac {
					// send Rejection. That Mac is already contained
					byteArray, err := json.Marshal(Rejection{Opcode: string(rejection)})
					if err != nil {
						panic(err)
					}

					fmt.Println(string(byteArray))

					_, err = c.Write(byteArray)
					if err != nil {
						panic(err)
					}
					return
				}
			}

			// send Acceptance
			byteArray, err := json.Marshal(Acceptance{Opcode: string(acceptance)})
			if err != nil {
				panic(err)
			}

			fmt.Println(string(byteArray))

			_, err = c.Write(byteArray)
			if err != nil {
				panic(err)
			}
			return

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
