package ding

import "context"

// Ding interface
type Ding interface {
	// 发送text类型消息
	// content 消息内容 isAtAll 可选，@所有人时：true，否则为：false
	// atMobiles 可选，被@人的手机号(在content里添加@人的手机号)
	SendText(ctx context.Context, content string, isAtAll bool, atMobiles ...string) error
	// 发送Link类型消息
	// title：消息标题 text：消息内容，如果太长只会部分展示
	// msgurl：点击消息跳转的URL picurl：图片URL，可选
	SendLink(ctx context.Context, text, title, msgurl string, picurl ...string) error
	// 发送markdown类型消息
	// title：消息标题 text：消息内容，如果太长只会部分展示
	// isAtAll 可选，@所有人时：true，否则为：false
	// atMobiles 可选，被@人的手机号(在content里添加@人的手机号)
	SendMarkdown(ctx context.Context, title, text string, isAtAll bool, atMobiles ...string) error
	// 发送带有按钮菜单的 ActionCard 消息
	// alignment：同文档 btnOrientation 属性，按钮方向是否横向排列
	// hideAvatar：是否隐藏发送者的头像 false-正常发消息者头像，true-隐藏发消息者头像
	SendActionCardWithMenus(ctx context.Context, title, text string, btns []*ActionCardMenu, alignment, hideAvatar bool) error
	// 发送独立跳转 ActionCard 类型消息
	// title：首屏会话透出的展示内容 text：markdown格式的消息
	// btns：按钮的信息：title-按钮方案，actionURL-点击按钮触发的URL
	// alignment：同文档 btnOrientation 属性，按钮方向是否横向排列
	// hideAvatar：是否隐藏发送者的头像
	SendActionCard(ctx context.Context, title, text string, singleTitle, singleURL string, alignment, hideAvatar bool) error
	// 发送 FeedCard 类型消息
	// title 单条信息文本
	// messageURL 点击单条信息到跳转链接 picURL 单条信息后面图片的URL
	SendFeedCard(ctx context.Context, links ...*FeedCard) error

	// 设置是否静默模式，静默模式下，不发送任何消息
	SetSilenceMode(silence bool) Ding
}
