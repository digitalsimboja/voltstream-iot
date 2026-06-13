resource "aws_kinesis_stream" "telemetry" {
  name             = "${var.project}-${var.environment}-telemetry"
  shard_count      = var.shard_count
  retention_period = var.retention_hours

  stream_mode_details {
    stream_mode = "PROVISIONED"
    # TODO: switch to ON_DEMAND for auto-scaling at higher throughput
  }

  encryption_type = "KMS"
  kms_key_id      = "alias/aws/kinesis"

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}
