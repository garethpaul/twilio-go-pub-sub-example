terraform {
  backend "gcs" {
    bucket = "tf-state-gjones-webinar"
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 4.19.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 4.19.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.1.3"
    }
  }
}

provider "google" {
  project = var.project_id
}

provider "google-beta" {
  project = var.project_id
}
