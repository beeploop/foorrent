package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/beeploop/foorrent/internal/peer"
	"github.com/beeploop/foorrent/internal/piece"
	"github.com/beeploop/foorrent/internal/storage"
	"github.com/beeploop/foorrent/internal/tracker"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download torrent file",
	Long:  `Starts the download process of downloading torrent file`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")

		torrent, err := metadata.Read(file)
		if err != nil {
			log.Fatalf("Reading torrent failed: %s\n", err.Error())
		}

		peerID, err := peer.GeneratePeerID()
		if err != nil {
			panic(err.Error())
		}

		infoHash, err := torrent.InfoHash()
		if err != nil {
			log.Fatalf("Failed to get info hash: %s\n", err.Error())
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

		client, err := tracker.TrackerFactory(input)
		if err != nil {
			panic(err.Error())
		}

		response, err := client.Request(input)
		if err != nil {
			log.Fatalf("Failed to contact tracker: %s\n", err.Error())
		}

		var writer storage.Storage
		if torrent.IsSingleFileMode() {
			writer = storage.NewSingleFileStorage(torrent.Info.Name)
		} else {
			writer = storage.NewMultiFileStorage(torrent.Info.Name, torrent.FileMap())
		}
		writer.Init()

		pieceManager, err := piece.NewManager(torrent, writer)
		if err != nil {
			log.Fatalf("Failed to create piece manager: %s\n", err.Error())
		}

		peerManager := peer.NewManager(torrent, peerID, response.Peers, pieceManager)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		peerManager.Start(ctx)

		quitChan := make(chan os.Signal, 1)
		signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM|syscall.SIGKILL)
		fmt.Println("Running... Press CTRL+C to quit")

		ticker := time.NewTicker(time.Second * 1)
		go func() {
			for {
				<-ticker.C
				downloaded, total := pieceManager.Downloaded()
				peers, totalPeers := peerManager.ActivePeers()
				missingBlk, requestedBlk, doneBlk := pieceManager.BlockStats()
				log.Printf(
					"[ Peers %d/%d ] [ Pieces %d/%d ] [ Block (missing - requested - done) %d - %d - %d ]\n",
					peers,
					totalPeers,
					downloaded,
					total,
					missingBlk,
					requestedBlk,
					doneBlk,
				)
			}
		}()

		<-quitChan
		ticker.Stop()

		fmt.Println("Gracefully shutting down...")
	},
}

func init() {
	downloadCmd.PersistentFlags().String("file", "", "input torrent file")
	downloadCmd.MarkPersistentFlagRequired("file")

	rootCmd.AddCommand(downloadCmd)
}
