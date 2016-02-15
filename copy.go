package ioutil

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Copy(source, destination string) (err error) {
	source = trimTrailingSlash(source)
	destination = trimTrailingSlash(destination)
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

func ShadowCopy(source, destination string) (err error) {
	source = trimTrailingSlash(source)
	destination = trimTrailingSlash(destination)
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
			e = shadowCopyFile(pathFile, cpyPath)
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

func shadowCopyFile(source, destination string) (err error) {
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
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	return err
}

func copyFile(source, destination string) (err error) {
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
		cerr := out.Close()
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

func IsFileEmpty(filepath string) bool {
	fileinfo, err := os.Stat(filepath)
	if err != nil {
		return false
	}
	return fileinfo.Size() == 0
}

func trimTrailingSlash(path string) string {
	i := strings.LastIndex(path, "/")
	if i == len(path)-1 {
		return path[:i]
	}
	return path
}
