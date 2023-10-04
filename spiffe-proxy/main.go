package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	SpiffeEndpointSocket string
	Port                 int
}

func getJwtHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("JWT token placeholder"))
}

func main() {
	var cfg Config

	var rootCmd = &cobra.Command{Use: "spiffe-proxy"}

	rootCmd.PersistentFlags().StringVar(&cfg.SpiffeEndpointSocket, "SPIFFE_ENDPOINT_SOCKET", "", "SPIFFE Workload API endpoint socket")
	rootCmd.PersistentFlags().IntVar(&cfg.Port, "port", 8080, "Port to run the server on")

	// Bind flags to environment variables using Viper
	viper.BindPFlag("SPIFFE_ENDPOINT_SOCKET", rootCmd.PersistentFlags().Lookup("SPIFFE_ENDPOINT_SOCKET"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	// Set Viper to automatically look for environment variables
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	cobra.OnInitialize(func() {
		cfg.SpiffeEndpointSocket = viper.GetString("SPIFFE_ENDPOINT_SOCKET")
		cfg.Port = viper.GetInt("port")
	})

	var serverCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the SPIFFE proxy server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Server is starting on port %d\n", cfg.Port)
			fmt.Printf("Using SPIFFE Workload API socket path: %s\n", cfg.SpiffeEndpointSocket)

			http.HandleFunc("/api/getjwt", getJwtHandler)
			http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil)
		},
	}

	rootCmd.AddCommand(serverCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error executing root command:", err)
		os.Exit(1)
	}
}
