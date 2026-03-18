package i18n

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	errorsMap = map[string]map[int]string{}
	keysMap   = map[string]map[string]string{}
)

type ctxKey string

const LangKey ctxKey = "lang"

// Load 加在指定语言文件
func Load(lang string, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取语言文件失败 %s: %w", lang, err)
	}
	raw := map[string]any{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("解析语言文件失败 %s: %w", lang, err)
	}
	errMap := map[int]string{}
	keyMap := map[string]string{}
	for k, v := range raw {
		// error code
		if k == "error" {
			if m, ok := v.(map[interface{}]interface{}); ok {
				for codeKey, msg := range m {
					// 处理 key (可能是 int 或 float64)
					var code int
					switch c := codeKey.(type) {
					case int:
						code = c
					case float64:
						code = int(c)
					case string:
						if cInt, err := strconv.Atoi(c); err == nil {
							code = cInt
						} else {
							continue
						}
					default:
						fmt.Printf("跳过不支持的 key 类型: %T\n", codeKey)
						continue
					}
					// 处理 value
					errMap[code] = fmt.Sprint(msg)
				}
			}
			continue
		}
		flatten(keyMap, k, v)
	}
	errorsMap[lang] = errMap
	keysMap[lang] = keyMap
	return nil
}

// LoadLocales 扫描目录加载所有语言文件
func LoadLocales(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取语言目录失败: %w", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			continue
		}
		lang := strings.TrimSuffix(name, filepath.Ext(name))
		path := filepath.Join(dir, name)
		if err := Load(lang, path); err != nil {
			return err
		}
	}
	if len(errorsMap) == 0 && len(keysMap) == 0 {
		return fmt.Errorf("没有加载到任何语言文件")
	}
	return nil
}

// GetMessage 获取指定 Code 对应语言内容
func GetMessage(lang string, code int) string {
	fmt.Println("请求语言", lang, code)
	if langMap, ok := errorsMap[lang]; ok {
		if msg, ok := langMap[code]; ok {
			return msg
		}
	}
	// fallback zh
	if zhMap, ok := errorsMap["zh"]; ok {
		if msg, ok := zhMap[code]; ok {
			return msg
		}
	}
	return "unknown error"
}

// T 获取指定 Key 对应语言内容
func T(lang string, key string, args ...string) string {
	if len(args) > 0 {
		key = key + "." + strings.Join(args, ".")
	}
	if langMap, ok := keysMap[lang]; ok {
		if msg, ok := langMap[key]; ok {
			return msg
		}
	}
	// fallback zh
	if zhMap, ok := keysMap["zh"]; ok {
		if msg, ok := zhMap[key]; ok {
			return msg
		}
	}
	return key
}

// flatten 用于拼接 YAML 配置的 Key
func flatten(result map[string]string, prefix string, value any) {
	switch v := value.(type) {
	case string:
		result[prefix] = v
	case map[string]any:
		for k, val := range v {
			key := prefix + "." + k
			flatten(result, key, val)
		}
	}
}

// GetLang 获取当前上下文使用语言
func GetLang(ctx context.Context) string {
	lang, ok := ctx.Value(LangKey).(string)
	if !ok {
		return "zh"
	}
	return lang
}
