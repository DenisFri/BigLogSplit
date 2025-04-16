# Text File Splitter

This Go script splits large log files into smaller parts. This project uses the [Bubble Tea](https://github.com/charmbracelet/bubbletea) library for creating terminal-based UIs.

## Features

- **Split Large Files:** The program can split large log files into smaller parts, each of a specified maximum size.
- **Terminal Progress Bar:** The progress of the file splitting is displayed as a progress bar in the terminal.
- **File Size Metrics:** Shows processed and total file size in human-readable format (KB, MB, GB).
- **ETA Calculation:** Displays estimated time remaining to complete the file splitting.
- **Dynamic UI:** Uses a spinner for small files and a progress bar for larger files.
- **Completion Message:** Once the file is fully split, a completion message is shown, and the program waits for a key press before exiting.

## Installation

1. **Prerequisites:** Make sure you have Go installed on your machine. You can download and install Go from [here](https://golang.org/dl/).

2. **Clone the Repository:**

   ```bash
   git clone https://github.com/yourusername/yourrepository.git
   cd yourrepository
   ```

3. **Install Dependencies:**

   ```bash
   go get github.com/charmbracelet/bubbletea
   go get github.com/charmbracelet/bubbles/progress
   go get github.com/charmbracelet/bubbles/spinner
   ```
   
4. **Build the Program:**

   ```bash
   go build -o file-splitter
   ```

## Usage

1. **Prepare a `config.json` File:**

Create a `config.json` file in the root directory of the repository with the following structure:

```json
   {
       "filePath": "path/to/input/file.log",
       "maxSizeMB": 200,
       "outputFolder": "path/to/output/directory"
   }
   ```

   - `filePath`: The path to the input log file that needs to be split.
   - `maxSizeMB`: The maximum size of each split file in megabytes.
   - `outputFolder`: The path to the directory where the split files will be saved.

2. **Run the Program:**

   ```bash
    ./BigLogSplit
    ```
   
## Example

Given a `config.json` file with the following content:

```json
{
  "filePath": "C:\\Users\\yourusername\\Desktop\\LargeLogFile.log",
  "maxSizeMB": 100,
  "outputFolder": "C:\\Users\\yourusername\\Desktop\\SplitLogs"
}
```
The program will split `LargeLogFile.log` into parts of 100 MB each, saving them in the SplitLogs directory on your desktop.

## UI Features

- **Progress Bar:** For files larger than 10MB, a progress bar shows the percentage completed.
- **Spinner:** For smaller files, a spinner animation is displayed instead of a progress bar.
- **File Size Display:** Shows how much of the file has been processed (e.g., "1.25 MB of 5.00 MB").
- **ETA Display:** Shows the estimated time remaining to complete the operation.

## Customization

You can adjust the speed of the progress bar by modifying the `time.Sleep(10 * time.Millisecond)` line in the `splitFile` function within the `split.go` file. Adjust the duration to control how quickly the progress bar updates.

## Acknowledgements

- Bubble Tea: A Go framework for building terminal applications. Learn more at [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- Go: The Go programming language. Learn more at golang.org.
