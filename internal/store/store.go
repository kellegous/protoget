package store

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kellegous/poop"
	"github.com/kellegous/protoget/internal"
)

const ManifestFile = "protoget.yaml"

type Store struct {
	root string
}

func Open(root string) (*Store, error) {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err := os.MkdirAll(root, 0755); err != nil {
			return nil, err
		}
	}

	return &Store{root: root}, nil
}

func (s *Store) Ensure(ctx context.Context, dep internal.Dep) (*Bundle, error) {
	// this only applies if the ref is a hash
	dst := filepath.Join(s.root, dep.Path())
	if _, err := os.Stat(dst); err == nil {
		return &Bundle{dep: dep, path: dst}, nil
	}

	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, poop.Chain(err)
	}
	defer os.RemoveAll(dir)

	sha, err := getCloneTo(ctx, dep, dir)
	if err != nil {
		return nil, poop.Chain(err)
	}

	// now that we have the hash, re-check if the bundle exists
	dep = dep.WithRef(sha)
	dst = filepath.Join(s.root, dep.Path())
	if _, err := os.Stat(dst); err == nil {
		return &Bundle{dep: dep, path: dst}, nil
	}

	mf, err := readManifestFile(filepath.Join(dir, ManifestFile))
	if err != nil {
		return nil, poop.Chain(err)
	}

	tmp := dst + ".tmp"

	if err := mf.archiveTo(tmp, dir); err != nil {
		return nil, poop.Chain(err)
	}

	if err := os.Rename(tmp, dst); err != nil {
		return nil, poop.Chain(err)
	}

	return &Bundle{dep: dep, path: dst}, nil
}

func getCloneTo(
	ctx context.Context,
	dep internal.Dep,
	root string,
) (string, error) {
	// init the repo
	if err := exec.CommandContext(ctx, "git", "init", root).Run(); err != nil {
		return "", poop.Chain(err)
	}

	// add the remote
	c := exec.CommandContext(ctx, "git", "remote", "add", "origin", dep.URL())
	c.Dir = root
	if err := c.Run(); err != nil {
		return "", poop.Chain(err)
	}

	// shallow fetch the ref
	c = exec.CommandContext(ctx, "git", "fetch", "--depth", "1", "origin", dep.Ref())
	c.Dir = root
	if err := c.Run(); err != nil {
		return "", poop.Chain(err)
	}

	// checkout the ref
	c = exec.CommandContext(ctx, "git", "checkout", "FETCH_HEAD")
	c.Dir = root
	if err := c.Run(); err != nil {
		return "", poop.Chain(err)
	}

	c = exec.CommandContext(ctx, "git", "rev-parse", "HEAD")
	c.Dir = root
	out, err := c.Output()
	if err != nil {
		return "", poop.Chain(err)
	}

	return string(bytes.TrimSpace(out)), nil
}
