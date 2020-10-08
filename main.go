package main

import (
	"fmt"
	"io"
	_ "io"
	"log"
	"os"
	_ "path/filepath"
	"sort"
	"strconv"
	"strings"

)

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

func scanDir(out io.Writer, path string, printFiles bool, level int, isLastDir int) error {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	files, _ := file.Readdir(0)



	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	var isLastElem bool = false

	for num, f := range files {
		if printFiles || f.IsDir() {
			var len = len(files)
			isLastElem = isLast(num, len)
			var sym string = "├───";
			if isLastElem {
				sym = "└───"
			}

			var repeatSym string = strings.Repeat("│\t", level);
			if (isLastDir > 0) {
				repeatSym =  strings.Repeat("│\t", level - isLastDir) + strings.Repeat("\t", isLastDir);
			}

			var size = "";
			if (!f.IsDir()) {
				size = " (" + formatSize(f.Size()) + ")"
			}

			fmt.Fprintln(out, repeatSym + sym + f.Name() + size)
		}

		if (f.IsDir()) {
			if (isLastElem) {
				isLastDir++
			}
			err := scanDir(out, path + "/" +f.Name() , printFiles, level + 1, isLastDir)
			if err != nil {
				panic(err.Error())
			}
		}
	}

	return nil
}

func formatSize(size int64) string {
	var result = "empty"
	if size > 0 {
		result = strconv.FormatInt(size, 10) + "b"
	}
	return result
}

func isLast(num int, len int) bool {
	return num + 1 == len
}

