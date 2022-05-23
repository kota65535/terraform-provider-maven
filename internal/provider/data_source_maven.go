package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"os"
)

func dataSourceMavenArtifact() *schema.Resource {
	return &schema.Resource{
		Description: "Download an artifact from the maven repository.",
		ReadContext: dataSourceMavenArtifactRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group ID",
			},
			"artifact_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Artifact ID",
			},
			"version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Version",
			},
			"classifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Classifier",
			},
			"extension": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "jar",
				Description: "Extension of the artifact file",
			},
			"output_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Path of the directory where the artifact file is located",
			},
			"output_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Path of the artifact file",
			},
			"output_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				ForceNew:    true,
				Description: "Size of the artifact file",
			},
			"output_sha": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "SHA1 checksum of the artifact file",
			},
			"output_base64sha256": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "Base64 Encoded SHA256 checksum of the artifact file",
			},
			"output_md5": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "MD5 of the artifact file",
			},
		},
	}
}

func dataSourceMavenArtifactRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	groupId := d.Get("group_id").(string)
	artifactId := d.Get("artifact_id").(string)
	version := d.Get("version").(string)
	classifier := d.Get("classifier").(string)
	extension := d.Get("extension").(string)
	outputDir := d.Get("output_dir").(string)

	artifact := NewArtifact(groupId, artifactId, version, classifier, extension)
	repository := m.(*Repository)

	filePath, err := DownloadMavenArtifact(repository, artifact, outputDir)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check download is successful
	fi, err := os.Stat(filePath)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("output_path", filePath)
	d.Set("output_size", fi.Size())

	// Calculate hashes
	sha1, base64sha256, md5, err := GenerateHashes(filePath)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("output_sha", sha1)
	d.Set("output_base64sha256", base64sha256)
	d.Set("output_md5", md5)
	d.SetId(d.Get("output_sha").(string))

	return nil
}
