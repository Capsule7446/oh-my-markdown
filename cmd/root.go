package cmd

import (
	"github.com/Capsule7446/oh-my-markdown/internal/logger"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "omm",
	Short:            "批次處理 Markdown 的工具",
	SilenceErrors:    true,
	SilenceUsage:     true,
	PersistentPreRunE: initLogger,
}

func initLogger(cmd *cobra.Command, args []string) error {
	return logger.Init()
}

func Execute() error {
	return rootCmd.Execute()
}
