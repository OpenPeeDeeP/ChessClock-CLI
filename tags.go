package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/OpenPeeDeeP/ChessClock-CLI/chessclock"
	"github.com/spf13/cobra"
)

var tagsLogger = mainLogger.With().Str("cmd", "tags").Logger()
var tagsCmd = &cobra.Command{
	Use:      "tags [DATE(YYYY/MM/DD)]",
	Short:    "list all of the flags for the date (default is current date)",
	Args:     validateTagsCmdArgs,
	RunE:     tagsCmdRun,
	PreRunE:  startClient(tagsLogger),
	PostRunE: stopClient,
}

func tagsCmdRun(cmd *cobra.Command, args []string) error {
	var date int64
	var err error
	if len(args) < 1 {
		date = time.Now().Unix()
	} else {
		date, err = parseDate(args[0])
		if err != nil {
			tagsLogger.Error().Err(err).Msg("How did you get here")
			return err
		}
	}
	res, err := client.ListTags(context.Background(), &pb.ListTagsRequest{
		Date: date,
	})
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			if len(args) < 1 {
				tagsLogger.Error().Msg("No task started for today")
			} else {
				tagsLogger.Error().Msg("No task for that date")
			}
			return err
		}
		tagsLogger.Error().Err(err).Msg("Could not get a list of tags")
		return err
	}
	for _, tag := range res.GetTags() {
		fmt.Println(tag)
	}
	return nil
}

func validateTagsCmdArgs(cmd *cobra.Command, args []string) (err error) {
	if len(args) > 1 {
		err = errors.New("Too many arguments")
		tagsLogger.Error().Err(err).Msg("Incorrect Arguments")
		return err
	}
	if len(args) != 1 {
		return nil
	}
	if _, err = parseDate(args[0]); err != nil {
		tagsLogger.Error().Err(err).Msg("Incorrect Arguments")
		return err
	}
	return nil
}
