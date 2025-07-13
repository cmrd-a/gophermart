package config

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SnakesConfig struct {
	RunAddress           string `mapstructure:"RUN_ADDRESS"`
	DatabaseURI          string `mapstructure:"DATABASE_URI"`
	AccrualSystemAddress string `mapstructure:"ACCRUAL_SYSTEM_ADDRESS"`
	SaltSecret           string `mapstructure:"SALT_SECRET"`
	JWTSecret            string `mapstructure:"JWT_SECRET"`
}

var Config SnakesConfig

func InitConfig() {
	var rootCmd = &cobra.Command{
		Use: "gophermart",
		Run: func(cmd *cobra.Command, args []string) {
			if err := viper.Unmarshal(&Config); err != nil {
				slog.Error("Unable to decode config into struct", "error", err)
				return
			}
			slog.Info("Configuration loaded",
				"run_address", Config.RunAddress,
				"database_uri", maskURI(Config.DatabaseURI),
				"accrual_system_address", Config.AccrualSystemAddress,
			)
		},
	}
	rootCmd.Flags().StringP("address", "a", ":8080", "gophermart address and port")
	rootCmd.Flags().StringP("database", "d", "", "database URI")
	rootCmd.Flags().StringP("accrual", "r", "http://localhost:9090", "accrual address and port")
	viper.BindPFlag("RUN_ADDRESS", rootCmd.Flags().Lookup("address"))
	viper.BindPFlag("DATABASE_URI", rootCmd.Flags().Lookup("database"))
	viper.BindPFlag("ACCRUAL_SYSTEM_ADDRESS", rootCmd.Flags().Lookup("accrual"))
	viper.BindPFlags(rootCmd.Flags())
	viper.AutomaticEnv()
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Failed to execute command", "error", err)
		os.Exit(1)
	}
}

// maskURI masks sensitive information in database URI for logging
func maskURI(uri string) string {
	if uri == "" {
		return ""
	}
	return "***masked***"
}
