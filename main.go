package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/webklex/juck/app"
	"github.com/webklex/juck/log"
	"os"
)

var buildNumber string
var buildVersion string

func main() {
	a := app.NewApplication()

	flag.CommandLine.StringVar(&a.OutputDir, "output", a.OutputDir, "Directory to output from sourcemap to")
	flag.CommandLine.StringVar(&a.FileList, "file-list", a.FileList, "File path of a file containing a list of target source map file paths")
	flag.CommandLine.StringVar(&a.UrlList, "url-list", a.UrlList, "File path of a file containing a list of target source map urls")
	flag.CommandLine.StringVar(&a.SourceFile, "file", a.SourceFile, "Target sourcemap file path")
	flag.CommandLine.StringVar(&a.SourceUrl, "url", a.SourceUrl, "Target sourcemap url")
	flag.CommandLine.DurationVar(&a.Delay, "delay", a.Delay, "Delay between two requests. Only applies if --url-list is used")
	flag.CommandLine.BoolVar(&a.ForceDownload, "force", a.ForceDownload, "Force to download and overwrite local sourcemap")
	flag.CommandLine.BoolVar(&a.LocalOnly, "local", a.LocalOnly, "Only use local files. Don't perform any requests")
	flag.CommandLine.BoolVar(&a.Combined, "combined", a.Combined, "Combine all source files into one")
	flag.CommandLine.BoolVar(&a.DisableSSL, "disable-ssl", a.DisableSSL, "Don't verify the site's SSL certificate")
	flag.CommandLine.IntVar(&log.Mode, "log", log.Mode, "Set the log mode (0 = all, 1 = success, 2 = warning, 3 = statistic, 4 = error)")
	flag.CommandLine.BoolVar(&a.DangerouslyWritePaths, "dangerously-write-paths", a.DangerouslyWritePaths, "Write full paths. WARNING: Be careful here, you are pulling directories from an untrusted source")

	sv := flag.Bool("version", false, "Show version and exit")
	nc := flag.Bool("no-color", false, "Disable color output")
	flag.Parse()

	if *nc {
		color.NoColor = true // disables colorized output
	}

	if *sv {
		fmt.Printf("version: %s\nbuild number: %s\n", color.CyanString(buildVersion), color.CyanString(buildNumber))
		os.Exit(0)
	}

	if err := a.Run(); err != nil {
		log.Error(err)
	}
}
