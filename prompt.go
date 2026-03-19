package einolib

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func FormatMessages(ctx context.Context, formatType schema.FormatType, templates []schema.MessagesTemplate, variables map[string]any) ([]*schema.Message, error) {
	template := prompt.FromMessages(formatType, templates...)
	return template.Format(ctx, variables)
}
