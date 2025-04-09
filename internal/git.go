package internal

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kellegous/poop"
)

func CloneTo(
	ctx context.Context,
	dep *Dep,
	root string,
) (string, error) {
	dir := filepath.Join(root, dep.path, dep.ref)

	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return dir, nil
	}

	// create the directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", poop.Chain(err)
	}

	// init the repo
	if err := exec.CommandContext(ctx, "git", "init", dir).Run(); err != nil {
		return "", errors.Join(err, os.RemoveAll(dir))
	}

	// add the remote
	c := exec.CommandContext(ctx, "git", "remote", "add", "origin", dep.URL())
	c.Dir = dir
	if err := c.Run(); err != nil {
		return "", errors.Join(poop.ChainWith(err, "git remote add origin"), os.RemoveAll(dir))
	}

	// shallow fetch the ref
	c = exec.CommandContext(ctx, "git", "fetch", "--depth", "1", "origin", dep.ref)
	c.Dir = dir
	if err := c.Run(); err != nil {
		return "", errors.Join(poop.ChainWithf(err, "git fetch %s", dep.ref), os.RemoveAll(dir))
	}

	// checkout the ref
	c = exec.CommandContext(ctx, "git", "checkout", "FETCH_HEAD")
	c.Dir = dir
	if err := c.Run(); err != nil {
		return "", errors.Join(poop.ChainWith(err, "git checkout FETCH_HEAD"), os.RemoveAll(dir))
	}

	return dir, nil
}
