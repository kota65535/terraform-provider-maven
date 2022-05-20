package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Params struct {
	RepositoryUrl string
	Username      string
	Password      string
}

func New(params *Params) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"repository_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://repo1.maven.org/maven2",
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Default:   "",
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"maven_package": dataSourceMavenPackage(),
		},
		ConfigureContextFunc: configure(params),
	}
}

func configure(params *Params) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(cxt context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		if params == nil {
			params = &Params{}
		}
		params.RepositoryUrl = d.Get("repository_url").(string)
		params.Username = d.Get("username").(string)
		params.Password = d.Get("password").(string)
		return NewRepository(params.RepositoryUrl, params.Username, params.Password), nil
	}
}
