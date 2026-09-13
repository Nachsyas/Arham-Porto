"use client";

import React, { Suspense } from "react";
import { Canvas } from "@react-three/fiber";
import HologramScene from "./HologramScene";
import { useHologramScroll } from "./useHologramScroll";
import { HOLOGRAM_CONSTANTS, type HologramMode } from "./hologram.constants";

interface HologramCanvasProps {
  mode?: HologramMode;
  isAIActive?: boolean;
  orientationBias?: number;
}

export default function HologramCanvas({
  mode = "idle",
  isAIActive = false,
  orientationBias = 0,
}: HologramCanvasProps) {
  const { targetYaw, targetPitch, isReducedMotion, isMobile } = useHologramScroll();

  return (
    <div className="w-full h-full relative" aria-hidden="true">
      <Canvas
        camera={{
          fov: HOLOGRAM_CONSTANTS.CAMERA.FOV,
          position: HOLOGRAM_CONSTANTS.CAMERA.POSITION,
        }}
        dpr={isMobile ? HOLOGRAM_CONSTANTS.PERFORMANCE.DPR_MOBILE : HOLOGRAM_CONSTANTS.PERFORMANCE.DPR_DESKTOP}
        gl={{
          antialias: true,
          alpha: true,
          powerPreference: "high-performance",
        }}
        className="w-full h-full"
      >
        <Suspense fallback={null}>
          <HologramScene
            targetYaw={targetYaw}
            targetPitch={targetPitch}
            isReducedMotion={isReducedMotion}
            isMobile={isMobile}
            mode={mode}
            isAIActive={isAIActive}
            orientationBias={orientationBias}
          />
        </Suspense>
      </Canvas>
    </div>
  );
}
