package utils

import "fmt"

func BytesToMB(size int) string {
	mb := float64(size) / float64(1024) / float64(1024)
	return fmt.Sprintf("%0.2fMB", mb)
}
