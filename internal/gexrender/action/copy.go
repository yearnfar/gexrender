package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yearnfar/gexrender/internal/gexrender"
	"github.com/yearnfar/gexrender/internal/pkg/util"
)

// CopyFile 拷贝文件
type Copy struct{}

// Invoke 执行
func (c *Copy) Invoke(job *gexrender.Job, data []byte) (err error) {
	param := &copyParam{}
	err = json.Unmarshal(data, param)
	if err != nil {
		return
	}

	err = param.Valid(job)
	if err != nil {
		return
	}

	err = util.CopyFile(param.Input, param.Output)
	if err != nil {
		err = fmt.Errorf("copy file err: %w, src: %s, dest: %s", err, param.Input, param.Output)
		return
	}
	return
}

type copyParam struct {
	Input  string `json:"input"`
	Output string `json:"output"`
}

// Valid 校验
func (p *copyParam) Valid(job *gexrender.Job) (err error) {
	if p.Input == "" {
		err = errors.New("input can't be empty")
		return
	} else if !filepath.IsAbs(p.Input) {
		p.Input = filepath.Join(job.WorkPath, p.Input)
	}

	if !util.IsFile(p.Input) {
		err = fmt.Errorf("文件%s不存在", p.Input)
		return
	}

	if p.Output == "" {
		err = errors.New("output can't be empty")
		return
	} else if !filepath.IsAbs(p.Output) {
		p.Output = filepath.Join(job.WorkPath, p.Output)
	}

	if strings.HasSuffix(p.Output, "/") && !util.IsDir(p.Output) {
		err = os.MkdirAll(p.Output, 0755)
		if err != nil {
			err = fmt.Errorf("创建目录失败：%w", err)
			return
		}
	}
	return
}

func init() {
	gexrender.Register("@gexrender/action-copy", &Copy{})
}
