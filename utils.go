package main

import (
	"os"
	"strings"
	"syscall"
)

func indexOf[T comparable](DesiredValue T, array []T) int {
	for index, value := range array {
		if value == DesiredValue {
			return index
		}
	}
	return -1
}

/*
	Аргумент Flag это искомые флаг.
	Аргумент Value это значение если flag не был найден
*/
func getValueFlag(Flag, value string) string {
	IndexFlag := indexOf(Flag, os.Args)
	if IndexFlag == -1 {
		return value
	}
	return os.Args[IndexFlag+1]
}

func isFlag(arg string) bool {
	if strings.Contains(arg, "--") || strings.Contains(arg, "-") {
		return true
	}
	return false
}

func isWritableDir(path string) (isWritable bool, err error) {
	var stat syscall.Stat_t
	info, err := os.Stat(path)
	if err != nil {
		return false, ERROR_PATH_NOT_EXISTS
	}
	if !info.IsDir() {
		return false, ERROR_NOT_DIRECTORY
	}

	// Check if the user bit is enabled in file permission
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return false, ERROR_WRITE_BIT_NOT_SET
	}

	if err = syscall.Stat(path, &stat); err != nil {
		return false, ERROR_STAT
	}

	if uint32(os.Geteuid()) != stat.Uid {
		return false, ERROR_NO_PERMISSION
	}

	return true, nil
}
