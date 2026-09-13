"use client";

import React, { Component, useState, useEffect, type ReactNode } from "react";
import dynamic from "next/dynamic";
import HologramFallback from "./HologramFallback";
import { Layers } from "lucide-react";
import type { HologramMode } from "./hologram.constants";

interface ErrorBoundaryProps {
  children: ReactNode;
  fallback: ReactNode;
}

interface ErrorBoundaryState {
  hasError: boolean;
}

class HologramErrorBoundary extends Component<ErrorBoundaryProps, ErrorBoundaryState> {
  constructor(props: ErrorBoundaryProps) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(): ErrorBoundaryState {
    return { hasError: true };
  }

  componentDidCatch(error: Error) {
    // Log non-sensitive diagnostic in development
    if (process.env.NODE_ENV !== "production") {
      console.warn("[HologramErrorBoundary] WebGL scene failed, displaying fallback:", error.message);
    }
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback;
    }
    return this.props.children;
  }
}

// Dynamically import HologramCanvas with SSR disabled so Three.js loads asynchronously
const DynamicHologramCanvas = dynamic(() => import("./HologramCanvas"), {
  ssr: false,
  loading: () => <HologramLoadingState />,
});

function HologramLoadingState() {
  return (
    <div
      className="relative w-full h-full min-h-[380px] max-w-md aspect-square rounded-2xl border border-border bg-gradient-to-b from-surfaceElevated/60 to-surface/40 p-6 flex flex-col items-center justify-center shadow-2xl overflow-hidden"
      aria-label="Loading 3D scene"
    >
      <div className="relative w-64 h-64 flex items-center justify-center">
        <div className="absolute inset-0 rounded-full border border-dashed border-primary/20 animate-spin" />
        <div className="flex h-14 w-14 items-center justify-center rounded-xl bg-canvas border border-primary/30">
          <Layers className="h-7 w-7 text-primary/70 animate-pulse" />
        </div>
      </div>
      <div className="w-full mt-4 pt-3 border-t border-border/60 flex items-center justify-between text-[10px] font-code text-themeText-muted">
        <span>Initializing 3D Pipeline</span>
        <span className="text-primary font-medium">Readying WebGL</span>
      </div>
    </div>
  );
}

interface HologramStageProps {
  mode?: HologramMode;
  isAIActive?: boolean;
  orientationBias?: number;
}

export default function HologramStage({
  mode = "idle",
  isAIActive = false,
  orientationBias = 0,
}: HologramStageProps) {
  const [webGLSupported, setWebGLSupported] = useState<boolean | null>(null);

  useEffect(() => {
    // Check if WebGL is supported in browser
    try {
      if (typeof window === "undefined" || !window.WebGLRenderingContext) {
        setWebGLSupported(false);
        return;
      }
      const canvas = document.createElement("canvas");
      const gl = canvas.getContext("webgl") || canvas.getContext("experimental-webgl");
      setWebGLSupported(Boolean(gl));
    } catch {
      setWebGLSupported(false);
    }
  }, []);

  // During SSR or initial mount, render loading state
  if (webGLSupported === null) {
    return <HologramLoadingState />;
  }

  // If WebGL is not supported on this device, gracefully render static fallback
  if (!webGLSupported) {
    return <HologramFallback reason="WebGL unsupported" />;
  }

  return (
    <HologramErrorBoundary fallback={<HologramFallback reason="Render error" />}>
      <div
        className="relative w-full max-w-md aspect-square rounded-2xl border border-border bg-gradient-to-b from-surfaceElevated/60 to-surface/40 flex flex-col items-center justify-center shadow-2xl overflow-hidden group"
        role="region"
        aria-label="Decorative 3D seated wireframe avatar used as a portfolio visual."
      >
        {/* Technical Viewport Coordinate Markers */}
        <span className="absolute top-3 left-3 text-[10px] font-code text-primary/50 pointer-events-none z-10">┌ SYS:MATRIX</span>
        <span className="absolute top-3 right-3 text-[10px] font-code text-primary/50 pointer-events-none z-10">NODE:AP-01 ┐</span>
        <span className="absolute bottom-10 left-3 text-[10px] font-code text-themeText-mutedSoft pointer-events-none z-10">└ SYS:3D_STAGE</span>
        <span className="absolute bottom-10 right-3 text-[10px] font-code text-themeText-mutedSoft pointer-events-none z-10">SCALE:1.0 ┘</span>

        {/* 3D WebGL Canvas Viewport */}
        <div className="relative w-full h-full flex-1">
          <DynamicHologramCanvas
            mode={mode}
            isAIActive={isAIActive}
            orientationBias={orientationBias}
          />
        </div>

        {/* Technical Telemetry Strip */}
        <div className="w-full px-5 py-2.5 border-t border-border/60 bg-surface/80 backdrop-blur-xs flex items-center justify-between text-[10px] font-code text-themeText-muted z-10">
          <span className="flex items-center gap-1.5">
            <span className="h-1.5 w-1.5 rounded-full bg-status-success animate-pulse" />
            Operational
          </span>
          <span className="text-primary font-medium">3D Hologram Stage</span>
        </div>
      </div>
    </HologramErrorBoundary>
  );
}
