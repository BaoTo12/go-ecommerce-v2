# Titan Commerce Platform - Terraform Configuration
# AWS EKS + Supporting Infrastructure

terraform {
  required_version = ">= 1.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.23"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.11"
    }
  }

  backend "s3" {
    bucket         = "titan-commerce-terraform-state"
    key            = "infrastructure/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "TitanCommerce"
      Environment = var.environment
      ManagedBy   = "Terraform"
    }
  }
}

# ===== VARIABLES =====

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "environment" {
  description = "Environment name"
  type        = string
  default     = "production"
}

variable "cluster_name" {
  description = "EKS cluster name"
  type        = string
  default     = "titan-commerce-cluster"
}

variable "cluster_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.28"
}

variable "cell_count" {
  description = "Number of cells to deploy"
  type        = number
  default     = 10
}

# ===== VPC =====

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.1.2"

  name = "${var.cluster_name}-vpc"
  cidr = "10.0.0.0/16"

  azs             = ["${var.aws_region}a", "${var.aws_region}b", "${var.aws_region}c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24", "10.0.3.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24", "10.0.103.0/24"]

  enable_nat_gateway     = true
  single_nat_gateway     = var.environment != "production"
  enable_dns_hostnames   = true
  enable_dns_support     = true

  public_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                    = 1
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/internal-elb"           = 1
  }
}

# ===== EKS CLUSTER =====

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "19.17.2"

  cluster_name    = var.cluster_name
  cluster_version = var.cluster_version

  vpc_id     = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  cluster_endpoint_public_access = true

  eks_managed_node_groups = {
    # System nodes for critical workloads
    system = {
      name           = "system-nodes"
      instance_types = ["m6i.large"]
      min_size       = 3
      max_size       = 5
      desired_size   = 3

      labels = {
        role = "system"
      }

      taints = [{
        key    = "CriticalAddonsOnly"
        value  = "true"
        effect = "NO_SCHEDULE"
      }]
    }

    # Application nodes
    application = {
      name           = "app-nodes"
      instance_types = ["m6i.xlarge"]
      min_size       = 5
      max_size       = 100
      desired_size   = 10

      labels = {
        role = "application"
      }
    }

    # High-memory nodes for databases
    database = {
      name           = "db-nodes"
      instance_types = ["r6i.2xlarge"]
      min_size       = 3
      max_size       = 10
      desired_size   = 3

      labels = {
        role = "database"
      }

      taints = [{
        key    = "database"
        value  = "true"
        effect = "NO_SCHEDULE"
      }]
    }

    # Flash sale nodes (burstable)
    flashsale = {
      name           = "flashsale-nodes"
      instance_types = ["c6i.4xlarge"]
      min_size       = 0
      max_size       = 50
      desired_size   = 0

      labels = {
        role = "flashsale"
      }
    }
  }

  # Enable cluster addons
  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
    }
    aws-ebs-csi-driver = {
      most_recent = true
    }
  }
}

# ===== RDS POSTGRESQL =====

module "rds" {
  source  = "terraform-aws-modules/rds/aws"
  version = "6.2.0"

  identifier = "${var.cluster_name}-postgres"

  engine               = "postgres"
  engine_version       = "15.4"
  family               = "postgres15"
  major_engine_version = "15"
  instance_class       = "db.r6g.2xlarge"

  allocated_storage     = 100
  max_allocated_storage = 1000

  db_name  = "titancommerce"
  username = "titan_admin"
  port     = 5432

  multi_az               = var.environment == "production"
  db_subnet_group_name   = module.vpc.database_subnet_group_name
  vpc_security_group_ids = [aws_security_group.rds.id]

  backup_retention_period = 7
  skip_final_snapshot     = var.environment != "production"
  deletion_protection     = var.environment == "production"

  performance_insights_enabled = true
  monitoring_interval          = 60

  parameters = [
    {
      name  = "max_connections"
      value = "1000"
    }
  ]
}

resource "aws_security_group" "rds" {
  name_prefix = "${var.cluster_name}-rds-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [module.eks.cluster_security_group_id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# ===== ELASTICACHE REDIS =====

resource "aws_elasticache_replication_group" "redis" {
  replication_group_id = "${var.cluster_name}-redis"
  description          = "Titan Commerce Redis Cluster"

  node_type            = "cache.r6g.large"
  num_cache_clusters   = var.environment == "production" ? 3 : 1
  port                 = 6379
  parameter_group_name = "default.redis7.cluster.on"

  automatic_failover_enabled = var.environment == "production"
  multi_az_enabled           = var.environment == "production"

  subnet_group_name  = aws_elasticache_subnet_group.redis.name
  security_group_ids = [aws_security_group.redis.id]

  at_rest_encryption_enabled = true
  transit_encryption_enabled = true

  snapshot_retention_limit = 7
}

resource "aws_elasticache_subnet_group" "redis" {
  name       = "${var.cluster_name}-redis-subnet"
  subnet_ids = module.vpc.private_subnets
}

resource "aws_security_group" "redis" {
  name_prefix = "${var.cluster_name}-redis-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [module.eks.cluster_security_group_id]
  }
}

# ===== MSK KAFKA =====

resource "aws_msk_cluster" "kafka" {
  cluster_name           = "${var.cluster_name}-kafka"
  kafka_version          = "3.5.1"
  number_of_broker_nodes = var.environment == "production" ? 6 : 3

  broker_node_group_info {
    instance_type   = "kafka.m5.large"
    client_subnets  = module.vpc.private_subnets
    security_groups = [aws_security_group.kafka.id]

    storage_info {
      ebs_storage_info {
        volume_size = 100
      }
    }
  }

  encryption_info {
    encryption_in_transit {
      client_broker = "TLS"
      in_cluster    = true
    }
  }

  configuration_info {
    arn      = aws_msk_configuration.kafka.arn
    revision = aws_msk_configuration.kafka.latest_revision
  }
}

resource "aws_msk_configuration" "kafka" {
  name              = "${var.cluster_name}-kafka-config"
  kafka_versions    = ["3.5.1"]
  
  server_properties = <<PROPERTIES
auto.create.topics.enable=true
delete.topic.enable=true
num.partitions=12
default.replication.factor=3
min.insync.replicas=2
PROPERTIES
}

resource "aws_security_group" "kafka" {
  name_prefix = "${var.cluster_name}-kafka-"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port       = 9092
    to_port         = 9098
    protocol        = "tcp"
    security_groups = [module.eks.cluster_security_group_id]
  }
}

# ===== OUTPUTS =====

output "cluster_endpoint" {
  value = module.eks.cluster_endpoint
}

output "cluster_name" {
  value = module.eks.cluster_name
}

output "rds_endpoint" {
  value = module.rds.db_instance_endpoint
}

output "redis_endpoint" {
  value = aws_elasticache_replication_group.redis.primary_endpoint_address
}

output "kafka_bootstrap_brokers" {
  value = aws_msk_cluster.kafka.bootstrap_brokers_tls
}
