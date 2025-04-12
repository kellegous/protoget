package internal

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Dep struct {
	path string
	ref  string
}

func (d Dep) URL() string {
	return fmt.Sprintf("git@github.com:%s.git", d.path)
}

func (d Dep) Ref() string {
	return d.ref
}

func (d Dep) Path() string {
	return filepath.Join(d.path, d.ref)
}

func ParseDep(s string) (Dep, error) {
	parts := strings.SplitN(s, "@", 2)

	if len(parts) != 2 {
		return Dep{}, fmt.Errorf("invalid dep: %s", s)
	}

	if !strings.HasPrefix(parts[0], "github.com/") {
		return Dep{}, fmt.Errorf("invalid dep: %s", s)
	}

	return Dep{
		path: strings.TrimPrefix(parts[0], "github.com/"),
		ref:  parts[1],
	}, nil
}
