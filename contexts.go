package main

import (
	"fmt"
	"log"

	"github.com/go-clang/clang-v10/clang"
)

type TranspileContext struct {
	files map[string]*FileContext
}

type FileContext struct {
	file  File
	funcs []*Function
}

func NewTranspileContext() *TranspileContext {
	return &TranspileContext{
		files: make(map[string]*FileContext),
	}
}

func (tc *TranspileContext) get_mainPackage() *Package {
	return &Package{
		IsMain:    true,
		PkgPath:   ".",
		Functions: tc.get_functions(),
		Types:     map[string]*Type{},
		Vars:      map[string]*Var{},
	}
}

func (tc *TranspileContext) get_functions() map[string]*Function {
	functions := map[string]*Function{}
	for _, fc := range tc.files {
		for _, f := range fc.funcs {
			functions[f.Identity.Name] = f
		}
	}
	return functions
}

func (tc *TranspileContext) get_repo() *Repository {
	pkg := tc.get_mainPackage()
	module := Module{
		Name: ".",
		Dir:  ".",
		Packages: map[PkgPath]*Package{
			pkg.PkgPath: pkg,
		},
		Dependencies: map[string]string{},
		Files:        tc.get_files(),
	}
	repo := Repository{
		Name: "current_repo",
		Modules: map[string]*Module{
			module.Name: &module,
		},
	}
	return &repo
}

func (tc *TranspileContext) get_files() map[string]*File {
	files := map[string]*File{}
	for _, fc := range tc.files {
		files[fc.file.Name] = &fc.file
	}
	return files
}

func (tc *TranspileContext) AddFile(file File, tu *clang.TranslationUnit) {
	if _, exists := tc.files[file.Name]; !exists {
		tc.files[file.Name] = NewFileContext(file)
		tc.files[file.Name].analyze_tu(tu)
	}
}

func NewFileContext(file File) *FileContext {
	return &FileContext{
		file:  file,
		funcs: []*Function{},
	}
}

func (fc *FileContext) analyze_tu(tu *clang.TranslationUnit) error {
	fmt.Printf("tu: %s\n", tu.Spelling())
	var cursor clang.Cursor
	cursor = tu.TranslationUnitCursor()
	fc.file = File{
		Name:    tu.Spelling(),
		Imports: []string{},
	}

	cursor.Location().SpellingLocation()

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() {
			return clang.ChildVisit_Continue
		}

		fmt.Printf("%s: %s (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR())

		switch cursor.Kind() {
		case clang.Cursor_FunctionDecl:
			f := visit_function(cursor)
			if f != nil {
				fc.funcs = append(fc.funcs, f)
			}
			return clang.ChildVisit_Continue
		}
		return clang.ChildVisit_Continue
	})

	fmt.Printf("parsed.\n")

	return nil
}

func visit_function(cursor clang.Cursor) *Function {
	if cursor.IsNull() || cursor.Kind() != clang.Cursor_FunctionDecl {
		log.Fatalf("visit_function: not a function")
		return nil
	}

	fname := cursor.Spelling()
	identity := Identity{
		ModPath: ".",
		PkgPath: ".",
		Name:    fname,
	}

	params := []Identity{}
	for i := 0; i < int(cursor.NumArguments()); i++ {
		param := cursor.Argument(uint32(i))
		param_ident := Identity{
			ModPath: ".",
			PkgPath: ".",
			Name:    param.Spelling(),
		}
		params = append(params, param_ident)
	}

	resultTy := get_ty_identity(cursor.ResultType())
	results := []Identity{resultTy}
	f := Function{
		Exported:      true,
		IsMethod:      false,
		Identity:      identity,
		FileLine:      get_fileLine(cursor),
		Content:       get_func_content(cursor),
		Receiver:      nil,
		Params:        params,
		Results:       results,
		FunctionCalls: []Identity{},
		MethodCalls:   []Identity{},
		Types:         results, // now only int
		GlobalVars:    []Identity{},
		CompressData:  new(string),
	}

	// if j, err := json.MarshalIndent(f, "", "  "); err == nil {
	// 	fmt.Println("Event JSON:", string(j))
	// }

	return &f
}
