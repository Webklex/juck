package npm

import (
	"encoding/json"
	"fmt"
	"github.com/webklex/juck/utils"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Npm struct {
	registry string
	frontend string
}

type cacheItem struct {
	error    error
	response *RepositoryResponse
}

var cache = map[string]*cacheItem{}

func NewNpmRegistry() *Npm {
	return &Npm{
		registry: "https://registry.npmjs.org/",
	}
}

func (npm *Npm) Get(name string) (*RepositoryResponse, error) {
	if c, ok := cache[name]; ok {
		return c.response, c.error
	}
	u, err := url.Parse(npm.registry + url.QueryEscape(name))
	if err != nil {
		return registerCache(name, nil, err)
	}
	resp, err := npm.request(http.MethodGet, u, nil)
	if err != nil {
		return registerCache(name, nil, err)
	}

	var r RepositoryResponse
	return registerCache(name, &r, json.Unmarshal(resp, &r))
}

func (npm *Npm) Dependencies(name string, dependencies ...string) (result []string, err error) {
	pkg, err := npm.Get(name)
	if err != nil {
		return nil, err
	}

	result = append(dependencies, pkg.Name())
	for dependency, _ := range pkg.Dependencies() {
		if utils.InStringList(result, dependency) {
			continue
		}
		if entries, err2 := npm.Dependencies(dependency, result...); entries != nil {
			result = entries
		} else if err2 != nil {
			result = append(result, dependency)
		}
	}

	return
}

func (npm *Npm) request(method string, u *url.URL, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("npm: could not create request: %s\n", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("npm: error making http request: %s\n", err)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("npm: invalid response status: %d - %s\n", res.StatusCode, res.Status)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("npm: could not read response body: %s\n", err)
	}
	return resBody, nil
}

func registerCache(name string, r *RepositoryResponse, err error) (*RepositoryResponse, error) {
	cache[name] = &cacheItem{
		error:    err,
		response: r,
	}

	return cache[name].response, cache[name].error
}
