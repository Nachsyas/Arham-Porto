"use client";

import React, { useMemo, useRef } from "react";
import * as THREE from "three";
import { useFrame } from "@react-three/fiber";
import { HOLOGRAM_CONSTANTS } from "./hologram.constants";

interface HologramParticlesProps {
  count?: number;
  isReducedMotion?: boolean;
}

export default function HologramParticles({
  count = HOLOGRAM_CONSTANTS.PERFORMANCE.PARTICLES_DESKTOP,
  isReducedMotion = false,
}: HologramParticlesProps) {
  const pointsRef = useRef<THREE.Points>(null);

  // Generate sparse particle distribution in a cylindrical volume around hologram
  const [positions, initialY] = useMemo(() => {
    const pos = new Float32Array(count * 3);
    const initY = new Float32Array(count);

    for (let i = 0; i < count; i++) {
      const radius = 0.35 + Math.random() * 0.95;
      const angle = Math.random() * Math.PI * 2;
      const y = -0.8 + Math.random() * 1.8;

      pos[i * 3] = Math.cos(angle) * radius;
      pos[i * 3 + 1] = y;
      pos[i * 3 + 2] = Math.sin(angle) * radius;

      initY[i] = y;
    }

    return [pos, initY];
  }, [count]);

  const geometry = useMemo(() => {
    const geom = new THREE.BufferGeometry();
    geom.setAttribute("position", new THREE.BufferAttribute(positions, 3));
    return geom;
  }, [positions]);

  const material = useMemo(() => {
    return new THREE.PointsMaterial({
      color: HOLOGRAM_CONSTANTS.COLORS.HOLOGRAM,
      size: 0.025,
      transparent: true,
      opacity: 0.45,
      blending: THREE.AdditiveBlending,
    });
  }, []);

  // Subtle upward drift in useFrame
  useFrame((state) => {
    if (isReducedMotion || !pointsRef.current) return;

    const positionAttr = pointsRef.current.geometry.attributes.position as THREE.BufferAttribute;
    const time = state.clock.elapsedTime * 0.12;

    for (let i = 0; i < count; i++) {
      let currentY = positionAttr.getY(i);
      currentY += 0.0015; // Slow ambient upward drift
      if (currentY > 1.2) {
        currentY = -0.8;
      }
      positionAttr.setY(i, currentY);
    }

    positionAttr.needsUpdate = true;
  });

  React.useEffect(() => {
    return () => {
      geometry.dispose();
      material.dispose();
    };
  }, [geometry, material]);

  return <points ref={pointsRef} geometry={geometry} material={material} />;
}
