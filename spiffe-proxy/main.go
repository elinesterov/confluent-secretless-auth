package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"go.uber.org/zap"
)

type Config struct {
	SpiffeEndpointSocket string
	Port                 int
}

type OAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	// Optional
	ExpiresIn int64 `json:"expires_in,omitempty"`
}

func loggingMiddleware(next http.HandlerFunc, logger *zap.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Received request",
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
			zap.String("remote_addr", r.RemoteAddr),
		)

		// Call the next handler
		next(w, r)
	}
}

func getJwtHandler(w http.ResponseWriter, r *http.Request, jwtSource *workloadapi.JWTSource, logger *zap.Logger) {
	audience := "confluent.io"
	svid, err := jwtSource.FetchJWTSVID(r.Context(), jwtsvid.Params{
		Audience: audience,
	})
	if err != nil {
		logger.Error("Error fetching JWT SVID", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	expiresIn := int64(time.Until(svid.Expiry).Seconds())
	oauthResp := OAuthResponse{
		AccessToken: svid.Marshal(),
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
	}

	jsonResp, err := json.Marshal(oauthResp)
	if err != nil {
		logger.Error("Error serializing JSON", zap.Error(err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Could not create logger:", err)
		os.Exit(1)
	}
	defer logger.Sync()

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
			logger.Info("Starting server", zap.Int("port", cfg.Port), zap.String("spiffe_endpoint_socket", cfg.SpiffeEndpointSocket))

			var socketFilePath string
			if strings.HasPrefix(cfg.SpiffeEndpointSocket, "unix://") {
				socketFilePath = strings.TrimPrefix(cfg.SpiffeEndpointSocket, "unix://")
			} else {
				socketFilePath = cfg.SpiffeEndpointSocket
			}

			if _, err := os.Stat(socketFilePath); os.IsNotExist(err) {
				logger.Fatal("SPIFFE Workload API socket does not exist", zap.String("path", socketFilePath))
			}

			// Initialize the Workload API client
			client, err := workloadapi.New(context.Background(), workloadapi.WithAddr("unix://"+socketFilePath))
			if err != nil {
				logger.Fatal("Could not connect to SPIFFE Workload API", zap.Error(err))
			}
			defer client.Close()

			// JWT Source is a blocking call
			jwtSource, err := workloadapi.NewJWTSource(context.Background(), workloadapi.WithClient(client))
			if err != nil {
				logger.Fatal("Could not create JWT source", zap.Error(err))
			}
			defer jwtSource.Close()

			http.HandleFunc("/api/getjwt", loggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
				getJwtHandler(w, r, jwtSource, logger)
			}, logger))

			err = http.ListenAndServe(":"+strconv.Itoa(cfg.Port), nil)
			if err != nil {
				logger.Fatal("Failed to start HTTP server", zap.Error(err))
			}
		},
	}

	rootCmd.AddCommand(serverCmd)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Error executing root command", zap.Error(err))
	}
}
