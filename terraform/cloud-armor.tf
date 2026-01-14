# =============================================================================
# Cloud Armor Security Policy
# =============================================================================
# WAF (Web Application Firewall) for protecting the application.
#
# Pricing (Standard tier):
#   - $5/month per policy
#   - $0.75 per million requests
#
# Features:
#   - Rate limiting (DDoS protection)
#   - IP allowlist/denylist
#   - Geo-blocking
#   - Custom rules (SQL injection, XSS patterns)
#
# Docs: https://cloud.google.com/armor/docs/configure-security-policies
# =============================================================================

resource "google_compute_security_policy" "grud_policy" {
  name        = "grud-security-policy"
  description = "Cloud Armor security policy for GRUD application"
  project     = var.project_id

  # Default rule - allow all traffic
  rule {
    action   = "allow"
    priority = "2147483647"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    description = "Default rule - allow all"
  }

  # Rate limiting - prevent DDoS
  rule {
    action   = "throttle"
    priority = "1000"
    match {
      versioned_expr = "SRC_IPS_V1"
      config {
        src_ip_ranges = ["*"]
      }
    }
    description = "Rate limit - 100 requests per minute per IP"
    rate_limit_options {
      conform_action = "allow"
      exceed_action  = "deny(429)"
      rate_limit_threshold {
        count        = 100
        interval_sec = 60
      }
      enforce_on_key = "IP"
    }
  }

  # Block known bad IPs (optional - add IPs as needed)
  # rule {
  #   action   = "deny(403)"
  #   priority = "100"
  #   match {
  #     versioned_expr = "SRC_IPS_V1"
  #     config {
  #       src_ip_ranges = ["1.2.3.4/32"]  # Add malicious IPs here
  #     }
  #   }
  #   description = "Block known malicious IPs"
  # }

  # Block SQL injection attempts (basic pattern)
  rule {
    action   = "deny(403)"
    priority = "2000"
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('sqli-stable')"
      }
    }
    description = "Block SQL injection"
  }

  # Block XSS attempts
  rule {
    action   = "deny(403)"
    priority = "2001"
    match {
      expr {
        expression = "evaluatePreconfiguredExpr('xss-stable')"
      }
    }
    description = "Block XSS attacks"
  }
}

output "cloud_armor_policy_name" {
  value       = google_compute_security_policy.grud_policy.name
  description = "Cloud Armor security policy name for GCPBackendPolicy"
}
