package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dindin",
	Short: "dindin is a slack bot to automate dinner rotation",
	Long:  "dindin keeps track of meals eaten and meals cooked and makes sure that the members with the worst ratios are up to cook next. Additionally, it automatically sends out slack messages during cooking nights to keep track of who's eating",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
