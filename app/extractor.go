package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/webklex/juck/log"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Extractor struct {
	dir      string
	data     map[string]interface{}
	sources  []string
	contents []string
	combined bool
}

//
// NewExtractor
// @Description: Create anew Extractor instance
// @param dir string
// @return *Extractor
func NewExtractor(dir string) *Extractor {
	return &Extractor{
		dir:      dir,
		data:     map[string]interface{}{},
		sources:  make([]string, 0),
		contents: make([]string, 0),
		combined: false,
	}
}

//
// Extract
// @Description: Attempt to extract all files from a given source map
// @receiver e *Extractor
// @param filename string
// @return error
func (e *Extractor) Extract(filename string) (nodeModules []string, err error) {
	log.Info("Extracting: %s", filename)

	if err = e.load(filename); err != nil {
		return
	}

	if err = e.parseSources(); err != nil {
		return
	}
	if err = e.parseContents(); err != nil {
		return
	}

	if err = makeDirIfNotExist(path.Join(e.dir, "combined")); err != nil {
		return
	}
	if err = makeDirIfNotExist(path.Join(e.dir, "sources")); err != nil {
		return
	}

	targetFile := "combined.js"
	var tfh *os.File
	if e.combined {
		if tf, ok := e.data["file"]; ok && tf != "" {
			targetFile = tf.(string)
		}
		targetFile = path.Join(e.dir, "combined", SanitizePath(targetFile))
		if err := makeDirIfNotExist(filepath.Dir(targetFile)); err != nil {
			log.Error("Failed to create directory \"%s\": %s", targetFile, err.Error())
		}
		tfh, err = os.OpenFile(targetFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
		if err != nil {
			return
		}
	}

	sc := len(e.sources)
	cc := len(e.contents)

	log.Info("Discovered sources: %d", sc)
	log.Info("Discovered contents: %d", cc)

	if sc > cc {
		log.Warning("There are more sources than contents, filenames may not match content")
	} else if sc < cc {
		log.Warning("There are more contents than sources, filenames may not match content")
	}

	for i, content := range e.contents {
		sourcePath := path.Join(e.dir, "sources", fmt.Sprintf("undefined-%d.js", i))
		if i < sc {
			sourcePath = e.sources[i]
		}
		if content == "" {
			log.Warning("Skipping %s -  no content", sourcePath)
			continue
		}
		if ext := filepath.Ext(sourcePath); ext == "" {
			sourcePath = sourcePath + ".js"
		}

		if name := e.getModuleName(sourcePath); name != "" {
			nodeModules = append(nodeModules, name)
		}

		if err := makeDirIfNotExist(filepath.Dir(sourcePath)); err != nil {
			log.Error("Failed to create directory \"%s\": %s", sourcePath, err.Error())
		} else {
			e.saveSource(sourcePath, content, tfh)
		}
	}
	return
}

func (e *Extractor) getModuleName(sourcePath string) string {
	if i := strings.Index(sourcePath, "node_modules"); i >= 0 {
		if len(sourcePath) > i+13 {
			parts := strings.SplitN(sourcePath[i+13:], "/", 3)
			if len(parts) == 3 {
				author := parts[0]
				repository := parts[1]
				reference := parts[2]
				if author[0:1] == "@" {
					log.Success("Node module discovered: %s/%s (%s)", author, repository, reference)
					return strings.Join([]string{author, repository}, "/")
				}
				log.Success("Node module discovered: %s (%s/%s)", author, repository, reference)
				return author
			}
		}
	}
	return ""
}

//
// saveSource
// @Description: Save a given source file
// @receiver e *Extractor
// @param sourcePath string
// @param content string
// @param tfh *os.File
func (e *Extractor) saveSource(sourcePath, content string, tfh *os.File) {
	f, err := os.OpenFile(sourcePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)

	if err != nil {
		log.Error("Failed to open file \"%s\": %s", sourcePath, err.Error())
	} else {
		defer f.Close()
		data, err := ioutil.ReadAll(f)
		if err != nil {
			log.Error("Failed to read file \"%s\": %s", sourcePath, err.Error())
		} else {
			if strings.Contains(string(data), content) == false {
				if e.combined {
					if _, err := tfh.WriteString(fmt.Sprintf("\n/**\nRestored: %s\n**/\n\n%s\n\n", sourcePath, content)); err != nil {
						log.Error(err)
					}
				}
				if _, err := f.WriteString(content); err != nil {
					log.Error("Failed to write to file \"%s\": %s", sourcePath, err.Error())
					log.Error(err)
				} else {
					log.Success("Wrote to: %s", sourcePath)
				}
			} else {
				log.Info("Skipping %s -  content already known", sourcePath)
			}
		}
	}
}

//
// Combine
// @Description: Set the Combined flag - combine the output into one file
// @receiver e *Extractor
// @param state bool
func (e *Extractor) Combine(state bool) {
	e.combined = state
}

//
// parseSources
// @Description: Attempt to parse all sources specified within the webpack map
// @receiver e *Extractor
// @return error
func (e *Extractor) parseSources() error {
	_sources, ok := e.data["sources"]
	if !ok {
		return errors.New("sourcemap does not contain sources")
	}
	sources, ok := _sources.([]interface{})
	if !ok {
		return errors.New("sourcemap sources has an invalid format")
	}
	for _, s := range sources {
		if str, ok := s.(string); ok && str != "" {
			str = path.Join(e.dir, "sources", SanitizePath(str))
			e.sources = append(e.sources, str)
		}
	}
	return nil
}

//
// parseContents
// @Description: Attempt to parse all sourcesContent specified within the webpack map
// @receiver e *Extractor
// @return error
func (e *Extractor) parseContents() error {
	_sourcesContents, ok := e.data["sourcesContent"]
	if !ok {
		return errors.New("sourcemap does not contain sourcesContents")
	}
	sourcesContents, ok := _sourcesContents.([]interface{})
	if !ok {
		return errors.New("sourcemap sourcesContent has an invalid format")
	}
	for _, s := range sourcesContents {
		str, _ := s.(string)
		e.contents = append(e.contents, str)
	}
	return nil
}

//
// load
// @Description: Load a given webpack file
// @receiver e *Extractor
// @param filepath string
// @return error
func (e *Extractor) load(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	e.data = map[string]interface{}{}

	return json.Unmarshal(byteValue, &e.data)
}
