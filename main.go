package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Counts struct {
	Lines      int
	Characters int
	Files      int
	Functions  int
}

var fileTypes = map[string]*regexp.Regexp{
	".go":   regexp.MustCompile(`\bfunc\b`),
	".py":   regexp.MustCompile(`^\s*def\s+`),
	".js":   regexp.MustCompile(`\bfunction\b`),
	".ts":   regexp.MustCompile(`\bfunction\b`),
	".njs":  regexp.MustCompile(`\bfunction\b`),
	".tfx":  regexp.MustCompile(`\bfunction\b`),
	".itl":  regexp.MustCompile(`\bfunction\b`),
	".json": nil,
	".ini":  nil,
	".key":  nil,
	".md":   nil,
	".toml": nil,
	".mod":  nil,
}

func main() {
	root := "."

	counts := make(map[string]*Counts)
	var mu sync.Mutex

	fileChan := make(chan string, 100)
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for path := range fileChan {
			ext := strings.ToLower(filepath.Ext(path))
			if regex, ok := fileTypes[ext]; ok && regex != nil {
				processFile(path, ext, regex, counts, &mu)
			} else if ok {
				processFile(path, ext, nil, counts, &mu)
			}
		}
	}

	numWorkers := 10
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if _, supported := fileTypes[ext]; supported {
			fileChan <- path
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", root, err)
	}

	close(fileChan)
	wg.Wait()

	// Display the results
	fmt.Printf("%-10s %-10s %-15s %-10s %-10s\n", "LANGUAGE", "FILES", "LINES", "CHARS", "FUNCTIONS")
	fmt.Println(strings.Repeat("-", 60))
	totalCounts := &Counts{}
	for langExt, count := range counts {
		lang := strings.ToUpper(langExt[1:])
		fmt.Printf("%-10s %-10d %-15d %-10d %-10d\n", lang, count.Files, count.Lines, count.Characters, count.Functions)
		totalCounts.Files += count.Files
		totalCounts.Lines += count.Lines
		totalCounts.Characters += count.Characters
		totalCounts.Functions += count.Functions
	}
	fmt.Println(strings.Repeat("-", 60))
	fmt.Printf("%-10s %-10d %-15d %-10d %-10d\n", "TOTAL", totalCounts.Files, totalCounts.Lines, totalCounts.Characters, totalCounts.Functions)
	fmt.Println("\nhttps://github.com/XanaOG/Walker")
	fmt.Println("Please do not claim that you made this tool, or attempt to sell it.")
}

func processFile(path, ext string, regex *regexp.Regexp, counts map[string]*Counts, mu *sync.Mutex) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Failed to open file %s: %v\n", path, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	charCount := 0
	funcCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++
		charCount += len(line) + 1 // +1 for the newline character
		if regex != nil && regex.MatchString(line) {
			funcCount++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	if _, exists := counts[ext]; !exists {
		counts[ext] = &Counts{}
	}
	counts[ext].Files++
	counts[ext].Lines += lineCount
	counts[ext].Characters += charCount
	counts[ext].Functions += funcCount
}
