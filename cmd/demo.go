package cmd

import (
	"context"
	"fan2go-tui/internal"
	"fan2go-tui/internal/configuration"
	"fan2go-tui/internal/demo"
	"fan2go-tui/internal/logging"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	demoHost   string
	demoPort   int
	demoWithUi bool
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a simulated fan2go server with generated data",
	Long: `Start a mock fan2go API server that simulates sensors, curves, and fans with
realistic dynamic data.

Can be run standalone so that 'fan2go-tui' (or other clients) can connect to it,
or with the '--ui' flag to launch the TUI directly connected to the simulated server.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath := configuration.DetectAndReadConfigFile()
		logging.Info("Using configuration file at: %s", configPath)
		configuration.LoadConfig()

		host := demoHost
		port := demoPort

		if !cmd.Flags().Changed("host") && configuration.CurrentConfig.Api.Host != "" {
			host = configuration.CurrentConfig.Api.Host
		}
		if !cmd.Flags().Changed("port") && configuration.CurrentConfig.Api.Port != 0 {
			port = configuration.CurrentConfig.Api.Port
		}

		sim := demo.NewSimulator()
		server := demo.NewServer(host, port, sim)

		if err := server.Start(); err != nil {
			logging.Error("Failed to start demo server: %v", err)
			pterm.Error.Printfln("Failed to start demo server: %v", err)
			os.Exit(1)
		}

		actualPort := server.Port()
		addr := net.JoinHostPort(host, strconv.Itoa(actualPort))

		if demoWithUi {
			configuration.CurrentConfig.Api.Host = host
			configuration.CurrentConfig.Api.Port = actualPort

			defer func() {
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				_ = server.Shutdown(shutdownCtx)
			}()

			internal.RunApplication()
			return
		}

		pterm.DefaultHeader.WithFullWidth().Println("fan2go-tui demo server")
		pterm.Info.Printfln("Simulated fan2go server running on http://%s", addr)
		pterm.Println()
		pterm.DefaultSection.Println("Endpoints")
		pterm.DefaultBulletList.WithItems([]pterm.BulletListItem{
			{Level: 0, Text: fmt.Sprintf("GET http://%s/fan      (Fans status & config)", addr)},
			{Level: 0, Text: fmt.Sprintf("GET http://%s/curve    (Curves status & config)", addr)},
			{Level: 0, Text: fmt.Sprintf("GET http://%s/sensor   (Sensors status & config)", addr)},
		}).Render()
		pterm.Println()
		pterm.Info.Println("Simulating 4 fans, 4 curves, and 4 sensors with dynamic realistic data.")
		pterm.Info.Println("Tip: Run 'fan2go-tui' in another terminal to view the UI, or run 'fan2go-tui demo --ui'.")
		pterm.Println()
		pterm.Warning.Println("Press Ctrl+C to stop.")

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-sig

		pterm.Println()
		pterm.Info.Println("Shutting down simulated fan2go server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logging.Warning("Error stopping server: %v", err)
		}
		pterm.Info.Println("Server stopped.")
	},
}

func init() {
	demoCmd.Flags().StringVar(&demoHost, "host", "127.0.0.1", "Host to listen on")
	demoCmd.Flags().IntVarP(&demoPort, "port", "p", 9001, "Port to listen on")
	demoCmd.Flags().BoolVar(&demoWithUi, "ui", false, "Also launch the terminal UI connected to the demo server")

	rootCmd.AddCommand(demoCmd)
}
