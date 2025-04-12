package store

import "github.com/kellegous/protoget/internal"

type Bundle struct {
	dep  *internal.Dep
	path string
}

func (b *Bundle) CloneTo(dst string) error {
	// TODO(kellegous): This is what puts the protos in the right place.
	return nil
}
