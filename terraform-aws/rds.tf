# =============================================================================
# [4/6] RDS POSTGRESQL
# =============================================================================
# Creates the managed PostgreSQL instance.
#
# What gets created:
#   - DB subnet group (private subnets from VPC)
#   - Security group allowing port 5432 from EKS nodes only
#   - RDS instance: db.t4g.micro, PostgreSQL 15, 20GB gp3, single-AZ
#   - No backups, no deletion protection (dev/test)
#   - Master user: grud_admin (databases + app users created in [5/6])
#
# Depends on: [2/6] vpc.tf (subnets), [3/6] eks.tf (node security group)
# Used by:    [5/6] databases.tf (postgresql provider connects here)
# =============================================================================

resource "aws_db_subnet_group" "main" {
  name       = "${var.cluster_name}-db"
  subnet_ids = module.vpc.private_subnets

  tags = {
    Environment = "test"
    Project     = "grud"
  }
}

resource "aws_security_group" "rds" {
  name_prefix = "${var.cluster_name}-rds-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [module.eks.node_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Environment = "test"
    Project     = "grud"
  }
}

resource "aws_db_instance" "main" {
  identifier = "${var.cluster_name}-postgres"

  engine         = "postgres"
  engine_version = "15"
  instance_class = var.db_instance_class

  allocated_storage = 20
  storage_type      = "gp3"

  db_name  = "grud"
  username = "grud_admin"
  password = var.db_master_password

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  publicly_accessible    = false
  multi_az               = false
  backup_retention_period = 0
  skip_final_snapshot    = true
  deletion_protection    = false

  tags = {
    Environment = "test"
    Project     = "grud"
  }
}
