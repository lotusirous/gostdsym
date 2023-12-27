package gostdsym

import (
	"fmt"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"slices"

	"golang.org/x/tools/go/packages"
)

// LoadPackages returns all packages from a given pattern.
func LoadPackages(pattern string) ([]string, error) {
	pkgs, err := packages.Load(nil, pattern)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(pkgs))
	for i := 0; i < len(pkgs); i++ {
		out[i] = pkgs[i].PkgPath
	}
	return out, nil
}

// GetPackageSymbols extracts all exported symbols from a package.
func GetPackageSymbols(pattern string) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	buildPkg, err := build.Import(pattern, wd, build.ImportComment)
	if err != nil {
		return nil, err
	}

	syms, err := parsePackage(buildPkg)
	if err != nil {
		return nil, err
	}

	syms = slices.Compact(syms)
	for i := range syms {
		syms[i] = buildPkg.ImportPath + "." + syms[i]
	}
	syms = append(syms, buildPkg.ImportPath)
	return syms, nil
}

func parsePackage(pkg *build.Package) ([]string, error) {
	// include tells parser.ParseDir which files to include.
	// That means the file must be in the build package's GoFiles or CgoFiles
	// list only (no tag-ignored files, tests, swig or other non-Go files).
	include := func(info fs.FileInfo) bool {
		for _, name := range pkg.GoFiles {
			if name == info.Name() {
				return true
			}
		}
		return false
	}
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, pkg.Dir, include, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no source-code package in directory %s", pkg.Dir)
	}
	astPkg := pkgs[pkg.Name]

	// TODO: go/doc does not include typed constants in the constants
	// list, which is what we want. For instance, time.Sunday is of type
	// time.Weekday, so it is defined in the type but not in the
	// Consts list for the package. This prevents
	//	go doc time.Sunday
	// from finding the symbol. Work around this for now, but we
	// should fix it in go/doc.
	// A similar story applies to factory functions.
	docPkg := doc.New(astPkg, pkg.ImportPath, doc.AllDecls)

	typs := types(docPkg)
	vars := variables(docPkg.Vars)
	consts := constants(docPkg.Consts)
	fns := functions(docPkg.Funcs)
	return append(append(append(typs, vars...), consts...), fns...), nil
}

func types(docPkg *doc.Package) []string {
	var out []string
	for _, typ := range docPkg.Types {
		if !token.IsExported(typ.Name) {
			continue
		}
		out = append(out, typ.Name)
		for _, va := range typ.Vars {
			for _, name := range va.Names {
				if name == "_" {
					continue
				}
				out = append(out, name)
			}
		}
		for _, va := range typ.Consts {
			for _, name := range va.Names {
				if name == "_" {
					continue
				}
				out = append(out, name)
			}
		}
		for _, fn := range typ.Funcs {
			out = append(out, fn.Name)
		}
		for _, method := range typ.Methods {
			if !token.IsExported(method.Name) {
				continue
			}
			out = append(out, typ.Name+"."+method.Name)
		}
	}
	return out
}

func variables(vars []*doc.Value) []string {
	var out []string
	for _, va := range vars {
		for _, name := range va.Names {
			if name == "_" {
				continue
			}
			if token.IsExported(name) {
				out = append(out, name)
			}
		}
	}
	return out
}

func constants(consts []*doc.Value) []string {
	var out []string
	for _, va := range consts {
		for _, name := range va.Names {
			if name == "_" {
				continue
			}
			if token.IsExported(name) {
				out = append(out, name)
			}
		}
	}
	return out
}

func functions(funcs []*doc.Func) []string {
	var out []string
	for _, fn := range funcs {
		if token.IsExported(fn.Name) {
			out = append(out, fn.Name)
		}
	}
	return out
}
