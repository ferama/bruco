package common

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func FindPython() (string, error) {
	var path string
	var err error
	path, err = exec.LookPath("python3")
	if err != nil {
		path, err = exec.LookPath("python")
		if err != nil {
			return "", err
		}
	}
	return path, nil
}

func GetenvStr(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return v, fmt.Errorf("%s: env var is empty", key)
	}
	return v, nil
}

func GetenvBool(key string) (bool, error) {
	s, err := GetenvStr(key)
	if err != nil {
		return false, err
	}
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}
	return v, nil
}
