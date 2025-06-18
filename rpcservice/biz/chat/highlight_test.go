// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cloudwego/kitex/pkg/streaming"
	"github.com/Trae-AI/stream-to-river/internal/config"
	"github.com/Trae-AI/stream-to-river/rpcservice/biz/chat/coze"
	"github.com/Trae-AI/stream-to-river/rpcservice/dal/redis"
	"github.com/Trae-AI/stream-to-river/rpcservice/kitex_gen/words"
)

func TestHighlight(t *testing.T) {
	config.LoadConfig("../../")
	coze.InitCozeConfig(config.GetStringMapString("Coze"))
	highlightResourceChan := make(chan string, 100)
	highlightResultChan := make(chan []HighlightItem, 100)
	defer close(highlightResourceChan)
	go cozeHighlight(context.Background(), highlightResourceChan, highlightResultChan)
	go func() {
		for result := range highlightResultChan {
			fmt.Println(result)
		}
	}()
	highlightResourceChan <- "The ambiguous instructions made the task more comprehensive than expected."
	close(highlightResourceChan)
	time.Sleep(1 * time.Minute)
}

type MockStream struct {
}

func (m2 MockStream) SendMsg(ctx context.Context, m any) error {
	j, _ := json.MarshalIndent(m, "", " ")
	println(string(j))
	return nil
}

func (m2 MockStream) RecvMsg(ctx context.Context, m any) error {
	return nil
}

func (m2 MockStream) SetHeader(hd streaming.Header) error {
	return nil
}

func (m2 MockStream) SendHeader(hd streaming.Header) error {
	return nil
}

func (m2 MockStream) SetTrailer(hd streaming.Trailer) error {

	return nil
}

func TestHighlightWithMsgAggregation(t *testing.T) {
	redis.InitCache()
	stream := streaming.NewServerStreamingServer[words.ChatResp](&MockStream{})
	ChatHandler(context.Background(), &words.ChatReq{
		QueryMsg: "用英语介绍下字节跳动",
		//QueryMsg: "翻译一下下面这段话：1. 数据智能要大幅投入，是 AI Search 的基石和指挥棒，标注、评估、考题都要做到自动化，效率和规模提升百倍。\n2. 做好 listwise 选取、单 doc 摘要，目标是在不影响信息完备性的前提下，减少总输出 doc 数和单 doc 摘要长度。listwise 选取要注重增量价值和多样性。\n3. 精排端到端建模满意度，通过多目标 token 输出（相关/权威/时效/满意），提升可解释性和下游 listwise 选取能力，也利于清晰的团队分工。\n4. 召回、粗排全面大模型化，多语言问题引刃而解，可以学习自动标注、或者蒸馏 SuperRank。\n5. 权威性、时效性等方向，关键是全链路把对应目标、样本、特征搞对，提升大模型对站点、作者信息的理解，对 query 爆发速度的感知。\n6. 索引按垂类拆分，头部垂类做好筛选去重可大幅提升 ROI（文章+文库+企信占 82% 索引量）；重点垂类激进提升覆盖率（如官网、学术、金融等）。\n7. 成本必须激进优化，才能支撑激进的用户增长、Agent 使用增长（如 Deep Research、RL 训练等），全链路做好 cache。",
	}, stream)
}
