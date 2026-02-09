package helpers

import (
	"net/http"
	"strconv"
)

func GetOffsetLimit(r *http.Request, defaultLimit, maxLimit int) (int, int) {
	q := r.URL.Query()

	offset, _ := strconv.Atoi(q.Get("offset"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	if offset < 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return offset, limit
}
