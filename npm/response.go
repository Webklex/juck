package npm

import (
	"bytes"
	"encoding/json"
	"time"
)

type RepositoryResponse struct {
	ID                    string `json:"_id"`
	Rev                   string `json:"_rev"`
	RepositoryName        string `json:"name"`
	RepositoryDescription string `json:"description"`
	DistTags              struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
	RepositoryVersions map[string]Version   `json:"versions"`
	Readme             string               `json:"readme"`
	Maintainers        []User               `json:"maintainers"`
	Time               map[string]time.Time `json:"time"`
	Homepage           string               `json:"homepage"`
	Keywords           StringList           `json:"keywords"`
	Repository         Repository           `json:"repository"`
	RepositoryAuthor   User                 `json:"author"`
	Bugs               struct {
		Url string `json:"url"`
	} `json:"bugs"`
	RepositoryLicense string          `json:"license"`
	ReadmeFilename    string          `json:"readmeFilename"`
	Contributors      []User          `json:"contributors"`
	Users             map[string]bool `json:"users"`
}

type Version struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Version     string     `json:"version"`
	Homepage    string     `json:"homepage"`
	Author      User       `json:"author"`
	Repository  Repository `json:"repository"`
	Bugs        struct {
		Url string `json:"url"`
	} `json:"bugs"`
	License         string                 `json:"license"`
	Files           StringList             `json:"files"`
	Main            string                 `json:"main"`
	Engines         map[string]string      `json:"engines"`
	Scripts         map[string]string      `json:"scripts"`
	Dependencies    map[string]string      `json:"dependencies"`
	Verb            map[string]interface{} `json:"verb"`
	Keywords        StringList             `json:"keywords"`
	DevDependencies map[string]string      `json:"devDependencies"`
	GitHead         string                 `json:"gitHead"`
	Id              string                 `json:"_id"`
	Shasum          string                 `json:"_shasum"`
	From            string                 `json:"_from"`
	NpmVersion      string                 `json:"_npmVersion"`
	NodeVersion     string                 `json:"_nodeVersion"`
	NpmUser         User                   `json:"_npmUser"`
	Dist            struct {
		Shasum     string `json:"shasum"`
		Tarball    string `json:"tarball"`
		Integrity  string `json:"integrity"`
		Signatures []struct {
			Keyid string `json:"keyid"`
			Sig   string `json:"sig"`
		} `json:"signatures"`
	} `json:"dist"`
	Maintainers            []User `json:"maintainers"`
	NpmOperationalInternal struct {
		Host string `json:"host"`
		Tmp  string `json:"tmp"`
	} `json:"_npmOperationalInternal"`
	Directories StringList `json:"directories"`
}

type StringList []string

type Repository struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type User struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Email string `json:"email"`
}

func (r *RepositoryResponse) Id() string {
	return r.ID
}

func (r *RepositoryResponse) Author() string {
	return r.RepositoryAuthor.Name
}

func (r *RepositoryResponse) Name() string {
	return r.RepositoryName
}

func (r *RepositoryResponse) Description() string {
	return r.RepositoryDescription
}

func (r *RepositoryResponse) License() string {
	return r.RepositoryLicense
}

func (r *RepositoryResponse) Maintainer() string {
	return r.Maintainers[0].Name
}

func (r *RepositoryResponse) Url() string {
	return r.Repository.Url
}

func (r *RepositoryResponse) Type() string {
	return r.Repository.Type
}

func (r *RepositoryResponse) Versions() (versions []string) {
	for k, _ := range r.RepositoryVersions {
		versions = append(versions, k)
	}
	return
}

func (r *RepositoryResponse) Dependencies() map[string]string {
	return r.RepositoryVersions[r.DistTags.Latest].Dependencies
}

func (r *RepositoryResponse) Original() interface{} {
	return r
}

func (sl *StringList) UnmarshalJSON(data []byte) error {
	if bytes.Compare(data, []byte("{}")) != 0 {
		var v []string
		if err := json.Unmarshal(data, &v); err != nil {
			return err
		}
		*sl = v
	}
	return nil
}
