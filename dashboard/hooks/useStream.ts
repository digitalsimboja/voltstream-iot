import { useEffect, useRef, useState } from "react";

export interface BatteryPacket {
  battery_id: string;
  machine_type: string;
  voltage: number;
  temperature_celsius: number;
  state_of_charge_pct: number;
  current_amps: number;
  timestamp: string;
}

interface UseStreamOptions {
  maxBuffer?: number;
}

/**
 * Opens a WebSocket connection to the ingestion API and streams
 * incoming battery telemetry packets into a rolling buffer.
 */
export function useStream(
  url: string,
  { maxBuffer = 500 }: UseStreamOptions = {}
) {
  const [packets, setPackets] = useState<BatteryPacket[]>([]);
  const [connected, setConnected] = useState(false);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);

    ws.onmessage = (event) => {
      try {
        const packet: BatteryPacket = JSON.parse(event.data);
        setPackets((prev) => {
          const next = [packet, ...prev];
          return next.length > maxBuffer ? next.slice(0, maxBuffer) : next;
        });
      } catch {
        // drop malformed frame
      }
    };

    return () => ws.close();
  }, [url, maxBuffer]);

  return { packets, connected };
}
