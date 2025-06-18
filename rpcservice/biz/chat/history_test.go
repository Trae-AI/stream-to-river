// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package chat

import (
	"context"
	"testing"

	"github.com/Trae-AI/stream-to-river/rpcservice/dal/redis"
	"github.com/cloudwego/eino/schema"
)

func TestConversation_Init(t *testing.T) {
	type fields struct {
		ConversationId string
		UserId         string
		History        []*schema.Message
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test_init",
			fields: fields{
				ConversationId: "124",
			},
			args: args{ctx: context.Background()},
		},
	}
	redis.InitCache()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conversation{
				ConversationId: tt.fields.ConversationId,
				UserId:         tt.fields.UserId,
				History:        tt.fields.History,
			}
			c.Init(tt.args.ctx)
			c.AppendMessages(tt.args.ctx, "pi", "qi")
			resp := c.GetHistory(tt.args.ctx)
			println(resp)
		})
	}
}
