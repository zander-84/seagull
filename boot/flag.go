package boot

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func ParseCmd() (string, string, string, error) {
	var action string
	flag.StringVar(&action, "a", "", "行为")
	action = strings.TrimSpace(action)
	var in string
	flag.StringVar(&in, "i", "", "入参")
	in = strings.TrimSpace(in)

	var out string
	flag.StringVar(&out, "o", "", "保存位置")
	out = strings.TrimSpace(out)

	flag.Parse()

	if action == "" || in == "" {
		return "", "", "", errors.New("参数 -a 或者 -i 不能为空")
	}
	return action, in, out, nil
}

func readFile(fileName string) (string, error) {
	if err := isFile(fileName); err != nil {
		return "", err
	}

	d1, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(d1), nil
}

func isDir(fileName string) bool {
	fi, e := os.Stat(fileName)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

func isFile(fileName string) error {
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File %s does not exist\n", fileName)
		}
		return err
	}
	return err
}

func fileExist(fileName string) (bool, error) {
	if _, err := os.Stat(fileName); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			fmt.Println("Some other error:", err)
			return false, err
		}
	} else {
		//fmt.Println("File or directory exists")
		return true, nil
	}

}
