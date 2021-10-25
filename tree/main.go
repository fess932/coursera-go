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

	write(out, names, path, "")

	return nil
}

func write(out io.Writer, names []os.FileInfo, p, space string) {
	for i, v := range names {
		if i == len(names)-1 {
			fmt.Fprintf(out, space+"└───")
		} else {
			fmt.Fprintf(out, space+"├───")
		}

		fmt.Fprintf(out, "%s", v.Name())
		if !v.IsDir() {
			fmt.Fprintf(out, " (%vb)", v.Size())
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

			write(out, n, p, space+"\t")
		}
	}
}
