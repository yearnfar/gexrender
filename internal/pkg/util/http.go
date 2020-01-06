package util

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
)

// Download 下载文件
func Download(urlStr, filename string) (err error) {
	resp, err := http.Get(urlStr)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return ioutil.WriteFile(filename, data, 0666)
}

// DownloadTo 下载到目录
func DownloadTo(urlStr, saveDir string) (filename string, err error) {
	filename = filepath.Join(saveDir, filepath.Base(urlStr))
	return filename, Download(urlStr, filename)
}
