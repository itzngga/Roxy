package helper

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
