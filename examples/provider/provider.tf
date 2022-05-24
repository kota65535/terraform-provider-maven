// For maven central repository
provider "maven" {}

// For private maven repository
provider "maven" {
  repository_url = "https://maven.pkg.github.com/kota65535/maven-repository"
  username       = "username"
  password       = "password"
  alias          = "private"
}

