package types

import (
	"fmt"
	"strings"
)

var directions = [2]string{"asc", "desc"}

type Filter struct {
	Prefix string
	Params map[string]string
	Sort   string
}

func NewFilter(prefix, dftSort, query string, avlSort []string) Filter {
	sort := sortQuery(prefix, query, avlSort)

	if sort == "" {
		sort = dftSort
	}

	return Filter{
		Prefix: prefix,
		Params: make(map[string]string),
		Sort:   sort,
	}
}

func sortQuery(pfx, srt string, avl []string) string {
	var builder strings.Builder
	fields := strings.Split(srt, ",")

	for _, field := range fields {
		sortChunks := strings.Split(field, ":")

		if len(sortChunks) == 2 {
			col := ""
			direction := ""

			for _, d := range directions {
				if sortChunks[1] == d {
					direction = d
					break
				}
			}

			for _, c := range avl {
				if sortChunks[0] == c {
					col = c
					break
				}
			}

			if col != "" && direction != "" {
				if pfx == "" {
					builder.WriteString(fmt.Sprintf(",%s %s", col, direction))
				} else {
					builder.WriteString(fmt.Sprintf(",%s.%s %s", pfx, col, direction))
				}
			}
		}
	}

	result := builder.String()

	if result != "" {
		result = result[1:]
	}

	return result
}
