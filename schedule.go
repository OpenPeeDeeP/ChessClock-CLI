package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	pb "github.com/OpenPeeDeeP/ChessClock-CLI/chessclock"
	"github.com/spf13/cobra"
)

var scheduleLogger = mainLogger.With().Str("cmd", "schedule").Logger()
var scheduleCmd = &cobra.Command{
	Use:      "schedule [DATE]",
	Short:    "Show all of the tasks for the specified date (defaults to today)",
	Args:     validateScheduleCmdArgs,
	RunE:     scheduleCmdRun,
	PreRunE:  startClient(scheduleLogger),
	PostRunE: stopClient,
}

func scheduleCmdRun(cmd *cobra.Command, args []string) error {
	var date int64
	var err error
	if len(args) < 1 {
		date = time.Now().Unix()
	} else {
		date, err = parseDate(args[0])
		if err != nil {
			scheduleLogger.Error().Err(err).Msg("How did you get here")
			return err
		}
	}
	res, err := client.Schedule(context.Background(), &pb.ScheduleRequest{
		Date: date,
	})
	if err != nil {
		scheduleLogger.Error().Err(err).Msg("Could not get the schedule")
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	defer w.Flush()
	fmt.Fprint(w, "TIME(UTC)\tTAG\tDESCRIPTION\n")
	for _, task := range res.GetTasks() {
		tm := time.Unix(task.GetTimestamp(), 0).UTC()
		fmt.Fprintf(w, "%s\t%s\t%s\n", tm.Format("15:04:05"), task.GetTag(), task.GetDescription())
	}
	return nil
}

func validateScheduleCmdArgs(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 1 {
		err = errors.New("Too many arguments")
		scheduleLogger.Error().Err(err).Msg("Incorrect Arguments")
		return err
	}
	if len(args) != 1 {
		return nil
	}
	if _, err = parseDate(args[0]); err != nil {
		scheduleLogger.Error().Err(err).Msg("Incorrect Arguments")
		return err
	}
	return nil
}
