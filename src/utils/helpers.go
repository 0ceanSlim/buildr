package utils

import (
	"time"
)

// Helper function to prepend a directory path to a list of filenames
func PrependDir(dir string, files []string) []string {
	var fullPaths []string
	for _, file := range files {
		fullPaths = append(fullPaths, dir+file)
	}
	return fullPaths
}

func formatTimestamp(unixTime int64) string {
    t := time.Unix(unixTime, 0)
    return t.Format("Jan 02, 2006 15:04")
}
