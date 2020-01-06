package action

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/yearnfar/gexrender/internal/gexrender"
	"github.com/yearnfar/gexrender/internal/pkg/util"
)

// Encode 视频编码
type Encode struct{}

// Invoke 执行
func (e *Encode) Invoke(job *gexrender.Job, data []byte) (err error) {
	ffmpegBin := "ffmpeg"
	if runtime.GOOS == "windows" {
		ffmpegBin += ".exe"
	}

	ffmpegBin, err = exec.LookPath(ffmpegBin)
	if err != nil {
		err = errors.New("ffmpeg no in path")
		return
	}

	param := &encodeParam{}
	err = json.Unmarshal(data, param)
	if err != nil {
		return
	}

	err = param.Valid(job)
	if err != nil {
		return
	}

	// FFmpeg将其所有的日志数据输出到stderr，以留出stdout空间，以便将输出数据输送到其他程序或另一个Ffmpeg实例。
	// https://stackoverflow.com/questions/35169650/differentiate-between-error-and-standard-terminal-log-with-ffmpeg-nodejs
	args := []string{"-loglevel", "error"}

	switch param.Preset {
	case "mp4":
		args = append(args, "-i", param.Input)
		args = append(args, "-acodec", "aac")
		args = append(args, "-ab", "128k")
		args = append(args, "-ar", "44100")
		args = append(args, "-vcodec", "libx264")
		args = append(args, "-r", "25")
		args = append(args, "-pix_fmt", "yuv420p")
		args = append(args, "-y", param.Output)
	case "ogg":
		args = append(args, "-i", param.Input)
		args = append(args, "-acodec", "libvorbis")
		args = append(args, "-ab", "128k")
		args = append(args, "-ar", "44100")
		args = append(args, "-vcodec", "libtheora")
		args = append(args, "-r", "25")
		args = append(args, "-y", param.Output)
	case "webm":
		args = append(args, "-i", param.Input)
		args = append(args, "-acodec", "libvorbis")
		args = append(args, "-ab", "128k")
		args = append(args, "-ar", "44100")
		args = append(args, "-vcodec", "libvpx")
		args = append(args, "-b", "614400")
		args = append(args, "-aspect", "16:9")
		args = append(args, "-y", param.Output)
	case "mp3":
		args = append(args, "-i", param.Input)
		args = append(args, "-acodec", "libmp3lame")
		args = append(args, "-ab", "128k")
		args = append(args, "-ar", "44100")
		args = append(args, "-y", param.Output)
	case "m4a":
		args = append(args, "-i", param.Input)
		args = append(args, "-acodec", "aac")
		args = append(args, "-ab", "64k")
		args = append(args, "-ar", "44100")
		args = append(args, "-strict", "-2")
		args = append(args, "-y", param.Output)
	default:
		args = append(args, "-i", param.Input)
		args = append(args, "-y", param.Output)
	}

	var stderr bytes.Buffer
	cmd := exec.Command(ffmpegBin, args...)
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("ffmpeg转码失败：%w", err)
		return
	}

	if stderr.Len() > 0 {
		err = fmt.Errorf("ffmpeg转码失败, stderr：%s", stderr.String())
		return
	}
	return
}

type encodeParam struct {
	Preset string `json:"preset"`
	Input  string `json:"input"`
	Output string `json:"output"`
}

// Valid 数据校验
func (p *encodeParam) Valid(job *gexrender.Job) (err error) {
	if p.Input == "" {
		err = errors.New("输入视频地址不能为空")
		return
	} else if !filepath.IsAbs(p.Input) {
		p.Input = filepath.Join(job.WorkPath, p.Input)
	}

	if !util.IsFile(p.Input) {
		err = fmt.Errorf("文件%s不存在", p.Input)
		return
	}

	if p.Output == "" {
		err = errors.New("输出视频地址不能为空")
		return
	} else if !filepath.IsAbs(p.Output) {
		p.Output = filepath.Join(job.WorkPath, p.Output)
	}
	return
}

func init() {
	gexrender.Register("@gexrender/action-encode", &Encode{})
}
