package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

/*

TODO:
- test transaction to populate tx.db and
tx add --from=andrej --to=andrej --value=3
tx add --from=andrej --to=andrej --value=700
tx add --from=babayaga --to=andrej --value=2000
tx add --from=andrej --to=andrej --value=100 --data=reward
tx add --from=babayaga --to=andrej --value=1

- create functions descriptions, e.g:

// a does X
func a() {}

*/

func main() {
	var glancCmd = &cobra.Command{
		Use:   "glanc",
		Short: "gyors lanc - the fastest bchain ever!",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	glancCmd.AddCommand(versionCmd)
	glancCmd.AddCommand(balancesCmd())
	glancCmd.AddCommand(txCmd())

	err := glancCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
