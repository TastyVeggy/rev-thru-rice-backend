package utils

import (
	"fmt"
	"strings"
)

func BuildFiltersforQuery(filters map[string]any) (string, []any) {
	conditions := []string{}
	args := []any{}

	placeholderCounter := 1

	for key, value := range filters {
		switch v := value.(type) {
		case []int:
			placeholders := []string{}
			for i := range v {
				placeholders = append(placeholders, fmt.Sprintf("$%d", placeholderCounter))
				args = append(args, v[i])
				placeholderCounter++
			}
			conditions = append(conditions, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ", ")))
		case string:
			conditions = append(conditions, fmt.Sprintf("%s=%d", key, placeholderCounter))
			args = append(args, v)
			placeholderCounter++
		case int:
			conditions = append(conditions, fmt.Sprintf("%s=%d", key, placeholderCounter))
			args = append(args, v)
			placeholderCounter++
		}
	}

	var query = ""
	if len(conditions) > 0 {
		query = fmt.Sprintf("WHERE %s", strings.Join(conditions, " AND "))
		return query, args
	}
	return query, args
}
