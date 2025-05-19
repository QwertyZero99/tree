package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type config struct {
	Indent     string
	DirSyntax  string
	FileSyntax string
}

func main() {
	// Flags
	dir := flag.String("dir", ".", "directory to list")
	indentChar := flag.String("indent", "\t", "what string to indent with")
	dirSyntax := flag.String("dirSyntax", "$/", "syntax to put directories in, replaces $ with the dir name")
	fileSyntax := flag.String("FileSyntax", "$", "syntax to put files in, replaces $ with the file name")

	flag.Parse()

	c := config{
		Indent:     *indentChar,
		DirSyntax:  *dirSyntax,
		FileSyntax: *fileSyntax,
	}

	// Run the directory listing
	result := strDir(*dir, 0, c)
	fmt.Print(result)
	fmt.Println("Press Enter to exit...")
	_, _ = fmt.Scanln() // Waits for the user to press Enter
}

// strDir builds a string representation of a directory's contents recursively
func strDir(dir string, depth int, c config) string {
	entries := readEntries(dir)
	return strEntries(entries, dir, depth, c)
}

// strEntries formats directory entries with indentation
func strEntries(entries []os.DirEntry, basePath string, depth int, c config) string {
	var sb strings.Builder
	prefix := strings.Repeat(c.Indent, depth)

	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(basePath, name)

		if entry.IsDir() {
			sb.WriteString(prefix + strings.ReplaceAll(c.DirSyntax, "$", name) + "\n")
			subEntries := readEntries(fullPath)
			sb.WriteString(strEntries(subEntries, fullPath, depth+1, c))
		} else {
			sb.WriteString(prefix + strings.ReplaceAll(c.FileSyntax, "$", name) + "\n")
		}
	}
	return sb.String()
}

// readEntries safely reads directory entries
func readEntries(dir string) []os.DirEntry {
	entries, err := os.ReadDir(dir)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "unable to read directory %s: %v\n", dir, err)
		if err != nil {
			return nil
		}
		return nil
	}
	return entries
}
