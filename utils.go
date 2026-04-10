package einolib

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type PrintJSONOptions struct {
	Prefix string
	Indent bool
}

func NewPrintJSONOptions(prefix string, indent bool) *PrintJSONOptions {
	return &PrintJSONOptions{
		Prefix: prefix,
		Indent: indent,
	}
}

// 主要用于调试输出，错误可忽略
func PrintJSON(obj interface{}, options *PrintJSONOptions) error {
	var (
		content []byte
		err     error
	)
	if options == nil {
		options = NewPrintJSONOptions("", false)
	}
	if options.Indent {
		content, err = json.MarshalIndent(obj, "", "  ")
	} else {
		content, err = json.Marshal(obj)
	}
	if err != nil {
		return err
	}
	if options.Prefix != "" {
		fmt.Printf("%s %s\n", options.Prefix, string(content))
	} else {
		fmt.Printf("%s\n", string(content))
	}
	return nil
}

// 从环境变量加载值到字符串指针，支持自动拼接前缀并转大写
func BindVarFromEnv(target *string, key string, prefixes ...string) {
	if target == nil {
		return
	}
	if len(prefixes) > 0 && prefixes[0] != "" {
		key = fmt.Sprintf("%s_%s", prefixes[0], key)
	}
	envKey := strings.ToUpper(key)
	if val := os.Getenv(envKey); val != "" {
		*target = val
	}
}
