package ding

import "context"

// ActionCardMenu 类型
type ActionCardMenu struct {
	Title     string `json:"title"`     // 按钮菜单文字
	ActionURL string `json:"actionURL"` // 点击按钮触发的URL
}

// NewActionCardMenu creates action card data
func NewActionCardMenu(title, actionURL string) *ActionCardMenu {
	return &ActionCardMenu{
		Title:     title,
		ActionURL: actionURL,
	}
}

func (d *clientimpl) SendActionCardWithMenus(
	ctx context.Context,
	title, text string,
	btns []*ActionCardMenu,
	alignment, hideAvatar bool,
) error {
	var hideAvatarVar = "0"
	if hideAvatar {
		hideAvatarVar = "1"
	}

	var btnOrientationVar = "0"
	if alignment {
		btnOrientationVar = "1"
	}

	data := map[string]interface{}{
		"msgtype": "actionCard",
		"actionCard": map[string]interface{}{
			"title":          title,
			"text":           text,
			"hideAvatar":     hideAvatarVar,
			"btnOrientation": btnOrientationVar,
			"btns":           btns,
		},
	}
	return d.request(ctx, data)
}

func (d *clientimpl) SendActionCard(
	ctx context.Context,
	title, text string,
	singleTitle, singleURL string,
	alignment, hideAvatar bool,
) error {
	var hideAvatarVar = "0"
	if hideAvatar {
		hideAvatarVar = "1"
	}

	var btnOrientationVar = "0"
	if alignment {
		btnOrientationVar = "1"
	}

	data := map[string]interface{}{
		"msgtype": "actionCard",
		"actionCard": map[string]interface{}{
			"title":          title,
			"text":           text,
			"singleTitle":    singleTitle,
			"singleURL":      singleURL,
			"btnOrientation": btnOrientationVar,
			"hideAvatar":     hideAvatarVar,
		},
	}
	return d.request(ctx, data)
}
