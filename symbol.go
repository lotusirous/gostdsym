package gostdsym

import (
	"fmt"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"io/fs"
	"regexp"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"
)

var internalPkg = regexp.MustCompile(`(^|/)internal($|/)`)

func isSkipPackage(v string) bool {
	return internalPkg.MatchString(v) || strings.HasPrefix(v, "vendor/")
}

// LoadPackages returns a list of packages.
func LoadPackages(pattern string) ([]string, error) {
	pkgs, err := packages.Load(nil, pattern)
	if err != nil {
		return nil, err
	}
	out := make([]string, len(pkgs))
	for i := 0; i < len(pkgs); i++ {
		path := pkgs[i].PkgPath
		if isSkipPackage(path) {
			continue
		}
		out[i] = pkgs[i].PkgPath
	}
	return out, nil
}

// GetPackageSymbols extracts all exported symbols from a package.
func GetPackageSymbols(name, srcDir string) ([]string, error) {
	buildPkg, err := build.Import(name, srcDir, build.ImportComment)
	if err != nil {
		return nil, err
	}

	syms, err := buildSymbols(buildPkg)
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

func buildSymbols(pkg *build.Package) ([]string, error) {
	fset := token.NewFileSet()
	include := func(info fs.FileInfo) bool {
		for _, name := range pkg.GoFiles {
			if name == info.Name() {
				return true
			}
		}
		return false
	}
	pkgs, err := parser.ParseDir(fset, pkg.Dir, include, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	astPkg, ok := pkgs[pkg.Name]
	if !ok {
		return nil, fmt.Errorf("not found package name: %s", pkg.Name)
	}
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
