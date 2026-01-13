# =============================================================================
# SSL Policy for Gateway
# =============================================================================
# Enforces modern TLS settings for HTTPS traffic.
#
# Profile options:
#   - COMPATIBLE: Wide compatibility, supports older clients
#   - MODERN: Recommended, TLS 1.2+, secure cipher suites
#   - RESTRICTED: Most secure, TLS 1.2+, limited cipher suites
#   - CUSTOM: Define your own cipher suites
#
# Docs: https://cloud.google.com/load-balancing/docs/ssl-policies-concepts
# =============================================================================

resource "google_compute_ssl_policy" "grud" {
  name            = "grud-ssl-policy"
  profile         = "MODERN"
  min_tls_version = "TLS_1_2"
}
