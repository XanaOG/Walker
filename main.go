package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	root := "." // change this to the root directory of your project, leave as "." if you are running it in the root directory of your project.

	lineCount := make(map[string]int)
	charCount := make(map[string]int)
	fileCount := make(map[string]int)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		switch ext {
		case ".go", ".py", ".js", ".njs", ".tfx", ".itl", ".json", ".ini", ".key", ".md", ".toml": // add other file extensions as desired
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			lines := strings.Split(string(b), "\n")
			lineCount[ext] += len(lines)
			charCount[ext] += len(b)
			fileCount[ext]++
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("%-10s %-15s %-20s %-15s\n", "LANGUAGE", "LINE COUNT", "CHARACTER COUNT", "FILE COUNT")
		fmt.Println(strings.Repeat("-", 65))
		for lang, count := range lineCount {
			fmt.Printf("%-10s %-15d %-20d %-15d\n", strings.ToUpper(lang[1:]), count, charCount[lang], fileCount[lang])
		}
		fmt.Println(strings.Repeat("-", 65))
		fmt.Printf("%-10s %-15s %-20s %-15s\n", "TOTAL", sum(lineCount), sum(charCount), sum(fileCount))
		fmt.Printf("\nhttps://github.com/XanaOG/Walker\n")
       fmt.Printf("Please do not claim that you made the tool.\n")
	}
}

func sum(m map[string]int) string {
	total := 0
	for _, v := range m {
		total += v
	}
	return fmt.Sprintf("%d", total)
}
