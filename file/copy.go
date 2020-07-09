package file

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Copy copies src to dst.
func Copy(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	return tryCopy(src, dst, info)
}

func tryCopy(src, dst string, info os.FileInfo) error {
	if info.IsDir() {
		return directoryCopy(src, dst, info)
	}
	return fileCopy(src, dst, info)
}

func fileCopy(src, dst string, info os.FileInfo) error {
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer safeClose(f.Close)

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer safeClose(s.Close)

	_, err = io.Copy(f, s)
	return err
}

func directoryCopy(src, dst string, info os.FileInfo) error {
	if err := os.MkdirAll(dst, info.Mode()); err != nil {
		return err
	}

	infos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if err := tryCopy(filepath.Join(src, info.Name()), filepath.Join(dst, info.Name()), info); err != nil {
			return err
		}
	}

	return nil
}

func safeClose(fn func() error) {
	if err := fn(); err != nil {
		log.Println(err)
	}
}
