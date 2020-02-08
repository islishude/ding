package ding

import "context"

func (d *clientimpl) SendLink(ctx context.Context, text, title, msgurl string, picurl ...string) error {
	link := map[string]interface{}{
		"title":      title,
		"text":       text,
		"messageUrl": msgurl,
	}

	if len(picurl) > 0 {
		link["picUrl"] = picurl[0]
	}

	data := map[string]interface{}{
		"msgtype": "link",
		"link":    link,
	}
	return d.request(ctx, data)
}
