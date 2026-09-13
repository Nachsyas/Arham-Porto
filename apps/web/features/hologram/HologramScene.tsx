"use client";

import React, { useRef, useMemo } from "react";
import * as THREE from "three";
import { useFrame, useThree } from "@react-three/fiber";
import HologramModel from "./HologramModel";
import ProjectionPlatform from "./ProjectionPlatform";
import HologramParticles from "./HologramParticles";
import { HOLOGRAM_CONSTANTS, type HologramMode } from "./hologram.constants";

interface HologramSceneProps {
  targetYaw: number;
  targetPitch: number;
  isReducedMotion: boolean;
  isMobile: boolean;
  mode?: HologramMode;
  isAIActive?: boolean;
  orientationBias?: number;
}

export default function HologramScene({
  targetYaw,
  targetPitch,
  isReducedMotion,
  isMobile,
  mode = "idle",
  isAIActive = false,
  orientationBias = 0,
}: HologramSceneProps) {
  const modelGroupRef = useRef<THREE.Group>(null);
  const scanPlaneRef = useRef<THREE.Mesh>(null);
  const { size } = useThree();

  // Scan line geometry & material
  const { scanGeometry, scanMaterial } = useMemo(() => {
    const geom = new THREE.PlaneGeometry(1.6, 0.02);
    const mat = new THREE.MeshBasicMaterial({
      color: HOLOGRAM_CONSTANTS.COLORS.HOLOGRAM,
      transparent: true,
      opacity: 0.35,
      side: THREE.DoubleSide,
    });
    return { scanGeometry: geom, scanMaterial: mat };
  }, []);

  React.useEffect(() => {
    return () => {
      scanGeometry.dispose();
      scanMaterial.dispose();
    };
  }, [scanGeometry, scanMaterial]);

  // Main Render Loop Interpolation
  useFrame((state, delta) => {
    if (!modelGroupRef.current) return;

    // 1. Pointer Parallax (Desktop only, subtle)
    let pointerYaw = 0;
    let pointerPitch = 0;
    if (!isMobile && !isReducedMotion) {
      pointerYaw = state.pointer.x * 0.04;
      pointerPitch = -state.pointer.y * 0.02;
    }

    // 2. Compute final target rotations
    const finalYaw = isReducedMotion ? 0 : targetYaw + pointerYaw + orientationBias;
    const finalPitch = isReducedMotion ? 0 : targetPitch + pointerPitch;

    // 3. Smooth Damping with delta
    const damping = HOLOGRAM_CONSTANTS.ROTATION.DAMPING_FACTOR;
    modelGroupRef.current.rotation.y = THREE.MathUtils.damp(
      modelGroupRef.current.rotation.y,
      finalYaw,
      damping,
      delta
    );
    modelGroupRef.current.rotation.x = THREE.MathUtils.damp(
      modelGroupRef.current.rotation.x,
      finalPitch,
      damping,
      delta
    );

    // 4. Subtle Idle Breathing Motion
    if (!isReducedMotion) {
      const breathing = Math.sin(state.clock.elapsedTime * 1.6) * 0.012;
      modelGroupRef.current.position.y = breathing;
    } else {
      modelGroupRef.current.position.y = 0;
    }

    // 5. Subtle Vertical Hologram Scanning Plane
    if (scanPlaneRef.current && !isReducedMotion) {
      const scanCycle = (state.clock.elapsedTime * 0.4) % 1; // 0 to 1 cycle
      scanPlaneRef.current.position.y = -0.7 + scanCycle * 1.6;
      // Fade out at ends
      const edgeFade = Math.sin(scanCycle * Math.PI);
      scanMaterial.opacity = 0.15 + edgeFade * 0.25;
    }
  });

  return (
    <>
      {/* Ambient & Rim Lighting */}
      <ambientLight color="#07111C" intensity={1.2} />
      <directionalLight position={[2, 3, 2]} color={HOLOGRAM_CONSTANTS.COLORS.HOLOGRAM} intensity={0.6} />
      <directionalLight position={[-2, -1, -2]} color={HOLOGRAM_CONSTANTS.COLORS.PRIMARY} intensity={0.3} />

      {/* Main Animated Group (Hologram + Pedestal) */}
      <group ref={modelGroupRef} position={[0, 0, 0]}>
        <HologramModel mode={mode} isAIActive={isAIActive} />
      </group>

      {/* Subtle Horizontal Scanning Plane */}
      {!isReducedMotion && (
        <mesh
          ref={scanPlaneRef}
          geometry={scanGeometry}
          material={scanMaterial}
          rotation={[-Math.PI / 2, 0, 0]}
        />
      )}

      {/* Ground Projection Platform */}
      <ProjectionPlatform isAIActive={isAIActive} />

      {/* Sparse Ambient Depth Particles */}
      <HologramParticles
        count={isMobile ? HOLOGRAM_CONSTANTS.PERFORMANCE.PARTICLES_MOBILE : HOLOGRAM_CONSTANTS.PERFORMANCE.PARTICLES_DESKTOP}
        isReducedMotion={isReducedMotion}
      />
    </>
  );
}
