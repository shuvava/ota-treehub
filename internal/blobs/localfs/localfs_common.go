package localfs

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/shuvava/treehub/internal/utils/fshelper"
)

func copyContentAndClose(file *os.File, reader io.Reader) (int64, error) {
	defer func() { _ = file.Close() }()
	written, err := io.Copy(file, reader)
	if err != nil {
		return 0, err
	}
	return written, nil
}

// safeStoreStream store stream into tmp file and replace path by this tmp file
func safeStoreStream(path string, reader io.Reader) (size int64, err error) {
	parent := filepath.Dir(path)
	fname := filepath.Base(path)
	file, err := ioutil.TempFile(parent, fname)
	if err != nil {
		return 0, err
	}
	tmppath := file.Name()

	defer func() {
		if err != nil {
			_ = os.Remove(tmppath)
		}
	}()
	written, err := copyContentAndClose(file, reader)
	if fshelper.IsPathExist(path) {
		err = os.Remove(path)
		if err != nil {
			return 0, err
		}
	}
	err = os.Rename(tmppath, path)
	if err != nil {
		return 0, err
	}
	return written, nil
}
