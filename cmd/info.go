package cmd

import (
	"fmt"

	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Prints the torrent content information",
	Long:  `Reads the given torrent file and prints the content information`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		withPieces, _ := cmd.Flags().GetBool("with_pieces")

		torrent, err := metadata.Read(file)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("Announce: ", torrent.Announce)
		fmt.Println("Comment: ", torrent.Comment)
		fmt.Println("Creation Date: ", torrent.CreationDate)
		fmt.Println("Name: ", torrent.Info.Name)
		fmt.Println("Length: ", torrent.Info.Length)
		fmt.Println("Piece Length: ", torrent.Info.PieceLength)
		fmt.Println("Files:")
		for _, v := range torrent.Info.Files {
			fmt.Printf("\t %v - %d\n", v.Path, v.Length)
		}

		if withPieces {
			fmt.Println("Files: ", torrent.Info.Pieces)
		}
	},
}

func init() {
	infoCmd.Flags().Bool("with_pieces", false, "include the pieces in the printed data")

	rootCmd.AddCommand(infoCmd)
}
