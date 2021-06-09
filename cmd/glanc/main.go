package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

/*

TODO:
- Document types with golang // <StructName> description style
*/

func main() {
	var glancCmd = &cobra.Command{
		Use:   "glanc",
		Short: "gyors lanc - the fastest bchain ever!",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	glancCmd.AddCommand(versionCmd())
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
