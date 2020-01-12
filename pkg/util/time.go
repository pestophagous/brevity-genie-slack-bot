package util

import (
	"strconv"
	"time"
)

func MessageTimestampToUTC(ts string) time.Time {
	i, _ := strconv.ParseInt(ts[:10], 10, 64)
	return time.Unix(i, 0)
}
