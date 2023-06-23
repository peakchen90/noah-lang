package compiler

import (
	"errors"
	"os"
	"path/filepath"
)

type VirtualFS struct {
	Root         string
	isFileSystem bool
	files        map[string][]byte
}

func newVirtualFS(root string, isFileSystem bool) *VirtualFS {
	if isFileSystem {
		root, _ = filepath.Abs(root)
	} else if len(root) == 0 {
		root = "/virtual-fs"
	}

	return &VirtualFS{
		Root:         root,
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
