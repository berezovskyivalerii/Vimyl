package files

import "os"

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
	return files, nil
}
