# main.tf

# Configure GCP provider
provider "google" {
  project = var.project_id
  region  = var.region
}

# Enable required APIs
resource "google_project_service" "services" {
  for_each = toset([
    "cloudrun.googleapis.com",
    "sqladmin.googleapis.com",
    "containerregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "secretmanager.googleapis.com"
  ])
  service = each.key
  disable_on_destroy = false
}

# VPC Network
resource "google_compute_network" "vpc" {
  name                    = "testhero-vpc"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "subnet" {
  name          = "testhero-subnet"
  ip_cidr_range = "10.0.0.0/24"
  network       = google_compute_network.vpc.id
  region        = var.region
}

# Cloud SQL Instance
resource "google_sql_database_instance" "testhero" {
  name             = "testhero-db-instance"
  database_version = "POSTGRES_14"
  region           = var.region
  depends_on       = [google_project_service.services]

  settings {
    tier = var.database_tier
    
    backup_configuration {
      enabled    = true
      start_time = "02:00"  # Daily backup at 2 AM
    }

    ip_configuration {
      ipv4_enabled    = true
      private_network = google_compute_network.vpc.id
      
      authorized_networks {
        name  = "allow-cloud-run"
        value = "0.0.0.0/0"  # In production, restrict this to your application's IP range
      }
    }
  }

  deletion_protection = true
}

# Database and User
resource "google_sql_database" "database" {
  name     = "testhero"
  instance = google_sql_database_instance.testhero.name
}

resource "google_sql_user" "user" {
  name     = "testhero_app"
  instance = google_sql_database_instance.testhero.name
  password = var.database_password
}

# Service Account for Cloud Run
resource "google_service_account" "cloud_run_service_account" {
  account_id   = "testhero-service"
  display_name = "TestHero Service Account"
}

# IAM bindings for the service account
resource "google_project_iam_member" "cloud_run_service_account_bindings" {
  for_each = toset([
    "roles/cloudsql.client",
    "roles/secretmanager.secretAccessor"
  ])
  role    = each.key
  member  = "serviceAccount:${google_service_account.cloud_run_service_account.email}"
  project = var.project_id
}

# Cloud Storage Bucket
resource "google_storage_bucket" "assets" {
  name          = "${var.project_id}-testhero-assets"
  location      = var.region
  force_destroy = true

  uniform_bucket_level_access = true
}

# Secret Manager for sensitive values
resource "google_secret_manager_secret" "openai_api_key" {
  secret_id = "openai-api-key"
  replication {
    automatic = true
  }
}

resource "google_secret_manager_secret_version" "openai_api_key" {
  secret      = google_secret_manager_secret.openai_api_key.id
  secret_data = var.openai_api_key
}

# Cloud Run Service
resource "google_cloud_run_service" "testhero" {
  name     = "testhero"
  location = var.region

  template {
    spec {
      service_account_name = google_service_account.cloud_run_service_account.email
      
      containers {
        image = var.container_image

        env {
          name  = "DB_HOST"
          value = "/cloudsql/${google_sql_database_instance.testhero.connection_name}"
        }
        
        env {
          name  = "DB_USER"
          value = google_sql_user.user.name
        }

        env {
          name  = "DB_NAME"
          value = google_sql_database.database.name
        }

        env {
          name = "DB_PASSWORD"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.database_password.secret_id
              key  = "latest"
            }
          }
        }

        env {
          name = "OPENAI_API_KEY"
          value_from {
            secret_key_ref {
              name = google_secret_manager_secret.openai_api_key.secret_id
              key  = "latest"
            }
          }
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "10"
        "run.googleapis.com/cloudsql-instances" = google_sql_database_instance.testhero.connection_name
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Allow unauthenticated access to Cloud Run service
resource "google_cloud_run_service_iam_member" "public" {
  service  = google_cloud_run_service.testhero.name
  location = google_cloud_run_service.testhero.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}