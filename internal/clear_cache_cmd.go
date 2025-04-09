package internal

import (
	"os"

	"github.com/kellegous/poop"
	"github.com/spf13/cobra"
)

func clearCacheCmd(f *Flags) *cobra.Command {
	return &cobra.Command{
		Use:   "clear-cache",
		Short: "Clear the cache",
		Run: func(cmd *cobra.Command, args []string) {
			if err := os.RemoveAll(f.CacheDir); err != nil {
				poop.HitFan(err)
			}
		},
	}
}
