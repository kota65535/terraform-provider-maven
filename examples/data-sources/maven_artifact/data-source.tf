data "maven_artifact" "commons" {
  group_id    = "org.apache.commons"
  artifact_id = "commons-text"
  version     = "1.9"
  output_dir  = "${path.root}/.terraform/tmp"
}
