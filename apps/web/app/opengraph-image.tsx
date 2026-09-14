import { ImageResponse } from "next/og";

export const alt = "Nachsyas Arham Mumtaz Nashohi — Software Engineer Portfolio";
export const size = {
  width: 1200,
  height: 630,
};
export const contentType = "image/png";

export default function Image() {
  return new ImageResponse(
    (
      <div
        style={{
          background: "#0a0a0f",
          width: "100%",
          height: "100%",
          display: "flex",
          flexDirection: "column",
          justifyContent: "space-between",
          padding: 80,
          fontFamily: "sans-serif",
          border: "2px solid #1e293b",
        }}
      >
        <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
          <div
            style={{
              width: 12,
              height: 12,
              borderRadius: "50%",
              background: "#22c55e",
            }}
          />
          <span style={{ color: "#22c55e", fontSize: 20, fontWeight: 600, letterSpacing: "0.1em" }}>
            PORTFOLIO // EVIDENCE-FIRST SYSTEM
          </span>
        </div>

        <div style={{ display: "flex", flexDirection: "column", gap: 20 }}>
          <h1 style={{ color: "#f8fafc", fontSize: 56, fontWeight: 800, margin: 0, lineHeight: 1.1 }}>
            Nachsyas Arham Mumtaz Nashohi
          </h1>
          <p style={{ color: "#94a3b8", fontSize: 28, margin: 0 }}>
            Software Engineer — Systems, Backend Architecture & Grounded AI
          </p>
        </div>

        <div
          style={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            borderTop: "1px solid #1e293b",
            paddingTop: 30,
          }}
        >
          <span style={{ color: "#64748b", fontSize: 20 }}>Arham Porto · Ask Arham Copilot</span>
          <span style={{ color: "#22c55e", fontSize: 20, fontFamily: "monospace" }}>
            github.com/Nachsyas
          </span>
        </div>
      </div>
    ),
    {
      ...size,
    }
  );
}
