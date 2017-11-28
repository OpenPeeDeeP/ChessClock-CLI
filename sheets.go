package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/OpenPeeDeeP/ChessClock-CLI/chessclock"
	"github.com/spf13/cobra"
)

var sheetsLogger = mainLogger.With().Str("cmd", "sheets").Logger()
var sheetsCmd = &cobra.Command{
	Use:      "sheets",
	Short:    "list all of the timesheets",
	Args:     validateSheetsCmdArgs,
	RunE:     sheetsCmdRun,
	PreRunE:  startClient(sheetsLogger),
	PostRunE: stopClient,
}

func sheetsCmdRun(cmd *cobra.Command, args []string) error {
	sheets, err := client.ListTimeSheets(context.Background(), &pb.ListTimeSheetsRequest{})
	if err != nil {
		sheetsLogger.Error().Err(err).Msg("Could not get a list of time sheets")
		return err
	}
	for _, sheet := range sheets.GetDates() {
		date := time.Unix(sheet, 0).UTC()
		fmt.Printf("%04d/%02d/%02d\n", date.Year(), date.Month(), date.Day())
	}
	return nil
}

func validateSheetsCmdArgs(cmd *cobra.Command, args []string) (err error) {
	switch {
	case len(args) > 0:
		err = errors.New("Too many arguments")
		sheetsLogger.Error().Err(err).Msg("Incorrect Arguments")
		return err
	default:
		return nil
	}
}
