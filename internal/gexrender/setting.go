package gexrender

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"github.com/yearnfar/gexrender/internal/gexrender/assets"
	"github.com/yearnfar/gexrender/internal/pkg/util"
)

// Setting 参数
type Setting struct {
	Config                string // JSON配置
	ConfigFile            string // JSON配置文件
	Binary                string // AeRender文件位置
	WorkPath              string // 运行目录
	NoLicense             bool   // 不需要license
	SkipCleanup           bool   // 不清理缓存目录
	ForceCommandLinePatch bool   // 补丁
	Debug                 bool   // 调试模式
	LogFile               string // 日志输出文件
	MultiFrames           bool
	Reuse                 bool
	StopOnError           bool
	MaxMemoryPercent      int
	ImageCachePercent     int
}

// Setup 配置
func (s *Setting) Setup() (job *Job, err error) {
	s.autoBinary()

	job, err = CreateJob(s)
	if err != nil {
		return
	}

	err = validate.Struct(job)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		for _, err_ := range errs {
			log.Error(err_.Translate(trans))
		}
		return
	}

	if job.Template.IsRenderImageSequence() {
		job.ResultName = "result_[#####]." + job.Template.OutputExt
		job.Template.ImageSequence = true
	} else {
		job.ResultName = "result." + job.Template.GetOutputExt()
	}

	if (job.Actions == nil || job.Actions["postrender"] == nil) && !s.SkipCleanup {
		log.Infof(`[%s]-- W A R N I N G: --

You haven't provided any post-render actions!
After render is finished all the files inside temporary folder (INCLUDING your target video) will be removed.

To prevent this from happening, please add an action to "job.actions.postrender".
For more info checkout: https://github.com/inlife/gexrender#Actions

P.S. to prevent gexrender from removing temp file data, you also can please provide an argument:
--skip-cleanup (or skipCleanup: true if using programmatically)\n`, job.Uid)
	}

	job.WorkPath = path.Join(s.WorkPath, job.Uid)
	job.Output = path.Join(job.WorkPath, job.ResultName)

	// add license helper
	if !s.NoLicense {
		s.addLicense()
	}

	// attempt to patch the default
	// Scripts/commandLineRenderer.js
	err = s.patch()
	if err != nil {
		log.Fatal(err)
		return
	}
	return
}

// 增加license
func (s *Setting) addLicense() {
	filename := "ae_render_only_node.txt"
	home, _ := os.UserHomeDir()

	documents := path.Join(home, "Documents")
	adobe := path.Join(documents, "Adobe")

	license1 := path.Join(documents, filename)
	license2 := path.Join(adobe, filename)

	if !util.IsDir(adobe) {
		_ = os.MkdirAll(adobe, os.ModePerm)
	}

	for _, fp := range []string{license1, license2} {
		log.Infof("adding default render-only-node licenses for After Effects at: - %s", fp)

		if !util.IsFile(fp) {
			_ = ioutil.WriteFile(fp, []byte{}, os.ModePerm)
		}
	}
}

func (s *Setting) autoBinary() {
	binUser := s.Binary
	if binUser != "" && util.IsFile(binUser) {
		return
	}

	binAuto := s.autoFind()
	if binAuto == "" {
		log.Fatal(`you should provide a proper path to After Effects\' \"aerender\" binary`)
		return
	}

	log.Infof("using automatically determined directory of After Effects installation: - %s", binAuto)
	s.Binary = binAuto
}

// 自动查找 AeRender文件
func (s *Setting) autoFind() string {
	platform := runtime.GOOS

	if defaultPaths[platform] == nil {
		return ""
	}

	bin := "aerender"
	if platform == "windows" {
		bin += ".exe"
	}

	for _, dir := range defaultPaths[platform] {
		if binPath := path.Join(dir, bin); util.IsFile(binPath) {
			return binPath
		}
	}
	return ""
}

// Patch 补丁
func (s *Setting) patch() (err error) {
	targetScript := "commandLineRenderer.jsx"

	afterEffects := path.Dir(s.Binary)
	originFile := path.Join(afterEffects, "Scripts", "Startup", targetScript)
	backupFile := path.Join(afterEffects, "Backup.Scripts", "Startup", targetScript)

	data, err := ioutil.ReadFile(originFile)
	if err != nil {
		log.Errorf("read command line fail: %v", err)
		return
	}

	log.Info("checking After Effects command line renderer patch...")

	if bytes.Contains(data, []byte("gexrender")) {
		log.Info("command line patch already is in place")

		if s.ForceCommandLinePatch {
			log.Info("forced rewrite of command line patch")

			err = ioutil.WriteFile(originFile, []byte(assets.CommandLineRendererScript), 0755)
			if err != nil {
				log.Errorf("forced rewrite of command line patch fail: %v", err)
				return
			}
		}
		return nil
	} else {
		log.Infof("backing up original command line script to: %s", backupFile)

		_ = os.MkdirAll(filepath.Join(afterEffects, "Backup.Scripts", "Startup"), 0755)
		err = util.CopyFile(originFile, backupFile)
		if err != nil {
			log.Error("")
		}

		log.Info("patching the command line script")

		defer func() {
			if err != nil && errors.Is(err, os.ErrPermission) {
				log.Info("\n\n              -- E R R O R --\n")
				log.Info("you need to run application with admin priviledges once")
				log.Info("to install Adobe After Effects commandLineRenderer.jsx patch\n")

				if runtime.GOOS == "windows" {
					log.Info("reading/writing inside Program Files folder on windows is blocked")
					log.Info("please run gexrender with Administrator Privilidges only ONE TIME, to install the patch\n\n")
				} else {
					log.Info("you might need to try to run gexrender with \"sudo\" only ONE TIME to install the patch\n\n")
				}
			}
		}()

		err = os.Chmod(originFile, 0755)
		if err != nil {
			return
		}

		err = ioutil.WriteFile(originFile, []byte(assets.CommandLineRendererScript), 0755)
		return
	}
}

var defaultPaths = map[string][]string{
	"darwin": {
		"/Applications/Adobe After Effects 2022",
		"/Applications/Adobe After Effects 2021",
		"/Applications/Adobe After Effects CC 2021",
		"/Applications/Adobe After Effects 2020",
		"/Applications/Adobe After Effects CC 2019",
		"/Applications/Adobe After Effects CC 2018",
		"/Applications/Adobe After Effects CC 2017",
		"/Applications/Adobe After Effects CC 2016",
		"/Applications/Adobe After Effects CC 2015",
		"/Applications/Adobe After Effects CC",
	},
	"windows": {
		"C:\\Program Files\\Adobe\\Adobe After Effects 2022\\Support Files",
		"C:\\Program Files\\Adobe\\Adobe After Effects 2021\\Support Files",
		"C:\\Program Files\\Adobe\\Adobe After Effects 2020\\Support Files",
		"C:\\Program Files\\Adobe\\Adobe After Effects CC 2020\\Support Files",
		"C:\\Program Files\\Adobe\\Adobe After Effects CC 2019\\Support Files",
		"C:\\Program Files\\Adobe\\Adobe After Effects CC 2018\\Support Files",
		"C:\\Program Files\\Adobe\\Adobe After Effects CC 2017\\Support Files",
		"C:\\Program Files\\Adobe\\Adobe After Effects CC 2016\\Support Files",
		"C:\\Program Files\\Adobe\\Adobe After Effects CC 2015\\Support Files",
		"C:\\Program Files\\Adobe\\Adobe After Effects CC\\Support Files",
		"C:\\Program Files\\Adobe\\Adobe After Effects CC",

		"C:\\Program Files\\Adobe\\After Effects 2022\\Support Files",
		"C:\\Program Files\\Adobe\\After Effects 2021\\Support Files",
		"C:\\Program Files\\Adobe\\After Effects 2020\\Support Files",
		"C:\\Program Files\\Adobe\\After Effects CC 2020\\Support Files",
		"C:\\Program Files\\Adobe\\After Effects CC 2019\\Support Files",
		"C:\\Program Files\\Adobe\\After Effects CC 2018\\Support Files",
		"C:\\Program Files\\Adobe\\After Effects CC 2017\\Support Files",
		"C:\\Program Files\\Adobe\\After Effects CC 2016\\Support Files",
		"C:\\Program Files\\Adobe\\After Effects CC 2015\\Support Files",
		"C:\\Program Files\\Adobe\\After Effects CC\\Support Files",
		"C:\\Program Files\\Adobe\\After Effects CC",
	},
}
