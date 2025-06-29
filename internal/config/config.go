package config

import (
	"fmt"

	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SnakesConfig struct {
	RunAddress           string `mapstructure:"RUN_ADDRESS"`
	DatabaseURI          string `mapstructure:"DATABASE_URI"`
	AccrualSystemAddress string `mapstructure:"ACCRUAL_SYSTEM_ADDRESS"`
}

var Config SnakesConfig

func InitConfig() {
	var rootCmd = &cobra.Command{
		Use: "gophermart",
		Run: func(cmd *cobra.Command, args []string) {
			if err := viper.Unmarshal(&Config); err != nil {
				fmt.Printf("Unable to decode into struct, %v\n", err)
				return
			}
			fmt.Println("RunAddress:", Config.RunAddress)
			fmt.Println("DatabaseURI:", Config.DatabaseURI)
			fmt.Println("AccrualSystemAddress:", Config.AccrualSystemAddress)
		},
	}
	rootCmd.Flags().StringP("address", "a", ":8080", "gophermart address and port")
	rootCmd.Flags().StringP("database", "d", "", "database URI")
	rootCmd.Flags().StringP("accural", "r", ":9090", "accrual address and port")
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
}
