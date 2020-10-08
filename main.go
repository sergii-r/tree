package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type FileTreeLeaf struct {
	fileName  string
	fileSize  string
	isDir     bool
	isLast    bool
	lastLevel int
	level     int
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	scanDir(out, path, printFiles, 0, 0)
	return nil
}

func scanDir(out io.Writer, path string, printFiles bool, level int, lastLavel int) error {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	FileTree := make([]FileTreeLeaf, 0)

	files, _ := file.Readdir(0)

	for _, f := range files {
		if !f.IsDir() && !printFiles {
			continue
		}

		var leaf = FileTreeLeaf{
			fileName:  f.Name(),
			fileSize:  formatSize(f),
			isDir:     f.IsDir(),
			isLast:    false,
			lastLevel: lastLavel,
			level:     level,
		}

		FileTree = append(FileTree, leaf)
	}

	sort.SliceStable(FileTree, func(i, j int) bool {
		return FileTree[i].fileName < FileTree[j].fileName
	})

	if len(FileTree) > 0 {
		FileTree[len(FileTree)-1].isLast = true
	}

	for _, fileLeaf := range FileTree {
		var repeatSym = strings.Repeat("│\t", level)
		if fileLeaf.lastLevel > 0 {
			repeatSym = strings.Repeat("│\t", fileLeaf.level-fileLeaf.lastLevel) + strings.Repeat("\t", fileLeaf.lastLevel)
		}
		var sym = "├───"
		if fileLeaf.isLast {
			sym = "└───"
		}
		fmt.Fprintln(out, repeatSym+sym+fileLeaf.fileName+fileLeaf.fileSize)
		if fileLeaf.isDir {
			lastLavel := fileLeaf.lastLevel
			if fileLeaf.isLast {
				lastLavel = fileLeaf.lastLevel + 1
			}
			scanDir(out, path+"/"+fileLeaf.fileName, printFiles, level+1, lastLavel)
		}
	}

	return nil
}

func formatSize(file os.FileInfo) string {
	if file.IsDir() {
		return ""
	}
	var result = " (empty)"
	size := file.Size()
	if size > 0 {
		result = " (" + strconv.FormatInt(size, 10) + "b)"
	}
	return result
}
