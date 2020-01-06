package gexrender

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/yearnfar/gexrender/internal/pkg/util"
)

// Asset 可替换资源
type Asset struct {
	Type        string      `json:"type" validate:"required"` // 类型
	Src         string      `json:"src"  validate:"required"` // 文件地址
	Dest        string      `json:"dest"`                     // 下载文件地址
	Composition string      `json:"composition"`              // 合成
	LayerName   string      `json:"layer_name"`               // 图层名
	LayerIndex  int         `json:"layer_index"`              // 图层索引
	Value       interface{} `json:"value"`                    // 值
	Expression  interface{} `json:"expression"`               // 表达式
	Property    string      `json:"property"`                 // 属性
}

// Fetch 抓取资源
func (a *Asset) Fetch(saveDir string) error {
	dest, err := Fetch(a.Src, saveDir)
	if err != nil {
		return err
	}

	a.Dest = dest
	return nil
}

// Wrap 转换
func (a *Asset) Wrap() string {
	switch a.Type {
	case "video", "audio", "image":
		var scripts []string
		if a.Src != "" {
			scripts = append(scripts, a.wrapFootage())
		}
		if a.Property != "" {
			scripts = append(scripts, a.wrapData())
		}
		return strings.Join(scripts, "\n\n")
	case "data":
		return a.wrapData()
	case "script":
		return a.wrapScript()
	default:
		return ""
	}
}

// scripting wrappers
func (a *Asset) wrapFootage() string {
	return `(function() {
	gexrender.` + a.getMethod() + `(` + a.getComposition() + `, ` + a.getValue() + `, ` + `function(layer) {
		gexrender.replaceFootage(layer, '` + strings.ReplaceAll(a.Dest, "\\", "\\\\") + `');
	})
})();`
}

func (a *Asset) wrapData() string {
	return `(function() {
	gexrender.` + a.getMethod() + `(` + a.getComposition() + `, ` + a.getValue() + `, ` + `function(layer) {
		var parts = ` + util.Stringify(a.partsOfKeypath(a.Property)) + `;
		` + a.renderIf(a.Value, `var value = {"value": $value};`) + `
		` + a.renderIf(a.Expression, `var value = {"expression": $value};`) + `

		gexrender.changeValueForKeypath(layer, parts, value);
		return true;
	})
})();`
}

func (a *Asset) wrapScript() string {
	script, _ := ioutil.ReadFile(a.Dest)
	return `(function() {
	` + string(script) + `
})();`
}

func (a *Asset) getMethod() string {
	if a.LayerName != "" {
		return "selectLayersByName"
	} else {
		return "selectLayersByIndex"
	}
}

func (a *Asset) getValue() string {
	if a.LayerName != "" {
		return escape(a.LayerName)
	} else {
		return strconv.Itoa(a.LayerIndex)
	}
}

func (a *Asset) getComposition() string {
	if a.Composition != "" {
		return escape(a.Composition)
	} else {
		return "null"
	}
}

func (a *Asset) renderIf(val interface{}, str string) string {
	if val == nil {
		return ""
	}

	var encoded string
	if v, ok := val.(string); ok {
		encoded = escape(v)
	} else {
		encoded = util.Stringify(val)
	}

	return strings.ReplaceAll(str, "$value", encoded)
}

func (a *Asset) partsOfKeypath(keypath string) []string {
	parts := strings.Split(keypath, "->")
	if len(parts) > 1 {
		return parts
	} else {
		return strings.Split(keypath, ".")
	}
}

func escape(str string) string {
	d, _ := json.Marshal(str)
	s := string(d)
	s = s[1 : len(s)-1]
	return `'` + strings.ReplaceAll(s, `'`, `\'`) + `'`
}
