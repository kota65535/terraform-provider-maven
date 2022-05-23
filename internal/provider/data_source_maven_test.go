package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceMavenArtifactBasic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMavenArtifactBasicConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceExists("data.maven_artifact.basic"),
				),
			},
		},
	})
}

func testAccDataSourceMavenArtifactBasicConfig() string {
	return fmt.Sprintf(`
	data "maven_artifact" "basic" {
		group_id    = "org.apache.commons"
		artifact_id = "commons-text"
		version     = "1.9"
	}
	`)
}

func testAccCheckDataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID set")
		}
		return nil
	}
}
