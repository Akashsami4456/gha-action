package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"gha-publish-evidence/internal/evidence"

	"github.com/spf13/cobra"
)

var (
	cmd = &cobra.Command{
		Use:   "gha-publish-evidence-action",
		Short: "Publish the step level evidence to CloudBees Build Platform",
		Long:  "Publish the step level evidence to CloudBees Build Platform",
		RunE:  run,
	}
	cfg evidence.Config
)

func Execute() error {
	return cmd.Execute()
}

func init() {
	setDefaultValues(&cfg)
}

func setDefaultValues(cfg *evidence.Config) {

	cfg.Format = evidence.Markdown

}

func run(_ *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unknown arguments: %v", args)
	}
	newContext, cancel := context.WithCancel(context.Background())
	osChannel := make(chan os.Signal, 1)
	signal.Notify(osChannel, os.Interrupt)
	go func() {
		<-osChannel
		cancel()
	}()

	return cfg.Run(newContext)
}
