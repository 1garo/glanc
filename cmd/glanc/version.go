package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Major = "0"
const Minor = "1"
const Fix = "1"
const Verbal = "TX Add && Balances List"

// versionCmd -> output cli version
func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Describes version.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s.%s.%s-beta %s\n", Major, Minor, Fix, Verbal)
		},
	}
}
