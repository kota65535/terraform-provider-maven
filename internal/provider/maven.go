package provider

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

const DefaultMavenRepoUrl = "https://repo1.maven.org/maven2/"

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

func (a *Artifact) Url(r *Repository) string {
	return r.Url + a.Path() + a.FileName()
}

func (a *Artifact) Path() string {
	return fmt.Sprintf("%s/%s/%s/", strings.Replace(a.GroupId, ".", "/", -1), a.ArtifactId, a.Version)
}

func (a *Artifact) FileName() string {
	if a.Classifier != "" {
		return fmt.Sprintf("%s-%s-%s.%s", a.ArtifactId, a.Version, a.Classifier, a.Extension)
	} else {
		return fmt.Sprintf("%s-%s.%s", a.ArtifactId, a.Version, a.Extension)
	}
}

func DownloadMavenPackage(repository *Repository, artifact *Artifact, outputDir string) (string, error) {

	url := artifact.Url(repository)
	resp, err := httpGet(artifact.Url(repository), repository.Username, repository.Password)
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

	filepath := path.Join(outputDir, artifact.FileName())

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
