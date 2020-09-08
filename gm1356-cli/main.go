package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// root コマンド
var rootCmd = &cobra.Command{
	Use:   "gm1356-cli",
	Short: "GM1356 command line controller",
	Long:  `GM1356 command line controller`,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
