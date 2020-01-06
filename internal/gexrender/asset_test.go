package gexrender

import "testing"

func TestAsset_WrapData1(t *testing.T) {
	asset := &Asset{
		Type:      "data",
		LayerName: "MyNicePicture.jpg",
		Property:  "Position",
		Value:     []int{500, 100},
	}

	s := asset.Wrap()
	t.Logf("\n %s", s)
}

func TestAsset_WrapData2(t *testing.T) {
	asset := &Asset{
		Type:       "data",
		LayerName:  "my text field",
		Property:   "Source Text",
		Expression: "time > 100 ? 'Bye bye' : 'Hello world'",
	}

	s := asset.Wrap()
	t.Logf("\n %s", s)
}

func TestAsset_WrapData3(t *testing.T) {
	asset := &Asset{
		Type:      "data",
		LayerName: "my text field",
		Property:  "Source Text.font",
		Value:     "Arial-BoldItalicMT",
	}

	s := asset.Wrap()
	t.Logf("\n %s", s)
}

func TestAsset_WrapData4(t *testing.T) {
	asset := &Asset{
		Type:      "data",
		LayerName: "background",
		Property:  "Effects.Skin_Color.Color",
		Value:     []int{1, 0, 0},
	}

	s := asset.Wrap()
	t.Logf("\n %s", s)
}

func TestAsset_FootageAudio(t *testing.T) {
	asset := &Asset{
		Type:        "image",
		Composition: "Audio Layer",
		LayerName:   "banner2",
		Dest:        "file:///C:\\Users\\Administrator\\Workspace\\aerender\\templates\\part1\\project\\resource\\Replace picture\\banner2.png",
	}

	s := asset.Wrap()
	t.Logf("\n %s", s)
}
