package cmd

import (
	"fmt"
	"os"

	"github.com/paraizofelipe/fullcycle/stress-test/internal/loadtest"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "stress-test",
	Short: "CLI para testes de carga em serviços web",
	RunE:  run,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("url", "", "URL do serviço a ser testado (obrigatório)")
	rootCmd.Flags().Int("requests", 100, "Número total de requisições")
	rootCmd.Flags().Int("concurrency", 1, "Número de chamadas simultâneas")

	_ = rootCmd.MarkFlagRequired("url")

	_ = viper.BindPFlag("url", rootCmd.Flags().Lookup("url"))
	_ = viper.BindPFlag("requests", rootCmd.Flags().Lookup("requests"))
	_ = viper.BindPFlag("concurrency", rootCmd.Flags().Lookup("concurrency"))
}

func run(cmd *cobra.Command, args []string) error {
	url := viper.GetString("url")
	requests := viper.GetInt("requests")
	concurrency := viper.GetInt("concurrency")

	if requests <= 0 {
		return fmt.Errorf("--requests deve ser maior que 0")
	}
	if concurrency <= 0 {
		return fmt.Errorf("--concurrency deve ser maior que 0")
	}
	if concurrency > requests {
		concurrency = requests
	}

	fmt.Printf("Iniciando stress test...\n")
	fmt.Printf("  URL:         %s\n", url)
	fmt.Printf("  Requisições: %d\n", requests)
	fmt.Printf("  Concorrência: %d\n\n", concurrency)

	report := loadtest.Run(url, requests, concurrency)
	report.Print()

	return nil
}
