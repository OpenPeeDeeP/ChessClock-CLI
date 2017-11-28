package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "github.com/OpenPeeDeeP/ChessClock-CLI/chessclock"
	"github.com/ianschenck/envflag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var version string

var (
	connection      *grpc.ClientConn
	client          pb.ChessClockClient
	mainLogger      = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	daemonConString string
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	rootCmd.AddCommand(startCmd, stopCmd, tagsCmd, sheetsCmd, scheduleCmd, tallyCmd)
	envflag.StringVar(&daemonConString, "CCD_CONNECTION_STRING", "localhost:4242", "Connection string to the daemon")
}

func main() {
	envflag.Parse()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func startClient(log zerolog.Logger) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var err error
		log.Debug().Str("con", daemonConString).Msg("Connecting to daemon")
		connection, err = grpc.Dial(daemonConString, grpc.WithInsecure())
		if err != nil {
			log.Error().Err(err).Msg("Could not connect to the daemon")
			return err
		}
		client = pb.NewChessClockClient(connection)
		return nil
	}
}

func stopClient(cmd *cobra.Command, args []string) error {
	log.Debug().Msg("Disconnecting from the daemon")
	return connection.Close()
}

func parseDate(date string) (int64, error) {
	dateString := strings.Split(date, "/")
	if len(dateString) != 3 {
		return 0, errors.New("Date not in the proper format (YYYY/MM/DD)")
	}
	dateInts := make([]int, 0, 3)
	for _, ds := range dateString {
		i, err := strconv.Atoi(ds)
		if err != nil {
			return 0, errors.New("Date not in the proper format (YYYY/MM/DD)")
		}
		dateInts = append(dateInts, i)
	}
	dateTime := time.Date(dateInts[0], time.Month(dateInts[1]), dateInts[2], 0, 0, 0, 0, time.Local).Unix()
	return dateTime, nil
}
