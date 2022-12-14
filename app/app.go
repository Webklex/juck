package app

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/webklex/juck/log"
	"github.com/webklex/juck/npm"
	"github.com/webklex/juck/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Application struct {
	OutputDir             string
	SourceFile            string
	SourceUrl             string
	FileList              string
	UrlList               string
	Delay                 time.Duration
	ForceDownload         bool
	DisableSSL            bool
	LocalOnly             bool
	DangerouslyWritePaths bool
	Combined              bool
	sources               []string
}

//
// NewApplication
// @Description: Create a new Application instance
// @return *Application
func NewApplication() *Application {
	dir, _ := os.Getwd()
	return &Application{
		OutputDir:             path.Join(dir, "output"),
		SourceFile:            "",
		SourceUrl:             "",
		Delay:                 0,
		ForceDownload:         false,
		DisableSSL:            false,
		DangerouslyWritePaths: false,
		Combined:              false,
		LocalOnly:             false,
		sources:               make([]string, 0),
	}
}

//
// Run
// @Description: Run the application
// @receiver a *Application
// @return error
func (a *Application) Run() error {
	if err := a.verify(); err != nil {
		return err
	}

	log.Statistic("Verified sources: %d", len(a.sources))
	var coreModules []string
	for _, source := range a.sources {
		e := NewExtractor(a.OutputDir)
		e.Combine(a.Combined)
		if nm, err := e.Extract(source); err != nil {
			log.Error(err)
		} else {
			coreModules = append(coreModules, nm...)
		}
	}

	n := npm.NewNpmRegistry()
	coreModules = utils.UniqueStringList(coreModules)
	sort.Strings(coreModules)

	log.Statistic("Discovered node modules: %d", len(coreModules))

	fh, err := os.OpenFile(path.Join(a.OutputDir, "node_modules.txt"), os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer fh.Close()
	for _, name := range coreModules {
		if _, err := fh.WriteString(name + "\n"); err != nil {
			return err
		}
	}

	nodeModules := coreModules
	nmc := len(nodeModules)
	for _, name := range coreModules {
		log.Info("Analyzing %s", name)
		if dependencies, _ := n.Dependencies(name, nodeModules...); dependencies != nil {
			nodeModules = utils.UniqueStringList(dependencies)
			if delta := len(nodeModules) - nmc; delta > 0 {
				log.Info("\t%d new dependencies discovered", delta)
				nmc = nmc + delta
			}
		}
	}

	sort.Strings(nodeModules)
	log.Statistic("Discovered node dependencies: %d", len(nodeModules))

	fh, err = os.OpenFile(path.Join(a.OutputDir, "dependencies.txt"), os.O_TRUNC|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer fh.Close()
	for _, name := range nodeModules {
		if _, err := fh.WriteString(name + "\n"); err != nil {
			return err
		}
	}

	return nil
}

//
// verify
// @Description: Verify all options and settings / prepare the battlefield
// @receiver a *Application
// @return error
func (a *Application) verify() error {
	err := makeDirIfNotExist(a.OutputDir)
	if err != nil {
		return err
	}

	if a.UrlList != "" {
		list, err := a.loadList(a.UrlList)
		if err != nil {
			return err
		}
		a.downloadList(list)
	}
	if a.FileList != "" {
		list, err := a.loadList(a.FileList)
		if err != nil {
			return err
		}

		a.loadLocalList(list)
	}

	if a.SourceUrl != "" {
		// Disable delay - there is only one file to be downloaded
		a.Delay = 0
		a.downloadList([]string{a.SourceUrl})
	}

	if a.SourceFile != "" {
		a.loadLocalList([]string{a.SourceFile})
	} else if len(a.sources) == 0 {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			source := sc.Text()
			if _, err := url.ParseRequestURI(source); err == nil {
				a.downloadList([]string{source})
			} else {
				a.loadLocalList([]string{source})
			}
		}

		if err = sc.Err(); err != nil {
			return err
		}
	}

	if len(a.sources) == 0 {
		return errors.New("no target specified. please use --file, --url or stdin and provide at least one target")
	}
	a.sources = utils.UniqueStringList(a.sources)

	return nil
}

//
// loadList
// @Description: Load a given file into a string list
// @receiver a *Application
// @param filepath string
// @return []string
// @return error
func (a *Application) loadList(filepath string) ([]string, error) {
	byteValue, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	list := strings.Split(string(byteValue), "\n")

	if len(list) == 0 {
		return nil, errors.New("list file is empty")
	}

	return list, nil
}

func (a *Application) loadLocalList(list []string) {
	for _, filename := range list {
		if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
			log.Error(err)
			continue
		}
		a.sources = append(a.sources, filename)
	}
}

func (a *Application) downloadList(list []string) {
	for _, _url := range list {
		if _url == "" {
			continue
		}
		u, err := url.Parse(_url)
		if err != nil {
			log.Error(err)
			continue
		}
		if strings.HasSuffix(u.Path, ".js") || strings.HasSuffix(u.Path, ".css") {
			u.Path = u.Path + ".map"
		}
		if strings.HasSuffix(u.Path, ".map") {
			filename := path.Join(a.OutputDir, "sourcemaps", SanitizePath(filepath.Base(u.Path)))
			if err := a.download(u.String(), filename); err != nil {
				log.Error(err)
			} else {
				a.sources = append(a.sources, filename)
			}
		}
	}
}

//
// download
// @Description: Download a given source to a given target
// @receiver a *Application
// @param source string
// @param filepath string
// @return error
func (a *Application) download(source, target string) error {
	if _, err := os.Stat(target); err == nil {
		// File already exist
		if a.ForceDownload == false {
			log.Info("Local cache: %s", source)
			return nil
		}
	}
	if a.LocalOnly {
		return errors.New("local only mode is active")
	}
	log.Info("Downloading: %s", source)
	if err := makeDirIfNotExist(filepath.Dir(target)); err != nil {
		return err
	}
	// Create the file
	out, err := os.Create(target)
	if err != nil {
		_ = os.Remove(target)
		return err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Error(err)
		}
		if a.Delay > 0 {
			time.Sleep(a.Delay)
		}
	}(out)

	// Get the data
	resp, err := http.Get(source)
	if err != nil {
		_ = os.Remove(target)
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		_ = os.Remove(target)
		return fmt.Errorf("failed to download: %s - %s", source, resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		_ = os.Remove(target)
		return err
	}

	return nil
}

//
// makeDirIfNotExist
// @Description: Create all directories in a given path
// @param dirName string
// @return error
func makeDirIfNotExist(dirName string) error {
	err := os.MkdirAll(dirName, os.ModeDir)
	if err == nil {
		return nil
	}
	if os.IsExist(err) {
		// check that the existing path is a directory
		info, err := os.Stat(dirName)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	return err
}
