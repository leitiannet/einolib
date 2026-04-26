package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/adk/filesystem"
	"github.com/leitiannet/einolib"
	_ "github.com/leitiannet/einolib/operators"
)

func main() {
	ctx := context.Background()
	operator, err := einolib.NewOperator(ctx,
		einolib.WithOperatorType(einolib.OperatorTypeLocal),
		einolib.WithOperatorDescription("local operator example"),
	)
	if err != nil {
		panic(err)
	}
	fmt.Printf("operator type: %T\n", operator)

	dir := "examples/operator"
	files, err := operator.LsInfo(ctx, &filesystem.LsInfoRequest{Path: dir})
	if err != nil {
		panic(err)
	}
	fmt.Printf("files in %s:\n", dir)

	var mainPath string
	for _, file := range files {
		fmt.Printf("- %s\n", file.Path)
		if !file.IsDir && strings.EqualFold(filepath.Base(file.Path), "main.go") {
			if filepath.Dir(file.Path) == "." {
				mainPath = filepath.Join(dir, file.Path)
			} else {
				mainPath = file.Path
			}
		}
	}
	if mainPath == "" {
		fmt.Println("main.go not found")
		return
	}
	content, err := operator.Read(ctx, &filesystem.ReadRequest{FilePath: mainPath})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n%s content:\n%s\n", mainPath, content.Content)
}
