terraform {
  required_providers {
    bitrise = {
      source  = "terraform.local/local/bitrise"
      version = "1.0.0"
    }
  }
}

provider "bitrise" {}

resource "bitrise_app" "app" {
  token         = "EZgewzA9KET4uj4cFqoadeLiHwBMKV4orgmZ7kd3AGy_yiMKGBPt050u7KT7fFRd7otH3KGuDKBeftVj0pCxkw"
  repo_url      = "https://github.com/pgdevelopers/nates_bitrise_provider_app.git"
  git_repo_slug = "nates_bitrise_provider_app"
  title         = "nates-cool-sanity-check"
  project_type  = "ios"
  stack_id      = "osx-xcode-13.2.x"
  config        = "default-ios-config"
}