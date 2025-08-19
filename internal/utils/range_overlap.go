package utils

type Range struct {
	Start int64
	End   int64
}

func IsOverlapping(range1, range2 Range) bool {
	return !(range1.End < range2.Start || range2.End < range1.Start)
}
