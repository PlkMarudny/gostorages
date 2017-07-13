package gostorages

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io/ioutil"
	//"mime"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type VizGHStorage struct {
	*BaseStorage
	User       string
	Pass       string
	Expiration time.Duration
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
	var result []byte

	fileContent, err := file.ReadAll()
	if err != nil {
		return err
	}
	// one http client for all the requests
	client := &http.Client{}

	// upload the image
	sUrl := assembleSaveUrl(v.BaseURL, v.Location)
	headers := map[string]string{"Content-Type": "image/jpeg", "Slug": trimExt(filepath.Base(name))}
	result, err = restRequest(client, sUrl.String(), POST, v.User, v.Pass, headers, fileContent)
	if err != nil {
		return err
	}
	// process the response - extract the metadata url
	img := GHFeed{}
	err = xml.Unmarshal(result, &img)
	if err != nil {
		return err
	}
	metaUrl := getMetadataUrl(&img)
	if metaUrl == "" {
		return errors.New("No metadata url found in VizGH reply")
	}

	// get metadata from the url
	result, err = restRequest(client, metaUrl, GET, v.User, v.Pass, map[string]string{}, nil)
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
	metaJson, err := xml.Marshal(&meta)
	if err != nil {
		return err
	}
	// update the metadata in the VizGH
	headers = map[string]string{"Content-Type": "application/vnd.vizrt.payload+xml"}
	_, err = restRequest(client, metaUrl, PUT, v.User, v.Pass, headers, metaJson)
	if err != nil {
		return errors.New("File probably exists")
	}
	return nil
}

func restRequest(client *http.Client, url, method, u, p string, h map[string]string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(payload))
	if err != nil {
		return []byte{}, err
	}
	for k, v := range h {
		req.Header.Set(k, v)
	}
	req.SetBasicAuth(u, p)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return []byte{}, errors.New(fmt.Sprintf("Error code: %d, error: %v", resp.StatusCode, err))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, nil
}

func trimExt(filename string) string {
	extension := filepath.Ext(filename)
	return filename[0 : len(filename)-len(extension)]
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
func assembleSaveUrl(base, location string) *url.URL {
	b, _ := url.Parse(base)
	u, _ := url.Parse(strings.Join([]string{"folder", location}, "/"))
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
	// URL here is: IMAGE*/path/file
	// get location and traverse up until root
	var result []byte
	var err error
	var fullPath, partPath, parentFolderId string

	client := &http.Client{}
	metaUrl := assembleMetadataUrl(v.BaseURL, v.Location)
	// get metadata, find parent-id and name
	for {
		result, err = restRequest(client, metaUrl.String(), GET, v.User, v.Pass, map[string]string{}, nil)
		if err != nil {
			return ""
		}
		meta := GHPayload{}
		err = xml.Unmarshal(result, &meta)
		if err != nil {
			return ""
		}
		partPath, parentFolderId = findNameParentField(&meta)
		if parentFolderId == "00000000-0000-0000-0000-000000000000" {
			break
		}
		fullPath = strings.Join([]string{partPath, fullPath}, "/")
		metaUrl = assembleMetadataUrl(v.BaseURL, parentFolderId)
	}
	return strings.Join([]string{"IMAGE*", fullPath, trimExt(filename)}, "")
}

func findNameParentField(meta *GHPayload) (string, string) {
	var parId, name string
	for _, v := range meta.GHField {
		if v.AttrName == "parent-id" {
			parId = v.GHValue.Text
		}
		if v.AttrName == "name" {
			name = v.GHValue.Text
		}
	}
	return name, parId
}

func assembleMetadataUrl(base, folder string) *url.URL {
	b, _ := url.Parse(base)
	u, _ := url.Parse(path.Join("metadata", "folder", folder))
	metaUrl := b.ResolveReference(u)
	return metaUrl
}

func (v *VizGHStorage) HasBaseURL() bool {
	return v.BaseURL != ""
}
