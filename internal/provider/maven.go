package provider

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

const DefaultMavenRepoUrl = "https://repo1.maven.org/maven2/"

const SnapshotVersionSuffix = "-SNAPSHOT"

type Repository struct {
	Url      string
	Username string
	Password string
}

type Artifact struct {
	GroupId    string
	ArtifactId string
	Version    string
	Classifier string
	Extension  string
}

type MavenMetadata struct {
	Timestamp   string `xml:"versioning>snapshot>timestamp"`
	BuildNumber string `xml:"versioning>snapshot>buildNumber"`
}

func NewRepository(url, username, password string) *Repository {
	if url == "" {
		url = DefaultMavenRepoUrl
	}
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	return &Repository{
		Url:      url,
		Username: username,
		Password: password,
	}
}

func NewArtifact(groupId, artifactId, version, classifier, extension string) *Artifact {
	if extension == "" {
		extension = "jar"
	}
	return &Artifact{
		GroupId:    groupId,
		ArtifactId: artifactId,
		Version:    version,
		Classifier: classifier,
		Extension:  extension,
	}
}

func (a *Artifact) ParentUrl(r *Repository) string {
	return r.Url + a.Path()
}

func (a *Artifact) Url(r *Repository, snapshotVersion string) string {
	return a.ParentUrl(r) + a.FileName(snapshotVersion)
}

func (a *Artifact) Path() string {
	return fmt.Sprintf("%s/%s/%s/", strings.Replace(a.GroupId, ".", "/", -1), a.ArtifactId, a.Version)
}

func (a *Artifact) FileName(snapshotVersion string) string {
	version := a.Version
	if snapshotVersion != "" {
		version = version[0 : len(version)-len(SnapshotVersionSuffix)]
		version = fmt.Sprintf("%s-%s", version, snapshotVersion)
	}
	if a.Classifier != "" {
		return fmt.Sprintf("%s-%s-%s.%s", a.ArtifactId, version, a.Classifier, a.Extension)
	} else {
		return fmt.Sprintf("%s-%s.%s", a.ArtifactId, version, a.Extension)
	}
}

func DownloadMavenArtifact(repository *Repository, artifact *Artifact, outputDir string) (string, error) {

	snapshotVersion := ""
	if strings.HasSuffix(artifact.Version, SnapshotVersionSuffix) {
		metadataUrl := artifact.ParentUrl(repository) + "maven-metadata.xml"
		resp, err := httpGet(metadataUrl, repository.Username, repository.Password)
		if err != nil {
			return "", err
		}
		if 400 <= resp.StatusCode {
			return "", errors.New(fmt.Sprintf("status code %d returned. URL: %s", resp.StatusCode, metadataUrl))
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		metadata := MavenMetadata{}
		err = xml.Unmarshal(body, &metadata)
		if err != nil {
			return "", nil
		}
		snapshotVersion = fmt.Sprintf("%s-%s", metadata.Timestamp, metadata.BuildNumber)
	}
	url := artifact.Url(repository, snapshotVersion)
	resp, err := httpGet(url, repository.Username, repository.Password)
	if err != nil {
		return "", err
	}
	if 400 <= resp.StatusCode {
		return "", errors.New(fmt.Sprintf("status code %d returned. URL: %s", resp.StatusCode, url))
	}
	defer resp.Body.Close()

	// default is current directory
	if outputDir == "" {
		outputDir = "."
	}
	// ensure outputDir is directory
	if _, err := os.Stat(outputDir); err != nil {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return "", err
		}
	}

	filepath := path.Join(outputDir, artifact.FileName(""))

	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	return filepath, nil
}

func httpGet(url, user, pwd string) (*http.Response, error) {
	if user != "" && pwd != "" {
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(user, pwd)
		return client.Do(req)
	}

	return http.Get(url)
}
