package compiler

import (
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

type VirtualFS struct {
	Root         string
	PackageRoot  string
	isFileSystem bool
	files        map[string][]byte
}

func newVirtualFS(root string, isFileSystem bool) *VirtualFS {
	packageRoot := ""

	if isFileSystem {
		_root, err := filepath.Abs(root)
		if err != nil {
			panic(err)
		}
		root = _root

		current, err := user.Current()
		if err != nil {
			panic(err)
		}
		packageRoot = filepath.Join(current.HomeDir, ".noah/packages")
	} else {
		if len(root) == 0 {
			root = "/noah-virtual:root"
		}
		packageRoot = "/noah-virtual:packages"
	}

	return &VirtualFS{
		Root:         root,
		PackageRoot:  packageRoot,
		isFileSystem: isFileSystem,
		files:        make(map[string][]byte),
	}
}

func (v *VirtualFS) ReadFile(filename string) ([]byte, error) {
	if v.isFileSystem {
		return os.ReadFile(filename)
	}

	buffer, has := v.files[filename]
	if has {
		return buffer, nil
	}
	return nil, errors.New("No such file: " + filename)
}

func (v *VirtualFS) WriteFile(filename string, buffer []byte) error {
	if v.isFileSystem {
		return os.WriteFile(filename, buffer, os.ModePerm)
	}

	v.files[filename] = buffer
	return nil
}

func (v *VirtualFS) Remove(filename string) error {
	if v.isFileSystem {
		return os.Remove(filename)
	}

	delete(v.files, filename)
	return nil
}

func (v *VirtualFS) ExistFile(filename string) bool {
	if v.isFileSystem {
		s, err := os.Stat(filename)
		return err == nil && !s.IsDir()
	}

	_, has := v.files[filename]
	return has
}
