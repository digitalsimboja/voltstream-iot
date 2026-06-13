"""
AWS Glue streaming job — consumes battery telemetry from Kinesis,
cleans the records, and writes time-series data to Amazon Timestream.

TODO: deploy this script via the Glue module in infra/modules/glue (to be added).
"""

import sys
from awsglue.transforms import *
from awsglue.utils import getResolvedOptions
from awsglue.context import GlueContext
from awsglue.job import Job
from pyspark.context import SparkContext
from pyspark.sql import functions as F
from pyspark.sql.types import StructType, StructField, StringType, DoubleType, TimestampType

args = getResolvedOptions(sys.argv, ["JOB_NAME", "kinesis_stream_arn", "timestream_database", "timestream_table"])

sc = SparkContext()
glueContext = GlueContext(sc)
spark = glueContext.spark_session
job = Job(glueContext)
job.init(args["JOB_NAME"], args)

BATTERY_SCHEMA = StructType([
    StructField("battery_id", StringType(), False),
    StructField("machine_type", StringType(), True),
    StructField("voltage", DoubleType(), True),
    StructField("temperature_celsius", DoubleType(), True),
    StructField("state_of_charge_pct", DoubleType(), True),
    StructField("current_amps", DoubleType(), True),
    StructField("timestamp", TimestampType(), True),
])


def clean_record(df):
    """Drop records with null battery_id or out-of-range values."""
    return df.filter(
        F.col("battery_id").isNotNull()
        & F.col("voltage").between(0, 100)
        & F.col("temperature_celsius").between(-40, 120)
        & F.col("state_of_charge_pct").between(0, 100)
    )


def add_moving_average(df):
    """Append a 10-record rolling average temperature per battery."""
    from pyspark.sql.window import Window
    window = Window.partitionBy("battery_id").orderBy("timestamp").rowsBetween(-9, 0)
    return df.withColumn("temp_moving_avg_10", F.avg("temperature_celsius").over(window))


# TODO: read from Kinesis Data Stream using Glue streaming source
# data_frame = glueContext.create_data_frame_from_catalog(...)

# TODO: parse JSON payload against BATTERY_SCHEMA
# TODO: call clean_record() → add_moving_average() → write to Timestream

job.commit()
