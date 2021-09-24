package signaling

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"

	api "github.com/JakWai01/airdrip/pkg/api/websockets/v1"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

// This signaling protocol is heavily inspired by the weron project created by @pojntfx
// Take a look at the specification by clicking the following link: https://github.com/pojntfx/weron/blob/main/docs/signaling-protocol.txt#L12

func NewSignalingServer() *SignalingServer {
	return &SignalingServer{
		communities:    map[string][]string{},
		macs:           map[string]bool{},
		connections:    map[string]net.Conn{},
		candidateCache: []string{},
	}
}

func (s *SignalingServer) HandleConn(conn *websocket.Conn) {

	go func() {
		for {

			// Read message from connection
			_, data, err := conn.Read(context.Background())
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(data))

			// Parse message
			var v api.Message
			if err := json.Unmarshal(data, &v); err != nil {
				log.Fatal(err)
			}

			// Handle different message types
			switch v.Opcode {
			case api.OpcodeApplication:
				var application api.Application
				if err := json.Unmarshal(data, &application); err != nil {
					log.Fatal(err)
				}

				if _, ok := s.macs[application.Mac]; ok {
					// Send rejection. That mac is already contained

					// Check if this conn is correct
					if err := wsjson.Write(context.Background(), conn, api.NewRejection()); err != nil {
						log.Fatal(err)
					}
				}
			case api.OpcodeReady:
			case api.OpcodeOffer:
			case api.OpcodeAnswer:
			case api.OpcodeCandidate:
			case api.OpcodeExited:
			}
		}
	}()
}
