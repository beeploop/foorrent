package cmd

import (
	"fmt"

	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/beeploop/foorrent/internal/utils"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints the torrent content information",
	Long:  `Reads the given torrent file and prints the content information`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		withPieces, _ := cmd.Flags().GetBool("with_pieces")
		humanReadable, _ := cmd.Flags().GetBool("human_readable")

		torrent, err := metadata.Read(file)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Announce: ", torrent.Announce)
		fmt.Println("Comment: ", torrent.Comment)
		fmt.Println("Creation Date: ", torrent.CreationDate)
		fmt.Println("Name: ", torrent.Info.Name)

		if humanReadable {
			fmt.Println("Length: ", utils.BytesToMB(torrent.Info.Length))
		} else {
			fmt.Println("Length: ", torrent.Info.Length)
		}

		if humanReadable {
			fmt.Println("Piece Length: ", utils.BytesToMB(torrent.Info.PieceLength))
		} else {
			fmt.Println("Piece Length: ", torrent.Info.PieceLength)

		}

		fmt.Println("Files:")
		for _, v := range torrent.Info.Files {
			if humanReadable {
				fmt.Printf("\t %v - %s\n", v.Path, utils.BytesToMB(v.Length))
			} else {
				fmt.Printf("\t %v - %d\n", v.Path, v.Length)
			}
		}

		if withPieces {
			fmt.Println("Files: ", torrent.Info.Pieces)
		}
	},
}

func init() {
	infoCmd.Flags().Bool("with_pieces", false, "include the pieces in the printed data")
	infoCmd.Flags().Bool("human_readable", false, "format file size/length into human-readable string")

	rootCmd.AddCommand(infoCmd)
}
