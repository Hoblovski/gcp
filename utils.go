package main

import (
	"github.com/go-clang/clang-v10/clang"

	"fmt"
	"io"
	"log"
	"os"
)

func get_fileLine(cursor clang.Cursor) FileLine {
	file, line, _, _ := cursor.Location().SpellingLocation()
	fileLine := FileLine{
		file: file.Name(),
		line: int(line),
	}
	return fileLine
}

// TODO: more types than int
func get_ty_identity(ty clang.Type) Identity {
	identity := Identity{
		ModPath: ".",
		PkgPath: ".",
		Name:    "int",
	}
	return identity
}

func get_func_content(cursor clang.Cursor) string {
	if cursor.IsNull() || cursor.Kind() != clang.Cursor_FunctionDecl {
		log.Fatalf("get_func_content: not a function")
		return ""
	}

	fnExtent := cursor.Extent()
	fnFile, _, _, startOffset := fnExtent.Start().SpellingLocation()
	_, _, _, endOffset := fnExtent.End().SpellingLocation()
	res, err := read_file(fnFile.Name(), int64(startOffset), int64(endOffset))
	if err != nil {
		log.Fatalf("get_func_content: %s", err)
		return ""
	}
	return res
}

func read_file(fpath string, begin_offset int64, end_offset int64) (string, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()
	if begin_offset < 0 || end_offset < 0 || begin_offset > end_offset || end_offset > fileSize {
		return "", fmt.Errorf("invalid offset range: begin_offset=%d, end_offset=%d, file_size=%d", begin_offset, end_offset, fileSize)
	}
	readSize := end_offset - begin_offset
	_, err = file.Seek(begin_offset, io.SeekStart)
	if err != nil {
		return "", fmt.Errorf("failed to seek to begin_offset: %w", err)
	}
	buffer := make([]byte, readSize)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(buffer[:n]), nil
}
