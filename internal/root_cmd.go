package internal

import (
	"context"
	"os"
	"path/filepath"

	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

const manifest = "protoget.yaml"

func rootCmd() *cobra.Command {
	var flags Flags
	cmd := &cobra.Command{
		Use:  "protoget",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				dep, err := parseDep(arg)
				if err != nil {
					poop.HitFan(err)
				}

				path, err := CloneTo(context.Background(), &dep, flags.CacheDir)
				if err != nil {
					poop.HitFan(err)
				}

				manifest, err := ReadFile(filepath.Join(path, manifest))
				if err != nil {
					poop.HitFan(err)
				}

				if err := manifest.CopySources(flags.DestDir, path); err != nil {
					poop.HitFan(err)
				}
			}
		},
	}

	flags.Register(cmd.PersistentFlags())

	cmd.AddCommand(clearCacheCmd(&flags))

	return cmd
}

func Execute() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
