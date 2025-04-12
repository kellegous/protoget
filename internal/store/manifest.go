package store

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/kellegous/poop"
	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Name    string   `yaml:"name"`
	Sources []string `yaml:"sources"`
}

func (m *Manifest) archiveTo(
	dst string,
	root string,
) error {
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

	gw := gzip.NewWriter(w)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, source := range m.Sources {
		source = strings.TrimLeft(source, "/")
		src := filepath.Join(root, source)

		s, err := os.Stat(src)
		if err != nil {
			return poop.Chain(err)
		}

		if err := tw.WriteHeader(&tar.Header{
			Name: source,
			Mode: 0755,
			Size: s.Size(),
		}); err != nil {
			return poop.Chain(err)
		}

		if err := copyFileTo(tw, src); err != nil {
			return poop.Chain(err)
		}
	}
	return nil
}

func copyFileTo(w io.Writer, src string) error {
	r, err := os.Open(src)
	if err != nil {
		return poop.Chain(err)
	}
	defer r.Close()

	if _, err := io.Copy(w, r); err != nil {
		return poop.Chain(err)
	}
	return nil
}

func readManifest(r io.Reader) (*Manifest, error) {
	var m Manifest
	if err := yaml.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func readManifestFile(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return readManifest(f)
}
