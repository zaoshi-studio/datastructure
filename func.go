package datastructure

import "time"

func Align[integer Integer](x, m integer) integer {
	if m <= 0 {
		return x
	}
	return x - x%m
}

func TimeToMs(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
