package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	p := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, p, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	names, err := f.Readdir(0)
	if err != nil {
		return err
	}

	if !printFiles {
		names = onlyDirs(names)
	}

	writeFiles(out, printFiles, names, path, "")

	return nil
}

func onlyDirs(names []os.FileInfo) (nnames []os.FileInfo) {
	for _, v := range names {
		if v.IsDir() {
			nnames = append(nnames, v)
		}
	}
	return
}

func writeFiles(out io.Writer, printFiles bool, names []os.FileInfo, p, space string) {
	for i, v := range names {
		switch {
		case i == len(names)-1:
			fmt.Fprintf(out, space+"└───")
		default:
			fmt.Fprintf(out, space+"├───")
		}

		fmt.Fprintf(out, "%s", v.Name())
		if !v.IsDir() {
			if v.Size() == 0 {
				fmt.Fprintf(out, " (empty)")
			} else {
				fmt.Fprintf(out, " (%db)", v.Size())
			}
		}

		fmt.Fprintf(out, "\n")

		if v.IsDir() {
			p := path.Join(p, v.Name())
			f, err := os.Open(p)
			if err != nil {
				log.Println(err)
			}
			n, err := f.Readdir(0)
			if err != nil {
				log.Println(err)
			}

			if !printFiles {
				n = onlyDirs(n)
			}

			if i == len(names)-1 {
				writeFiles(out, printFiles, n, p, space+"\t")
			} else {
				writeFiles(out, printFiles, n, p, space+"│\t")
			}
		}
	}
}
