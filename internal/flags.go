package internal

import (
	"os/user"
	"path/filepath"

	"github.com/spf13/pflag"
)

type Flags struct {
	CacheDir string
	DestDir  string
}

func (f *Flags) Register(s *pflag.FlagSet) {
	s.StringVar(
		&f.CacheDir,
		"cache-directory",
		mustHaveDefaultCacheDir(),
		"cache directory")

	s.StringVar(
		&f.DestDir,
		"destination-directory",
		"./external",
		"destination directory")
}

func mustHaveDefaultCacheDir() string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(u.HomeDir, ".cache/protoget")
}
