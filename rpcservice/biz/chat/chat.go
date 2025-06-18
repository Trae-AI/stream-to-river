// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package chat

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
	"github.com/cloudwego/gopkg/concurrency/gopool"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/Trae-AI/stream-to-river/rpcservice/biz/llm"
	llmconfig "github.com/Trae-AI/stream-to-river/rpcservice/conf/llm_config"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
	"github.com/Trae-AI/stream-to-river/rpcservice/utils"
)

// Warning thresholds for latency
var chatTimeToFirstModelChunkWarningThreshold = time.Second
var chatTimeToFirstHighlightWarningThreshold = 2 * time.Second
var chatTimeToFirstResponseWarningThreshold = 2 * time.Second

// isValidConv checks if the conversation response is valid.
// It first checks if the response is in the refused string list.
// Then it waits for the review result for 300 milliseconds.
//
// Parameters:
//   - ctx: The context for the request, used for logging and cancellation.
//   - response: The response string to be checked.
//   - reviewEnd: A channel to receive the review result.
//
// Returns:
//   - bool: true if the response is valid, false otherwise.
func isValidConv(ctx context.Context, response string, reviewEnd chan bool) bool {
	if utils.IsInRefusedString(response) {
		klog.CtxWarnf(ctx, "response is refused, resp=%s", response)
		return false
	}

	select {
	case reviewPass := <-reviewEnd:
		return reviewPass
	// Most cases should return from the above branch as the service p99 is 300ms.
	case <-time.After(300 * time.Millisecond):
		klog.CtxWarnf(ctx, "[chat] review pass not received in time, waited for 300ms")
	}

	return true
}

// ChatHandler handles the chat request with streaming.
// It initializes the conversation, creates a model stream, and processes model chunks.
// It also handles highlight resources and sends responses via the stream.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request-scoped values.
//   - req: A pointer to the ChatReq struct containing the request information.
//   - stream: A WordService_ChatServer interface for handling streaming responses.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the process.
func ChatHandler(ctx context.Context, req *words.ChatReq, stream words.WordService_ChatServer) (err error) {
	startTime := time.Now()
	conversation := NewConversation(req.ConversationId)
	conversation.Init(ctx)

	reviewEnd := make(chan bool, 1)
	gopool.CtxGo(ctx, func() {
		reviewEnd <- PromptReview(ctx, req.QueryMsg)
	})
	modelMsg, err := ArkModelStreamMsg(ctx, conversation.GetHistory(ctx), req.QueryMsg)
	if err != nil {
		return kerrors.NewBizStatusError(-1, err.Error())
	}
	firstModelChunk := false
	firstHighlightRecv := false
	firstResponse := false

	highlightResourceChan := make(chan string, 100)
	highlightResultChan := make(chan []HighlightItem, 100)
	highLightEnd := make(chan struct{}, 1)

	var modelMsgSendDone int32 = 0
	defer func() {
		atomic.StoreInt32(&modelMsgSendDone, 1)
	}()

	gopool.CtxGo(ctx, func() {
		cozeHighlight(ctx, highlightResourceChan, highlightResultChan)
	})
	gopool.CtxGo(ctx, func() {
		for highlightItems := range highlightResultChan {
			if !firstHighlightRecv {
				firstHighlightRecv = true
				klog.CtxInfof(ctx, "[chat] first highlight resource received, time: %s", time.Since(startTime))
				if time.Since(startTime) > chatTimeToFirstHighlightWarningThreshold {
					klog.CtxWarnf(ctx, "[chat] time to first highlight resource is too long, time: %s", time.Since(startTime))
				}
			}
			metaInfo, _ := sonic.MarshalString([]map[string]interface{}{
				{
					"type":  1,
					"items": highlightItems,
				},
			})
			errSend := stream.Send(ctx, &words.ChatResp{Extra: map[string]string{
				"meta_info": metaInfo,
			}})
			if errSend != nil && atomic.LoadInt32(&modelMsgSendDone) != 1 {
				klog.CtxErrorf(ctx, "send highlight item failed, err=%v", errSend)
			}
		}
		highLightEnd <- struct{}{}
	})

	var mu sync.Mutex // protect respFullText
	respFullText := ""
	defer func() {
		close(highlightResourceChan)
		// wait for highlight End
		waitHighlightTimeStart := time.Now()
		select {
		case <-highLightEnd:
			klog.CtxInfof(ctx, "[chat] highlight end received, time: %s", time.Since(waitHighlightTimeStart))
		case <-time.After(2 * time.Second):
			klog.CtxWarnf(ctx, "[chat] highlight end not received in time, waited for 2 seconds")
		}

		mu.Lock()
		if isValidConv(ctx, respFullText, reviewEnd) {
			appendErr := conversation.AppendMessages(ctx, req.GetQueryMsg(), respFullText)
			if appendErr != nil {
				klog.CtxErrorf(ctx, "Conversation.AppendMessages failed, err=%v", appendErr)
			}
		}
		mu.Unlock()

		klog.CtxInfof(ctx, "[chat] total latency: %s", time.Since(startTime))
	}()

	for {
		chunk, e := modelMsg.Recv()
		if !firstModelChunk {
			firstModelChunk = true
			klog.CtxInfof(ctx, "[chat] first model chunk received, time: %s", time.Since(startTime))
			if time.Since(startTime) > chatTimeToFirstModelChunkWarningThreshold {
				klog.CtxWarnf(ctx, "[chat] time to first model chunk is too long, time: %s", time.Since(startTime))
			}
		}
		if e != nil {
			if e == io.EOF {
				break
			}
			klog.CtxErrorf(ctx, "arkModel Stream recv failed, err=%v", e)
			return
		}
		mu.Lock()
		respFullText += chunk.Content
		mu.Unlock()
		if err = stream.Send(ctx, &words.ChatResp{Msg: chunk.Content}); err != nil {
			klog.CtxErrorf(ctx, "ChatService stream send failed, err=%v", err)
			return
		}
		if !firstResponse {
			firstResponse = true
			klog.CtxInfof(ctx, "[chat] first response sent, time: %s", time.Since(startTime))
			if time.Since(startTime) > chatTimeToFirstResponseWarningThreshold {
				klog.CtxWarnf(ctx, "[chat] time to first response is too long, time: %s", time.Since(startTime))
			}
		}
		highlightResourceChan <- chunk.Content
	}

	return
}

// ArkModelStreamMsg creates a streaming message from the Ark model.
// It initializes the Ark model and creates a streaming message based on the conversation history and user message.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request-scoped values.
//   - history: A slice of pointers to Message structs representing the conversation history.
//   - userMsg: The user's query message.
//
// Returns:
//   - *schema.StreamReader[*schema.Message]: A pointer to the stream reader for the model messages.
//   - error: An error object if an unexpected error occurs during the process.
func ArkModelStreamMsg(ctx context.Context, history []*schema.Message, userMsg string) (outStream *schema.StreamReader[*schema.Message], err error) {
	arkModel, err := llm.GetArkModel()
	if err != nil {
		err = fmt.Errorf("arkModel init failed, err=%v", err)
		klog.CtxErrorf(ctx, "%s", err)
		return
	}

	inMsg := append([]*schema.Message{
		schema.SystemMessage(llmconfig.ChatPE),
	}, history...)
	inMsg = append(inMsg, schema.UserMessage(userMsg))

	chatMsg, err := arkModel.Stream(ctx, inMsg)

	if err != nil {
		err = fmt.Errorf("arkModel create Stream failed, err=%v", err)
		klog.CtxErrorf(ctx, err.Error())
		return nil, err
	}
	return chatMsg, nil
}
