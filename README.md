# Walker - Advanced Code Analysis Tool

## Features

### Comprehensive Analysis
- **Multi-language Support**: 35+ programming languages and file formats
- **Detailed Statistics**: Lines of code, comments, blank lines, functions, classes, file sizes
- **Smart Comment Detection**: Language-specific comment pattern recognition
- **Function & Class Counting**: Accurate detection using language-specific patterns
- **File Size Analysis**: Total codebase size with human-readable formatting

### Beautiful Output
- **Colorized Terminal Output**: Eye-catching colored terminal display
- **Progress Bars**: Real-time analysis progress with file count and speed
- **Multiple Output Formats**: Table view (default) and JSON export
- **Top Files Ranking**: See your largest files at a glance
- **Summary Statistics**: Code ratio, average lines per function, and more

### Performance & Efficiency
- **Concurrent Processing**: Multi-threaded analysis for maximum speed
- **Smart File Filtering**: Automatic exclusion of build artifacts and temporary files
- **Memory Efficient**: Optimized for large codebases
- **Fast Pattern Matching**: Regex-based analysis for accurate results

### Flexible Configuration
- **Command Line Flags**: Extensive customization options
- **Include/Exclude Patterns**: Fine-tune what gets analyzed
- **Custom Root Directory**: Analyze any directory
- **Configurable Output**: Control what information is displayed

## Supported Languages

| Language | Extensions | Functions | Classes | Comments |
|----------|------------|-----------|---------|----------|
| **Go** | `.go` | ✅ | ✅ | ✅ |
| **Python** | `.py`, `.pyw`, `.pyx` | ✅ | ✅ | ✅ |
| **JavaScript** | `.js`, `.jsx`, `.mjs`, `.cjs` | ✅ | ✅ | ✅ |
| **TypeScript** | `.ts`, `.tsx` | ✅ | ✅ | ✅ |
| **Java** | `.java` | ✅ | ✅ | ✅ |
| **C** | `.c`, `.h` | ✅ | ❌ | ✅ |
| **C++** | `.cpp`, `.cc`, `.cxx`, `.hpp`, `.hxx` | ✅ | ✅ | ✅ |
| **C#** | `.cs` | ✅ | ✅ | ✅ |
| **Rust** | `.rs` | ✅ | ✅ | ✅ |
| **PHP** | `.php`, `.phtml` | ✅ | ✅ | ✅ |
| **Ruby** | `.rb`, `.rbw` | ✅ | ✅ | ✅ |
| **Swift** | `.swift` | ✅ | ✅ | ✅ |
| **Kotlin** | `.kt`, `.kts` | ✅ | ✅ | ✅ |
| **Dart** | `.dart` | ✅ | ✅ | ✅ |
| **Scala** | `.scala`, `.sc` | ✅ | ✅ | ✅ |
| **Lua** | `.lua` | ✅ | ❌ | ✅ |
| **Perl** | `.pl`, `.pm`, `.perl` | ✅ | ❌ | ✅ |
| **R** | `.r`, `.R`, `.Rmd` | ✅ | ❌ | ✅ |
| **MATLAB** | `.m`, `.mlx` | ✅ | ❌ | ✅ |
| **Julia** | `.jl` | ✅ | ❌ | ✅ |
| **Haskell** | `.hs`, `.lhs` | ✅ | ❌ | ✅ |
| **Erlang** | `.erl`, `.hrl` | ✅ | ❌ | ✅ |
| **Elixir** | `.ex`, `.exs` | ✅ | ❌ | ✅ |
| **F#** | `.fs`, `.fsx`, `.fsi` | ✅ | ❌ | ✅ |
| **OCaml** | `.ml`, `.mli` | ✅ | ❌ | ✅ |
| **Assembly** | `.asm`, `.s`, `.S` | ❌ | ❌ | ✅ |
| **Shell** | `.sh`, `.bash`, `.zsh`, `.fish` | ✅ | ❌ | ✅ |
| **Vim Script** | `.vim`, `.vimrc` | ✅ | ❌ | ✅ |
| **Batch** | `.bat`, `.cmd` | ✅ | ❌ | ✅ |
| **PowerShell** | `.ps1`, `.psm1`, `.psd1` | ✅ | ❌ | ✅ |
| **Groovy** | `.groovy`, `.gradle` | ✅ | ✅ | ✅ |
| **HTML** | `.html`, `.htm`, `.xhtml` | ❌ | ❌ | ✅ |
| **CSS** | `.css`, `.scss`, `.sass`, `.less` | ❌ | ❌ | ✅ |
| **SQL** | `.sql` | ❌ | ❌ | ✅ |
| **YAML** | `.yml`, `.yaml` | ❌ | ❌ | ✅ |
| **JSON** | `.json` | ❌ | ❌ | ❌ |
| **XML** | `.xml`, `.xsd`, `.xsl` | ❌ | ❌ | ✅ |
| **Markdown** | `.md`, `.markdown` | ❌ | ❌ | ❌ |
| **TOML** | `.toml` | ❌ | ❌ | ✅ |
| **INI** | `.ini`, `.cfg`, `.conf` | ❌ | ❌ | ✅ |
| **Dockerfile** | `Dockerfile`, `.dockerfile` | ❌ | ❌ | ✅ |
| **Terraform** | `.tf`, `.tfvars` | ❌ | ❌ | ✅ |
| **GraphQL** | `.graphql`, `.gql` | ❌ | ❌ | ✅ |
| **Protobuf** | `.proto` | ❌ | ❌ | ✅ |
| **CMake** | `.cmake`, `CMakeLists.txt` | ✅ | ❌ | ✅ |
| **Makefile** | `Makefile`, `.mk`, `makefile` | ❌ | ❌ | ✅ |
| **Properties** | `.properties`, `.env` | ❌ | ❌ | ✅ |

## Installation

### Prerequisites
- Go 1.21 or higher

### Build from Source
```bash
git clone https://github.com/XanaOG/Walker.git
cd Walker
go mod tidy
go build -o walker main.go
```

### Run Directly
```bash
go run main.go [flags]
```

##  Usage

### Basic Usage
```bash
# Analyze current directory
./walker

# Analyze specific directory
./walker -path /path/to/project

# Show progress bar (default: enabled)
./walker -progress=true

# Disable progress bar
./walker -progress=false
```

### Output Formats
```bash
# Table format (default)
./walker -format table

# JSON format
./walker -format json

# JSON with top files
./walker -format json -top 20
```

### Filtering Options
```bash
# Exclude patterns
./walker -exclude "*.test.go,*_test.py,node_modules"

# Include only specific patterns
./walker -include "*.go,*.py,*.js"

# Analyze with custom exclusions
./walker -exclude ".git,dist,build"
```

### Display Options
```bash
# Show top 20 files by line count
./walker -top 20

# Show detailed file statistics
./walker -detailed

# Group results by directory
./walker -by-dir
```

### Advanced Examples
```bash
# Comprehensive analysis with custom settings
./walker -path ./src -top 15 -exclude "vendor,node_modules" -format table

# JSON export for CI/CD
./walker -format json -progress=false > code-stats.json

# Quick analysis without progress bar
./walker -progress=false -top 5
```

##  Sample Output

### Table Format
```
 Walker - Advanced Code Analysis Tool
Analyzing codebase at: .

 Code Analysis Results
Generated on: 2024-01-15 14:30:45

LANGUAGE     FILES      LINES     CODE   COMMENTS    BLANK    CHARS    FUNCS  CLASSES
───────────────────────────────────────────────────────────────────────────────────────
Go               5       2547     1893       234      420    76543       89       12
TypeScript      12       1834     1245       289      300    52341       45        8
JavaScript       8       1432      987       145      300    41234       38        6
Python           4        876      654        87      135    23456       28        4
JSON             3        234      234         0        0     8765        0        0
YAML             2        123       98        15       10     3456        0        0
───────────────────────────────────────────────────────────────────────────────────────
TOTAL           34       7046     5111       770     1165   205795      200       30

 Top 10 Files by Lines:
 1. src/main.go                                      547 lines    16234 chars
 2. frontend/src/components/Dashboard.tsx            423 lines    12456 chars
 3. api/handlers/users.go                           312 lines     9876 chars
 4. frontend/src/utils/helpers.ts                   298 lines     8765 chars
 5. config/database.go                              267 lines     7654 chars

 Summary:
   Total Size: 200.8 kB
   Code Ratio: 72.5%
   Avg Lines/Function: 25.6

 https://github.com/XanaOG/Walker
   Enhanced by AI - Please respect the original author
```

### JSON Format
```json
{
  "generated_at": "2024-01-15T14:30:45Z",
  "languages": {
    "Go": {
      "files": 5,
      "lines": 2547,
      "code_lines": 1893,
      "comment_lines": 234,
      "blank_lines": 420,
      "characters": 76543,
      "functions": 89,
      "classes": 12,
      "size": 76543
    }
  },
  "summary": {
    "total_files": 34,
    "total_lines": 7046,
    "total_code_lines": 5111,
    "total_comments": 770,
    "total_blank": 1165,
    "total_chars": 205795,
    "total_functions": 200,
    "total_classes": 30,
    "total_size": 205795,
    "code_ratio": 72.5
  }
}
```

##  Command Line Options

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `-path` | string | `.` | Root directory to analyze |
| `-format` | string | `table` | Output format (table, json) |
| `-progress` | bool | `true` | Show progress bar |
| `-top` | int | `10` | Show top N files by lines |
| `-detailed` | bool | `false` | Show detailed file statistics |
| `-by-dir` | bool | `false` | Group results by directory |
| `-exclude` | string | | Comma-separated exclusion patterns |
| `-include` | string | | Comma-separated inclusion patterns |

##  Default Exclusions

Walker automatically excludes common non-source directories and files:
- Version control: `.git`, `.svn`, `.hg`, `.bzr`
- Dependencies: `node_modules`, `vendor`, `target`, `build`, `dist`
- IDEs: `.idea`, `.vscode`, `.vs`
- Binaries: `*.exe`, `*.dll`, `*.so`, `*.dylib`, `*.jar`, `*.war`, `*.class`
- Python: `*.pyc`, `*.pyo`, `__pycache__`
- System: `.DS_Store`, `Thumbs.db`

##  Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

---

**⭐ If you find this tool useful, please give it a star!**