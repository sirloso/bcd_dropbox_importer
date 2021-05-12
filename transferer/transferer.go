package transferer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/karrick/godirwalk"
)

// map[users][orders][rolls][photos]
type Files map[string]map[string]map[string][]string
type Orders map[string]map[string][]string
type Rolls map[string][]string

var fromPath string
var toPath string

// Transfer : entry point to transfer the files
func Transfer(from string, to string, verbose bool, rename bool) {
	fmt.Printf("starting transfer from %s to %s\n", from, to)

	pathFrom, err := filepath.Abs(from)
	handleErr(err)
	fromPath = pathFrom

	pathTo, err := filepath.Abs(to)
	handleErr(err)
	toPath = pathTo

	fromFolders := make(Files)
	toFolders := make(Files)

	// process source
	err = godirwalk.Walk(pathFrom, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			appendToList(from, &fromFolders, osPathname)
			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	handleErr(err)

	// process from
	err = godirwalk.Walk(pathTo, &godirwalk.Options{
		Callback: func(osPathname string, de *godirwalk.Dirent) error {
			appendToList(to, &toFolders, osPathname)
			return nil
		},
		Unsorted: true, // (optional) set true for faster yet non-deterministic enumeration (see godoc)
	})
	handleErr(err)

	// create go chanel to communicate with ui
	var wg sync.WaitGroup
	// var chans []chan float32

	users := make(map[string]chan float32)

	for user := range fromFolders {
		userChan := make(chan float32)
		wg.Add(1)
		// chans = append(chans, userChan)
		users[user] = userChan
	}
	fmt.Println("starin", fromFolders)
	// send info to ui
	// ui(users, verbose, rename)

	// process data bb
	for user := range fromFolders {
		go process(&wg, users[user], user, fromFolders, toFolders, rename)
	}

	wg.Wait()
	fmt.Println("done")
}

func appendToList(topLevel string, list *Files, item string) {
	files := strings.Split(item, "/")
	fl := len(files)

	if strings.Contains(topLevel, files[fl-2]) {
		//add user
		(*list)[files[fl-1]] = make(Orders)
	}

	if strings.Contains(topLevel, files[fl-3]) {
		// add orders
		(*list)[files[fl-2]][files[fl-1]] = make(Rolls)
	}

	if strings.Contains(topLevel, files[fl-4]) {
		// add rolls
		(*list)[files[fl-3]][files[fl-2]][files[fl-1]] = make([]string, 0)
	}

	if strings.Contains(topLevel, files[fl-5]) {
		// add photo
		photos := (*list)[files[fl-4]][files[fl-3]][files[fl-2]]
		photos = append(photos, files[fl-1])
		(*list)[files[fl-4]][files[fl-3]][files[fl-2]] = photos
	}
}

func markAsUsed(path string) {
	fileArr := strings.Split(path, "/")
	oldFileName := fileArr[len(fileArr)-1]

	fileArr[len(fileArr)-1] = fmt.Sprintf("%s_%s", "processed", oldFileName)
	newFile := strings.Join(fileArr, "/")

	err := os.Rename(oldFileName, newFile)
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
		// panic("path error")
	}
}

func process(wg *sync.WaitGroup, c chan float32, user string, from Files, to Files, rename bool) {
	fmt.Println("processing", user)
	defer wg.Done()

	rolls := float32(0.0)
	processed := float32(0.0)

	for order := range from[user] {
		rolls += float32(len(from[user][order]))
	}

	// if not present create folder in destination
	userPath := fmt.Sprintf("%s/%s", toPath, user)
	_, err := os.Stat(userPath)
	if os.IsNotExist(err) {
		// create folder in destination
		createFolder(userPath)
	}
	fmt.Println("processing user", user)

	// itterate over orders
	for order := range from[user] {
		orderPath := fmt.Sprintf("%s/%s/", userPath, order)
		_, err := os.Stat(orderPath)
		// if order not present in destination; create order
		fmt.Println("processing order", order)
		if os.IsNotExist(err) {
			createFolder(orderPath)
		}

		// itterate over rolls
		for roll := range from[user][order] {
			// if roll not present in destination
			rollPath := fmt.Sprintf("%s/%s/", orderPath, roll)
			_, err := os.Stat(rollPath)
			fmt.Println("processing roll", roll)
			if os.IsNotExist(err) {
				// create roll
				createFolder(rollPath)
				// copy all photos to roll
				files := from[user][order][roll]
				for _, file := range files {
					toFilePath := fmt.Sprintf("%s/%s/%s/%s", user, order, roll, file)
					fromFilePath := fmt.Sprintf("%s/%s", rollPath, file)

					in, _ := os.Open(fromFilePath)
					out, _ := os.Create(toFilePath)

					_, err = io.Copy(out, in)

					handleErr(err)
					in.Close()
					out.Close()
				}

				if rename {
					markAsUsed(rollPath)
				}
				processed += 1.0
				// c <- (processed / total)
			} else {
				// else skip
				if rename {
					markAsUsed(rollPath)
				}
				processed += 1.0
				// c <- (processed / total)
				break
			}
		}

		if rename {
			markAsUsed(orderPath)
		}
	}

	if rename {
		markAsUsed(userPath)
	}
}

func ui(users map[string]chan float32, verbose bool, rename bool) {
	// ui
}

func createFolder(path string) {
	os.Mkdir(path, 0755)
}
