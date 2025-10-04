package web

import (
	"net/http"
	"strings"
)

func detectContentType(b []byte) string {
	return strings.SplitN(http.DetectContentType(b), ";", 2)[0]
}
