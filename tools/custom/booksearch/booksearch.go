// 图书检索示例自定义工具
package booksearch

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/leitiannet/einolib"
)

const (
	SearchBookToolName = "search_book"
	SearchBookToolDesc = "Search books based on user preferences"
)

type BookSearchInput struct {
	Genre     string `json:"genre" jsonschema:"description=Preferred book genre,enum=fiction,enum=sci-fi,enum=mystery,enum=biography,enum=business"`
	MaxPages  int    `json:"max_pages" jsonschema:"description=Maximum page length (0 for no limit)"`
	MinRating int    `json:"min_rating" jsonschema:"description=Minimum user rating (0-5 scale)"`
}

type BookSearchOutput struct {
	Books []string `json:"books"`
}

func searchBook(ctx context.Context, input *BookSearchInput) (*BookSearchOutput, error) {
	// search code
	// ...
	return &BookSearchOutput{Books: []string{"God's blessing on this wonderful world!"}}, nil
}

type BookSearchToolConfig struct{}

func NewBookSearchToolConfig(bookSearchToolOptions ...BookSearchToolOption) *BookSearchToolConfig {
	bookSearchToolConfig := &BookSearchToolConfig{}
	einolib.ApplyOptions(bookSearchToolConfig, bookSearchToolOptions)
	return bookSearchToolConfig
}

type BookSearchToolOption func(*BookSearchToolConfig)

func NewSearchBookTool(ctx context.Context, toolConfig *einolib.ToolConfig, bookSearchToolConfig *BookSearchToolConfig) ([]tool.BaseTool, error) {
	t, err := utils.InferTool(SearchBookToolName, SearchBookToolDesc, searchBook)
	if err != nil {
		return nil, err
	}
	return []tool.BaseTool{t}, nil
}

func createSearchBookTool(ctx context.Context, toolConfig *einolib.ToolConfig, specificConfig interface{}) ([]tool.BaseTool, error) {
	bookSearchToolConfig, err := einolib.ParseSpecificConfig(specificConfig, func() *BookSearchToolConfig { return NewBookSearchToolConfig() })
	if err != nil {
		return nil, err
	}
	return NewSearchBookTool(ctx, toolConfig, bookSearchToolConfig)
}

func init() {
	if err := einolib.RegisterToolConstructFunc(einolib.ToolTypeCustom, SearchBookToolName, createSearchBookTool, (*BookSearchToolConfig)(nil)); err != nil {
		einolib.GetLogger().Errorf("register tool %s failed: %v", SearchBookToolName, err)
	}
}
