package internal

import (
	"context"
	"fan2go-tui/internal/configuration"
	"fan2go-tui/internal/logging"
	"fan2go-tui/internal/ui"
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	"github.com/pterm/pterm"
)

func RunApplication() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var g run.Group
	{
		if configuration.CurrentConfig.Profiling.Enabled {
			g.Add(func() error {
				mux := http.NewServeMux()
				mux.HandleFunc("/debug/pprof/", pprof.Index)
				mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
				mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
				mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
				mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

				go func() {
					logging.Info("Starting profiling webserver...")
					pterm.Info.Printfln("Starting profiling webserver...")

					profilingConfig := configuration.CurrentConfig.Profiling
					address := fmt.Sprintf("%s:%d", profilingConfig.Host, profilingConfig.Port)
					err := http.ListenAndServe(address, mux)
					logging.Error("Error running profiling webserver: %v", err)
					pterm.Error.Printfln("Error running profiling webserver: %v", err)
				}()

				<-ctx.Done()
				logging.Info("Stopping profiling webserver...")
				return nil
			}, func(err error) {
				if err != nil {
					logging.Warning("Error stopping parca webserver: %v", err.Error())
					pterm.Warning.Printfln("Error stopping parca webserver: %v", err.Error())
				} else {
					logging.Debug("Webservers stopped.")
					pterm.Debug.Printfln("parca webserver stopped.")
				}
			})
		}
	}
	{
		g.Add(func() error {
			pterm.Info.Printfln("Launching UI...")
			logging.Info("Launching UI...")
			return ui.CreateUi(true).Run()
		}, func(err error) {
			if err != nil {
				logging.Warning("Error stopping UI: %v", err.Error())
				pterm.Warning.Printfln("Error stopping UI: %v", err.Error())
			} else {
				logging.Debug("UI stopped.")
				pterm.Debug.Printfln("Received SIGTERM signal, exiting...")
			}
		})
	}
	{
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		g.Add(func() error {
			<-sig
			logging.Info("Received SIGTERM signal, exiting...")
			pterm.Info.Printfln("Received SIGTERM signal, exiting...")

			return nil
		}, func(err error) {
			defer close(sig)
			cancel()
		})
	}

	if err := g.Run(); err != nil {
		logging.Error("%v", err)
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		logging.Info("Done.")
		pterm.Info.Printfln("Done.")
		os.Exit(0)
	}
}
