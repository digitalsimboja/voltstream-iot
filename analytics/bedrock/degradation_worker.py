"""
Async degradation analysis worker — queries Amazon Timestream for per-battery
SoC trends and uses AWS Bedrock to classify degradation risk.

Run as a scheduled Lambda or ECS task.
TODO: wire into an EventBridge schedule via infra/modules/scheduler (to be added).
"""

import boto3
import json
import os
from datetime import datetime, timezone

TIMESTREAM_DATABASE = os.environ.get("TIMESTREAM_DATABASE", "voltstream-dev")
TIMESTREAM_TABLE = os.environ.get("TIMESTREAM_TABLE", "battery-metrics")
BEDROCK_MODEL_ID = os.environ.get("BEDROCK_MODEL_ID", "amazon.titan-text-express-v1")
SNS_TOPIC_ARN = os.environ.get("SNS_TOPIC_ARN", "")

query_client = boto3.client("timestream-query", region_name="eu-north-1")
bedrock_client = boto3.client("bedrock-runtime", region_name="eu-north-1")
sns_client = boto3.client("sns", region_name="eu-north-1")


def fetch_soc_trend(battery_id: str, hours: int = 6) -> list[dict]:
    """Query Timestream for recent SoC readings for a specific battery."""
    query = f"""
        SELECT battery_id, time, measure_value::double AS soc_pct
        FROM "{TIMESTREAM_DATABASE}"."{TIMESTREAM_TABLE}"
        WHERE measure_name = 'state_of_charge_pct'
          AND battery_id = '{battery_id}'
          AND time > ago({hours}h)
        ORDER BY time ASC
    """
    # TODO: paginate results for batteries with dense history
    response = query_client.query(QueryString=query)
    return response.get("Rows", [])


def evaluate_degradation(battery_id: str, trend: list[dict]) -> dict:
    """Ask Bedrock to classify the SoC trend as normal, degrading, or critical."""
    prompt = f"""
You are an industrial battery health expert. Analyse this State of Charge (SoC) trend
for battery {battery_id} and classify its degradation risk.

SoC readings (oldest → newest): {json.dumps(trend)}

Respond with JSON only:
{{
  "battery_id": "{battery_id}",
  "risk_level": "normal | degrading | critical",
  "summary": "one sentence explanation",
  "recommended_action": "one sentence"
}}
"""
    # TODO: switch to Claude claude-haiku-4-5-20251001 for lower latency and cost
    response = bedrock_client.invoke_model(
        modelId=BEDROCK_MODEL_ID,
        body=json.dumps({"inputText": prompt}),
        contentType="application/json",
        accept="application/json",
    )
    result = json.loads(response["body"].read())
    return json.loads(result.get("results", [{}])[0].get("outputText", "{}"))


def publish_alert(assessment: dict) -> None:
    """Publish a degradation alert to SNS if risk is degrading or critical."""
    if assessment.get("risk_level") not in ("degrading", "critical"):
        return
    if not SNS_TOPIC_ARN:
        return
    sns_client.publish(
        TopicArn=SNS_TOPIC_ARN,
        Subject=f"[VoltStream] Battery degradation alert — {assessment['battery_id']}",
        Message=json.dumps(assessment, indent=2),
    )


def handler(event: dict, context) -> dict:
    """Lambda / ECS task entry point. Expects event = {"battery_ids": [...]}"""
    battery_ids = event.get("battery_ids", [])
    results = []

    for bid in battery_ids:
        trend = fetch_soc_trend(bid)
        if not trend:
            continue
        assessment = evaluate_degradation(bid, trend)
        publish_alert(assessment)
        results.append(assessment)

    return {"assessed": len(results), "timestamp": datetime.now(timezone.utc).isoformat()}
