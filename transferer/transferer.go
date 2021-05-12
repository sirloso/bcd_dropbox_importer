package transferer

import (
	"fmt"
	"path/filepath"
)

func Transfer(from string, to string, verbose bool) {
	pathFrom, _ := filepath.Abs(from)
	pathTo, _ := filepath.Abs(to)
	fmt.Println(pathFrom, pathTo, verbose)
}
