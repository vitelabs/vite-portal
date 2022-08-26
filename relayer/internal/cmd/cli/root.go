package cli

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/vitelabs/vite-portal/relayer/internal/app"
	"github.com/vitelabs/vite-portal/relayer/internal/cmd"
	"github.com/vitelabs/vite-portal/relayer/internal/types"
	"github.com/vitelabs/vite-portal/shared/pkg/version"
)

var (
	debug      bool
	profile    bool
	configPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s", types.AppName),
	Short: fmt.Sprintf("%s relays data requests and responses to and from Vite full nodes.", types.AppName),
	Long: fmt.Sprintf(`%s is the core component of VitePortal and responsible for relaying
	data requests and responses to and from Vite full nodes.`, types.AppName),
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cmd.Exit("Execute error", err)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "sets log level to debug")
	startCmd.Flags().BoolVar(&profile, "profile", false, "expose cpu & memory profiling")
	startCmd.Flags().StringVar(&configPath, "config", "", "path to the configuration file")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(versionCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: fmt.Sprintf("Starts %s daemon", types.AppName),
	Long:  fmt.Sprintf(`Starts the %s daemon, picks up the config from %s`, types.AppName, types.DefaultConfigFilename),
	Run: func(command *cobra.Command, args []string) {
		a, err := app.InitApp(debug, configPath)
		if err != nil {
			cmd.Exit("init error", err)
		}

		err = a.Start(profile)
		if err != nil {
			cmd.Exit("start error", err)
		}

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
			a.Shutdown()
			os.Exit(0)
		}()

		fmt.Printf("%s started\n", types.AppName)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get current version",
	Long:  `Retrieves the version`,
	Run: func(command *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", version.PROJECT_BUILD_VERSION)
	},
}
