"use client";

import React from "react";
import { Layers } from "lucide-react";

interface HologramFallbackProps {
  reason?: string;
}

export default function HologramFallback({ reason }: HologramFallbackProps) {
  return (
    <div
      className="relative w-full h-full min-h-[380px] max-w-md aspect-square rounded-2xl border border-border bg-gradient-to-b from-surfaceElevated/60 to-surface/40 p-6 flex flex-col items-center justify-center shadow-2xl overflow-hidden"
      role="img"
      aria-label="Static 3D seated wireframe fallback placeholder used as a portfolio visual."
    >
      {/* Technical Viewport Coordinate Markers */}
      <span className="absolute top-3 left-3 text-[10px] font-code text-primary/50">┌ SYS:STATIC</span>
      <span className="absolute top-3 right-3 text-[10px] font-code text-primary/50">NODE:AP-01 ┐</span>
      <span className="absolute bottom-3 left-3 text-[10px] font-code text-themeText-mutedSoft">└ WEBGL:FALLBACK</span>
      <span className="absolute bottom-3 right-3 text-[10px] font-code text-themeText-mutedSoft">MODE:2D_VECTOR ┘</span>

      {/* Static Vector Wireframe Geometric Seated Silhouette */}
      <div className="relative w-64 h-64 flex items-center justify-center">
        {/* Concentric rings */}
        <div className="absolute inset-0 rounded-full border border-dashed border-primary/30" />
        <div className="absolute inset-8 rounded-full border border-hologram/30 shadow-[0_0_20px_rgba(97,216,255,0.15)]" />
        <div className="absolute inset-16 rounded-full border border-primary/40 bg-surfaceStrong/40" />

        {/* Central Geometric Icon */}
        <div className="relative z-10 flex flex-col items-center justify-center p-6 text-center space-y-3">
          <div className="flex h-14 w-14 items-center justify-center rounded-xl bg-canvas border border-primary/50 shadow-[0_0_24px_rgba(38,184,255,0.35)]">
            <Layers className="h-7 w-7 text-primary" />
          </div>
          <div className="space-y-0.5">
            <p className="text-xs font-code font-bold tracking-wider text-hologram">
              SYSTEM ARCHITECTURE
            </p>
            <p className="text-[11px] font-body text-themeText-muted">
              Full-Stack & Backend Systems
            </p>
          </div>
        </div>
      </div>

      {/* Technical Telemetry Strip */}
      <div className="w-full mt-4 pt-3 border-t border-border/60 flex items-center justify-between text-[10px] font-code text-themeText-muted">
        <span className="flex items-center gap-1.5">
          <span className="h-1.5 w-1.5 rounded-full bg-status-success" />
          Static Mode
        </span>
        <span className="text-primary font-medium">Standard Render</span>
      </div>
    </div>
  );
}
