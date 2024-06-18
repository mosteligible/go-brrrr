package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go_brrr",
	Short: "run load test on various api endpoints",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting load test..")
	},
}

var SrvName, IndexName, QueryType, VectorField, SemanticConf, Query string
var Top, Concurrency int
var Rate, testDuration float64
var LocalEmbeddings bool

func init() {
	rootCmd.PersistentFlags().IntVar(&Concurrency, "concurrency", 0, "number of concurrent requests")
	rootCmd.PersistentFlags().Float64Var(&Rate, "rate", 0.0, "concurrent request repetition interval")
	rootCmd.PersistentFlags().Float64Var(&testDuration, "test-duration", 10.0, "number of seconds to keep the load going")
	rootCmd.MarkPersistentFlagRequired("concurrency")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
