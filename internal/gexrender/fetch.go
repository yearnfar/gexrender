package gexrender

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/yearnfar/gexrender/internal/pkg/util"
)

// Fetch 抓取资源
func Fetch(rawUrl, saveDir string) (destFile string, err error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return
	}

	protocol := u.Scheme
	destName := filepath.Base(rawUrl)

	// remove possible query search string params
	if idx := strings.LastIndex(destName, "?"); idx != -1 {
		destName = destName[0:idx]
	}

	// prevent same name file collisions
	if util.IsFile(filepath.Join(saveDir, destName)) {
		destExt, randStr := filepath.Ext(destName), util.RandomString(6)
		destName = fmt.Sprintf("%s-%s%s", strings.TrimRight(destName, destExt), randStr, destExt)
	}

	destFile = filepath.Join(saveDir, destName)

	switch protocol {
	case "data":
		// TODO
		return
	case "http",
		"https":
		log.Infof("http down: %s", rawUrl)
		err = util.Download(rawUrl, destFile)
		return

	case "file":
		log.Infof("copy file: %s", rawUrl)
		var src string
		src, err = util.UriToPath(rawUrl)
		if err != nil {
			return
		}

		err = util.CopyFile(src, destFile)
		if err != nil {
			err = fmt.Errorf("复制文件%s失败：%w", src, err)
			return
		}
		return

	default:
		err = errors.New("unknown protocol")
		return
	}
}
