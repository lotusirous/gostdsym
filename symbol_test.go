package gostdsym

import (
	"reflect"
	"sort"
	"testing"
)

func TestAll(t *testing.T) {
	_, err := LoadPackages("std")
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range []struct {
		in   string
		deli string
		want []string
	}{
		{in: "cmp", deli: ".", want: []string{"cmp", "cmp.Compare", "cmp.Less", "cmp.Or", "cmp.Ordered"}},
		{in: "html/template",
			deli: ".",
			want: []string{
				"html/template",
				"html/template.CSS",
				"html/template.ErrAmbigContext",
				"html/template.ErrBadHTML",
				"html/template.ErrBranchEnd",
				"html/template.ErrEndContext",
				"html/template.ErrJSTemplate",
				"html/template.ErrNoSuchTemplate",
				"html/template.ErrOutputContext",
				"html/template.ErrPartialCharset",
				"html/template.ErrPartialEscape",
				"html/template.ErrPredefinedEscaper",
				"html/template.ErrRangeLoopReentry",
				"html/template.ErrSlashAmbig",
				"html/template.Error",
				"html/template.Error.Error",
				"html/template.ErrorCode",
				"html/template.FuncMap",
				"html/template.HTML",
				"html/template.HTMLAttr",
				"html/template.HTMLEscape",
				"html/template.HTMLEscapeString",
				"html/template.HTMLEscaper",
				"html/template.IsTrue",
				"html/template.JS",
				"html/template.JSEscape",
				"html/template.JSEscapeString",
				"html/template.JSEscaper",
				"html/template.JSStr",
				"html/template.Must",
				"html/template.New",
				"html/template.OK",
				"html/template.ParseFS",
				"html/template.ParseFiles",
				"html/template.ParseGlob",
				"html/template.Srcset",
				"html/template.Template",
				"html/template.Template.AddParseTree",
				"html/template.Template.Clone",
				"html/template.Template.DefinedTemplates",
				"html/template.Template.Delims",
				"html/template.Template.Execute",
				"html/template.Template.ExecuteTemplate",
				"html/template.Template.Funcs",
				"html/template.Template.Lookup",
				"html/template.Template.Name",
				"html/template.Template.New",
				"html/template.Template.Option",
				"html/template.Template.Parse",
				"html/template.Template.ParseFS",
				"html/template.Template.ParseFiles",
				"html/template.Template.ParseGlob",
				"html/template.Template.Templates",
				"html/template.URL",
				"html/template.URLQueryEscaper",
				"html/template.eatAttrName",
				"html/template.errorf",
				"html/template.parseFS",
				"html/template.parseFiles",
				"html/template.parseGlob",
			}},
		{
			in:   "container/list",
			deli: ".",
			want: []string{
				"container/list",
				"container/list.Element",
				"container/list.Element.Next",
				"container/list.Element.Prev",
				"container/list.List",
				"container/list.New",
				"container/list.List.Back",
				"container/list.List.Front",
				"container/list.List.Init",
				"container/list.List.InsertAfter",
				"container/list.List.InsertBefore",
				"container/list.List.Len",
				"container/list.List.MoveAfter",
				"container/list.List.MoveBefore",
				"container/list.List.MoveToBack",
				"container/list.List.MoveToFront",
				"container/list.List.PushBack",
				"container/list.List.PushBackList",
				"container/list.List.PushFront",
				"container/list.List.PushFrontList",
				"container/list.List.Remove",
			},
		},
		{
			in:   "container/list",
			deli: "#",
			want: []string{
				"container/list",
				"container/list#Element",
				"container/list#Element.Next",
				"container/list#Element.Prev",
				"container/list#List",
				"container/list#New",
				"container/list#List.Back",
				"container/list#List.Front",
				"container/list#List.Init",
				"container/list#List.InsertAfter",
				"container/list#List.InsertBefore",
				"container/list#List.Len",
				"container/list#List.MoveAfter",
				"container/list#List.MoveBefore",
				"container/list#List.MoveToBack",
				"container/list#List.MoveToFront",
				"container/list#List.PushBack",
				"container/list#List.PushBackList",
				"container/list#List.PushFront",
				"container/list#List.PushFrontList",
				"container/list#List.Remove",
			},
		},
		{
			in:   "context",
			deli: ".",
			want: []string{
				"context",
				"context.AfterFunc",
				"context.Background",
				"context.CancelCauseFunc",
				"context.CancelFunc",
				"context.Canceled",
				"context.Cause",
				"context.Context",
				"context.DeadlineExceeded",
				"context.TODO",
				"context.WithCancel",
				"context.WithCancelCause",
				"context.WithDeadline",
				"context.WithDeadlineCause",
				"context.WithTimeout",
				"context.WithTimeoutCause",
				"context.WithValue",
				"context.WithoutCancel",
			},
		},
		{
			in:   "errors",
			deli: ".",
			want: []string{
				"errors",
				"errors.Is",
				"errors.Join",
				"errors.New",
				"errors.Unwrap",
				"errors.ErrUnsupported",
				"errors.As",
			},
		},
	} {
		got, err := GetPackageSymbols(test.in, test.deli)
		if err != nil {
			t.Fatalf("want no error for MustExtract, got: %v", err)
		}
		sort.Strings(got)
		sort.Strings(test.want)
		if len(got) != len(test.want) || !reflect.DeepEqual(got, test.want) {
			t.Errorf("extractSymbol(%s) = %v, want = %v", test.in, got, test.want)
		}
	}
}
