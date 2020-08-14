package gexrender

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	log "github.com/sirupsen/logrus"
	"github.com/yearnfar/gexrender/internal/gexrender/assets"
)

// Job 任务
type Job struct {
	Uid        string   `json:"-"` // 唯一id
	WorkPath   string   `json:"-"` // 工作目录
	ScriptFile string   `json:"-"` // 脚本文件
	ResultName string   `json:"-"` // 输出文件名
	Output     string   `json:"-"` // 输出文件
	Setting    *Setting `json:"-"` // 参数配置

	Template *Template            `json:"template" validate:"required"` // 模板
	Assets   []*Asset             `json:"assets" validate:"required"`   // 可替换资源
	Actions  map[string][]*Action `json:"actions" validate:"-"`         // Action

}

var renderTimeRegex = regexp.MustCompile(`PROGRESS:  Total Time Elapsed: (\d+) Seconds`)
var renderErrorRegex = regexp.MustCompile(`Error: gexrender:(.*)`)
var aeErrorRegex = regexp.MustCompile(`aerender ERROR:(.*)`)

// CreateJob 创建任务
func CreateJob(setting *Setting) (job *Job, err error) {
	var data []byte
	if setting.ConfigFile != "" {
		data, err = ioutil.ReadFile(setting.ConfigFile)
		if err != nil {
			return
		}
	} else {
		data = []byte(setting.Config)
	}

	job = &Job{}
	err = json.Unmarshal(data, job)
	if err != nil {
		return
	}

	uid, err := gonanoid.Nanoid()
	if err != nil {
		return
	}

	job.Uid = uid
	job.Setting = setting
	return
}

// Run 执行
func (j *Job) Run() (err error) {
	_ = os.MkdirAll(j.WorkPath, os.ModePerm)
	defer j.CleanUp()

	err = j.Fetch()
	if err != nil {
		log.Error(err)
		return
	}

	err = j.Invoke("prerender")
	if err != nil {
		log.Error(err)
		return
	}

	err = j.CreateScriptFile()
	if err != nil {
		log.Error(err)
		return
	}

	err = j.Render()
	if err != nil {
		log.Error(err)
		return
	}

	err = j.Invoke("postrender")
	if err != nil {
		log.Error(err)
		return
	}
	return
}

// Invoke 调用Action
func (j *Job) Invoke(actionType string) (err error) {
	log.Infof("[%s] invoker %s", j.Uid, actionType)

	actions, ok := j.Actions[actionType]
	if !ok {
		return
	}

	for _, action := range actions {
		data, _ := action.Parameter.MarshalJSON()
		err = Invoke(action.Action, j, data)
		if err != nil {
			log.Errorf("[%s] invoker err: %v", j.Uid, err)
			return
		}
	}
	return
}

// Fetch 获取资源
func (j *Job) Fetch() (err error) {
	err = j.Template.Fetch(j.WorkPath)
	if err != nil {
		return
	}

	for _, asset := range j.Assets {
		if asset.Src != "" {
			err = asset.Fetch(j.WorkPath)
			if err != nil {
				return
			}
		}
	}
	return
}

// CreateScriptFile 生成脚本文件
func (j *Job) CreateScriptFile() error {
	log.Infof("[%s] running script assemble...", j.Uid)

	var data []string
	for _, asset := range j.Assets {
		data = append(data, asset.Wrap())
	}

	// write out assembled custom script file in the workpath
	j.ScriptFile = path.Join(j.WorkPath, fmt.Sprintf("gexrender-%s-script.jsx", j.Uid))

	script := assets.GexRenderScript
	script = strings.ReplaceAll(script, "/*COMPOSITION*/", j.Template.Composition)
	script = strings.ReplaceAll(script, "/*USERSCRIPT*/", strings.Join(data, "\n\n"))

	return ioutil.WriteFile(j.ScriptFile, []byte(script), os.ModePerm)
}

// Render 执行渲染
func (j *Job) Render() (err error) {
	startTime := time.Now()
	log.Infof("[%s] rendering job...", j.Uid)
	var param []string

	param = append(param, "-project", j.Template.Dest)
	param = append(param, "-comp", j.Template.Composition)
	param = append(param, "-output", j.Output)

	if j.Template.OutputModule != "" {
		param = append(param, "-OMtemplate", j.Template.OutputModule)
	}
	if j.Template.SettingsTemplate != "" {
		param = append(param, "-RStemplate", j.Template.SettingsTemplate)
	}
	if j.Template.FrameStart != "" {
		param = append(param, "-s", j.Template.FrameStart)
	}
	if j.Template.FrameEnd != "" {
		param = append(param, "-e", j.Template.FrameEnd)
	}
	if j.Template.IncrementFrame != "" {
		param = append(param, "-i", j.Template.IncrementFrame)
	}
	if j.ScriptFile != "" {
		param = append(param, "-r", j.ScriptFile)
	}

	if j.Setting.MultiFrames {
		param = append(param, "-mp")
	}
	if j.Setting.Reuse {
		param = append(param, "-reuse")
	}
	if j.Template.ContinueOnMissing {
		param = append(param, "-continueOnMissingFootage")
	}
	if j.Setting.ImageCachePercent > 0 || j.Setting.MaxMemoryPercent > 0 {
		imageCachePercent := j.Setting.ImageCachePercent
		if imageCachePercent == 0 {
			imageCachePercent = 50
		}

		maxMemoryPercent := j.Setting.MaxMemoryPercent
		if maxMemoryPercent == 0 {
			maxMemoryPercent = 50
		}

		param = append(param, "-mem_usage", strconv.Itoa(imageCachePercent), strconv.Itoa(maxMemoryPercent))
	}

	log.Infof(`[%s] spawning aerender process: %s %s`, j.Uid, j.Setting.Binary, strings.Join(param, " "))

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(j.Setting.Binary, param...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err = cmd.Run()

	var output []string
	if stdout.Len() > 0 {
		output = append(output, stdout.String())
	}
	if stderr.Len() > 0 {
		output = append(output, stderr.String())
	}

	outputStr := strings.Join(output, "")
	logPath := filepath.Join(j.WorkPath, fmt.Sprintf(`../aerender-%s.log`, j.Uid))
	_ = ioutil.WriteFile(logPath, []byte(outputStr), 0644)

	if err != nil {
		log.Errorf("[%s]Error starting aerender process:: %v", j.Uid, err)
	}

	if err != nil && j.Setting.StopOnError {
		log.Errorf("[%s]stop on error!", j.Uid)
		return
	}

	// 渲染成功
	if strings.Contains(outputStr, "Finished composition") {
		timeMatches := renderTimeRegex.FindStringSubmatch(outputStr)
		if len(timeMatches) == 2 {
			log.Infof(`[%s] rendering took ~%d sec.`, j.Uid, timeMatches[1])
		} else {
			log.Infof(`[%s] rendering took ~%.0f sec.`, j.Uid, time.Since(startTime).Seconds())
		}
		return
	} else {
		log.Infof(`[%s] rendering took ~%.0f sec.`, j.Uid, time.Since(startTime).Seconds())

		errMatches := renderErrorRegex.FindStringSubmatch(outputStr)
		if len(errMatches) > 0 {
			return errors.New(strings.TrimSpace(errMatches[1]))
		}

		aeErrMatches := aeErrorRegex.FindStringSubmatch(outputStr)
		if len(aeErrMatches) > 0 {
			return errors.New(strings.TrimSpace(aeErrMatches[1]))
		}

		return errors.New("渲染失败")
	}
}

func (j *Job) Parse(s string) {

}

// CleanUp 清理数据
func (j *Job) CleanUp() {
	if !j.Setting.SkipCleanup {
		log.Infof("[%s] cleaning up...", j.Uid)
		_ = os.RemoveAll(j.WorkPath)
	} else {
		log.Infof("[%s] skipping the clean up...", j.Uid)
	}
}
