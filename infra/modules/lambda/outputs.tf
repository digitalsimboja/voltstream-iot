output "function_arn" {
  value = aws_lambda_function.anomaly_detector.arn
}

output "function_name" {
  value = aws_lambda_function.anomaly_detector.function_name
}
