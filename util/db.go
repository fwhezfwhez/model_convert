package util

import (
	"strconv"
	"strings"
)

// Input counts per page, page index and all records count returns db limit ,offset
// Example:
// input "20", "1", 1000,
// returns limit=20, offset=0
func ToLimitOffset(sizeIn string, indexIn string, count int) (limit int, offset int) {
	size, _ := strconv.Atoi(sizeIn)
	index, _ := strconv.Atoi(indexIn)
	if size == 0 {
		size = 10
	}

	if index == 0 {
		index = 1
	}
	//1
	if count == 0 {
		return size, 0
	}
	var pageMax int
	//1%10
	if count%size == 0 {
		pageMax = count / size
	} else {
		//1
		pageMax = count/size + 1
	}
	//1<=1
	if pageMax <= index {
		index = pageMax
	}
	offset = size * (index - 1)

	if offset == -10 {
		offset = 0
	}
	return size, offset
}

// GenerateOrderBy
// input "id" returns "id asc"
// input "-id" returns "id desc"
func GenerateOrderBy(orderBy string) string {
	var tmp = make([]string, 0, 2)
	orders := strings.Split(orderBy, ",")
	for _, order := range orders {
		order = strings.TrimSpace(order)
		if strings.HasPrefix(order, "-") {
			tmp = append(tmp, order[1:]+" DESC")
		} else {
			tmp = append(tmp, order[:]+" ASC")
		}
	}
	return strings.Join(tmp, ",") + " "
}
