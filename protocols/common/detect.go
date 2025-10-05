package common

import (
	"net/http"
	"strings"
)

func DetectContentType(b []byte) string {
	return strings.SplitN(http.DetectContentType(b), ";", 2)[0]
}
