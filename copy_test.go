package ioutil

import (
	frand "crypto/rand"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
)

func testCreateRandomFile(dirname string) (string, error) {
	buf := make([]byte, rand.Intn(4096))
	_, err := frand.Read(buf)
	if err != nil {
		return "", err
	}
	if dirname == "" {
		dirname = os.TempDir()
	}
	tmpFile, err := ioutil.TempFile(dirname, "testCreateRandomFile")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()
	_, err = tmpFile.Write(buf)
	if err != nil {
		return "", err
	}
	err = tmpFile.Sync()
	if err != nil {
		return "", err
	}
	return tmpFile.Name(), nil
}

func TestCompareFileSameFile(t *testing.T) {
	file1, err := testCreateRandomFile("")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file1)
	file2, err := testCreateRandomFile("")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file2)
	err = copyFile(file1, file2)
	if err != nil {
		t.Fatal(err)
	}
	n, err := CompareFile(file1, file2)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("'%s' and '%s' should have the same content.\n", file1, file2)
	}
}

func TestCompareFileDifferentFile(t *testing.T) {
	file1, err := testCreateRandomFile("")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file1)
	file2, err := testCreateRandomFile("")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file2)
	n, err := CompareFile(file1, file2)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatalf("'%s' and '%s' should have the same content.\n", file1, file2)
	}
}

func TestCopyFileToFile(t *testing.T) {
	file1, err := testCreateRandomFile("")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file1)
	file2 := path.Join(os.TempDir(), "TestCopyFileToFile"+strconv.Itoa(rand.Int()))
	defer os.Remove(file2)
	err = Copy(file1, file2)
	if err != nil {
		t.Fatal(err)
	}
	n, err := CompareFile(file1, file2)
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("'%s' and '%s' should have the same content.\n", file1, file2)
	}
}

func TestCopyInexistingFile(t *testing.T) {
	file1 := path.Join(os.TempDir(), "TestCopyInexistingFile"+strconv.Itoa(rand.Int()))
	file2 := path.Join(os.TempDir(), "TestCopyInexistingFile"+strconv.Itoa(rand.Int()))
	err := Copy(file1, file2)
	if err == nil {
		t.Fatal("You shouldn't be able to copy a file that doesn't exist !")
	}
}

func TestCopyFileToDirectory(t *testing.T) {
	file, err := testCreateRandomFile("")
	defer os.Remove(file)
	dir, err := ioutil.TempDir(os.TempDir(), "TestCopyFileToDirectory")
	defer os.RemoveAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	err = Copy(file, dir)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(path.Join(dir, path.Base(file)))
	if err != nil {
		t.Fatal(err)
	}
	n, err := CompareFile(file, path.Join(dir, path.Base(file)))
	if err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("'%s' and '%s' should have the same content.\n", file, path.Join(dir, path.Base(file)))
	}
}

func TestCopyDirectoryToFile(t *testing.T) {
	file, err := testCreateRandomFile("")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file)
	dir, err := ioutil.TempDir(os.TempDir(), "TestCopyDirectoryToFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	err = Copy(dir, file)
	if err == nil {
		t.Fatal("You shouldn't be able to copy a directory into a file !")
	}
}

func TestCopyDirectoryToNonExistingDirectory(t *testing.T) {
	maindir, err := ioutil.TempDir(os.TempDir(), "source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(maindir)
	dir, err := ioutil.TempDir(maindir, "dir")
	if err != nil {
		t.Fatal(err)
	}
	subdir, err := ioutil.TempDir(dir, "subdir")
	if err != nil {
		t.Fatal(err)
	}
	files := make([]string, 25)
	for i := 0; i < len(files); i++ {
		n := rand.Int()
		if n%3 == 0 {
			files[i], err = testCreateRandomFile(maindir)
		} else if n%3 == 1 {
			files[i], err = testCreateRandomFile(dir)
		} else {
			files[i], err = testCreateRandomFile(subdir)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	cpyDir := path.Join(os.TempDir(), "destination"+strconv.Itoa(rand.Int()))
	err = Copy(maindir, cpyDir)
	defer os.RemoveAll(cpyDir)
	if err != nil {
		t.Fatal(err)
	}
	filepath.Walk(maindir, func(pathFile string, info os.FileInfo, err error) error {
		cpyPath := pathFile[len(maindir):]
		cpyPath = path.Join(cpyDir, cpyPath)
		// Check if file or directory has been created
		_, err = os.Stat(cpyPath)
		if err != nil {
			t.Fatal(err)
		}
		if !info.IsDir() {
			n, err := CompareFile(pathFile, cpyPath)
			if err != nil {
				t.Fatal(err)
			}
			if n != 0 {
				t.Fatalf("'%s' and '%s' should have the same content.\n", pathFile, cpyPath)
			}
		}
		return nil
	})
}

func TestCopyDirectoryToExistingDirectory(t *testing.T) {
	maindir, err := ioutil.TempDir(os.TempDir(), "source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(maindir)
	dir, err := ioutil.TempDir(maindir, "dir")
	if err != nil {
		t.Fatal(err)
	}
	subdir, err := ioutil.TempDir(dir, "subdir")
	if err != nil {
		t.Fatal(err)
	}
	files := make([]string, 25)
	for i := 0; i < len(files); i++ {
		n := rand.Int()
		if n%3 == 0 {
			files[i], err = testCreateRandomFile(maindir)
		} else if n%3 == 1 {
			files[i], err = testCreateRandomFile(dir)
		} else {
			files[i], err = testCreateRandomFile(subdir)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	cpyDir, err := ioutil.TempDir(os.TempDir(), "destination")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(cpyDir)
	err = Copy(maindir, cpyDir)
	if err != nil {
		t.Fatal(err)
	}
	filepath.Walk(maindir, func(pathFile string, info os.FileInfo, err error) error {
		cpyPath := pathFile[len(path.Dir(maindir)):]
		cpyPath = path.Join(cpyDir, cpyPath)
		// Check if file or directory has been created
		_, err = os.Stat(cpyPath)
		if err != nil {
			t.Fatal(err)
		}
		if !info.IsDir() {
			n, err := CompareFile(pathFile, cpyPath)
			if err != nil {
				t.Fatal(err)
			}
			if n != 0 {
				t.Fatalf("'%s' and '%s' should have the same content.\n", pathFile, cpyPath)
			}
		}
		return nil
	})
}

func TestShadowCopyFileToFile(t *testing.T) {
	file1, err := testCreateRandomFile("")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file1)
	file2 := path.Join(os.TempDir(), "TestShadowCopyFileToFile"+strconv.Itoa(rand.Int()))
	defer os.Remove(file2)
	err = ShadowCopy(file1, file2)
	if err != nil {
		t.Fatal(err)
	}
	if !IsFileEmpty(file2) {
		t.Fatalf("`%s` should be empty", file2)
	}
}

func TestShadowCopyInexistingFile(t *testing.T) {
	file1 := path.Join(os.TempDir(), "TestShadowCopyInexistingFile"+strconv.Itoa(rand.Int()))
	file2 := path.Join(os.TempDir(), "TestShadowCopyInexistingFile"+strconv.Itoa(rand.Int()))
	err := ShadowCopy(file1, file2)
	if err == nil {
		t.Fatal("You shouldn't be able to copy a file that doesn't exist !")
	}
}

func TestShadowCopyFileToDirectory(t *testing.T) {
	file, err := testCreateRandomFile("")
	defer os.Remove(file)
	dir, err := ioutil.TempDir(os.TempDir(), "TestShadowCopyFileToDirectory")
	defer os.RemoveAll(dir)
	if err != nil {
		t.Fatal(err)
	}
	err = ShadowCopy(file, dir)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(path.Join(dir, path.Base(file)))
	if err != nil {
		t.Fatal(err)
	}
	if !IsFileEmpty(path.Join(dir, path.Base(file))) {
		t.Fatalf("`%s` should be empty", path.Join(dir, path.Base(file)))
	}
}

func TestShadowCopyDirectoryToFile(t *testing.T) {
	file, err := testCreateRandomFile("")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file)
	dir, err := ioutil.TempDir(os.TempDir(), "TestShadowCopyDirectoryToFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	err = ShadowCopy(dir, file)
	if err == nil {
		t.Fatal("You shouldn't be able to copy a directory into a file !")
	}
}

func TestShadowCopyDirectoryToNonExistingDirectory(t *testing.T) {
	maindir, err := ioutil.TempDir(os.TempDir(), "source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(maindir)
	dir, err := ioutil.TempDir(maindir, "dir")
	if err != nil {
		t.Fatal(err)
	}
	subdir, err := ioutil.TempDir(dir, "subdir")
	if err != nil {
		t.Fatal(err)
	}
	files := make([]string, 25)
	for i := 0; i < len(files); i++ {
		n := rand.Int()
		if n%3 == 0 {
			files[i], err = testCreateRandomFile(maindir)
		} else if n%3 == 1 {
			files[i], err = testCreateRandomFile(dir)
		} else {
			files[i], err = testCreateRandomFile(subdir)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	cpyDir := path.Join(os.TempDir(), "destination"+strconv.Itoa(rand.Int()))
	err = ShadowCopy(maindir, cpyDir)
	defer os.RemoveAll(cpyDir)
	if err != nil {
		t.Fatal(err)
	}
	filepath.Walk(maindir, func(pathFile string, info os.FileInfo, err error) error {
		cpyPath := pathFile[len(maindir):]
		cpyPath = path.Join(cpyDir, cpyPath)
		// Check if file or directory has been created
		_, err = os.Stat(cpyPath)
		if err != nil {
			t.Fatal(err)
		}
		if !info.IsDir() {
			if !IsFileEmpty(cpyPath) {
				t.Fatalf("`%s` should be empty", cpyPath)
			}
		}
		return nil
	})
}

func TestShadowCopyDirectoryToExistingDirectory(t *testing.T) {
	maindir, err := ioutil.TempDir(os.TempDir(), "source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(maindir)
	dir, err := ioutil.TempDir(maindir, "dir")
	if err != nil {
		t.Fatal(err)
	}
	subdir, err := ioutil.TempDir(dir, "subdir")
	if err != nil {
		t.Fatal(err)
	}
	files := make([]string, 25)
	for i := 0; i < len(files); i++ {
		n := rand.Int()
		if n%3 == 0 {
			files[i], err = testCreateRandomFile(maindir)
		} else if n%3 == 1 {
			files[i], err = testCreateRandomFile(dir)
		} else {
			files[i], err = testCreateRandomFile(subdir)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
	cpyDir, err := ioutil.TempDir(os.TempDir(), "destination")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(cpyDir)
	err = ShadowCopy(maindir, cpyDir)
	if err != nil {
		t.Fatal(err)
	}
	filepath.Walk(maindir, func(pathFile string, info os.FileInfo, err error) error {
		cpyPath := pathFile[len(path.Dir(maindir)):]
		cpyPath = path.Join(cpyDir, cpyPath)
		// Check if file or directory has been created
		_, err = os.Stat(cpyPath)
		if err != nil {
			t.Fatal(err)
		}
		if !info.IsDir() {
			if !IsFileEmpty(cpyPath) {
				t.Fatalf("`%s` should be empty", cpyPath)
			}
		}
		return nil
	})
}

func TestIsFileEmptyWithEmptyFile(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "TestIsFileEmptyWithEmptyFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if !IsFileEmpty(file.Name()) {
		t.Fatalf("`%s` should be an empty file", file.Name())
	}
}

func TestIsFileEmptyWithNonEmptyFile(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "TestIsFileEmptyWithNonEmptyFile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())
	file.Write([]byte("Hello World!"))
	file.Sync()
	file.Close()
	if IsFileEmpty(file.Name()) {
		t.Fatalf("`%s` should not be an empty file", file.Name())
	}
}
