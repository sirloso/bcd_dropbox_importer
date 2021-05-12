package transferer

import (
	"path/filepath"

	"github.com/karrick/godirwalk"
)

// map[users][orders][rolls][photos]
type Files map[string]map[string]map[string]string

// Transfer : entry point to transfer the files
func Transfer(from string, to string, verbose bool) {
	pathFrom, err := filepath.Abs(from)
	handleErr(err)
	pathTo, err := filepath.Abs(to)
	handleErr(err)

	var fromFolders Files
	var toFolders Files

	// process source
	err = godirwalk.Walk(pathFrom, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			appendToList(&fromFolders, osPathname)
			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	handleErr(err)

	// process from
	err = godirwalk.Walk(pathTo, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			appendToList(&toFolders, osPathname)
			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	handleErr(err)

	err = startTransfer(&fromFolders, &toFolders)
	handleErr(err)
}

func appendToList(list *Files, item string) {

}

// MarkAsUsed: mark the file/folder/photo as processed
func MarkAsUsed(path string) {
	// rename file/folder/photo
}

func handleErr(err error) {
	if err != nil {
		panic("path error")
	}
}

func startTransfer(fromFolders *Files, toFolder *Files) error {
	return nil
}
