// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/Trae-AI/stream-to-river/rpcservice/dal/redis"
)

// Global variables defining configuration constants
var (
	// keyTemplate is the template for generating Redis keys based on conversation ID
	keyTemplate = "d_conv_%s"
	// maxTurn defines the maximum number of conversation turns to retain
	maxTurn = 2
	// conversationExpireTime sets the expiration time for conversations in Redis
	conversationExpireTime = 4 * time.Hour
)

// Conversation represents a chat conversation, containing conversation ID, user ID, and chat history
type Conversation struct {
	ConversationId string            `json:"conversation_id"`
	UserId         string            `json:"user_id"`
	History        []*schema.Message `json:"history"`
}

// cacheKey generates a Redis key for the conversation based on its ID
// Returns:
//   - string: The generated Redis key
func (c *Conversation) cacheKey() string {
	return fmt.Sprintf(keyTemplate, c.ConversationId)
}

// Init initializes the conversation by retrieving its history from Redis
// Parameters:
//   - ctx: The context for logging and potential cancellation
func (c *Conversation) Init(ctx context.Context) {
	// Attempt to get the conversation from Redis
	conv, exist := redis.Cache.Get(c.cacheKey())
	if !exist {
		// Log if the conversation is not found
		klog.CtxInfof(ctx, "conversation not found, id=%s", c.ConversationId)
		return
	}

	// Type assert the retrieved value to Conversation
	conversation, ok := conv.(*Conversation)
	if !ok {
		// Log if type assertion fails
		klog.CtxInfof(ctx, "conversation type assert failed, id=%s", c.ConversationId)
		return
	}
	// Update the conversation history
	c.History = conversation.History
}

// GetHistory retrieves the chat history of the conversation
// This function is non - strongly dependent and does not return an error
// Parameters:
//   - ctx: The context for the operation
//
// Returns:
//   - []*schema.Message: A slice of pointers to Message structs representing the chat history
func (c *Conversation) GetHistory(ctx context.Context) []*schema.Message {
	return c.History
}

// AppendMessages appends new user and assistant messages to the conversation history
// It also truncates the history to ensure it does not exceed the maximum number of turns
// Finally, it updates the conversation in Redis
// Parameters:
//   - ctx: The context for logging and potential cancellation
//   - query: The user's query message
//   - resp: The assistant's response message
//
// Returns:
//   - error: An error object if an unexpected error occurs during Redis operation
func (c *Conversation) AppendMessages(ctx context.Context, query string, resp string) error {
	// Append user message
	c.History = append(c.History, &schema.Message{
		Role:    schema.User,
		Content: query,
	})
	// Append assistant message
	c.History = append(c.History, &schema.Message{
		Role:    schema.Assistant,
		Content: resp,
	})

	// Truncate history to not exceed maxTurn
	if len(c.History) > maxTurn*2 {
		c.History = c.History[len(c.History)-maxTurn*2:]
	}

	// Update the conversation in Redis
	redis.Cache.Set(c.cacheKey(), c, conversationExpireTime)

	return nil
}

// NewConversation creates a new Conversation instance with the given conversation ID
// Parameters:
//   - conversationId: The unique identifier for the conversation
//
// Returns:
//   - *Conversation: A pointer to the newly created Conversation instance
func NewConversation(conversationId string) *Conversation {
	return &Conversation{ConversationId: conversationId}
}
