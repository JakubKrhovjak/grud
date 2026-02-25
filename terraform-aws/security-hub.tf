# =============================================================================
# Security Hub — Central Security Dashboard
# =============================================================================
# Aggregates findings from GuardDuty, Inspector, IAM Access Analyzer, etc.
# Runs automated security checks against AWS best practices.
# First 30 days free, then ~$1-3/month for small accounts.
#
# View findings: AWS Console → Security Hub → Findings / Security standards
# =============================================================================

resource "aws_securityhub_account" "main" {}

resource "aws_securityhub_standards_subscription" "aws_best_practices" {
  standards_arn = "arn:aws:securityhub:eu-central-1::standards/aws-foundational-security-best-practices/v/1.0.0"

  depends_on = [aws_securityhub_account.main]
}
