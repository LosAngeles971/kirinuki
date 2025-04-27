package mosaic

import (
	"time"

	"github.com/LosAngeles971/kirinuki/business/helpers"
)

func BuildFromLocalFile(name, filename string, targets []string) (*File, error) {
	k := NewFile(name, targets)
	return k, nil
}