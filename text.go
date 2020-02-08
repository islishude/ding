package ding

import "context"

func (d *clientimpl) SendText(ctx context.Context, content string, isAtAll bool, atMobiles ...string) error {
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
	}

	if isAtAll {
		data["at"] = map[string]interface{}{
			"isAtAll": true,
		}
	} else if len(atMobiles) > 0 {
		data["at"] = map[string]interface{}{
			"atMobiles": atMobiles,
		}
	}

	return d.request(ctx, data)
}
