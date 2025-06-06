# Text File Splitter

This Go script splits large log files into smaller parts. This project uses the [Bubble Tea](https://github.com/charmbracelet/bubbletea) library for creating terminal-based UIs.

## Features

- **Split Large Files:** The program can split large log files into smaller parts, each of a specified maximum size.
- **Terminal Progress Bar:** The progress of the file splitting is displayed as a progress bar in the terminal.
- **File Size Metrics:** Shows processed and total file size in human-readable format (KB, MB, GB).
- **ETA Calculation:** Displays estimated time remaining to complete the file splitting.
- **Dynamic UI:** Uses a spinner for small files and a progress bar for larger files.
- **Completion Message:** Once the file is fully split, a completion message is shown, and the program waits for a key press before exiting.
- **Log Analysis:** Analyze log files for error and custom patterns, creating a detailed report.
- **Pattern Filtering:** Include or exclude lines matching specific patterns from the output files.

## Installation

1. **Prerequisites:** Make sure you have Go installed on your machine. You can download and install Go from [here](https://golang.org/dl/).

2. **Clone the Repository:**

   ```bash
   git clone https://github.com/DenisFri/BigLogSplit.git
   cd BigLogSplit
   ```

3. **Install Dependencies:**

   ```bash
   go get github.com/charmbracelet/bubbletea
   go get github.com/charmbracelet/bubbles/progress
   go get github.com/charmbracelet/bubbles/spinner
   go get github.com/charmbracelet/lipgloss
   ```
   
4. **Build the Program:**

   ```bash
   go build -o bin/BigLogSplit ./cmd/BigLogSplit
   ```

## Usage

1. **Prepare a `config.json` File:**

Create a `config.json` file in the root directory of the repository with the following structure:

```json
{
  "filePath": "path/to/input/file.log",
  "maxSizeMB": 200,
  "outputFolder": "path/to/output/directory",
  "analysis": {
    "enabled": true,
    "errorPatterns": [
      "\\bERROR\\b",
      "\\bException\\b",
      "\\bFAILED\\b"
    ],
    "customPatterns": [
      "\\bAPI\\b",
      "\\bUSER\\:[0-9]+\\b"
    ],
    "outputFile": "analysis_report.md"
  },
  "filtering": {
    "mode": "none",
    "patterns": [
      "\\bDEBUG\\b",
      "\\bTRACE\\b"
    ]
  }
}
```

Fields explanation:
   - `filePath`: The path to the input log file that needs to be split.
   - `maxSizeMB`: The maximum size of each split file in megabytes.
   - `outputFolder`: The path to the directory where the split files will be saved.
   - `analysis`: Configuration for log analysis features.
     - `enabled`: Whether to enable log analysis.
     - `errorPatterns`: List of regex patterns to identify errors in logs.
     - `customPatterns`: List of regex patterns for custom analysis.
     - `outputFile`: Path where the analysis report will be saved.
   - `filtering`: Configuration for line filtering.
     - `mode`: Filtering mode ("none", "include", or "exclude").
     - `patterns`: List of regex patterns to match for filtering.

2. **Run the Program:**

   ```bash
   ./bin/BigLogSplit
   ```
   
## Example

Given a `config.json` file with the following content:

```json
{
  "filePath": "C:\\Users\\yourusername\\Desktop\\LargeLogFile.log",
  "maxSizeMB": 100,
  "outputFolder": "C:\\Users\\yourusername\\Desktop\\SplitLogs",
  "analysis": {
    "enabled": true,
    "errorPatterns": ["\\bERROR\\b", "\\bFAILED\\b"],
    "customPatterns": ["\\bUSER\\b"],
    "outputFile": "analysis_report.md"
  },
  "filtering": {
    "mode": "exclude",
    "patterns": ["\\bDEBUG\\b"]
  }
}
```

The program will:
1. Split `LargeLogFile.log` into parts of 100 MB each
2. Save the split files in the SplitLogs directory
3. Exclude any lines containing "DEBUG"
4. Generate an analysis report of error patterns and user mentions

## UI Features

- **Progress Bar:** For files larger than 10MB, a progress bar shows the percentage completed.
- **Spinner:** For smaller files, a spinner animation is displayed instead of a progress bar.
- **File Size Display:** Shows how much of the file has been processed (e.g., "1.25 MB of 5.00 MB").
- **ETA Display:** Shows the estimated time remaining to complete the operation.

## Log Analysis

The log analysis feature scans the file for patterns defined in the configuration and generates a report containing:

- Total number of lines and bytes processed
- Number of matches for each error pattern
- Number of matches for each custom pattern
- Sample lines matching each pattern (up to 10 examples per pattern)

The report is generated in Markdown format and saved to the specified output file.

## Filtering Options

The filtering feature allows you to control which lines are included in the output files:

- **None:** No filtering is applied (default)
- **Include:** Only include lines matching at least one of the specified patterns
- **Exclude:** Exclude any lines matching any of the specified patterns

## Customization

You can adjust the speed of the progress bar by modifying the `time.Sleep(10 * time.Millisecond)` line in the `SplitFile` function within the `split.go` file. Adjust the duration to control how quickly the progress bar updates.

## Acknowledgements

- Bubble Tea: A Go framework for building terminal applications. Learn more at [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- Go: The Go programming language. Learn more at golang.org.
