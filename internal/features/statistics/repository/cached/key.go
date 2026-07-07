package statistics_cached_repository

import (
	"strconv"
	"time"
)

func cacheKey(userId *int, from, to *time.Time) string {
	return "statistics:tasks:" + intPtrToken(userId) + ":" + timePtrToken(from) + ":" + timePtrToken(to)
}

func intPtrToken(v *int) string {
	if v == nil {
		return "nil"
	}

	return strconv.Itoa(*v)
}

func timePtrToken(v *time.Time) string {
	if v == nil {
		return "nil"
	}

	return v.Format(time.RFC3339)
}
