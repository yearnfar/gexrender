package util

import (
	"testing"
)

func TestUriToPath(t *testing.T) {
	uris := []string{
		"file:///c:%5CUsers%5CAdminstrator%5CDocuments%5Cabc.text",
		"file:///c%7C/Users%5CAdminstrator%5CDocuments%5Cabc.text",
		"file:///Users/yearnfar/Workspace/go/src/github.com/yearnfar/gexrender/gexrender.go",
	}

	for _, uri := range uris {
		uri2, err := UriToPath(uri)
		if err != nil {
			t.Fatal(err)
		}

		t.Log(uri2)
	}
}
