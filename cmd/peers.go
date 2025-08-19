package cmd

import (
	"fmt"
	"net"
	"strconv"

	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/beeploop/foorrent/internal/peer"
	"github.com/beeploop/foorrent/internal/tracker"
	"github.com/spf13/cobra"
)

var peersCmd = &cobra.Command{
	Use:   "peers",
	Short: "Prints the interval and list of peers",
	Long:  `Contacts the tracker and retrieve the interval and list of peers`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")

		fmt.Println("reading .torrent file")
		torrent, err := metadata.Read(file)
		if err != nil {
			panic(err.Error())
		}

		infoHash, err := torrent.InfoHash()
		if err != nil {
			panic(err.Error())
		}

		peerID, err := peer.GeneratePeerID()
		if err != nil {
			panic(err.Error())
		}

		input := tracker.TrackerInput{
			Announce:   torrent.Announce,
			InfoHash:   infoHash,
			PeerID:     peerID,
			Port:       tracker.DEFAULT_PORT,
			Left:       torrent.TotalSize(),
			Downloaded: 0,
			Uploaded:   0,
		}

		fmt.Println("requesting data from tracker...")
		result, err := tracker.Request(input)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Interval: ", result.Interval)
		fmt.Println("Peers:")
		for _, peer := range result.Peers {
			fmt.Printf("\t - %s\n", net.JoinHostPort(peer.IP.String(), strconv.Itoa(int(peer.Port))))
		}
	},
}

func init() {
	rootCmd.AddCommand(peersCmd)
}
