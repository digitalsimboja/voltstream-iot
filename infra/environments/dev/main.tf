provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "voltstream"
      Environment = "dev"
      ManagedBy   = "terraform"
    }
  }
}

module "vpc" {
  source             = "../../modules/vpc"
  project            = "voltstream"
  environment        = "dev"
  vpc_cidr           = var.vpc_cidr
  availability_zones = var.availability_zones
}

module "kinesis" {
  source      = "../../modules/kinesis"
  project     = "voltstream"
  environment = "dev"
  shard_count = 1
}

module "iam" {
  source             = "../../modules/iam"
  project            = "voltstream"
  environment        = "dev"
  kinesis_stream_arn = module.kinesis.stream_arn
}

module "lambda" {
  source             = "../../modules/lambda"
  project            = "voltstream"
  environment        = "dev"
  execution_role_arn = module.iam.lambda_exec_role_arn
  kinesis_stream_arn = module.kinesis.stream_arn
}

module "timestream" {
  source      = "../../modules/timestream"
  project     = "voltstream"
  environment = "dev"
}
