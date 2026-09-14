import { ImageResponse } from "next/og";

export const size = {
  width: 32,
  height: 32,
};
export const contentType = "image/png";

export default function Icon() {
  return new ImageResponse(
    (
      <div
        style={{
          fontSize: 20,
          background: "#0a0a0f",
          width: "100%",
          height: "100%",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          color: "#22c55e",
          fontWeight: 700,
          fontFamily: "monospace",
          borderRadius: 6,
          border: "1px solid #1e293b",
        }}
      >
        A
      </div>
    ),
    {
      ...size,
    }
  );
}
