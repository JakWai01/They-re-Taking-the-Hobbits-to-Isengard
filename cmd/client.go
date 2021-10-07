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
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Start a signaling client.",
	RunE: func(cmd *cobra.Command, args []string) error {

		fatal := make(chan error)
		done := make(chan struct{})

		client := signaling.NewSignalingClient()

		socket := viper.GetString(laddrKey)
		fmt.Println(socket)
		go func() {

			go client.HandleConn(socket, viper.GetString(communityKey))

		}()

		for {
			select {
			case err := <-fatal:
				log.Fatal(err)
			case <-done:
				os.Exit(0)
			}
		}
	},
}

func init() {
	clientCmd.PersistentFlags().String(laddrKey, "localhost", "Listen address")
	clientCmd.PersistentFlags().String(communityKey, "a", "Community to join")

	// Bind env variables
	if err := viper.BindPFlags(clientCmd.PersistentFlags()); err != nil {
		log.Fatal("could not bind flags:", err)
	}
	viper.SetEnvPrefix("airdrip")
	viper.AutomaticEnv()
}
