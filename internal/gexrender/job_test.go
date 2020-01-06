package gexrender

import (
	"io/ioutil"
	"testing"
)

func TestCreateJob(t *testing.T) {
	setting := &Setting{
		ConfigFile: "./testdata/myjob.json",
	}

	job, err := CreateJob(setting)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(job)
}

func TestJob_CreateScriptFile(t *testing.T) {
	workPath, _ := ioutil.TempDir("", "")

	setting := &Setting{
		ConfigFile: "./testdata/myjob.json",
		WorkPath:   workPath,
	}

	job, err := CreateJob(setting)
	if err != nil {
		t.Fatal(err)
		return
	}

	err = job.CreateScriptFile()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(job.ScriptFile)
}
