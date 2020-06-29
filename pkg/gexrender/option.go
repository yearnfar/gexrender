package gexrender

import (
	"github.com/yearnfar/gexrender/internal/gexrender"
)

// Option 设置
type Option func(s *gexrender.Setting)

// WithDebug 调试模式
func WithDebug() Option {
	return func(s *gexrender.Setting) {
		s.Debug = true
	}
}

// WithLogFile 日志输出到文件
func WithLogFile(file string) Option {
	return func(s *gexrender.Setting) {
		s.LogFile = file
	}
}

// WithConfigFile 提供一个包含作业的json文件的相对或绝对路径代替从参数中使用json。
func WithConfigFile(filename string) Option {
	return func(s *gexrender.Setting) {
		s.ConfigFile = filename
	}
}

// WithConfig 参数中使用json
func WithConfig(data string) Option {
	return func(s *gexrender.Setting) {
		s.Config = data
	}
}

// WithBinary 手动指定aerender二进制文件的路径，您可以将其留空以依赖于自动查找。
func WithBinary(bin string) Option {
	return func(s *gexrender.Setting) {
		s.Binary = bin
	}
}

// WithWorkPath 手动覆盖gexrender默认工作目录tmpdir/gexrender。
func WithWorkPath(workPath string) Option {
	return func(s *gexrender.Setting) {
		s.WorkPath = workPath
	}
}

// WithNoLicense 禁止创建ae_render_only_node.txt文件(默认启用)，该文件允许免费使用试用版的Adobe After Effects。
func WithNoLicense() Option {
	return func(s *gexrender.Setting) {
		s.NoLicense = true
	}
}

// WithSkipCleanup 强制在呈现完成后保留临时数据。
func WithSkipCleanup() Option {
	return func(s *gexrender.Setting) {
		s.SkipCleanup = true
	}
}

// WithForceCommandLinePatch 强制（重新）安装commandLineRenderer.jsx补丁。
func WithForceCommandLinePatch() Option {
	return func(s *gexrender.Setting) {
		s.ForceCommandLinePatch = true
	}
}

// WithMultiFrames （来自Adobe官网）： 根据系统配置和首选项设置，可以创建更多进程来同时渲染多个帧。(参见内存和多处理首选项)。
func WithMultiFrames() Option {
	return func(s *gexrender.Setting) {
		s.MultiFrames = true
	}
}

// WithReuse （来自Adobe官网）：重用当前正在运行的After Effects实例(如果找到的话)来执行gexrender。
// 当使用一个已经运行的实例时，当渲染完成时，aerender将首选项保存到磁盘，但在效果完成后不会退出。
// 如果没有使用这个参数，即使一个After Effects已经在运行，aerender也会启动一个新的实例. 当呈现完成时，它退出该实例，并且不保存首选项。
func WithReuse() Option {
	return func(s *gexrender.Setting) {
		s.Reuse = true
	}
}

// WithReuse 如果发生处理/呈现错误，则强制aerender停止，否则aerender将报告错误并继续工作。
func WithStopOnError() Option {
	return func(s *gexrender.Setting) {
		s.StopOnError = true
	}
}

// WithMaxMemoryPercent （来自Adobe官网）： 指定After Effects可以使用的内存的总百分比。
// 对于这两个值，如果安装的RAM小于给定的数量(n gb)，则该值是安装的RAM的百分比，否则为n的百分比。
// n的值为32位Windows的2gb、64位Windows的4gb和Mac OS的3.5 GB。
func WithMaxMemoryPercent(percent int) Option {
	return func(s *gexrender.Setting) {
		s.MaxMemoryPercent = percent
	}
}

// WithReuWithImageCachePercent （来自Adobe官网）：指定用于缓存已渲染图片和镜头的最大内存百分比。
func WithImageCachePercent(percent int) Option {
	return func(s *gexrender.Setting) {
		s.ImageCachePercent = percent
	}
}
