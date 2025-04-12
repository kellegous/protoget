package cmd

import (
	"os"

	"github.com/kellegous/poop"
	"github.com/kellegous/protoget/internal"
	"github.com/kellegous/protoget/internal/store"
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	var flags Flags
	cmd := &cobra.Command{
		Use:  "protoget",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			s, err := store.Open(flags.CacheDir)
			if err != nil {
				poop.HitFan(err)
			}

			for _, arg := range args {
				dep, err := internal.ParseDep(arg)
				if err != nil {
					poop.HitFan(err)
				}

				b, err := s.Ensure(cmd.Context(), &dep)
				if err != nil {
					poop.HitFan(err)
				}

				if err := b.CloneTo(flags.DestDir); err != nil {
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
