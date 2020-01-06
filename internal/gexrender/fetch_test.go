package gexrender

import (
	"io/ioutil"
	"testing"
)

func TestFetch(t *testing.T) {
	urls := []string{
		"https://www.yearnfar.com/wp-content/themes/yzipicc/images/logo.png",
		"file:///C:\\Users\\Administrator\\Workspace\\go\\src\\github.com\\yearnfar\\gexrender\\pkg\\gexrender\\testdata\\resource\\video1.mp4",
	}

	saveDir, _ := ioutil.TempDir("", "gexrender")

	for _, url := range urls {
		saveFile, err := Fetch(url, saveDir)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(saveFile)
	}
}
