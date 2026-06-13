resource "aws_timestreamwrite_database" "telemetry" {
  database_name = "${var.project}-${var.environment}"

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}

resource "aws_timestreamwrite_table" "battery_metrics" {
  database_name = aws_timestreamwrite_database.telemetry.database_name
  table_name    = "battery-metrics"

  retention_properties {
    memory_store_retention_period_in_hours  = var.memory_retention_hours
    magnetic_store_retention_period_in_days = var.magnetic_retention_days
  }

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}

resource "aws_timestreamwrite_table" "anomalies" {
  database_name = aws_timestreamwrite_database.telemetry.database_name
  table_name    = "anomalies"

  retention_properties {
    memory_store_retention_period_in_hours  = 24
    magnetic_store_retention_period_in_days = 365
  }

  tags = {
    Project     = var.project
    Environment = var.environment
  }
}
