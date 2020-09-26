package filesystem

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetRootDirectory() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("unable to get runtime caller to find out a root directory")
	}

	dir := filepath.Dir(file)

	parentDirectories := 3
	sep := string(os.PathSeparator)

	splits := strings.Split(dir, sep)
	root := strings.Join(splits[0:len(splits)-parentDirectories], sep)

	return root + sep, nil
}
