import { BatteryPacket } from "../hooks/useStream";

interface Props {
  packets: BatteryPacket[];
}

interface Alert {
  battery_id: string;
  machine_type: string;
  reason: string;
  timestamp: string;
  severity: "critical" | "warning";
}

function extractAlerts(packets: BatteryPacket[]): Alert[] {
  const alerts: Alert[] = [];
  for (const p of packets) {
    if (p.temperature_celsius > 80) {
      alerts.push({
        battery_id: p.battery_id,
        machine_type: p.machine_type,
        reason: `Temperature ${p.temperature_celsius}°C exceeds 80°C threshold`,
        timestamp: p.timestamp,
        severity: "critical",
      });
    } else if (p.state_of_charge_pct < 10) {
      alerts.push({
        battery_id: p.battery_id,
        machine_type: p.machine_type,
        reason: `SoC ${p.state_of_charge_pct}% below 10% minimum`,
        timestamp: p.timestamp,
        severity: "critical",
      });
    }
  }
  return alerts.slice(0, 20);
}

/**
 * Renders a live list of the most recent anomaly alerts derived from the stream.
 */
export default function AlertPanel({ packets }: Props) {
  const alerts = extractAlerts(packets);

  if (alerts.length === 0) {
    return (
      <div className="rounded-lg border border-green-200 bg-green-50 p-4 text-sm text-green-700">
        All systems nominal — no active alerts.
      </div>
    );
  }

  return (
    <div className="rounded-lg border border-gray-200 bg-white shadow-sm">
      <h3 className="border-b border-gray-100 px-4 py-3 text-sm font-semibold text-gray-700">
        Active Alerts ({alerts.length})
      </h3>
      <ul className="divide-y divide-gray-100">
        {alerts.map((a, i) => (
          <li key={i} className="flex items-start gap-3 px-4 py-3">
            <span
              className={`mt-0.5 h-2 w-2 rounded-full flex-shrink-0 ${
                a.severity === "critical" ? "bg-red-500" : "bg-yellow-400"
              }`}
            />
            <div className="min-w-0">
              <p className="text-xs font-semibold text-gray-800">
                {a.battery_id} · {a.machine_type}
              </p>
              <p className="text-xs text-gray-600">{a.reason}</p>
              <p className="mt-0.5 text-xs text-gray-400">
                {new Date(a.timestamp).toLocaleTimeString()}
              </p>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
