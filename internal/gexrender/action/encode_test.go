package action

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/yearnfar/gexrender/internal/gexrender"
)

func TestEncode_Invoke(t *testing.T) {
	workPath, _ := ioutil.TempDir("", "gexrender")
	job := &gexrender.Job{WorkPath: workPath}
	input, _ := filepath.Abs("./testdata/video1.mp4")

	param := encodeParam{
		Preset: "mp4",
		Input:  input,
		Output: "video1.mp4",
	}

	data, _ := json.Marshal(param)

	e := &Encode{}
	err := e.Invoke(job, data)
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(workPath)
}
