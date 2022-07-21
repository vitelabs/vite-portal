package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/vitelabs/vite-portal/internal/app"
	"github.com/vitelabs/vite-portal/internal/cmd"
	"github.com/vitelabs/vite-portal/internal/rpc"
	"github.com/vitelabs/vite-portal/internal/types"
)

var (
	debug   bool
	profile bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "portal-relayer",
	Short: "portal-relayer relays data requests and responses to and from Vite full nodes.",
	Long: `portal-relayer is the core component of VitePortal and responsible for relaying
	data requests and responses to and from Vite full nodes.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cmd.Exit("Execute error", err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "sets log level to debug")
	startCmd.Flags().BoolVar(&profile, "profile", false, "expose cpu & memory profiling")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(stopCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: fmt.Sprintf("Starts %s daemon", app.AppName),
	Long:  fmt.Sprintf(`Starts the %s daemon, picks up the config from the assigned <datadir>`, app.AppName),
	Run: func(command *cobra.Command, args []string) {
		if err := app.InitApp(debug); err != nil {
			cmd.Exit("start error", err)
		}

		go rpc.StartHttpRpc(types.GlobalConfig.RpcHttpPort, types.GlobalConfig.RpcTimeout, debug, profile)

		// trap kill signals
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
			os.Kill,
			os.Interrupt)

		defer func() {
			sig := <-signalChannel
			fmt.Printf("Exit signal %s received\n", sig)
			app.Shutdown()
			os.Exit(0)
		}()

		fmt.Printf("%s started\n", app.AppName)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get current version",
	Long:  `Retrieves the version`,
	Run: func(command *cobra.Command, args []string) {
		fmt.Printf("AppVersion: %s\n", app.AppVersion)
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: fmt.Sprintf("Stops %s daemon", app.AppName),
	Long:  fmt.Sprintf(`Stops the %s daemon`, app.AppName),
	Run: func(command *cobra.Command, args []string) {
		fmt.Printf("%s stopped\n", app.AppName)
	},
}
