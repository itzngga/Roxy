package util

import (
	"fmt"
	"github.com/google/uuid"
	waTypes "go.mau.fi/whatsmeow/types"
	"net/url"
	"strconv"
	"strings"
)

func RemoveDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func CreateUid() string {
	id := uuid.New().String()
	return id[:len(id)*1/4]
}

func ParseArgs(str string) []string {
	return strings.Split(str, "")
}

func IsValidUrl(s string) bool {
	if len(s) == 0 {
		return false
	}

	url, err := url.Parse(s)
	if err != nil || url.Scheme == "" {
		return false
	}

	if url.Host == "" && url.Fragment == "" && url.Opaque == "" {
		return false
	}

	return true
}

func StringIsOnSlice(target string, slice []string) bool {
	inSlice := false
	for _, i := range slice {
		if target == i {
			inSlice = true
			break
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

func ParseAllJid(jid any) (pJid waTypes.JID, err error) {
	switch uJid := jid.(type) {
	case string:
		result, ok := ParseJID(uJid)
		if !ok {
			return pJid, fmt.Errorf("error: failed to parse jid : %s\n", jid)
		}
		pJid = result
	case waTypes.JID:
		pJid = uJid
	default:
		return pJid, fmt.Errorf("error: unsupported jid types : %s\n", jid)
	}
	return pJid.ToNonAD(), nil
}

func ParseGroupJid(jid any) (pJid waTypes.JID, err error) {
	switch uJid := jid.(type) {
	case string:
		result, ok := ParseJID(uJid)
		if !ok {
			return pJid, fmt.Errorf("error: failed to parse jid : %s\n", jid)
		} else if result.Server != waTypes.GroupServer {
			return pJid, fmt.Errorf("error: given jid is not group jid : %s\n", jid)
		}
		pJid = result
	case waTypes.JID:
		if uJid.Server != waTypes.GroupServer {
			return pJid, fmt.Errorf("error: given jid is not group jid : %s\n", jid)
		}
	default:
		return pJid, fmt.Errorf("error: unsupported jid types : %s\n", jid)
	}
	return pJid.ToNonAD(), nil
}

func ParseUserJid(jid any) (pJid waTypes.JID, err error) {
	switch uJid := jid.(type) {
	case string:
		result, ok := ParseJID(uJid)
		if !ok {
			return pJid, fmt.Errorf("error: failed to parse jid : %s\n", jid)
		} else if result.Server != waTypes.DefaultUserServer {
			return pJid, fmt.Errorf("error: given jid is not user jid : %s\n", jid)
		}
		pJid = result
	case waTypes.JID:
		if uJid.Server != waTypes.DefaultUserServer {
			return pJid, fmt.Errorf("error: given jid is not user jid : %s\n", jid)
		}
	default:
		return pJid, fmt.Errorf("error: unsupported jid types : %s\n", jid)
	}
	return pJid.ToNonAD(), nil
}

func JIDToString(jidStr any) (result string, err error) {
	switch uJid := jidStr.(type) {
	case string:
		jid, ok := ParseJID(uJid)
		if !ok {
			return "", fmt.Errorf("error: failed to parse jid: %s", jid)
		} else if jid.Server != waTypes.DefaultUserServer && jid.Server != waTypes.GroupServer {
			return "", fmt.Errorf("error: given jid is not user or group jid: %s", jidStr)
		}
		return strconv.FormatUint(jid.ToNonAD().UserInt(), 10), nil
	case waTypes.JID:
		if uJid.Server != waTypes.DefaultUserServer && uJid.Server != waTypes.GroupServer {
			return "", fmt.Errorf("error: given jid is not user or group jid: %s", jidStr)
		}
		return strconv.FormatUint(uJid.ToNonAD().UserInt(), 10), nil
	default:
		return "", fmt.Errorf("error: unsupported jid types: %s", jidStr)
	}
}
