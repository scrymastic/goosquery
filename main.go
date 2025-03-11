package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/scrymastic/goosquery/collector"
	"github.com/scrymastic/goosquery/networking"
	"github.com/scrymastic/goosquery/system"
	"github.com/scrymastic/goosquery/utility"
	"github.com/sirupsen/logrus"
)

// Version information set by the build process
var (
	Version   = "dev"
	BuildTime = "unknown"
)

var log = logrus.New()

// setupLogging configures the logger with colors and formatting
func setupLogging(toFile bool, outputDir string) {
	// Configure console output
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		PadLevelText:    true,
		DisableQuote:    true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	// If logging to file is requested, create a multi-writer
	if toFile && outputDir != "" {
		// Create a log file
		logFile, err := os.Create(filepath.Join(outputDir, "goosquery.log"))
		if err != nil {
			log.Warnf("Failed to create log file: %v. Logs will only be displayed in console.", err)
		} else {
			// Create a hook to write to the file with a different formatter
			fileFormatter := &logrus.TextFormatter{
				DisableColors:   true,
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
			}

			// Create a hook to write to the file with a different formatter
			fileHook := &FileHook{
				Writer:    logFile,
				Formatter: fileFormatter,
			}
			log.AddHook(fileHook)
		}
	}
}

// FileHook is a hook for logrus that writes to a file
type FileHook struct {
	Writer    *os.File
	Formatter logrus.Formatter
}

// Fire writes the log entry to the file
func (hook *FileHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

// Levels returns the levels this hook is enabled for
func (hook *FileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// CollectorFunc is a generic type for all Gen functions
type CollectorFunc interface{}

// CollectorResult stores the result of a collector function
type CollectorResult struct {
	Name      string
	Data      interface{}
	Error     error
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

// BenchmarkResult stores benchmark information for a collector
type BenchmarkResult struct {
	Name            string        `json:"name"`
	Iterations      int           `json:"iterations"`
	TotalDuration   time.Duration `json:"total_duration"`
	AverageDuration time.Duration `json:"average_duration"`
	MinDuration     time.Duration `json:"min_duration"`
	MaxDuration     time.Duration `json:"max_duration"`
	SuccessRate     float64       `json:"success_rate"`
}

// RunCollector executes a collector function and returns its result with timing information
func RunCollector(name string, fn CollectorFunc, args ...interface{}) CollectorResult {
	log.Infof("Starting collector: %s", name)

	result := CollectorResult{
		Name:      name,
		StartTime: time.Now(),
	}

	// Use reflection to call the function with variable arguments
	fnValue := reflect.ValueOf(fn)

	// Prepare arguments
	var argValues []reflect.Value
	for _, arg := range args {
		argValues = append(argValues, reflect.ValueOf(arg))
	}

	// Call the function
	returnValues := fnValue.Call(argValues)

	// Process return values (data and error)
	if len(returnValues) >= 1 {
		result.Data = returnValues[0].Interface()
	}

	if len(returnValues) >= 2 && !returnValues[1].IsNil() {
		result.Error = returnValues[1].Interface().(error)
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	if result.Error != nil {
		log.Errorf("Collector %s failed: %v", name, result.Error)
	} else {
		log.Infof("Collector %s completed in %v", name, result.Duration.Milliseconds())
	}

	return result
}

// BenchmarkCollector runs a collector multiple times to measure its performance
func BenchmarkCollector(name string, fn CollectorFunc, iterations int, args ...interface{}) BenchmarkResult {
	log.Infof("Benchmarking collector: %s (%d iterations)", name, iterations)

	var durations []time.Duration
	successCount := 0

	for i := 0; i < iterations; i++ {
		result := RunCollector(fmt.Sprintf("%s (iteration %d/%d)", name, i+1, iterations), fn, args...)

		if result.Error == nil {
			successCount++
			durations = append(durations, result.Duration)
		}
	}

	// Calculate benchmark statistics
	benchmark := BenchmarkResult{
		Name:       name,
		Iterations: iterations,
	}

	if len(durations) > 0 {
		// Sort durations for min/max
		sort.Slice(durations, func(i, j int) bool {
			return durations[i] < durations[j]
		})

		var totalDuration time.Duration
		for _, d := range durations {
			totalDuration += d
		}

		benchmark.TotalDuration = totalDuration
		benchmark.AverageDuration = totalDuration / time.Duration(len(durations))
		benchmark.MinDuration = durations[0]
		benchmark.MaxDuration = durations[len(durations)-1]
	}

	benchmark.SuccessRate = float64(successCount) / float64(iterations) * 100.0

	log.Infof("Benchmark results for %s: Avg=%v, Min=%v, Max=%v, Success=%.1f%%",
		name, benchmark.AverageDuration, benchmark.MinDuration, benchmark.MaxDuration, benchmark.SuccessRate)

	return benchmark
}

// SaveResultToJSON saves the collector result to a JSON file
func SaveResultToJSON(result CollectorResult, outputDir string) error {
	if result.Error != nil {
		return fmt.Errorf("cannot save result with error: %v", result.Error)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create file
	filename := filepath.Join(outputDir, fmt.Sprintf("%s.json", result.Name))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Marshal data to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result.Data); err != nil {
		return fmt.Errorf("failed to encode data to JSON: %v", err)
	}

	log.Debugf("[SAVED] Saved result to %s", filename)
	return nil
}

// SaveBenchmarkToJSON saves benchmark results to a JSON file
func SaveBenchmarkToJSON(benchmarks []BenchmarkResult, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create file
	filename := filepath.Join(outputDir, "benchmarks.json")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create benchmark file: %v", err)
	}
	defer file.Close()

	// Marshal data to JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(benchmarks); err != nil {
		return fmt.Errorf("failed to encode benchmarks to JSON: %v", err)
	}

	log.Infof("[SAVED] Saved benchmark results to %s", filename)
	return nil
}

// RunAllCollectors runs all collector functions and saves results to JSON files
func RunAllCollectors(outputDir string) {
	// Setup logging with colors and formatting
	setupLogging(true, outputDir)

	log.Info("[START] Starting GoOSQuery data collection")

	// Initialize all collectors
	collector.Initialize()

	// Import packages to trigger their init functions
	_ = networking.GenARPCache
	_ = system.GenOSVersion
	_ = utility.GenTime

	// Get all collectors from the registry
	allCollectors := collector.GetAll()
	log.Infof("[INFO] Found %d collectors in the registry", len(allCollectors))

	// Run all collectors and save results
	var results []CollectorResult
	var skippedCount int
	var successCount int
	var failureCount int

	log.Info("[INFO] Starting collector execution...")

	// Calculate total collectors for progress reporting
	totalCollectors := len(allCollectors)

	for i, c := range allCollectors {
		// Skip collectors that require arguments
		requiresArg, argType, argDesc, err := collector.RequiresArgument(c.Name)
		if err != nil {
			log.Errorf("[ERROR] Error checking if collector %s requires arguments: %v", c.Name, err)
			continue
		}

		if requiresArg {
			log.Debugf("[SKIP] Skipping collector %s because it requires arguments: %s (%s)", c.Name, argDesc, argType)
			skippedCount++
			continue
		}

		// Show progress
		progress := float64(i+1) / float64(totalCollectors) * 100
		log.Infof("[RUNNING] [%d/%d] Running collector: %s (%.1f%% complete)", i+1, totalCollectors, c.Name, progress)

		// Run the collector using the Execute function
		startTime := time.Now()
		result, err := collector.Execute(c.Name)
		endTime := time.Now()
		duration := endTime.Sub(startTime)

		if err != nil {
			log.Errorf("[FAILED] Failed to execute collector %s: %v", c.Name, err)
			failureCount++
			results = append(results, CollectorResult{
				Name:      c.Name,
				Error:     err,
				StartTime: startTime,
				EndTime:   endTime,
				Duration:  duration,
			})
			continue
		}

		log.Infof("[SUCCESS] Collector %s completed in %v", c.Name, duration)
		successCount++

		collectorResult := CollectorResult{
			Name:      c.Name,
			Data:      result,
			StartTime: startTime,
			EndTime:   endTime,
			Duration:  duration,
		}
		results = append(results, collectorResult)

		if err := SaveResultToJSON(collectorResult, outputDir); err != nil {
			log.Warnf("[WARNING] Failed to save %s result: %v", c.Name, err)
		} else {
			log.Debugf("[SAVED] Saved result for %s", c.Name)
		}
	}

	// Generate summary report
	summary := struct {
		TotalCollectors int
		SuccessCount    int
		FailureCount    int
		SkippedCount    int
		TotalDuration   time.Duration
		CollectorStats  []struct {
			Name     string
			Duration time.Duration
			Success  bool
		}
	}{
		TotalCollectors: len(allCollectors),
		SuccessCount:    successCount,
		FailureCount:    failureCount,
		SkippedCount:    skippedCount,
	}

	for _, result := range results {
		success := result.Error == nil
		summary.TotalDuration += result.Duration
		summary.CollectorStats = append(summary.CollectorStats, struct {
			Name     string
			Duration time.Duration
			Success  bool
		}{
			Name:     result.Name,
			Duration: result.Duration,
			Success:  success,
		})
	}

	// Save summary to JSON
	summaryFile, err := os.Create(filepath.Join(outputDir, "summary.json"))
	if err != nil {
		log.Errorf("[ERROR] Failed to create summary file: %v", err)
	} else {
		defer summaryFile.Close()
		encoder := json.NewEncoder(summaryFile)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(summary); err != nil {
			log.Errorf("[ERROR] Failed to encode summary to JSON: %v", err)
		}
	}

	log.Infof("[COMPLETE] Data collection completed. Results: %d successful, %d failed, %d skipped. Total duration: %v",
		successCount, failureCount, skippedCount, summary.TotalDuration)
}

// RunBenchmarks runs benchmarks for all collector functions
func RunBenchmarks(outputDir string, iterations int) {
	// Setup logging with colors and formatting
	setupLogging(true, outputDir)

	log.Info("[START] Starting GoOSQuery benchmarks")

	// Initialize all collectors
	collector.Initialize()

	// Import packages to trigger their init functions
	_ = networking.GenARPCache
	_ = system.GenOSVersion
	_ = utility.GenTime

	// Get all collectors from the registry
	allCollectors := collector.GetAll()
	log.Infof("[INFO] Found %d collectors in the registry", len(allCollectors))

	// Run benchmarks
	var benchmarks []BenchmarkResult
	var skippedCount int

	log.Infof("[INFO] Starting benchmark execution with %d iterations per collector...", iterations)

	// Calculate total collectors for progress reporting
	totalCollectors := len(allCollectors)

	for i, c := range allCollectors {
		// Skip collectors that require arguments
		requiresArg, argType, argDesc, err := collector.RequiresArgument(c.Name)
		if err != nil {
			log.Errorf("[ERROR] Error checking if collector %s requires arguments: %v", c.Name, err)
			continue
		}

		if requiresArg {
			log.Debugf("[SKIP] Skipping benchmark for collector %s because it requires arguments: %s (%s)", c.Name, argDesc, argType)
			skippedCount++
			continue
		}

		// Show progress
		progress := float64(i+1) / float64(totalCollectors) * 100
		log.Infof("[RUNNING] [%d/%d] Benchmarking collector: %s (%.1f%% complete)", i+1, totalCollectors, c.Name, progress)

		// Run the benchmark using the Execute function
		var durations []time.Duration
		successCount := 0

		for j := 0; j < iterations; j++ {
			log.Debugf("  [ITER] Iteration %d/%d", j+1, iterations)
			startTime := time.Now()
			result, err := collector.Execute(c.Name)
			endTime := time.Now()
			duration := endTime.Sub(startTime)

			if err == nil && result != nil {
				successCount++
				durations = append(durations, duration)
				log.Debugf("  [SUCCESS] Iteration %d completed in %v", j+1, duration)
			} else {
				log.Debugf("  [FAILED] Iteration %d failed: %v", j+1, err)
			}
		}

		// Calculate benchmark statistics
		benchmark := BenchmarkResult{
			Name:       c.Name,
			Iterations: iterations,
		}

		if len(durations) > 0 {
			// Sort durations for min/max
			sort.Slice(durations, func(i, j int) bool {
				return durations[i] < durations[j]
			})

			var totalDuration time.Duration
			for _, d := range durations {
				totalDuration += d
			}

			benchmark.TotalDuration = totalDuration
			benchmark.AverageDuration = totalDuration / time.Duration(len(durations))
			benchmark.MinDuration = durations[0]
			benchmark.MaxDuration = durations[len(durations)-1]
			benchmark.SuccessRate = float64(successCount) / float64(iterations)

			log.Infof("[RESULT] Benchmark results for %s: Avg: %v, Min: %v, Max: %v, Success rate: %.1f%%",
				c.Name, benchmark.AverageDuration, benchmark.MinDuration, benchmark.MaxDuration, benchmark.SuccessRate*100)
		} else {
			log.Warnf("[WARNING] No successful iterations for %s", c.Name)
			benchmark.SuccessRate = 0
		}

		benchmarks = append(benchmarks, benchmark)
	}

	// Save benchmark results
	if err := SaveBenchmarkToJSON(benchmarks, outputDir); err != nil {
		log.Errorf("[ERROR] Failed to save benchmark results: %v", err)
	}

	log.Infof("[COMPLETE] Benchmarks completed. Processed %d collectors (%d skipped).",
		len(benchmarks), skippedCount)
}

// RunCollectorsWithArgs runs collectors that require arguments with default values
func RunCollectorsWithArgs(outputDir string) {
	log.Info("[START] Running collectors that require arguments with default values")

	// Create a map of collector names to default argument values
	defaultArgs := map[string]interface{}{
		"authenticode":       "C:\\Windows\\System32\\notepad.exe",
		"hash":               "C:\\Windows\\System32\\notepad.exe",
		"registry":           "HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion",
		"process_memory_map": uint32(os.Getpid()), // Use the current process ID
		"file":               "C:\\Windows\\System32\\notepad.exe",
		"curl":               []interface{}{"https://www.google.com", "GoOSQuery/1.0"},
	}

	// Run each collector with its default argument
	totalCollectors := len(defaultArgs)
	successCount := 0
	failureCount := 0

	i := 0
	for name, arg := range defaultArgs {
		i++
		// Show progress
		progress := float64(i) / float64(totalCollectors) * 100

		// Format the argument for display
		var argDisplay string
		if argSlice, ok := arg.([]interface{}); ok {
			argStrings := make([]string, len(argSlice))
			for i, a := range argSlice {
				argStrings[i] = fmt.Sprintf("%v", a)
			}
			argDisplay = strings.Join(argStrings, ", ")
		} else {
			argDisplay = fmt.Sprintf("%v", arg)
		}

		log.Infof("[RUNNING] [%d/%d] Running collector %s with argument: %s (%.1f%% complete)",
			i, totalCollectors, name, argDisplay, progress)

		var result interface{}
		var err error

		startTime := time.Now()

		// Check if the argument is a slice (for multiple arguments)
		if argSlice, ok := arg.([]interface{}); ok {
			result, err = collector.ExecuteWithMultipleArgs(name, argSlice)
		} else {
			result, err = collector.ExecuteWithArg(name, arg)
		}

		endTime := time.Now()
		duration := endTime.Sub(startTime)

		if err != nil {
			log.Errorf("[FAILED] Failed to execute collector %s with argument %v: %v", name, argDisplay, err)
			failureCount++
			continue
		}

		log.Infof("[SUCCESS] Collector %s completed in %v", name, duration)
		successCount++

		collectorResult := CollectorResult{
			Name:      name,
			Data:      result,
			StartTime: startTime,
			EndTime:   endTime,
			Duration:  duration,
		}

		if err := SaveResultToJSON(collectorResult, outputDir); err != nil {
			log.Warnf("[WARNING] Failed to save %s result: %v", name, err)
		} else {
			log.Debugf("[SAVED] Saved result for %s", name)
		}
	}

	log.Infof("[COMPLETE] Completed running collectors with arguments. Results: %d successful, %d failed.",
		successCount, failureCount)
}

// compressReportsFolder compresses the reports folder using Go's archive/zip package
func compressReportsFolder(folderPath string) error {
	log.Infof("[COMPRESS] Compressing folder: %s", folderPath)

	// Create the zip file name based on the folder path
	zipFileName := folderPath + ".zip"

	// Create a new zip file
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %v", err)
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the folder and add files to the zip
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Create a relative path for the file in the zip
		relPath, err := filepath.Rel(folderPath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %v", err)
		}

		log.Debugf("[COMPRESS] Adding file to zip: %s", relPath)

		// Open the file to be added
		fileToZip, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %v", path, err)
		}
		defer fileToZip.Close()

		// Create a new file header
		fileHeader, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create file header: %v", err)
		}

		// Set the name to the relative path
		fileHeader.Name = relPath

		// Use deflate compression method
		fileHeader.Method = zip.Deflate

		// Create the file in the zip
		writer, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			return fmt.Errorf("failed to create file in zip: %v", err)
		}

		// Copy the file content to the zip
		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			return fmt.Errorf("failed to write file content to zip: %v", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk through folder: %v", err)
	}

	log.Infof("[COMPRESS] Successfully compressed folder to: %s", zipFileName)
	return nil
}

// displayHelp shows usage information
func displayHelp() {
	greenColor := "\033[32m"
	resetColor := "\033[0m"
	boldColor := "\033[1m"

	fmt.Printf("%s%sUSAGE:%s\n", greenColor, boldColor, resetColor)
	fmt.Printf("  %sgoosquery [OUTPUT_DIR] [OPTIONS]%s\n\n", greenColor, resetColor)

	fmt.Printf("%s%sOPTIONS:%s\n", greenColor, boldColor, resetColor)
	fmt.Printf("  %s--help%s        Display this help message and exit\n\n", greenColor, resetColor)

	fmt.Printf("%s%sEXAMPLES:%s\n", greenColor, boldColor, resetColor)
	fmt.Printf("  goosquery                     # Run with default output directory 'reports'\n")
	fmt.Printf("  goosquery custom_dir          # Specify custom output directory\n")
}

func main() {
	// Configure basic logging for startup
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		DisableQuote:  true,
	})
	log.SetOutput(os.Stdout)

	// Check for help flag first
	for _, arg := range os.Args {
		if arg == "--help" || arg == "-h" {
			displayBanner()
			displayHelp()
			return
		}
	}

	// Print directly to ensure we see output
	fmt.Println("Starting GoOSQuery...")

	// Get output directory from command line or use default
	outputDir := "reports"
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "--") {
		outputDir = os.Args[1]
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Errorf("[ERROR] Failed to create output directory: %v", err)
		return
	}

	// Display startup banner
	displayBanner()

	// Run normal data collection
	RunAllCollectors(outputDir)

	// Run collectors that require arguments with default values
	RunCollectorsWithArgs(outputDir)

	// Compress the reports folder after completion
	if err := compressReportsFolder(outputDir); err != nil {
		log.Errorf("[ERROR] Failed to compress reports folder: %v", err)
	} else {
		log.Infof("[SUCCESS] Reports folder compressed successfully")
	}
}

// displayBanner shows a nice ASCII art banner
func displayBanner() {
	// ANSI color codes
	greenColor := "\033[32m"
	resetColor := "\033[0m"

	banner := `
   ____ _____  ____  _________ ___  _____  _______  __
  / __ '/ __ \/ __ \/ ___/ __ '/ / / / _ \/ ___/ / / /
 / /_/ / /_/ / /_/ (__  ) /_/ / /_/ /  __/ /  / /_/ / 
 \__, /\____/\____/____/\__, /\__,_/\___/_/   \__, /  
/____/                    /_/                /____/   -- scrymastic --                 
`
	fmt.Print(greenColor)
	fmt.Println(banner)
	fmt.Printf("GoOSQuery v%s - Windows System Information Collector\n", Version)
	fmt.Printf("Build Time: %s\n%s\n", BuildTime, resetColor)
}
