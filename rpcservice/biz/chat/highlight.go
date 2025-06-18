// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package chat

import (
	"context"
	"strings"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/coze-dev/coze-go"

	cozeClient "github.com/Trae-AI/stream-to-river/rpcservice/biz/chat/coze"
)

type HighlightItem struct {
	Text  string `json:"text"`
	Start int64  `json:"start"`
	End   int64  `json:"end"`
}

// highlight processes the incoming resources, generates a list of words using LLM,
// checks word details, and sends highlight items to the result channel.
// Although the current return results do not guarantee order,
// the code processes them based on index, so it's acceptable.
//
// Parameters:
//   - ctx: The context for the request, used for cancellation, deadlines, and passing request-scoped values.
//   - highlightResourceChan: A channel for receiving incoming resources.
//   - highlightResultChan: A channel for sending generated highlight items.
func cozeHighlightWords(ctx context.Context, text string) ([]string, error) {
	cli := cozeClient.GetCozeClient()

	req := &coze.RunWorkflowsReq{
		WorkflowID: cozeClient.CozeConf.WorkflowID,
		Parameters: map[string]any{
			"input": text,
		},
		IsAsync: false,
	}

	resp, err := cli.Workflows.Runs.Create(ctx, req)
	if err != nil {
		klog.CtxErrorf(ctx, "failed to run coze workflow, err: %v", err)
		return nil, err
	}

	type result struct {
		Output []struct {
			Word string `json:"word"`
		} `json:"output"`
	}

	var res result
	err = sonic.UnmarshalString(resp.Data, &res)
	if err != nil {
		klog.CtxErrorf(ctx, "failed to unmarshal coze workflow response, err: %v", err)
		return nil, err
	}

	var ret []string
	for _, item := range res.Output {
		ret = append(ret, item.Word)
	}

	return ret, nil
}

type syncMap struct {
	sync.Mutex
	m map[string]bool
}

func (sm *syncMap) Set(key string) {
	sm.Lock()
	sm.m[key] = true
	sm.Unlock()
}

func (sm *syncMap) Get(key string) bool {
	sm.Lock()
	v := sm.m[key]
	sm.Unlock()
	return v
}

func (sm *syncMap) SetOrGet(key string) bool {
	var ret bool

	sm.Lock()
	ret = !sm.m[key]
	if ret {
		sm.m[key] = true
	}
	sm.Unlock()

	return ret
}

type indexPair struct {
	start int64
	end   int64
}

func (i indexPair) notOverlap(ip indexPair) bool {
	return i.end <= ip.start || i.start >= ip.end
}

type indexDedup struct {
	sync.Mutex
	ips []indexPair
}

func anyOf[T any](s []T, f func(T) bool) bool {
	for _, item := range s {
		if f(item) {
			return true
		}
	}

	return false
}

func (i *indexDedup) Append(ip indexPair) bool {
	var ret bool
	i.Lock()
	ret = !anyOf(i.ips, func(ip_ indexPair) bool {
		return !ip.notOverlap(ip_)
	})

	if ret {
		i.ips = append(i.ips, ip)
	}
	i.Unlock()

	return ret
}

func cozeHighlight(ctx context.Context, highlightResourceChan chan string, highlightResultChan chan []HighlightItem) {

	var sentence string
	var wg sync.WaitGroup
	wordDedup := &syncMap{
		m: make(map[string]bool),
	}

	idxDedup := &indexDedup{}
	var sentenceOffset int

	highlightFunc := func(sentence string, offset int64) {
		wg.Add(1)
		go func() {
			defer wg.Done()

			words_, err := cozeHighlightWords(ctx, sentence)
			if err != nil {
				klog.CtxErrorf(ctx, "failed to highlight sentence, err: %v", err)
				return
			}

			items := make([]HighlightItem, 0, len(words_))
			lowerSentence := strings.ToLower(sentence)
			for _, w := range words_ {
				if wordDedup.SetOrGet(strings.ToLower(w)) {
					index := strings.Index(lowerSentence, w)
					if index >= 0 {
						start := int64(index) + offset
						end := start + int64(len(w))
						ip := indexPair{
							start: start,
							end:   end,
						}

						if idxDedup.Append(ip) {
							items = append(items, HighlightItem{
								Text:  sentence[index : index+len(w)],
								Start: start,
								End:   end,
							})
						} else {
							klog.CtxInfof(ctx, "index invalid, word=%s, start_index=%d, end_index=%d", w, start, end)
						}
					}
				}
			}

			if len(items) > 0 {
				highlightResultChan <- items
			}
		}()

		sentenceOffset += len(sentence)
	}

	for w := range highlightResourceChan {
		if strings.HasSuffix(sentence, "。") || strings.HasSuffix(sentence, ". ") {
			highlightFunc(sentence, int64(sentenceOffset))
			sentence = w
			continue
		}

		if strings.HasSuffix(sentence, "\n") || strings.HasPrefix(w, "\n") {
			highlightFunc(sentence, int64(sentenceOffset))
			sentence = w
			continue
		}

		if strings.HasSuffix(sentence, ".") && strings.HasPrefix(w, " ") {
			highlightFunc(sentence, int64(sentenceOffset))
			sentence = w
			continue
		}

		if len(sentence) > 50 && (strings.HasSuffix(sentence, " ") || strings.HasPrefix(w, " ")) {
			highlightFunc(sentence, int64(sentenceOffset))
			sentence = w
			continue
		}

		sentence += w
	}

	if len(sentence) > 0 {
		highlightFunc(sentence, int64(sentenceOffset))
	}

	wg.Wait()
	close(highlightResultChan)
}
