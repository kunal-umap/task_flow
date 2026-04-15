package utils

import (
	"net/http"
	"strconv"
)

func GetPagination(r *http.Request) (limit int, offset int) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, _ := strconv.Atoi(pageStr)
	limit, _ = strconv.Atoi(limitStr)

	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset = (page - 1) * limit

	return
}
