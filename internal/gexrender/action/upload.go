package action

import (
	"github.com/yearnfar/gexrender/internal/gexrender"
)

// Upload 上传
type Upload struct{}

// Invoke 执行
func (u *Upload) Invoke(job *gexrender.Job, data []byte) (err error) {
	return
}

func init() {
	gexrender.Register("@gexrender/action-upload", &Upload{})
}
