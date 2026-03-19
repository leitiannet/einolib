package einolib

import (
	"encoding/json"
	"fmt"
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
