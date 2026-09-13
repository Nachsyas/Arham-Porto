"use client";

import React, { useMemo } from "react";
import * as THREE from "three";
import { HOLOGRAM_CONSTANTS } from "./hologram.constants";

interface ProjectionPlatformProps {
  isAIActive?: boolean;
}

export default function ProjectionPlatform({ isAIActive = false }: ProjectionPlatformProps) {
  const primaryColor = isAIActive ? HOLOGRAM_CONSTANTS.COLORS.PRIMARY : HOLOGRAM_CONSTANTS.COLORS.HOLOGRAM;
  const faintColor = HOLOGRAM_CONSTANTS.COLORS.HOLOGRAM_FAINT;

  // Build concentric technical rings
  const { outerRing, innerRing, centerRing, radialLines } = useMemo(() => {
    // Ring 1: Outer base circle (radius 1.15)
    const outer = new THREE.RingGeometry(1.12, 1.14, 64);
    // Ring 2: Middle dashed base (radius 0.82)
    const inner = new THREE.RingGeometry(0.81, 0.83, 48);
    // Ring 3: Center focal circle (radius 0.45)
    const center = new THREE.RingGeometry(0.44, 0.45, 32);

    // Crosshair radial ticks
    const points: THREE.Vector3[] = [];
    const tickLen = 0.12;
    // 4 cardinal ticks
    points.push(new THREE.Vector3(-1.18, 0, 0), new THREE.Vector3(-1.18 + tickLen, 0, 0));
    points.push(new THREE.Vector3(1.18, 0, 0), new THREE.Vector3(1.18 - tickLen, 0, 0));
    points.push(new THREE.Vector3(0, 0, -1.18), new THREE.Vector3(0, 0, -1.18 + tickLen));
    points.push(new THREE.Vector3(0, 0, 1.18), new THREE.Vector3(0, 0, 1.18 - tickLen));

    const radialGeom = new THREE.BufferGeometry().setFromPoints(points);

    return {
      outerRing: outer,
      innerRing: inner,
      centerRing: center,
      radialLines: radialGeom,
    };
  }, []);

  const ringMaterial = useMemo(() => {
    return new THREE.MeshBasicMaterial({
      color: primaryColor,
      transparent: true,
      opacity: 0.4,
      side: THREE.DoubleSide,
    });
  }, [primaryColor]);

  const faintMaterial = useMemo(() => {
    return new THREE.MeshBasicMaterial({
      color: faintColor,
      transparent: true,
      opacity: 0.25,
      side: THREE.DoubleSide,
    });
  }, [faintColor]);

  const tickMaterial = useMemo(() => {
    return new THREE.LineBasicMaterial({
      color: primaryColor,
      transparent: true,
      opacity: 0.6,
    });
  }, [primaryColor]);

  React.useEffect(() => {
    return () => {
      outerRing.dispose();
      innerRing.dispose();
      centerRing.dispose();
      radialLines.dispose();
      ringMaterial.dispose();
      faintMaterial.dispose();
      tickMaterial.dispose();
    };
  }, [outerRing, innerRing, centerRing, radialLines, ringMaterial, faintMaterial, tickMaterial]);

  return (
    <group position={[0, -0.88, 0]}>
      {/* Horizontal floor rings (rotated flat on XZ plane) */}
      <mesh geometry={outerRing} material={ringMaterial} rotation={[-Math.PI / 2, 0, 0]} />
      <mesh geometry={innerRing} material={faintMaterial} rotation={[-Math.PI / 2, 0, 0]} />
      <mesh geometry={centerRing} material={ringMaterial} rotation={[-Math.PI / 2, 0, 0]} />

      {/* Cardinal crosshairs on the platform */}
      <lineSegments geometry={radialLines} material={tickMaterial} />

      {/* Subtle upward ambient light column beneath the plinth */}
      <pointLight position={[0, 0.2, 0]} color={primaryColor} intensity={0.45} distance={1.8} />
    </group>
  );
}
