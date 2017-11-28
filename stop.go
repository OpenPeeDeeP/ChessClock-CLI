package main

import (
	"context"
	"errors"
	"strings"
	"time"

	pb "github.com/OpenPeeDeeP/ChessClock-CLI/chessclock"
	"github.com/spf13/cobra"
)

var reasons = make(map[string]pb.StopRequest_Reason)

var stopLogger = mainLogger.With().Str("cmd", "stop").Logger()
var stopCmd = &cobra.Command{
	Use:      "stop [eod|break|lunch]",
	Short:    "Stop all tasks for the specified reasons",
	Args:     validateStopCmdArgs,
	RunE:     stopCmdRun,
	PreRunE:  startClient(stopLogger),
	PostRunE: stopClient,
}

func init() {
	reasons["eod"] = pb.StopRequest_EndOfDay
	reasons["break"] = pb.StopRequest_Break
	reasons["lunch"] = pb.StopRequest_Lunch
}

func stopCmdRun(cmd *cobra.Command, args []string) error {
	_, err := client.Stop(context.Background(), &pb.StopRequest{
		Timestamp: time.Now().Unix(),
		Reason:    reasons[strings.ToLower(args[0])],
	})
	if err != nil {
		stopLogger.Error().Err(err).Msg("Could not stop the previous task")
		return err
	}
	return nil
}

func validateStopCmdArgs(cmd *cobra.Command, args []string) (err error) {
	switch {
	case len(args) > 1:
		err = errors.New("Too many arguments")
		stopLogger.Error().Err(err).Msg("Incorrect Arguments")
		return err
	case len(args) < 1:
		err = errors.New("Must sepcify a reason for stopping (eod, break, lunch)")
		stopLogger.Error().Err(err).Msg("Incorrect Arguments")
		return err
	}
	if _, ok := reasons[strings.ToLower(args[0])]; !ok {
		err = errors.New("Reason given is not a valid reason (eod, break, lunch)")
		stopLogger.Error().Err(err).Msg("Incorrect Arguments")
		return err
	}
	return nil
}
