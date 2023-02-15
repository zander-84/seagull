package zap

import (
	"os"
	"strings"
)

func getPrefixPath(dir string, filename string) string {
	return strings.TrimRight(dir, "/") + "/" + filename
}

func openOrCreate(path string, fileName string) (*os.File, error) {
	file, err := os.OpenFile(getPrefixPath(path, fileName), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func openOrCreateWithAction(path string, fileName string, f func(f *os.File)) error {
	file, err := openOrCreate(path, fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	f(file)
	return nil
}
