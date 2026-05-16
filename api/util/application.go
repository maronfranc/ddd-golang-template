package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/maronfranc/poc-golang-ddd/domain/dto"
)

const QS_PAGE = "page"
const QS_LIMIT = "limit"
const LIMIT = 10

func GetQueryInt(r *http.Request, name string, defaultV int) int {
	qv := r.URL.Query().Get(name)
	if qv == "" {
		return defaultV
	}
	v, err := strconv.Atoi(qv)
	if err != nil {
		v = defaultV
	}
	return v
}

func BuildPagination(route string, total, page, limit int) *dto.Paginated {
	// SEE: https://stackoverflow.com/a/17974
	// SEE: https://stackoverflow.com/questions/17944/how-to-round-up-the-result-of-integer-division/17974
	totalPage := int((total + limit - 1) / limit)
	prevPage := page - 1
	nextPage := page + 1
	var prevLink, nextLink string
	if prevPage > 1 {
		prevLink = fmt.Sprintf("%s?page=%d&limit=%d", route, prevPage, limit)
	}
	if nextPage < totalPage {
		nextLink = fmt.Sprintf("%s?page=%d&limit=%d", route, nextPage, limit)
	}
	return &dto.Paginated{
		TotalRecord: total,
		TotalPage:   totalPage,
		PrevLink:    prevLink,
		NextLink:    nextLink,
	}
}

func EncodeResponseError(w http.ResponseWriter, msg string, code int) {
	m := dto.ResponseMessage{Message: msg}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(m)
}
