package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

type FileStats struct {
	Path         string
	Lines        int
	CodeLines    int
	CommentLines int
	BlankLines   int
	Characters   int
	Functions    int
	Classes      int
	Size         int64
}

type LanguageStats struct {
	Files        int
	Lines        int
	CodeLines    int
	CommentLines int
	BlankLines   int
	Characters   int
	Functions    int
	Classes      int
	Size         int64
	FileStats    []FileStats
}

type Config struct {
	Root         string
	OutputFormat string
	ShowProgress bool
	Exclude      []string
	Include      []string
	TopFiles     int
	Detailed     bool
	ByDirectory  bool
}

type LanguageConfig struct {
	Extensions       []string
	FunctionPattern  *regexp.Regexp
	ClassPattern     *regexp.Regexp
	CommentPatterns  []*regexp.Regexp
	StringDelimiters []string
}

var languages = map[string]LanguageConfig{
	"Go": {
		Extensions:      []string{".go"},
		FunctionPattern: regexp.MustCompile(`^\s*func\s+(\w+|\([^)]*\)\s*\w+)\s*\(`),
		ClassPattern:    regexp.MustCompile(`^\s*type\s+\w+\s+(struct|interface)`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"Python": {
		Extensions:      []string{".py", ".pyw", ".pyx"},
		FunctionPattern: regexp.MustCompile(`^\s*def\s+\w+\s*\(`),
		ClassPattern:    regexp.MustCompile(`^\s*class\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
			regexp.MustCompile(`^\s*""".*?"""`),
			regexp.MustCompile(`^\s*'''.*?'''`),
		},
	},
	"JavaScript": {
		Extensions:      []string{".js", ".jsx", ".mjs", ".cjs"},
		FunctionPattern: regexp.MustCompile(`^\s*(function\s+\w+|const\s+\w+\s*=\s*\(|let\s+\w+\s*=\s*\(|var\s+\w+\s*=\s*\(|\w+\s*:\s*function|\w+\s*=>\s*)`),
		ClassPattern:    regexp.MustCompile(`^\s*class\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"TypeScript": {
		Extensions:      []string{".ts", ".tsx"},
		FunctionPattern: regexp.MustCompile(`^\s*(function\s+\w+|const\s+\w+\s*=\s*\(|let\s+\w+\s*=\s*\(|export\s+function|\w+\s*:\s*\(|\w+\s*=>\s*)`),
		ClassPattern:    regexp.MustCompile(`^\s*(export\s+)?(abstract\s+)?class\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"Java": {
		Extensions:      []string{".java"},
		FunctionPattern: regexp.MustCompile(`^\s*(public|private|protected|static|\s)*\s+\w+\s+\w+\s*\(`),
		ClassPattern:    regexp.MustCompile(`^\s*(public|private|protected)?\s*(abstract\s+)?(class|interface)\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"C": {
		Extensions:      []string{".c", ".h"},
		FunctionPattern: regexp.MustCompile(`^\s*\w+\s+\w+\s*\(`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"C++": {
		Extensions:      []string{".cpp", ".cc", ".cxx", ".hpp", ".hxx"},
		FunctionPattern: regexp.MustCompile(`^\s*(\w+\s+)*\w+\s+\w+\s*\(`),
		ClassPattern:    regexp.MustCompile(`^\s*(class|struct)\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"C#": {
		Extensions:      []string{".cs"},
		FunctionPattern: regexp.MustCompile(`^\s*(public|private|protected|internal|static|\s)*\s+\w+\s+\w+\s*\(`),
		ClassPattern:    regexp.MustCompile(`^\s*(public|private|protected|internal)?\s*(abstract\s+)?(class|interface|struct)\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"Rust": {
		Extensions:      []string{".rs"},
		FunctionPattern: regexp.MustCompile(`^\s*(pub\s+)?fn\s+\w+`),
		ClassPattern:    regexp.MustCompile(`^\s*(pub\s+)?(struct|enum|trait)\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"PHP": {
		Extensions:      []string{".php", ".phtml"},
		FunctionPattern: regexp.MustCompile(`^\s*(public|private|protected)?\s*function\s+\w+`),
		ClassPattern:    regexp.MustCompile(`^\s*(abstract\s+)?(class|interface|trait)\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`^\s*#`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"Ruby": {
		Extensions:      []string{".rb", ".rbw"},
		FunctionPattern: regexp.MustCompile(`^\s*def\s+\w+`),
		ClassPattern:    regexp.MustCompile(`^\s*class\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"Swift": {
		Extensions:      []string{".swift"},
		FunctionPattern: regexp.MustCompile(`^\s*(private|public|internal)?\s*func\s+\w+`),
		ClassPattern:    regexp.MustCompile(`^\s*(public|private|internal)?\s*(class|struct|protocol)\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"Kotlin": {
		Extensions:      []string{".kt", ".kts"},
		FunctionPattern: regexp.MustCompile(`^\s*(private|public|internal|protected)?\s*fun\s+\w+`),
		ClassPattern:    regexp.MustCompile(`^\s*(public|private|internal|protected)?\s*(class|interface|object)\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"Shell": {
		Extensions:      []string{".sh", ".bash", ".zsh", ".fish"},
		FunctionPattern: regexp.MustCompile(`^\s*\w+\s*\(\s*\)\s*\{`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"HTML": {
		Extensions: []string{".html", ".htm", ".xhtml"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`<!--.*?-->`),
		},
	},
	"CSS": {
		Extensions: []string{".css", ".scss", ".sass", ".less"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`/\*.*?\*/`),
			regexp.MustCompile(`^\s*//`), // SCSS/Sass comments
		},
	},
	"SQL": {
		Extensions: []string{".sql"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*--`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"YAML": {
		Extensions: []string{".yml", ".yaml"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"JSON": {
		Extensions: []string{".json"},
	},
	"XML": {
		Extensions: []string{".xml", ".xsd", ".xsl"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`<!--.*?-->`),
		},
	},
	"Markdown": {
		Extensions: []string{".md", ".markdown"},
	},
	"TOML": {
		Extensions: []string{".toml"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"INI": {
		Extensions: []string{".ini", ".cfg", ".conf"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*[#;]`),
		},
	},
	"Dart": {
		Extensions:      []string{".dart"},
		FunctionPattern: regexp.MustCompile(`^\s*(static\s+)?\w+\s+\w+\s*\(`),
		ClassPattern:    regexp.MustCompile(`^\s*(abstract\s+)?class\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"Scala": {
		Extensions:      []string{".scala", ".sc"},
		FunctionPattern: regexp.MustCompile(`^\s*def\s+\w+`),
		ClassPattern:    regexp.MustCompile(`^\s*(class|object|trait)\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"Lua": {
		Extensions:      []string{".lua"},
		FunctionPattern: regexp.MustCompile(`^\s*(local\s+)?function\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*--`),
		},
	},
	"Perl": {
		Extensions:      []string{".pl", ".pm", ".perl"},
		FunctionPattern: regexp.MustCompile(`^\s*sub\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"R": {
		Extensions:      []string{".r", ".R", ".Rmd"},
		FunctionPattern: regexp.MustCompile(`^\s*\w+\s*<-\s*function`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"MATLAB": {
		Extensions:      []string{".m", ".mlx"},
		FunctionPattern: regexp.MustCompile(`^\s*function\s+.*=\s*\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*%`),
		},
	},
	"Julia": {
		Extensions:      []string{".jl"},
		FunctionPattern: regexp.MustCompile(`^\s*function\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"Haskell": {
		Extensions:      []string{".hs", ".lhs"},
		FunctionPattern: regexp.MustCompile(`^\s*\w+\s*::`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*--`),
		},
	},
	"Erlang": {
		Extensions:      []string{".erl", ".hrl"},
		FunctionPattern: regexp.MustCompile(`^\s*\w+\s*\(`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*%`),
		},
	},
	"Elixir": {
		Extensions:      []string{".ex", ".exs"},
		FunctionPattern: regexp.MustCompile(`^\s*def\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"F#": {
		Extensions:      []string{".fs", ".fsx", ".fsi"},
		FunctionPattern: regexp.MustCompile(`^\s*let\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`\(\*.*?\*\)`),
		},
	},
	"OCaml": {
		Extensions:      []string{".ml", ".mli"},
		FunctionPattern: regexp.MustCompile(`^\s*let\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`\(\*.*?\*\)`),
		},
	},
	"Assembly": {
		Extensions: []string{".asm", ".s", ".S"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*[#;]`),
		},
	},
	"Vim": {
		Extensions:      []string{".vim", ".vimrc"},
		FunctionPattern: regexp.MustCompile(`^\s*function!?\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*"`),
		},
	},
	"Batch": {
		Extensions:      []string{".bat", ".cmd"},
		FunctionPattern: regexp.MustCompile(`^\s*:\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*rem\s`),
			regexp.MustCompile(`^\s*::`),
		},
	},
	"PowerShell": {
		Extensions:      []string{".ps1", ".psm1", ".psd1"},
		FunctionPattern: regexp.MustCompile(`^\s*function\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"Dockerfile": {
		Extensions: []string{"Dockerfile", ".dockerfile"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"Terraform": {
		Extensions: []string{".tf", ".tfvars"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
	"GraphQL": {
		Extensions: []string{".graphql", ".gql"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"Protobuf": {
		Extensions: []string{".proto"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
		},
	},
	"CMake": {
		Extensions:      []string{".cmake", "CMakeLists.txt"},
		FunctionPattern: regexp.MustCompile(`^\s*function\s*\(`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"Makefile": {
		Extensions: []string{"Makefile", ".mk", "makefile"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*#`),
		},
	},
	"Properties": {
		Extensions: []string{".properties", ".env"},
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*[#!]`),
		},
	},
	"Groovy": {
		Extensions:      []string{".groovy", ".gradle"},
		FunctionPattern: regexp.MustCompile(`^\s*def\s+\w+`),
		ClassPattern:    regexp.MustCompile(`^\s*class\s+\w+`),
		CommentPatterns: []*regexp.Regexp{
			regexp.MustCompile(`^\s*//`),
			regexp.MustCompile(`/\*.*?\*/`),
		},
	},
}

var (
	defaultExcludes = []string{
		".git", ".svn", ".hg", ".bzr",
		"node_modules", "vendor", "target", "build", "dist",
		".idea", ".vscode", ".vs", "*.exe", "*.dll", "*.so", "*.dylib",
		"*.jar", "*.war", "*.class", "*.pyc", "*.pyo", "__pycache__",
		".DS_Store", "Thumbs.db",
	}
)

func main() {
	config := parseFlags()

	if config.ShowProgress {
		fmt.Println(color.CyanString("Walker - Code Analysis Tool"))
		fmt.Printf("%s\n", color.New(color.FgHiBlack).Sprintf("Analyzing codebase at: %s", config.Root))
		fmt.Println()
	}

	stats, err := analyzeCodebase(config)
	if err != nil {
		fmt.Printf("Error analyzing codebase: %v\n", err)
		os.Exit(1)
	}

	switch config.OutputFormat {
	case "json":
		outputJSON(stats)
	case "table":
		fallthrough
	default:
		outputTable(stats, config)
	}
}

func parseFlags() Config {
	var config Config

	flag.StringVar(&config.Root, "path", ".", "Root directory to analyze")
	flag.StringVar(&config.OutputFormat, "format", "table", "Output format (table, json)")
	flag.BoolVar(&config.ShowProgress, "progress", true, "Show progress bar")
	flag.IntVar(&config.TopFiles, "top", 10, "Show top N files by lines")
	flag.BoolVar(&config.Detailed, "detailed", false, "Show detailed file statistics")
	flag.BoolVar(&config.ByDirectory, "by-dir", false, "Group results by directory")

	var excludeStr, includeStr string
	flag.StringVar(&excludeStr, "exclude", "", "Comma-separated list of patterns to exclude")
	flag.StringVar(&includeStr, "include", "", "Comma-separated list of patterns to include")

	flag.Parse()

	if excludeStr != "" {
		config.Exclude = strings.Split(excludeStr, ",")
	}
	config.Exclude = append(config.Exclude, defaultExcludes...)

	if includeStr != "" {
		config.Include = strings.Split(includeStr, ",")
	}

	return config
}

func analyzeCodebase(config Config) (map[string]*LanguageStats, error) {
	stats := make(map[string]*LanguageStats)
	var mu sync.Mutex

	extToLang := make(map[string]string)
	for lang, langConfig := range languages {
		for _, ext := range langConfig.Extensions {
			extToLang[ext] = lang
		}
	}

	fileChan := make(chan string, 1000)
	var wg sync.WaitGroup

	var totalFiles int
	filepath.Walk(config.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if shouldProcessFile(path, config) {
			totalFiles++
		}
		return nil
	})

	var bar *progressbar.ProgressBar
	if config.ShowProgress && totalFiles > 0 {
		bar = progressbar.NewOptions(totalFiles,
			progressbar.OptionSetDescription("Analyzing files..."),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "█",
				SaucerHead:    "█",
				SaucerPadding: "░",
				BarStart:      "[",
				BarEnd:        "]",
			}),
			progressbar.OptionShowCount(),
			progressbar.OptionShowIts(),
			progressbar.OptionSetWidth(50),
		)
	}

	worker := func() {
		defer wg.Done()
		for path := range fileChan {
			ext := strings.ToLower(filepath.Ext(path))
			if lang, ok := extToLang[ext]; ok {
				fileStats := analyzeFile(path, languages[lang])

				mu.Lock()
				if stats[lang] == nil {
					stats[lang] = &LanguageStats{FileStats: make([]FileStats, 0)}
				}
				langStats := stats[lang]
				langStats.Files++
				langStats.Lines += fileStats.Lines
				langStats.CodeLines += fileStats.CodeLines
				langStats.CommentLines += fileStats.CommentLines
				langStats.BlankLines += fileStats.BlankLines
				langStats.Characters += fileStats.Characters
				langStats.Functions += fileStats.Functions
				langStats.Classes += fileStats.Classes
				langStats.Size += fileStats.Size
				langStats.FileStats = append(langStats.FileStats, fileStats)
				mu.Unlock()
			}

			if bar != nil {
				bar.Add(1)
			}
		}
	}

	numWorkers := 20
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go worker()
	}

	err := filepath.Walk(config.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() || !shouldProcessFile(path, config) {
			return nil
		}
		fileChan <- path
		return nil
	})

	close(fileChan)
	wg.Wait()

	if bar != nil {
		bar.Finish()
		fmt.Println()
	}

	return stats, err
}

func shouldProcessFile(path string, config Config) bool {
	for _, pattern := range config.Exclude {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return false
		}
		if strings.Contains(path, pattern) {
			return false
		}
	}

	if len(config.Include) > 0 {
		for _, pattern := range config.Include {
			if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
				return true
			}
		}
		return false
	}

	ext := strings.ToLower(filepath.Ext(path))
	for _, langConfig := range languages {
		for _, supportedExt := range langConfig.Extensions {
			if ext == supportedExt {
				return true
			}
		}
	}

	return false
}

func analyzeFile(path string, langConfig LanguageConfig) FileStats {
	file, err := os.Open(path)
	if err != nil {
		return FileStats{Path: path}
	}
	defer file.Close()

	info, _ := file.Stat()
	stats := FileStats{
		Path: path,
		Size: info.Size(),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stats.Lines++
		stats.Characters += len(line) + 1

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			stats.BlankLines++
		} else if isCommentLine(trimmed, langConfig.CommentPatterns) {
			stats.CommentLines++
		} else {
			stats.CodeLines++
			if langConfig.FunctionPattern != nil && langConfig.FunctionPattern.MatchString(line) {
				stats.Functions++
			}
			if langConfig.ClassPattern != nil && langConfig.ClassPattern.MatchString(line) {
				stats.Classes++
			}
		}
	}

	return stats
}

func isCommentLine(line string, patterns []*regexp.Regexp) bool {
	for _, pattern := range patterns {
		if pattern.MatchString(line) {
			return true
		}
	}
	return false
}

func outputTable(stats map[string]*LanguageStats, config Config) {
	if len(stats) == 0 {
		color.Yellow("No supported code files found!")
		return
	}

	color.Cyan("\nCode Analysis Results")
	fmt.Printf("%s\n", color.New(color.FgHiBlack).Sprintf("Generated on: %s", time.Now().Format("2006-01-02 15:04:05")))
	fmt.Println()

	type langSort struct {
		name  string
		stats *LanguageStats
	}
	var sorted []langSort
	for lang, langStats := range stats {
		sorted = append(sorted, langSort{lang, langStats})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].stats.Lines > sorted[j].stats.Lines
	})

	// Print clean, well-formatted table
	fmt.Printf("%-15s %8s %12s %12s %12s %8s %12s %8s %10s\n",
		"LANGUAGE", "FILES", "LINES", "CODE", "COMMENTS", "BLANK", "CHARS", "FUNCS", "CLASSES")

	fmt.Println(strings.Repeat("─", 120))

	var totals LanguageStats
	for _, item := range sorted {
		lang, langStats := item.name, item.stats

		fmt.Printf("%-15s %8d %12d %12d %12d %8d %12d %8d %10d\n",
			lang,
			langStats.Files,
			langStats.Lines,
			langStats.CodeLines,
			langStats.CommentLines,
			langStats.BlankLines,
			langStats.Characters,
			langStats.Functions,
			langStats.Classes)

		totals.Files += langStats.Files
		totals.Lines += langStats.Lines
		totals.CodeLines += langStats.CodeLines
		totals.CommentLines += langStats.CommentLines
		totals.BlankLines += langStats.BlankLines
		totals.Characters += langStats.Characters
		totals.Functions += langStats.Functions
		totals.Classes += langStats.Classes
		totals.Size += langStats.Size
	}

	fmt.Println(strings.Repeat("─", 120))
	fmt.Printf("%-15s %8d %12d %12d %12d %8d %12d %8d %10d\n",
		"TOTAL",
		totals.Files,
		totals.Lines,
		totals.CodeLines,
		totals.CommentLines,
		totals.BlankLines,
		totals.Characters,
		totals.Functions,
		totals.Classes)

	if config.TopFiles > 0 {
		showTopFiles(stats, config.TopFiles)
	}

	// Show summary
	fmt.Printf("\n Summary:\n")
	fmt.Printf("   Total Size: %s\n", formatBytes(totals.Size))
	fmt.Printf("   Code Ratio: %.1f%%\n", float64(totals.CodeLines)/float64(totals.Lines)*100)
	if totals.Functions > 0 {
		fmt.Printf("   Avg Lines/Function: %.1f\n", float64(totals.CodeLines)/float64(totals.Functions))
	}

	fmt.Printf("\n %s\n", color.BlueString("https://github.com/XanaOG/Walker"))
	fmt.Printf("   %s\n", color.New(color.FgHiBlack).Sprint("Please respect the original author"))
}

func showTopFiles(stats map[string]*LanguageStats, topN int) {
	fmt.Printf("\n Top %d Files by Lines:\n", topN)

	var allFiles []FileStats
	for _, langStats := range stats {
		allFiles = append(allFiles, langStats.FileStats...)
	}

	sort.Slice(allFiles, func(i, j int) bool {
		return allFiles[i].Lines > allFiles[j].Lines
	})

	if len(allFiles) > topN {
		allFiles = allFiles[:topN]
	}

	for i, file := range allFiles {
		fmt.Printf("%2d. %-55s %10d lines %12d chars\n",
			i+1,
			truncateString(file.Path, 55),
			file.Lines,
			file.Characters)
	}
}

func outputJSON(stats map[string]*LanguageStats) {
	output := struct {
		GeneratedAt time.Time                 `json:"generated_at"`
		Languages   map[string]*LanguageStats `json:"languages"`
		Summary     map[string]interface{}    `json:"summary"`
	}{
		GeneratedAt: time.Now(),
		Languages:   stats,
	}

	// Calculate summary
	var totals LanguageStats
	for _, langStats := range stats {
		totals.Files += langStats.Files
		totals.Lines += langStats.Lines
		totals.CodeLines += langStats.CodeLines
		totals.CommentLines += langStats.CommentLines
		totals.BlankLines += langStats.BlankLines
		totals.Characters += langStats.Characters
		totals.Functions += langStats.Functions
		totals.Classes += langStats.Classes
		totals.Size += langStats.Size
	}

	output.Summary = map[string]interface{}{
		"total_files":      totals.Files,
		"total_lines":      totals.Lines,
		"total_code_lines": totals.CodeLines,
		"total_comments":   totals.CommentLines,
		"total_blank":      totals.BlankLines,
		"total_chars":      totals.Characters,
		"total_functions":  totals.Functions,
		"total_classes":    totals.Classes,
		"total_size":       totals.Size,
		"code_ratio":       float64(totals.CodeLines) / float64(totals.Lines) * 100,
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
