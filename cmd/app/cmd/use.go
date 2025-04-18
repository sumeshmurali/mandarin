/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sumeshmurali/mandarin/internal/config"
	"github.com/sumeshmurali/mandarin/internal/orchestrator"
)

// useCmd represents the use command
var useCmd = &cobra.Command{
	Use:   "use [path to config file]",
	Short: "Use a predefined configuration",
	Long: `Use a predefined configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a path to the config file.")
			return
		}
		configPath := args[0]

		config, err := config.ParseConfiguration(configPath)
		if err != nil {
		    fmt.Println("Error parsing configuration:", err)
		    return
		}
		fmt.Printf("Using configuration from: %s\n", configPath)

		orchestrator.Run(config)
		// fmt.Printf("Loaded configuration: %+v\n", config)
	},
}

func init() {
	rootCmd.AddCommand(useCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// useCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// useCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
