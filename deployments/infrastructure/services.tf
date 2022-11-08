resource "google_project_service" "main" {
  disable_dependent_services = true

  for_each = toset(["storage.googleapis.com", "pubsub.googleapis.com", "run.googleapis.com", "secretmanager.googleapis.com"])

  service = each.key
}

resource "google_project_service_identity" "pubsub" {
  provider = google-beta

  service = "pubsub.googleapis.com"
}

