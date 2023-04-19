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
  title         = "nates-cool-flutter-again"
  project_type  = "flutter"
  stack_id      = "osx-xcode-14.2.x-ventura"
  config        = "flutter-config-test-app-both"
}