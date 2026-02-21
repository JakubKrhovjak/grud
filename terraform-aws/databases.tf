# =============================================================================
# [5/6] DATABASES & USERS
# =============================================================================
# Connects to the RDS instance from [4/6] and creates application databases
# and users. Uses the cyrilgdn/postgresql provider.
#
# What gets created:
#   - Database "university" — used by student-service
#   - Database "projects"   — used by project-service
#   - User "student_user"   — owner of university DB
#   - User "project_user"   — owner of projects DB
#
# NOTE: This step requires network access from your machine to RDS.
#       RDS is in private subnets, so this only works during terraform apply
#       when running from within the VPC or via VPN/bastion.
#       Alternative: run database init from a K8s Job after cluster is up.
#
# Depends on: [4/6] rds.tf (RDS instance must be running)
# =============================================================================

provider "postgresql" {
  host     = aws_db_instance.main.address
  port     = aws_db_instance.main.port
  username = aws_db_instance.main.username
  password = var.db_master_password
  database = "grud"
  sslmode  = "require"
}

# Databases
resource "postgresql_database" "university" {
  name  = "university"
  owner = postgresql_role.student_user.name
}

resource "postgresql_database" "projects" {
  name  = "projects"
  owner = postgresql_role.project_user.name
}

# Users
resource "postgresql_role" "student_user" {
  name     = "student_user"
  login    = true
  password = var.db_password_student
}

resource "postgresql_role" "project_user" {
  name     = "project_user"
  login    = true
  password = var.db_password_project
}
