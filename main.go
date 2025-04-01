package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/go-clang/clang-v10/clang"
)

var fname = flag.String("fname", "", "the file to analyze")

func main() {
	tu, err := parse_file_with_args(os.Args[1:])
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		os.Exit(1)
	}
	defer tu.Dispose()

	tc := NewTranspileContext()
	tc.AddFile(File{Name: *fname}, tu)
	repo := tc.get_repo()

	var jsonStr []byte
	jsonStr, err = json.MarshalIndent(repo, "", "  ")
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		os.Exit(1)
	} else {
		fmt.Printf("GOT REPO:\n%s\n", jsonStr)
	}
}

func parse_file_with_args(args []string) (*clang.TranslationUnit, error) {
	if err := flag.CommandLine.Parse(args); err != nil {
		fmt.Printf("ERROR: %s", err)
		return nil, err
	}

	fmt.Printf(":: fname: %s\n", *fname)
	fmt.Printf(":: args: %v\n", flag.Args())

	if *fname == "" {
		flag.Usage()
		fmt.Printf("please provide a file name to analyze\n")
		return nil, errors.New("no file name provided")
	}

	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()

	tuArgs := []string{}
	if len(flag.Args()) > 0 && flag.Args()[0] == "-" {
		tuArgs = make([]string, len(flag.Args()[1:]))
		copy(tuArgs, flag.Args()[1:])
	}

	tu := idx.ParseTranslationUnit(*fname, tuArgs, nil, 0)
	return &tu, nil
}
