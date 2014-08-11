package http_crawler

import (
	"net/http"
)

type CrawlerResponse struct {
	Body   []byte
	Header http.Header
}

func (r *CrawlerResponse) IsHTML() bool {
	return http.DetectContentType(r.Body) == "text/html; charset=utf-8"
}
