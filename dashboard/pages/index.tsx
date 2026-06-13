import Head from "next/head";
import FleetMap from "../components/FleetMap";
import TelemetryChart from "../components/TelemetryChart";
import AlertPanel from "../components/AlertPanel";
import { useStream } from "../hooks/useStream";

const WS_URL = process.env.NEXT_PUBLIC_WS_URL ?? "ws://localhost:8080/ws";

export default function Home() {
  const { packets, connected } = useStream(WS_URL, { maxBuffer: 500 });

  return (
    <>
      <Head>
        <title>VoltStream IoT — Fleet Dashboard</title>
      </Head>

      <div className="min-h-screen bg-gray-50 p-6">
        {/* Header */}
        <header className="mb-6 flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-gray-900">VoltStream IoT</h1>
            <p className="text-sm text-gray-500">Industrial Battery Telemetry Platform</p>
          </div>
          <div className="flex items-center gap-2 text-sm">
            <span
              className={`h-2.5 w-2.5 rounded-full ${connected ? "bg-green-500" : "bg-red-500"}`}
            />
            <span className="text-gray-600">{connected ? "Live" : "Disconnected"}</span>
            <span className="ml-2 text-gray-400">{packets.length} packets received</span>
          </div>
        </header>

        {/* Alert panel */}
        <section className="mb-6">
          <AlertPanel packets={packets} />
        </section>

        {/* Fleet map */}
        <section className="mb-6">
          <FleetMap packets={packets} />
        </section>

        {/* Live telemetry chart */}
        <section>
          <TelemetryChart packets={packets} />
        </section>
      </div>
    </>
  );
}
