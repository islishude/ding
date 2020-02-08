package ding

import "context"

// FeedCard 类型
type FeedCard struct {
	Title  string `json:"title"`      // 单条信息文本
	MsgURL string `json:"messageURL"` // 点击单条信息到跳转链接
	PicURL string `json:"picURL"`     // 单条信息后面图片的URL
}

// NewFeedCard 创建 FeedCard 类型
func NewFeedCard(title, msgurl, picurl string) *FeedCard {
	return &FeedCard{
		Title:  title,
		MsgURL: msgurl,
		PicURL: picurl,
	}
}

func (d *clientimpl) SendFeedCard(ctx context.Context, links ...*FeedCard) error {
	data := map[string]interface{}{
		"msgtype": "feedCard",
		"feedCard": map[string][]*FeedCard{
			"links": links,
		},
	}
	return d.request(ctx, data)
}
