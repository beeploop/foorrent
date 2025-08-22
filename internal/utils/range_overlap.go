package utils

type Range struct {
	Start int64
	End   int64
}

func IsOverlapping(r1, r2 Range) bool {
	return !(r1.End <= r2.Start || r2.End <= r1.Start)
}
