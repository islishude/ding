package ding

import "context"

func (d *clientimpl) SendMarkdown(ctx context.Context, title, text string, isAtAll bool, atMobiles ...string) error {
	data := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]interface{}{
			"title": title,
			"text":  text,
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
