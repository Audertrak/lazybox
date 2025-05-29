package fs

import (
	"lazybox/internal/ir"
	"lazybox/internal/util"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func Scan(path string) (*ir.FileInfo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return scanEntry(absPath, absPath)
}

func scanEntry(entryPath, basePath string) (*ir.FileInfo, error) {
	fi, err := os.Lstat(entryPath)
	if err != nil {
		return &ir.FileInfo{
			Name:         filepath.Base(entryPath),
			Path:         relPath(entryPath, basePath),
			AbsolutePath: entryPath,
			Error:        err.Error(),
		}, nil
	}

	fileType := getFileType(fi)
	symlinkTarget := ""
	if fileType == ir.FileTypeSymlink {
		target, err := os.Readlink(entryPath)
		if err == nil {
			symlinkTarget = target
		}
	}

	owner, group := getOwnerGroup(fi)
	var createTime *time.Time
	if ctime := getCreateTime(fi); !ctime.IsZero() {
		createTime = &ctime
	}

	info := &ir.FileInfo{
		Name:          filepath.Base(entryPath),
		Path:          relPath(entryPath, basePath),
		AbsolutePath:  entryPath,
		Type:          fileType,
		Size:          fi.Size(),
		Mode:          fi.Mode().String(),
		Owner:         owner,
		Group:         group,
		ModTime:       fi.ModTime(),
		CreateTime:    createTime,
		SymlinkTarget: symlinkTarget,
	}

	if fileType == ir.FileTypeDirectory {
		info.IsGitRepo = util.IsGitRepo(entryPath)
		if info.IsGitRepo {
			info.GitRemotes = util.GetGitRemotes(entryPath)
		}
		entries, err := os.ReadDir(entryPath)
		if err != nil {
			info.Error = err.Error()
			return info, nil
		}
		for _, entry := range entries {
			// Optionally skip hidden files except .git
			if strings.HasPrefix(entry.Name(), ".") && entry.Name() != ".git" {
				continue
			}
			childPath := filepath.Join(entryPath, entry.Name())
			childInfo, _ := scanEntry(childPath, basePath)
			info.Contents = append(info.Contents, childInfo)
		}
	} else if fileType == ir.FileTypeFile {
		info.Extension = filepath.Ext(entryPath)
	}

	return info, nil
}

func relPath(path, base string) string {
	rel, err := filepath.Rel(base, path)
	if err != nil {
		return path
	}
	return rel
}

func getFileType(fi os.FileInfo) ir.FileType {
	mode := fi.Mode()
	switch {
	case mode.IsRegular():
		return ir.FileTypeFile
	case mode.IsDir():
		return ir.FileTypeDirectory
	case mode&os.ModeSymlink != 0:
		return ir.FileTypeSymlink
	default:
		return ir.FileTypeOther
	}
}

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
