package gexrender

import (
	"runtime"

	"github.com/yearnfar/gexrender/internal/pkg/sliceutil"
)

// Template 素材模板
type Template struct {
	Src         string `json:"src" validate:"required"`         // 文件
	Composition string `json:"composition" validate:"required"` // 合成
	Dest        string `json:"dest"`                            // 输出

	FrameStart     string `json:"frame_start"`     // -s 开始帧
	FrameEnd       string `json:"frame_end"`       // -e 结束
	IncrementFrame string `json:"increment_frame"` // -i

	ContinueOnMissing bool   `json:"continue_on_missing"` //
	SettingsTemplate  string `json:"settings_template"`   // - RStemplate
	OutputModule      string `json:"output_module"`       // - OMtemplate
	OutputExt         string `json:"output_ext"`          // 输出视频后缀
	ImageSequence     bool   `json:"image_sequence"`      // 图片序列帧
}

// Fetch 抓取资源
func (t *Template) Fetch(saveDir string) error {
	dest, err := Fetch(t.Src, saveDir)
	if err != nil {
		return err
	}

	t.Dest = dest
	return nil
}

// IsRenderImageSequence 是否渲染序列帧
func (t *Template) IsRenderImageSequence() bool {
	return sliceutil.InStrings(t.OutputExt, []string{"jpeg", "jpg", "png"})
}

// GetOutputExt 渲染视频后缀
func (t *Template) GetOutputExt() string {
	if t.OutputExt != "" {
		return t.OutputExt
	}

	if runtime.GOOS == "darwin" {
		return "mov"
	} else {
		return "avi"
	}
}
