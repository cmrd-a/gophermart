package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cmrd-a/gophermart/internal/api"
	"github.com/cmrd-a/gophermart/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "gophermart",
		Run: func(cmd *cobra.Command, args []string) {
			if err := viper.Unmarshal(&config.Config); err != nil {
				fmt.Printf("Unable to decode into struct, %v\n", err)
				return
			}
			fmt.Println("RunAddress:", config.Config.RunAddress)
			fmt.Println("DatabaseURI:", config.Config.DatabaseURI)
			fmt.Println("AccrualSystemAddress:", config.Config.AccrualSystemAddress)
		},
	}
	rootCmd.Flags().StringP("address", "a", ":9090", "gophermart address and port")
	rootCmd.Flags().StringP("database", "d", "", "database URI")
	rootCmd.Flags().StringP("accural", "r", ":8080", "accrual address and port")
	viper.BindPFlag("RUN_ADDRESS", rootCmd.Flags().Lookup("address"))
	viper.BindPFlag("DATABASE_URI", rootCmd.Flags().Lookup("database"))
	viper.BindPFlag("ACCRUAL_SYSTEM_ADDRESS", rootCmd.Flags().Lookup("accural"))

	// Bind flags with Viper
	viper.BindPFlags(rootCmd.Flags())

	viper.AutomaticEnv()

	// Bind environment variables
	viper.BindEnv("RUN_ADDRESS")
	viper.BindEnv("DATABASE_URI")
	viper.BindEnv("ACCRUAL_SYSTEM_ADDRESS")

	// Set the configuration file name and path
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".") // Search in the working directory

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
	}

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	r := api.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":9090")
	if err != nil {
		log.Fatal(err.Error())
	}

}
