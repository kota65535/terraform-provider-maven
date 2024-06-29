package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	// Set descriptions to support markdown syntax, this will be used in document generation
	// and the language server.
	schema.DescriptionKind = schema.StringMarkdown

	// Customize the content of descriptions when output. For example you can add defaults on
	// to the exported descriptions if present.
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"repository_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "URL of the maven repository.",
				DefaultFunc: schema.EnvDefaultFunc("MAVEN_REPOSITORY_URL", "https://repo1.maven.org/maven2"),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Username to authenticate against the private maven repository.",
				DefaultFunc: schema.EnvDefaultFunc("MAVEN_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password to authenticate against the private maven repository.",
				DefaultFunc: schema.EnvDefaultFunc("MAVEN_PASSWORD", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"maven_artifact": dataSourceMavenArtifact(),
		},
		ConfigureContextFunc: configure(),
	}
}

func configure() func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(cxt context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		repositoryUrl := d.Get("repository_url").(string)
		username := d.Get("username").(string)
		password := d.Get("password").(string)
		return NewRepository(repositoryUrl, username, password), nil
	}
}
