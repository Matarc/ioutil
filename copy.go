package ioutil

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

func Copy(source, destination string) (err error) {
	srcInfo, srcErr := os.Stat(source)
	if srcErr != nil {
		return srcErr
	}
	dstInfo, dstErr := os.Stat(destination)
	if srcInfo.IsDir() && dstInfo != nil && !dstInfo.IsDir() && !os.IsNotExist(dstErr) {
		return fmt.Errorf("Cannot overwrite non-directory '%s' with directory '%s'", destination, source)
	}
	dstIsDir := false
	if dstInfo != nil {
		dstIsDir = dstInfo.IsDir()
	}
	return filepath.Walk(source, func(pathFile string, info os.FileInfo, e error) error {
		if !info.IsDir() && info.Mode().IsRegular() {
			// Copy file
			cpyPath := destination
			if dstInfo != nil && dstInfo.IsDir() {
				if srcInfo.IsDir() {
					if dstIsDir {
						cpyPath = pathFile[len(path.Dir(source)):]
					} else {
						cpyPath = pathFile[len(source):]
					}
				} else {
					cpyPath = path.Base(pathFile)
				}
				cpyPath = path.Join(destination, cpyPath)
			}
			e = copyFile(pathFile, cpyPath)
			if e != nil {
				err = e
				return err
			}
		} else if info.IsDir() {
			// Copy directory
			dirpath := destination
			if dstInfo != nil && dstInfo.IsDir() {
				if dstIsDir {
					dirpath = path.Join(destination, pathFile[len(path.Dir(source)):])
				} else {
					dirpath = path.Join(destination, pathFile[len(source):])
				}
			}
			e = os.Mkdir(dirpath, 0755)
			if e != nil {
				err = e
				return err
			}
			dstInfo, e = os.Stat(dirpath)
			if e != nil {
				err = e
				return err
			}
		}
		return nil
	})
}

func copyFile(source, destination string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer func() {
		cerr := in.Close()
		if err == nil {
			err = cerr
		}
	}()
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Sync()
}

func CompareFile(file1, file2 string) (int, error) {
	buf1, err := ioutil.ReadFile(file1)
	if err != nil {
		return 0, err
	}
	buf2, err := ioutil.ReadFile(file2)
	if err != nil {
		return 0, err
	}
	return bytes.Compare(buf1, buf2), nil
}
