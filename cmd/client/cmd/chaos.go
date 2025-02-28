package cmd

import (
	"github.com/TheCrabilia/chaos-shortener/internal/client"
	"github.com/spf13/cobra"
)

var (
	latencyRate float64
	errorRate   float64
	outageRate  float64
)

var chaosCmd = &cobra.Command{
	Use: "chaos",
	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(server)
		c.UpdateChaosSettings(&client.UpdateChaosSettingsOpts{
			LatencyRate: latencyRate,
			ErrorRate:   errorRate,
			OutageRate:  outageRate,
		})
	},
}

func init() {
	chaosCmd.Flags().StringVarP(&server, "server", "S", "http://localhost:8080", "server to connect to")
	chaosCmd.Flags().Float64VarP(&latencyRate, "latency-rate", "l", 0.0, "latency rate")
	chaosCmd.Flags().Float64VarP(&errorRate, "error-rate", "e", 0.0, "error rate")
	chaosCmd.Flags().Float64VarP(&outageRate, "outage-rate", "o", 0.0, "outage rate")
}
