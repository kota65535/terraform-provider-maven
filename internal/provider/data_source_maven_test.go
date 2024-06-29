package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceMavenArtifact(t *testing.T) {
	td, cwd := setup(t)
	defer tearDown(t, td, cwd)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMavenArtifactMinimalConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccFilesExists("commons-text-1.9.jar", "."),
					resource.TestCheckResourceAttr("data.maven_artifact.minimal", "output_sha", "ba6ac8c2807490944a0a27f6f8e68fb5ed2e80e2"),
				),
			},
			{
				Config: testAccDataSourceMavenArtifactAllConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccFilesExists("commons-text-1.9-javadoc.jar", "out"),
					resource.TestCheckResourceAttr("data.maven_artifact.all", "output_sha", "599bd81a3ceb32ec09c066fb0cc2005e05996f48"),
				),
			},
		},
	})
}

func TestAccDataSourceMavenArtifactSnapshot(t *testing.T) {
	td, cwd := setup(t)
	defer tearDown(t, td, cwd)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceMavenArtifactSnapshotConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccFilesExists("commons-text-1.10.1-SNAPSHOT.jar", "."),
					resource.TestCheckResourceAttrSet("data.maven_artifact.snapshot", "output_sha"),
				),
			},
		},
	})
}

func testAccDataSourceMavenArtifactMinimalConfig() string {
	return fmt.Sprintf(`
	data "maven_artifact" "minimal" {
		group_id    = "org.apache.commons"
		artifact_id = "commons-text"
		version     = "1.9"
	}
	`)
}

func testAccDataSourceMavenArtifactAllConfig() string {
	return fmt.Sprintf(`
	data "maven_artifact" "all" {
		group_id    = "org.apache.commons"
		artifact_id = "commons-text"
		version     = "1.9"
		classifier  = "javadoc"
		output_dir  = "out"
	}
	`)
}

func testAccDataSourceMavenArtifactSnapshotConfig() string {
	return fmt.Sprintf(`
	provider "maven" {
		repository_url = "https://repository.apache.org/content/repositories/snapshots"
	}

	data "maven_artifact" "snapshot" {
		group_id    = "org.apache.commons"
		artifact_id = "commons-text"
		version     = "1.10.1-SNAPSHOT"
	}
	`)
}

func testAccFilesExists(filename string, dir string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := os.Stat(filepath.Join(dir, filename))
		if err != nil {
			return err
		}
		return nil
	}
}
