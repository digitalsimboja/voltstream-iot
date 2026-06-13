output "database_name" {
  value = aws_timestreamwrite_database.telemetry.database_name
}

output "metrics_table_name" {
  value = aws_timestreamwrite_table.battery_metrics.table_name
}

output "anomalies_table_name" {
  value = aws_timestreamwrite_table.anomalies.table_name
}
