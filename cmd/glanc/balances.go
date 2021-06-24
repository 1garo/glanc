package main

import (
	"fmt"
	"os"

	"github.com/1garo/glanc/database"
	"github.com/spf13/cobra"
)

// balancesCmd -> create balances cli command and it's configs
func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interact with balances (list...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	balancesCmd.AddCommand(balancesListCmd())

	return balancesCmd
}

// balancesListCmd -> list all balances cmd
func balancesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists all balances.",
		Run: func(cmd *cobra.Command, args []string) {
			state, err := database.NewStateFromDisk()
			if err != nil {
				fmt.Println(err)
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			var firstFourSnap []byte
			snap := state.LatestBlockHash()
			firstFourSnap = snap[:4]

			fmt.Printf("Accounts balances: %x\n", firstFourSnap)
			fmt.Println("__________________")
			fmt.Println("")

			for account, balance := range state.Balances {
				fmt.Printf("%s: %d\n", account, balance)
			}
		},
	}
}
