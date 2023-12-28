package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/lotusirous/gostdsym"
)

func main() {
	log.SetPrefix("stdsym: ")
	log.SetFlags(0)
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: stdsym [options]\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func run() error {
	web := flag.Bool("web", false, "enable pkgsite format with the '#' delimiter")
	flag.Usage = usage
	flag.Parse()
	deli := "."
	if *web {
		deli = "#"
	}
	stdPattern := "std"
	pkgs, err := gostdsym.LoadPackages(stdPattern)
	if err != nil {
		return err
	}
	w := os.Stdout
	buf := bufio.NewWriter(w)
	for _, pattern := range pkgs {
		if isSkipPackage(pattern) {
			continue
		}
		out, err := gostdsym.GetPackageSymbols(pattern, deli)
		if err != nil {
			return err
		}
		for _, sym := range out {
			buf.WriteString(sym + "\n")
		}
	}
	return buf.Flush()
}

var internalPkg = regexp.MustCompile(`(^|/)internal($|/)`)

func isSkipPackage(v string) bool {
	return internalPkg.MatchString(v) || strings.HasPrefix(v, "vendor/") && v != ""
}
