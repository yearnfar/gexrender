package gexrender

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	log "github.com/sirupsen/logrus"
)

var validate *validator.Validate
var trans ut.Translator

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)

	trans, _ = ut.New(zh.New()).GetTranslator("zh")

	validate = validator.New()
	err := zh_translations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		log.Panic(err)
	}
}

// Render 渲染
func Render(setting *Setting) (err error) {
	if setting.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if setting.LogFile != "" {
		logFile, err := os.OpenFile(setting.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			log.Error(err)
			return err
		}

		log.SetOutput(logFile)
	}

	if setting.WorkPath == "" {
		setting.WorkPath, err = ioutil.TempDir("", "gexrender")
		if err != nil {
			log.Error(err)
			return err
		}
	} else if !filepath.IsAbs(setting.WorkPath) {
		setting.WorkPath, err = filepath.Abs(setting.WorkPath)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	if setting.Config == "" && setting.ConfigFile == "" {
		err = errors.New("you need to provide a gexrender job json as an argument")
		return
	}

	job, err := setting.Setup()
	if err != nil {
		log.Error(err)
		return
	}

	err = job.Run()
	if err != nil {
		log.Error(err)
		return
	}
	return
}
