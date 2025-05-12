terraform {
  required_providers {
    coder = {
      source = "coder/coder"
    }
  }
}

module "jetbrains_gateway" {}
