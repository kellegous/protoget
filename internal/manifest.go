package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Name    string   `yaml:"name"`
	Sources []string `yaml:"sources"`
}

func Read(r io.Reader) (*Manifest, error) {
	var m Manifest
	if err := yaml.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

func ReadFile(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(f)
}

func copyFile(src string, dst string) error {
	dir := filepath.Dir(dst)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, r)
	return err
}

func (m *Manifest) CopySources(dstDir string, srcDir string) error {
	for _, rel := range m.Sources {
		rel = strings.TrimLeft(rel, "/")
		src := filepath.Join(srcDir, rel)
		dst := filepath.Join(dstDir, rel)

		fmt.Printf("%s -> %s\n", src, dst)
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}
	return nil
}
