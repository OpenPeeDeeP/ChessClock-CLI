package main

import (
	"context"
	"fmt"
	"os"

	pb "github.com/OpenPeeDeeP/ChessClock-CLI/chessclock"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var rootLogger = mainLogger.With().Str("cmd", "root").Logger()
var rootCmd = &cobra.Command{
	Short:            "Chess Clock to keep track of tasks",
	SilenceErrors:    true,
	SilenceUsage:     true,
	TraverseChildren: true,
	RunE:             rootCmdRun,
	PersistentPreRun: handlePersitentFlags,
}

var (
	isVerbose   bool
	showVersion bool
)

func init() {
	rootCmd.Use = os.Args[0]
	rootCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "Print out in verbose mode")
	rootCmd.Flags().BoolVar(&showVersion, "version", false, "Display the version of the cli and daemon")
}

func rootCmdRun(cmd *cobra.Command, args []string) error {
	if showVersion {
		startClient(rootLogger)(cmd, args)
		defer stopClient(cmd, args)
		ver, err := client.Version(context.Background(), &pb.VersionRequest{})
		if err != nil {
			rootLogger.Error().Err(err).Msg("Could not get version of the daemon")
			return err
		}
		fmt.Printf("CLI version: %s\n", version)
		fmt.Printf("Daemon version: %s\n", ver.GetVersion())
		return nil
	}
	return cmd.Help()
}

func handlePersitentFlags(cmd *cobra.Command, args []string) {
	if isVerbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}
