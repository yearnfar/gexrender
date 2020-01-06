package gexrender

import (
	"github.com/yearnfar/gexrender/internal/gexrender"
	_ "github.com/yearnfar/gexrender/internal/gexrender/action"
)

// Render 渲染
func Render(options ...Option) (err error) {
	setting := &gexrender.Setting{}
	for _, option := range options {
		option(setting)
	}

	return gexrender.Render(setting)
}
