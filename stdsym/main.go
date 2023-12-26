package main

import (
	"bytes"
	"log"
	"os"

	"github.com/lotusirous/gostdsym"
)

func main() {
	log.SetPrefix("stdsym: ")
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	pkgs, err := gostdsym.LoadPackages("std")
	if err != nil {
		return err
	}

	w := os.Stdout
	var buf bytes.Buffer
	for _, v := range pkgs {
		out, err := gostdsym.GetPackageSymbols(v, cwd)
		if err != nil {
			return err
		}
		for _, sym := range out {
			buf.WriteString(sym + "\n")
		}
	}
	_, err = buf.WriteTo(w)
	return err
}
