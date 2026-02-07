package utils

import (
	"strconv"
	"strings"
)

func FormatPrice(p float64) string {
	return strconv.FormatFloat(p, 'f', 2, 64)
}

func Slugify(name string) string {
	s := strings.ToLower(name)
	s = strings.ReplaceAll(s, " ", "_")
	return s
}
