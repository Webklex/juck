package app

import (
	"errors"
	"fmt"
	"github.com/webklex/juck/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
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
	DisableSSL            bool
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
		DisableSSL:            false,
		DangerouslyWritePaths: false,
		Combined:              false,
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

	for _, source := range a.sources {
		e := NewExtractor(a.OutputDir)
		e.Combine(a.Combined)
		if err := e.Extract(source); err != nil {
			log.Error(err)
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
		} else if len(list) == 0 {
			return errors.New("url list file is empty")
		}
		for i, _url := range list {
			u, err := url.Parse(_url)
			if err != nil {
				return err
			}
			if strings.HasSuffix(u.Path, ".js") {
				u.Path = u.Path + ".map"
			}
			filename := path.Join(a.OutputDir, fmt.Sprintf("source.%d.js.map", i))
			if err := a.download(u.String(), filename); err != nil {
				log.Error(err)
			} else {
				a.sources = append(a.sources, filename)
			}
			if a.Delay > 0 {
				time.Sleep(a.Delay)
			}
		}
	}
	if a.FileList != "" {
		list, err := a.loadList(a.FileList)
		if err != nil {
			return err
		} else if len(list) == 0 {
			return errors.New("file list file is empty")
		}

		for _, filename := range list {
			if _, err := os.Stat(a.SourceFile); errors.Is(err, os.ErrNotExist) {
				log.Error(err)
			} else {
				a.sources = append(a.sources, filename)
			}
		}
	}

	if a.SourceUrl != "" {
		// Download source
		u, err := url.Parse(a.SourceUrl)
		if err != nil {
			return err
		}
		if strings.HasSuffix(u.Path, ".js") {
			u.Path = u.Path + ".map"
		}

		a.SourceFile = path.Join(a.OutputDir, "source.js.map")
		if err := a.download(u.String(), a.SourceFile); err != nil {
			return err
		}
	}

	if a.SourceFile != "" {
		if _, err := os.Stat(a.SourceFile); errors.Is(err, os.ErrNotExist) {
			return err
		}
		a.sources = []string{a.SourceFile}
	} else if len(a.sources) == 0 {
		return errors.New("no target specified. please use --file or --url and define a valid target")
	}

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

	return strings.Split(string(byteValue), "\n"), nil
}

//
// download
// @Description: Download a given source to a given target
// @receiver a *Application
// @param source string
// @param filepath string
// @return error
func (a *Application) download(source, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		_ = os.Remove(filepath)
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(source)
	if err != nil {
		_ = os.Remove(filepath)
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		_ = os.Remove(filepath)
		return fmt.Errorf("failed to download: %s - %s", source, resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		_ = os.Remove(filepath)
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
