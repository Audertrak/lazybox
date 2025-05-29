//go:build windows
// +build windows

package fs

import (
	"os"
	"time"
)

func getOwnerGroup(fi os.FileInfo) (string, string) {
	return "", ""
}

func getCreateTime(fi os.FileInfo) time.Time {
	return time.Time{}
}
