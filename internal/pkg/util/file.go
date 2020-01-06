package util

import (
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// IsFile 文件是否存在
func IsFile(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

// IsDir 检查是否目录
func IsDir(dir string) bool {
	info, err := os.Stat(dir)
	return err == nil && info.IsDir()
}

// CopyFile 复制文件
func CopyFile(src, dest string) error {
	if !IsFile(src) {
		return errors.New("file not exists")
	}

	if IsDir(dest) {
		dest = filepath.Join(dest, filepath.Base(src))
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer srcFile.Close()
	dstFile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}

//  UriToPath File URI to Path function.
func UriToPath(uri string) (file string, err error) {
	if !strings.HasPrefix(uri, "file://") {
		return "", errors.New("must pass in a file:// URI to convert to a file path")
	}

	rest, _ := url.QueryUnescape(uri[7:])
	firstSlash := strings.Index(rest, "/")
	host, path := rest[0:firstSlash], rest[firstSlash+1:]

	// 2.  Scheme Definition
	// As a special case, <host> can be the string "localhost" or the empty
	// string; this is interpreted as "the machine from which the URL is
	// being interpreted".
	if host == "localhost" {
		host = ""
	}

	if host != "" {
		host = string(filepath.Separator) + string(filepath.Separator) + host
	}

	// 3.2  Drives, drive letters, mount points, file system root
	// Drive letters are mapped into the top of a file URI in various ways,
	// depending on the implementation; some applications substitute
	// vertical bar ("|") for the colon after the drive letter, yielding
	// "file:///c|/tmp/test.txt".  In some cases, the colon is left
	// unchanged, as in "file:///c:/tmp/test.txt".  In other cases, the
	// colon is simply omitted, as in "file:///c/tmp/test.txt".
	path = regexp.MustCompile(`^(.+)\|`).ReplaceAllString(path, `$1:`)

	// for Windows, we need to invert the path separators from what a URI uses
	if filepath.Separator == '\\' {
		path = strings.ReplaceAll(path, "/", "\\")
	}

	if !regexp.MustCompile(`^.+:`).MatchString(path) {
		path = string(filepath.Separator) + path
	}

	return host + path, nil
}
