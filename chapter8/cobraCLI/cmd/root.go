package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "details",
	Short: "This project takes student informations",
	Long:  `A long string about description`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := cmd.PersistentFlags().Lookup("name").Value
		age := cmd.PersistentFlags().Lookup("age").Value
		log.Printf("Hello %s (%s years), welcome to the command line worlds", name, age)
	},
}

// Execute is the cobra logic starting point
func Execute() {
	rootCmd.PersistentFlags().StringP("name", "n", "stranger", "Name of the student")
	rootCmd.PersistentFlags().IntP("age", "a", 25, "Age of the student")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
