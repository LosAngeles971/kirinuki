package storage

import (
	"os"
)

func DeleteLocalFile(filename string) {
	if stat, err := os.Stat(filename); err == nil {
		if !stat.IsDir() {
			os.Remove(filename)
		}
	}
}