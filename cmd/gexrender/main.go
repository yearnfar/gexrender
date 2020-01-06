package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/yearnfar/gexrender/internal/gexrender"
	_ "github.com/yearnfar/gexrender/internal/gexrender/action"
)

func main() {
	app := &cli.App{
		Name:    "gexrender-cli",
		Usage:   "gexrender standalone renderer",
		Version: "1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "f",
				Aliases: []string{"file"},
				Usage:   "提供一个包含作业的json文件的相对或绝对路径代替从参数中使用json。",
			},
			&cli.StringFlag{
				Name:    "b",
				Aliases: []string{"binary"},
				Usage:   "手动指定aerender二进制文件的路径，您可以将其留空以依赖于自动查找。",
			},
			&cli.StringFlag{
				Name:    "w",
				Aliases: []string{"workpath"},
				Usage:   "手动覆盖gexrender默认工作目录tmpdir/gexrender。",
			},
			&cli.BoolFlag{
				Name:  "stop-on-error",
				Usage: "如果发生处理/呈现错误，则强制aerender停止，否则aerender将报告错误并继续工作。",
			},
			&cli.BoolFlag{
				Name:  "no-license",
				Usage: "禁止创建ae_render_only_node.txt文件(默认启用)，该文件允许免费使用试用版的Adobe After Effects。",
			},
			&cli.BoolFlag{
				Name:  "force-patch",
				Usage: "强制（重新）安装commandLineRenderer.jsx补丁。",
			},
			&cli.BoolFlag{
				Name:  "debug",
				Usage: "启用aerender的转储命令和其他调试内容。",
			},
			&cli.BoolFlag{
				Name:  "skip-cleanup",
				Usage: "强制在呈现完成后保留临时数据。",
			},
			&cli.BoolFlag{
				Name:  "multi-frames",
				Usage: "（来自Adobe官网）： 根据系统配置和首选项设置，可以创建更多进程来同时渲染多个帧。(参见内存和多处理首选项)。",
			},
			&cli.BoolFlag{
				Name:  "max-memory-percent",
				Usage: `（来自Adobe官网）： 指定After Effects可以使用的内存的总百分比。对于这两个值，如果安装的RAM小于给定的数量(n gb)，则该值是安装的RAM的百分比，否则为n的百分比。n的值为32位Windows的2gb、64位Windows的4gb和Mac OS的3.5 GB。`,
			},
			&cli.BoolFlag{
				Name:  "image-cache-percent",
				Usage: "（来自Adobe官网）：指定用于缓存已渲染图片和镜头的最大内存百分比。",
			},
			&cli.BoolFlag{
				Name:  "reuse",
				Usage: `（来自Adobe官网）：重用当前正在运行的After Effects实例(如果找到的话)来执行gexrender。当使用一个已经运行的实例时，当渲染完成时，aerender将首选项保存到磁盘，但在效果完成后不会退出。如果没有使用这个参数，即使一个After Effects已经在运行，aerender也会启动一个新的实例. 当呈现完成时，它退出该实例，并且不保存首选项。`,
			},
		},
		Action: func(c *cli.Context) error {
			setting := &gexrender.Setting{
				ConfigFile:            c.String("f"),
				WorkPath:              c.String("w"),
				StopOnError:           c.Bool("stop-on-error"),
				NoLicense:             c.Bool("no-license"),
				ForceCommandLinePatch: c.Bool("force-patch"),
				Debug:                 c.Bool("debug"),
				SkipCleanup:           c.Bool("skip-cleanup"),
				MultiFrames:           c.Bool("multi-frames"),
				Reuse:                 c.Bool("reuse"),
				MaxMemoryPercent:      c.Int("max-memory-percent"),
				ImageCachePercent:     c.Int("image-cache-percent"),
			}

			if c.String("f") == "" {
				if c.NArg() == 0 {
					log.Fatal("you need to provide a gexrender job json as an argument")
				} else {
					setting.Config = c.Args().First()
				}
			}

			err := gexrender.Render(setting)
			if err != nil {
				log.Error(err)
				return err
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
