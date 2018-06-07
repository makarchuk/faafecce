package cmd

import (
	"fmt"
	"os"

	"github.com/makarchuk/faafecce/faafecce"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Faafecce",
	Short: "Face mirroring tool",
	Run: func(cmd *cobra.Command, args []string) {
		filename, err := cmd.PersistentFlags().GetString("file")
		if err != nil {
			panic(err)
		}
		outfile, err := cmd.PersistentFlags().GetString("write")
		if err != nil {
			panic(err)
		}
		// if middle {
		// 	middler = faafecce.Middle
		// } else {
		// 	middler = faafecce.Face
		// }
		middler := faafecce.Middle

		err = faafecce.Transform(middler, filename, "", outfile)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("file", "f", "", "File with face picture")
	rootCmd.PersistentFlags().StringP("name", "n", "", "Name to apply")
	rootCmd.PersistentFlags().StringP("write", "w", "", "Output file")
	rootCmd.PersistentFlags().BoolP("middle", "m", false,
		"If applied no CV will be used and image will be mirrored over geometrical center")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
