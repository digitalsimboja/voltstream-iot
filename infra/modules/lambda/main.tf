resource "aws_lambda_function" "anomaly_detector" {
  function_name = "${var.project}-${var.environment}-anomaly-detector"
  role          = var.execution_role_arn
  handler       = "bootstrap"
  runtime       = "provided.al2023" # Go Lambda uses custom runtime
  filename      = var.lambda_zip_path
  timeout       = 30
  memory_size   = 128

  environment {
    variables = {
      ENVIRONMENT = var.environment
      # TODO: add TIMESTREAM_DATABASE and TIMESTREAM_TABLE vars
      # TODO: add SNS_ALERT_TOPIC_ARN for anomaly notifications
    }
  }

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}

resource "aws_lambda_event_source_mapping" "kinesis" {
  event_source_arn               = var.kinesis_stream_arn
  function_name                  = aws_lambda_function.anomaly_detector.arn
  starting_position              = "LATEST"
  batch_size                     = 100
  bisect_batch_on_function_error = true

  # TODO: configure destination for failed batch records (SQS DLQ)
}
