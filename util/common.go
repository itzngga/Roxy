package util

import (
	"github.com/google/uuid"
	"strings"
)

func CreateUid() string {
	id := uuid.New().String()
	return id[:len(id)*1/4]
}

func CreateButtonID(uid, cmd string) string {
	return uid + "-" + cmd
}

func ParseButtonID(uid, id string) string {
	if strings.HasPrefix(id, uid+"-") {
		return id[len(uid)+1:]
	}
	return ""
}

func ParseArgs(str string) []string {
	return strings.Split(str, "")
}

func StringIsOnSlice(target string, slice []string) bool {
	inSlice := false
	for _, i := range slice {
		if target == i {
			inSlice = true
		}
	}
	return inSlice
}

func RemoveElementByIndex[T []any](slice []T, index int) []T {
	sliceLen := len(slice)
	sliceLastIndex := sliceLen - 1
	if index != sliceLastIndex {
		slice[index] = slice[sliceLastIndex]
	}
	return slice[:sliceLastIndex]
}
