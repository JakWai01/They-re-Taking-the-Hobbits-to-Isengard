package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/JakWai01/airdrip/pkg/signaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	laddrKey     = "laddr"
	communityKey = "community"
	macKey       = "mac"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start a signaling client.",
	RunE: func(cmd *cobra.Command, args []string) error {

		fatal := make(chan error)
		done := make(chan struct{})

		client := signaling.NewSignalingClient()

		socket := viper.GetString(laddrKey) + ":8080"
		fmt.Println(socket)
		go func() {

			go client.HandleConn(socket, viper.GetString(communityKey), viper.GetString(macKey))

		}()

		for {
			select {
			case err := <-fatal:
				panic(err)
			case <-done:
				os.Exit(0)
			}
		}
	},
}

func init() {
	clientCmd.PersistentFlags().String(laddrKey, "localhost", "Listen address")
	clientCmd.PersistentFlags().String(communityKey, "a", "Community to join")
	clientCmd.PersistentFlags().String(macKey, "124", "Mac to identify you as a unique host")

	// Bind env variables
	if err := viper.BindPFlags(clientCmd.PersistentFlags()); err != nil {
		log.Fatal("could not bind flags:", err)
	}
	viper.SetEnvPrefix("airdrip")
	viper.AutomaticEnv()
}
