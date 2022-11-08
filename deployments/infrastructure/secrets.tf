resource "google_secret_manager_secret" "sendgrid" {
  secret_id = "sendgrid"
  replication {
    user_managed {
      replicas {
        location = local.region
      }
    }
  }

  depends_on = [
    google_project_service.main
  ]

  
}
