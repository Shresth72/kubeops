# namespace
resource "kubernetes_namespace" "k8s_cluster" {
  metadata {
    name = "k8s_ns_by_tf"
  }
}

# deployment
resource "kubernetes_deployment" "k8s_deployment" {
  metadata {
    name = "terraform_example"
    labels = {
      test = "MyApp"
    }
    namespace = "k8s_ns_by_tf"
  }

  spec {
    replicas = 2
    selector {
      match_labels = {
        test = "MyApp"
      }
    }

    template {
      metadata {
        labels = {
          test = "MyApp"
        }
      }
      spec {
        container {
          image = "nginx:1.7.9"
          name  = "example"
          resources {
            limits = {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "50Mi"
            }
          }
        }
      }
    }
  }
}
