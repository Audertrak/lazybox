//go:build !windows
// +build !windows

package fs

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"
)

func getOwnerGroup(fi os.FileInfo) (string, string) {
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return "", ""
	}
	uid := stat.Uid
	gid := stat.Gid
	u, _ := user.LookupId(strconv.FormatUint(uint64(uid), 10))
	g, _ := user.LookupGroupId(strconv.FormatUint(uint64(gid), 10))
	owner := ""
	group := ""
	if u != nil {
		owner = u.Username
	}
	if g != nil {
		group = g.Name
	}
	return owner, group
}

func getCreateTime(fi os.FileInfo) time.Time {
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}
	}
	sec := stat.Ctim.Sec
	nsec := stat.Ctim.Nsec
	return time.Unix(sec, nsec)
}
