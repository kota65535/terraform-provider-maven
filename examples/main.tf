variable github_token {
  type = string
}

terraform {
  required_providers {
    maven = {
      version = "~> 0.0.5"
      source  = "kota65535/maven"
    }
  }
  required_version = ">= 1.1"
}

// Download apache-commons package from maven central
data "maven_package" "commons" {
  group_id    = "org.apache.commons"
  artifact_id = "commons-text"
  version     = "1.9"
}

// Configure another provider for private maven repository
provider "maven" {
  repository_url = "https://maven.pkg.github.com/kota65535/maven-repository"
  username       = "kota65535"
  password       = var.github_token
  alias          = "private"
}

// Download a package from private maven repository
data "maven_package" "private" {
  group_id    = "com.kota65535"
  artifact_id = "hello-maven"
  version     = "0.0.20"
  output_dir  = "out"

  provider = maven.private
}

// Extract files from the package
resource null_resource unzip_package {
  triggers = {
    hash = filebase64sha256(data.maven_package.commons.output_path)
  }
  provisioner "local-exec" {
    command = "unzip ${data.maven_package.commons.output_path} -d ${path.root}/.terraform/tmp"
  }
}

output files {
  value = fileset("${path.root}/.terraform/tmp", "**")
}
