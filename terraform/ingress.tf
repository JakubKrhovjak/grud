# Static IP for Ingress
# This ensures the external IP doesn't change on redeploy

# Static IP for application Ingress (student-service API)
# prevent_destroy ensures IP is preserved after terraform destroy
resource "google_compute_global_address" "ingress_ip" {
  name        = "grud-ingress-ip"
  description = "Static IP for GRUD application Ingress"

  lifecycle {
    prevent_destroy = true
  }
}

# Static IP for Grafana Ingress
resource "google_compute_global_address" "grafana_ip" {
  name        = "grud-grafana-ip"
  description = "Static IP for Grafana Ingress"

  lifecycle {
    prevent_destroy = true
  }
}
