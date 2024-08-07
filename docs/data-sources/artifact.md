---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "maven_artifact Data Source - terraform-provider-maven"
subcategory: ""
description: |-
  Download an artifact from the maven repository.
---

# maven_artifact (Data Source)

Download an artifact from the maven repository.

## Example Usage

```terraform
// Download apache-commons package from maven central
data "maven_artifact" "commons" {
  group_id    = "org.apache.commons"
  artifact_id = "commons-text"
  version     = "1.9"
  output_dir  = "${path.root}/.terraform/tmp"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `artifact_id` (String) Artifact ID.
- `group_id` (String) Group ID.
- `version` (String) Version.

### Optional

- `classifier` (String) Classifier.
- `extension` (String) Extension of the artifact file. Defaults to `jar`.
- `output_dir` (String) Path of the directory where the artifact file is located.

### Read-Only

- `id` (String) The ID of this resource.
- `output_base64sha256` (String) Base64 Encoded SHA256 checksum of the artifact file.
- `output_md5` (String) MD5 of the artifact file.
- `output_path` (String) Path of the artifact file.
- `output_sha` (String) SHA1 checksum of the artifact file.
- `output_size` (Number) Size of the artifact file.
