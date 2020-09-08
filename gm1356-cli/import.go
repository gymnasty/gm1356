package main

import (
	"fmt"

	"github.com/gymnasty/gm1356"
	"github.com/spf13/cobra"
)

func init() {
	var importCmd = &cobra.Command{
		Use:   "import",
		Short: "import recorded measure data",
		Long:  `import recorded measure data`,
		Run: func(cmd *cobra.Command, args []string) {
			const eventBufferSize = 128
			driver, err := gm1356.Open(eventBufferSize)
			if err != nil {
				panic(err)
			}
			defer driver.Close()

			c := driver.EventChannel()
			go func() {
				for {
					event := <-c
					fmt.Println(event.String())
				}
			}()

			driver.Import()
		},
	}
	rootCmd.AddCommand(importCmd)
}
