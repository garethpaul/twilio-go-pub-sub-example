resource "google_service_account" "ordinary_processor" {
  account_id = "ordinary-processor"
}

resource "google_service_account" "dlq_processor" {
  account_id = "dlq-processor"
}

resource "google_cloud_run_service" "ordinary_processor" {
  name     = "ordinary-processor"
  location = "us-central1"

  metadata {
    annotations = {
      "run.googleapis.com/ingress" = "internal"
    }
  }

  template {
    spec {

      service_account_name = google_service_account.ordinary_processor.email

      containers {
        image = "gcr.io/${var.project_id}/message-processor"

        ports {
          container_port = 8080
        }

        env {
          name = "TWILIO_ACCOUNT_SID"
          value = var.twilio_account_sid
        }

        env {
          name = "TWILIO_AUTH_TOKEN"
          value = var.twilio_auth_token
        }

        env {
          name = "TWILIO_FROM_NUMBER"
          value = var.twilio_from_number
        }

        env {
          name  = "SUBSCRIPTION_NAME"
          value = local.ordinary_subscription_name
        }

        env {
          name  = "BUCKET_NAME"
          value = google_storage_bucket.processed_ok.name
        }
      }
    }
  }

  depends_on = [
    google_project_service.main
  ]
}

resource "google_cloud_run_service_iam_member" "ordinary_invoker" {
  project  = google_cloud_run_service.ordinary_processor.project
  location = google_cloud_run_service.ordinary_processor.location
  service  = google_cloud_run_service.ordinary_processor.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.ordinary_processor.email}"
}

resource "google_cloud_run_service" "dlq_processor" {
  name     = "dlq-processor"
  location = "us-central1"

  metadata {
    annotations = {
      "run.googleapis.com/ingress" = "internal"
    }
  }
  
  template {
    spec {
      service_account_name = google_service_account.dlq_processor.email

      containers {
        image = "gcr.io/${var.project_id}/message-processor"

        ports {
          container_port = 8080
        }

        env {
          name  = "BUCKET_NAME"
          value = google_storage_bucket.processed_error.name
        }
      }
    }
  }

  depends_on = [
    google_project_service.main
  ]
}

resource "google_cloud_run_service_iam_member" "dlq_invoker" {
  project  = google_cloud_run_service.dlq_processor.project
  location = google_cloud_run_service.dlq_processor.location
  service  = google_cloud_run_service.dlq_processor.name
  role     = "roles/run.invoker"
  member   = "serviceAccount:${google_service_account.dlq_processor.email}"
}
