import { BatteryPacket } from "../hooks/useStream";

interface Props {
  packets: BatteryPacket[];
}

function statusColor(packet: BatteryPacket): string {
  if (packet.temperature_celsius > 80) return "bg-red-500";
  if (packet.state_of_charge_pct < 10) return "bg-orange-400";
  if (packet.state_of_charge_pct < 25) return "bg-yellow-400";
  return "bg-green-500";
}

function statusLabel(packet: BatteryPacket): string {
  if (packet.temperature_celsius > 80) return "CRITICAL — OVERHEAT";
  if (packet.state_of_charge_pct < 10) return "CRITICAL — LOW SOC";
  if (packet.state_of_charge_pct < 25) return "WARNING";
  return "NORMAL";
}

/**
 * Renders a card grid showing the last-known state of each machine in the fleet.
 * Each card's accent colour reflects the battery health status.
 */
export default function FleetMap({ packets }: Props) {
  // Deduplicate by battery_id — keep the latest packet per machine.
  const latest = new Map<string, BatteryPacket>();
  for (const p of packets) {
    if (!latest.has(p.battery_id)) latest.set(p.battery_id, p);
  }
  const machines = Array.from(latest.values());

  return (
    <div>
      <h3 className="mb-3 text-sm font-semibold text-gray-700">
        Fleet Status — {machines.length} machines
      </h3>
      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-5">
        {machines.map((p) => (
          <div
            key={p.battery_id}
            className="rounded-lg border border-gray-200 bg-white p-3 shadow-sm"
          >
            <div className="flex items-center gap-2">
              <span className={`h-2.5 w-2.5 rounded-full ${statusColor(p)}`} />
              <span className="text-xs font-mono font-semibold text-gray-800">
                {p.battery_id}
              </span>
            </div>
            <p className="mt-1 text-xs text-gray-500">{p.machine_type}</p>
            <p className="mt-2 text-xs text-gray-600">
              {p.temperature_celsius}°C · {p.state_of_charge_pct}% SoC
            </p>
            <p className={`mt-1 text-xs font-semibold ${statusColor(p).replace("bg-", "text-")}`}>
              {statusLabel(p)}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
