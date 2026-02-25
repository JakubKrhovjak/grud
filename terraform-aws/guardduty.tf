# =============================================================================
# GuardDuty — Threat Detection
# =============================================================================
# Continuously monitors AWS account for suspicious activity:
#   - CloudTrail logs (unusual API calls)
#   - VPC Flow Logs (port scanning, malicious IPs)
#   - DNS logs (malware domains)
#   - EKS audit logs (privilege escalation)
#
# First 30 days free, then ~$2-5/month for small clusters.
# Findings visible in AWS Console → GuardDuty → Findings.
# =============================================================================

resource "aws_guardduty_detector" "main" {
  enable                       = true
  finding_publishing_frequency = "FIFTEEN_MINUTES"

  datasources {
    kubernetes {
      audit_logs {
        enable = true
      }
    }
  }

  tags = {
    Environment = "test"
    Project     = "grud"
  }
}
