package main

import (
	"fmt"
	"time"

	"github.com/gymnasty/gm1356"
	"github.com/spf13/cobra"
)

func init() {
	var measureCmd = &cobra.Command{
		Use:   "measure",
		Short: "measure sound level",
		Long:  `measure sound level`,
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

			for {
				if err := driver.Measure(); err != nil {
					panic(err)
				}
				time.Sleep(time.Second)
			}
		},
	}
	rootCmd.AddCommand(measureCmd)
}
