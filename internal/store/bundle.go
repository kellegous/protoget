package store

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"

	"github.com/kellegous/poop"

	"github.com/kellegous/protoget/internal"
)

type Bundle struct {
	dep  *internal.Dep
	path string
}

func (b *Bundle) CloneTo(dst string) error {
	r, err := os.Open(b.path)
	if err != nil {
		return poop.Chain(err)
	}
	defer r.Close()

	gr, err := gzip.NewReader(r)
	if err != nil {
		return poop.Chain(err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return poop.Chain(err)
		}

		if err := cloneFileTo(filepath.Join(dst, hdr.Name), tr); err != nil {
			return poop.Chain(err)
		}
	}

	return nil
}

func cloneFileTo(dst string, r io.Reader) error {
	dir := filepath.Dir(dst)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return poop.Chain(err)
		}
	}

	w, err := os.Create(dst)
	if err != nil {
		return poop.Chain(err)
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return poop.Chain(err)
}
