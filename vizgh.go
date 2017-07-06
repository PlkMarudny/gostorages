package gostorages

import (
	"time"
	"strings"
	"net/http"
	"bytes"
	"mime"
	"net/url"
	"path/filepath"
)

type VizGHStorage struct {
	*BaseStorage
	User		string
	Pass		string
}

type VizGHFile struct {

}

func NewVizGHStorage(url string, folder string, user string, pass string) *VizGHStorage {
	return &VizGHStorage{
		NewBaseStorage(folder, url),
		user,
		pass,
	}
}

// Save saves a file at the given path; the path is the uuid of a folder
// we want to save a file to
func (v *VizGHStorage) Save(name string, file File) error {
	content, err := file.ReadAll()
	if err != nil {
		return err
	}

	saveUrl := assembleSaveUrl(v.BaseURL, v.Location, name)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", saveUrl.String(), bytes.NewReader(content))
	req.Header.Set("Slug", filepath.Base(name))
	req.Header.Set("Content-Type", mime.TypeByExtension(filepath.Ext(name)))
	req.SetBasicAuth(v.User, v.Pass)

	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// url in a form http://server/uuid/name, where uuid describes the folder
func assembleSaveUrl(base, location, name string) *url.URL {
	b, _ := url.Parse(base)
	u, _ := url.Parse(strings.Join([]string{"folder", location, name}, "/"))
	return b.ResolveReference(u)
}

func (v *VizGHStorage) Path(filepath string) string {
	panic("implement me")
}

func (v *VizGHStorage) Exists(filepath string) bool {
	panic("implement me")
}

func (v *VizGHStorage) Delete(filepath string) error {
	panic("implement me")
}

func (v *VizGHStorage) Open(filepath string) (File, error) {
	panic("implement me")
}

func (v *VizGHStorage) ModifiedTime(filepath string) (time.Time, error) {
	panic("implement me")
}

func (v *VizGHStorage) Size(filepath string) int64 {
	panic("implement me")
}

func (v *VizGHStorage) URL(filename string) string {
	if v.HasBaseURL() {
		return strings.Join([]string{v.BaseURL, "file", filename}, "/")
	}
	return ""
}

func (v *VizGHStorage) HasBaseURL() bool {
	return v.BaseURL != ""
}



