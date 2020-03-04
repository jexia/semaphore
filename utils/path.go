package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// A FileInfo describes a file
type FileInfo struct {
	os.FileInfo
	Path string
}

// ReadDir reads the given path and returns all available files matching the given extention
func ReadDir(path string, recursive bool, ext string) (files []FileInfo, _ error) {
	list, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, file := range list {
		if file.IsDir() && recursive {
			result, err := ReadDir(filepath.Join(path, file.Name()), recursive, ext)
			if err != nil {
				return nil, err
			}

			files = append(files, result...)
			continue
		}

		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ext {
			continue
		}

		files = append(files, FileInfo{
			FileInfo: file,
			Path:     path,
		})
	}

	return files, nil
}
