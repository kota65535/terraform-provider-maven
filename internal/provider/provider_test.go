package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"maven": func() (*schema.Provider, error) {
		return New(&Params{
			RepositoryUrl: "https://repo1.maven.org/maven2",
			Username:      "",
			Password:      "",
		}), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New(nil).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}
