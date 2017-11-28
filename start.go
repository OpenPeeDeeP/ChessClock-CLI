package main

import (
	"context"
	"errors"
	"time"

	pb "github.com/OpenPeeDeeP/ChessClock-CLI/chessclock"
	"github.com/spf13/cobra"
)

var startLogger = mainLogger.With().Str("cmd", "tags").Logger()
var startCmd = &cobra.Command{
	Use:      "start [TAG]",
	Short:    "Start a specific task",
	Args:     validateStartCmdArgs,
	RunE:     startCmdRun,
	PreRunE:  startClient(startLogger),
	PostRunE: stopClient,
}

var (
	startDescription string
	startDuration    string
)

func init() {
	startCmd.Flags().StringVarP(&startDescription, "description", "d", "", "Description for the task")
	startCmd.Flags().StringVarP(&startDuration, "prior", "p", "", "Start a task that was started in the past (IE 30m)")
}

func startCmdRun(cmd *cobra.Command, args []string) error {
	t := time.Now()
	if startDuration != "" {
		dur, _ := time.ParseDuration(startDuration)
		t = t.Add(-1 * dur)
	}
	_, err := client.Start(context.Background(), &pb.StartRequest{
		Timestamp:   t.Unix(),
		Tag:         args[0],
		Description: startDescription,
	})
	if err != nil {
		startLogger.Error().Err(err).Msg("Could not start task")
		return err
	}
	return nil
}

func validateStartCmdArgs(cmd *cobra.Command, args []string) (err error) {
	switch {
	case len(args) > 1:
		err = errors.New("Too many arguments")
		startLogger.Error().Err(err).Msg("Incorrect Arguments")
		return err
	case len(args) < 1:
		err = errors.New("Must sepcify a tag for the task")
		startLogger.Error().Err(err).Msg("Incorrect Arguments")
		return err
	}

	if startDuration != "" {
		if _, err := time.ParseDuration(startDuration); err != nil {
			startLogger.Error().Err(err).Msg("Unable to parse duration")
			return err
		}
	}
	return nil
}
