package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"maven": func() (*schema.Provider, error) {
		return New(), nil
	},
}

type ProviderFixture struct {
	RepositoryUrl string
	Username      string
	Password      string
	Env           map[string]string
}

func (r *ProviderFixture) rawConfig() map[string]interface{} {
	rawConfig := map[string]interface{}{}
	if r.RepositoryUrl != "" {
		rawConfig["repository_url"] = r.RepositoryUrl
	}
	if r.Username != "" {
		rawConfig["username"] = r.Username
	}
	if r.Password != "" {
		rawConfig["password"] = r.Password
	}
	return rawConfig
}

func configureProvider(t *testing.T, fixture *ProviderFixture) *Repository {
	for k, v := range fixture.Env {
		_ = os.Setenv(k, v)
	}

	p := New()

	if err := p.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}

	ctx := context.Background()
	diags := p.Configure(ctx, terraform.NewResourceConfigRaw(fixture.rawConfig()))
	if len(diags) > 0 {
		issues := []string{}
		for _, d := range diags {
			issues = append(issues, d.Summary)
		}
		log.Fatalf(strings.Join(issues, ", "))
	}
	return p.Meta().(*Repository)
}

func TestProvider(t *testing.T) {
	repository := configureProvider(t, &ProviderFixture{})

	assert.Equal(t, "https://repo1.maven.org/maven2/", repository.Url)
	assert.Equal(t, "", repository.Username)
	assert.Equal(t, "", repository.Password)
}

func TestProviderWithConfig(t *testing.T) {
	repository := configureProvider(t, &ProviderFixture{
		RepositoryUrl: "foo",
		Username:      "bar",
		Password:      "baz",
	})

	assert.Equal(t, "foo/", repository.Url)
	assert.Equal(t, "bar", repository.Username)
	assert.Equal(t, "baz", repository.Password)
}

func TestProviderWithEnvs(t *testing.T) {
	repository := configureProvider(t, &ProviderFixture{
		Env: map[string]string{
			"MAVEN_REPOSITORY_URL": "hoge",
			"MAVEN_USERNAME":       "piyo",
			"MAVEN_PASSWORD":       "fuga",
		},
	})

	assert.Equal(t, "hoge/", repository.Url)
	assert.Equal(t, "piyo", repository.Username)
	assert.Equal(t, "fuga", repository.Password)
}
