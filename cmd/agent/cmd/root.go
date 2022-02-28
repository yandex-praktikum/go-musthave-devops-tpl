package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultServerAddress  = "127.0.0.1:8080"
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
	defaultServerTimeout  = 1 * time.Second
)

var (
	rootCmd = &cobra.Command{
		Use:   "agent",
		Short: "Simple metrics agent for learning purposes",
		Long:  `Start the agnet and enjoy a lot of metrics!`,
	}
	ServerAddress  string
	PollInterval   time.Duration
	ReportInterval time.Duration
	ServerTimeout  time.Duration
)

func init() {
	rootCmd.Flags().StringVarP(&ServerAddress, "address", "a", defaultServerAddress,
		"Pair of ip:port to connect to")

	rootCmd.Flags().DurationVarP(&ServerTimeout, "timeout", "t", defaultServerTimeout,
		"Timeout for server connection")

	rootCmd.Flags().DurationVarP(&ReportInterval, "report", "r", defaultReportInterval,
		"Number of seconds to periodically report metrics")

	rootCmd.Flags().DurationVarP(&PollInterval, "poll", "p", defaultPollInterval,
		"Number of seconds to periodically get metrics")
}

func Execute() error {
	return rootCmd.Execute()
}
