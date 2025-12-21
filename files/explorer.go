package files

import (
	"os"
	"time"
)

type VirtualEntry struct {
	name  string
	isDir bool
}

func (v VirtualEntry) Name() string       { return v.name }
func (v VirtualEntry) Size() int64        { return 0 }
func (v VirtualEntry) Mode() os.FileMode  { return os.ModeDir }
func (v VirtualEntry) ModTime() time.Time { return time.Now() }
func (v VirtualEntry) IsDir() bool        { return v.isDir }
func (v VirtualEntry) Sys() interface{}   { return nil }

// Returs a list of files by the path
func ListFiles(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	entries, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	files := make([]os.FileInfo, 0, 10)
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry)
		}
	}
	return files, nil
}

// Returs a list of dirs by the path
func ListDirs(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	entries, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	files := make([]os.FileInfo, 0, 10)
	for _, entry := range entries {
		if entry.IsDir() {
			files = append(files, entry)
		}
	}
	return addGeneral(files), nil
}

func addGeneral(list []os.FileInfo) []os.FileInfo {
	general := VirtualEntry{name: "General", isDir: true}
	return append([]os.FileInfo{general}, list...)
}
