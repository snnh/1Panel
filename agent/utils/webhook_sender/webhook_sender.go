package webhook_sender

import (
	"strings"
)

func NormalizeToText(s string) string {
	r := strings.NewReplacer(
		"\r\n", "\n",
		"\r", "\n",
		"<br/>", "\n",
		"<br />", "\n",
		"<br>", "\n",
		"</p>", "\n",
		"<p>", "",
		"</div>", "\n",
		"<div>", "",
		"&nbsp;", " ",
	)
	return r.Replace(s)
}
