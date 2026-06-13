import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import { BatteryPacket } from "../hooks/useStream";

interface Props {
  packets: BatteryPacket[];
  batteryId?: string;
}

/**
 * Renders a rolling time-series chart for temperature, SoC, and voltage
 * for a single battery or the most recently active one.
 */
export default function TelemetryChart({ packets, batteryId }: Props) {
  const filtered = batteryId
    ? packets.filter((p) => p.battery_id === batteryId)
    : packets.slice(0, 50);

  const data = [...filtered].reverse().map((p) => ({
    time: new Date(p.timestamp).toLocaleTimeString(),
    temp: p.temperature_celsius,
    soc: p.state_of_charge_pct,
    voltage: p.voltage,
  }));

  return (
    <div className="rounded-lg border border-gray-200 bg-white p-4 shadow-sm">
      <h3 className="mb-3 text-sm font-semibold text-gray-700">
        {batteryId ?? "Live Telemetry"} — Temperature / SoC / Voltage
      </h3>
      <ResponsiveContainer width="100%" height={260}>
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="time" tick={{ fontSize: 11 }} />
          <YAxis tick={{ fontSize: 11 }} />
          <Tooltip />
          <Legend />
          <Line type="monotone" dataKey="temp" stroke="#ef4444" dot={false} name="Temp °C" />
          <Line type="monotone" dataKey="soc" stroke="#3b82f6" dot={false} name="SoC %" />
          <Line type="monotone" dataKey="voltage" stroke="#22c55e" dot={false} name="Voltage V" />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
