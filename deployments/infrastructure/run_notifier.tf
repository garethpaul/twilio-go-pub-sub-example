resource "google_service_account" "notifier" {
  account_id = "notifier"
}

resource "google_cloud_run_service" "notifier" {
  name     = "notifier"
  location = "us-central1"

  metadata {
    annotations = {
      "run.googleapis.com/ingress" = "internal"
    }
  }

    template {
    spec {

      service_account_name = google_service_account.notifier.email

      containers {
        image = "gcr.io/${var.project_id}/deadletter-notifier"

        ports {
          container_port = 8080
        }

        env {
          name = "NOTIFIER_EMAIL"
          value = var.notifier_email
        }
        
        env {
          name = "SENDGRID_API_KEY"
          value = var.sendgrid
        }
      }
    }
  }
}

resource "google_cloud_run_service_iam_member" "notifier_invoker" {
  project  = google_cloud_run_service.notifier.project
  location = google_cloud_run_service.notifier.location
  service  = google_cloud_run_service.notifier.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.notifier.email}"
}
