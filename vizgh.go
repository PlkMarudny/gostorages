package gostorages

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type VizGHStorage struct {
	*BaseStorage
	User       string
	Pass       string
	Expiration time.Duration
}

type VizGHFile struct {
}

func NewVizGHStorage(url string, folder string, user string, pass string, expireAfter int) *VizGHStorage {
	return &VizGHStorage{
		NewBaseStorage(folder, url),
		user,
		pass,
		time.Duration(time.Duration(expireAfter) * 24 * time.Hour),
	}
}

// Save saves a file at the given path; the path is the uuid of a folder
// we want to save a file to
func (v *VizGHStorage) Save(name string, file File) error {
	content, err := file.ReadAll()
	if err != nil {
		return err
	}

	// upload the image
	sUrl := assembleSaveUrl(v.BaseURL, v.Location, name)

	client := &http.Client{}
	req, _ := http.NewRequest("POST", sUrl.String(), bytes.NewReader(content))
	req.Header.Set("Slug", filepath.Base(name))
	req.Header.Set("Content-Type", mime.TypeByExtension(filepath.Ext(name)))
	req.SetBasicAuth(v.User, v.Pass)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}
	// unmarshal the response to get metadata url
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	req.Body.Close()
	img := GHFeed{}
	err = xml.Unmarshal(result, &img)
	if err != nil {
		return err
	}

	metaUrl := getMetadataUrl(&img)
	if metaUrl == "" {
		return errors.New("No metadata url found in VizGH reply")
	}

	// read metadata
	req2, err := http.NewRequest("GET", metaUrl, nil)
	if err != nil {
		return err
	}
	req2.SetBasicAuth(v.User, v.Pass)
	resp, err = client.Do(req2)
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}
	defer req2.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	meta := GHPayload{}
	err = xml.Unmarshal(result, &meta)
	if err != nil {
		return err
	}

	// modify expiration date
	err = modifyExpirationDate(&meta)
	if err != nil {
		return err
	}

	// marshal the xml containing the metadata
	output, err := xml.Marshal(&meta)
	if err != nil {
		return err
	}

	// update the metadata in the VizGH
	req3, _ := http.NewRequest("PUT", metaUrl, bytes.NewReader(output))
	req3.Header.Set("Content-Type", "application/vnd.vizrt.payload+xml")
	req3.SetBasicAuth(v.User, v.Pass)
	resp, err = client.Do(req3)
	defer req3.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
	}

func modifyExpirationDate(feed *GHPayload) error {
	for _, v := range feed.GHField {
		if v.AttrName == "date-expired" {
			t := time.Now().Add(time.Duration(time.Duration(7) * 24 * time.Hour))
			v.GHValue.Text = t.UTC().Format(time.RFC3339)
			return nil
		}
	}
	return errors.New("No \"date-expired\" tag was found in the metadata vdf")
}

func getMetadataUrl(feed *GHFeed) string {
	var metadataUrl string
	for i := range feed.GHEntry.GHLink {
		if feed.GHEntry.GHLink[i].AttrRel == "describedby" {
			// get metadata of the image created
			metadataUrl = feed.GHEntry.GHLink[i].AttrHref
		}
	}
	return metadataUrl
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
