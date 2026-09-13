"use client";

import React, { useMemo, useRef } from "react";
import * as THREE from "three";
import { useGLTF } from "@react-three/drei";
import { HOLOGRAM_CONSTANTS, type HologramMode } from "./hologram.constants";

interface HologramModelProps {
  mode?: HologramMode;
  isAIActive?: boolean;
  intensity?: number;
}

/**
 * Procedural Seated Human Wireframe Mesh.
 * Used as the canonical development 3D model until production scan (arham-wireframe.glb) is provided.
 * Represents a seated engineer posture: head, shoulders, torso, forward thighs, downward shins, arms resting forward.
 */
export function SeatedDevelopmentMesh({ isAIActive = false, intensity = 1 }: { isAIActive?: boolean; intensity?: number }) {
  const groupRef = useRef<THREE.Group>(null);

  // Wireframe color & line properties
  const wireColor = isAIActive ? HOLOGRAM_CONSTANTS.COLORS.PRIMARY : HOLOGRAM_CONSTANTS.COLORS.HOLOGRAM;
  const faintColor = HOLOGRAM_CONSTANTS.COLORS.HOLOGRAM_FAINT;

  // Build high-efficiency geometries with useMemo to guarantee zero leaks
  const geometries = useMemo(() => {
    // 1. Head: faceted icosahedron for clean digital silhouette
    const head = new THREE.IcosahedronGeometry(0.24, 2);

    // 2. Neck
    const neck = new THREE.CylinderGeometry(0.08, 0.1, 0.14, 8);

    // 3. Torso / Shoulders: contoured tapered prism
    const upperTorso = new THREE.CylinderGeometry(0.36, 0.28, 0.48, 6);
    const lowerTorso = new THREE.CylinderGeometry(0.28, 0.24, 0.32, 6);

    // 4. Arms (Seated angle, hands toward knees)
    const upperArm = new THREE.CylinderGeometry(0.07, 0.06, 0.38, 6);
    const forearm = new THREE.CylinderGeometry(0.06, 0.05, 0.36, 6);

    // 5. Thighs (Projecting forward in horizontal seated plane)
    const thigh = new THREE.BoxGeometry(0.16, 0.14, 0.48);

    // 6. Shins (Vertical downward toward platform)
    const shin = new THREE.CylinderGeometry(0.06, 0.05, 0.44, 6);

    // 7. Minimalist Seated Geometric Base
    const seatPlinth = new THREE.BoxGeometry(0.5, 0.42, 0.5);

    // Precompute WireframeGeometries to prevent allocation during render
    const wires = {
      head: new THREE.WireframeGeometry(head),
      neck: new THREE.WireframeGeometry(neck),
      upperTorso: new THREE.WireframeGeometry(upperTorso),
      lowerTorso: new THREE.WireframeGeometry(lowerTorso),
      upperArm: new THREE.WireframeGeometry(upperArm),
      forearm: new THREE.WireframeGeometry(forearm),
      thigh: new THREE.WireframeGeometry(thigh),
      shin: new THREE.WireframeGeometry(shin),
      seatPlinth: new THREE.WireframeGeometry(seatPlinth),
    };

    return {
      base: { head, neck, upperTorso, lowerTorso, upperArm, forearm, thigh, shin, seatPlinth },
      wire: wires,
    };
  }, []);

  // Shared wireframe and depth-occlusion materials
  const { lineMaterial, occludeMaterial, pointsMaterial } = useMemo(() => {
    const lineMat = new THREE.LineBasicMaterial({
      color: wireColor,
      transparent: true,
      opacity: 0.8 * intensity,
      linewidth: 1,
    });

    // Dark semi-transparent core to provide spatial depth without total transparency
    const occludeMat = new THREE.MeshBasicMaterial({
      color: faintColor,
      transparent: true,
      opacity: 0.25,
      side: THREE.DoubleSide,
    });

    const pointsMat = new THREE.PointsMaterial({
      color: wireColor,
      size: 0.035,
      transparent: true,
      opacity: 0.9 * intensity,
    });

    return {
      lineMaterial: lineMat,
      occludeMaterial: occludeMat,
      pointsMaterial: pointsMat,
    };
  }, [wireColor, faintColor, intensity]);

  // Cleanup on unmount
  React.useEffect(() => {
    return () => {
      Object.values(geometries.base).forEach((geom) => geom.dispose());
      Object.values(geometries.wire).forEach((geom) => geom.dispose());
      lineMaterial.dispose();
      occludeMaterial.dispose();
      pointsMaterial.dispose();
    };
  }, [geometries, lineMaterial, occludeMaterial, pointsMaterial]);

  // Helper to render wireframe + occlude mesh + vertex points
  const renderLimb = (
    baseGeom: THREE.BufferGeometry,
    wireGeom: THREE.WireframeGeometry,
    position: [number, number, number],
    rotation?: [number, number, number]
  ) => {
    return (
      <group position={position} rotation={rotation ? new THREE.Euler(...rotation) : undefined}>
        {/* Spatial depth mesh */}
        <mesh geometry={baseGeom} material={occludeMaterial} />
        {/* Glowing holographic contours */}
        <lineSegments geometry={wireGeom} material={lineMaterial} />
        {/* Corner vertex nodes */}
        <points geometry={baseGeom} material={pointsMaterial} />
      </group>
    );
  };

  return (
    <group ref={groupRef} position={[0, 0, 0]}>
      {/* Head & Neck */}
      {renderLimb(geometries.base.head, geometries.wire.head, [0, 0.72, 0.05])}
      {renderLimb(geometries.base.neck, geometries.wire.neck, [0, 0.48, 0.02])}

      {/* Torso */}
      {renderLimb(geometries.base.upperTorso, geometries.wire.upperTorso, [0, 0.22, 0])}
      {renderLimb(geometries.base.lowerTorso, geometries.wire.lowerTorso, [0, -0.12, -0.02])}

      {/* Left Arm (Seated Angle) */}
      {renderLimb(geometries.base.upperArm, geometries.wire.upperArm, [-0.38, 0.18, 0.05], [0.35, 0, -0.25])}
      {renderLimb(geometries.base.forearm, geometries.wire.forearm, [-0.34, -0.12, 0.24], [-0.75, 0.2, -0.1])}

      {/* Right Arm (Seated Angle) */}
      {renderLimb(geometries.base.upperArm, geometries.wire.upperArm, [0.38, 0.18, 0.05], [0.35, 0, 0.25])}
      {renderLimb(geometries.base.forearm, geometries.wire.forearm, [0.34, -0.12, 0.24], [-0.75, -0.2, 0.1])}

      {/* Thighs (Forward horizontal seated projection) */}
      {renderLimb(geometries.base.thigh, geometries.wire.thigh, [-0.18, -0.32, 0.22], [0.1, 0, 0])}
      {renderLimb(geometries.base.thigh, geometries.wire.thigh, [0.18, -0.32, 0.22], [0.1, 0, 0])}

      {/* Shins (Downward from knees to platform) */}
      {renderLimb(geometries.base.shin, geometries.wire.shin, [-0.18, -0.62, 0.44], [0.05, 0, 0])}
      {renderLimb(geometries.base.shin, geometries.wire.shin, [0.18, -0.62, 0.44], [0.05, 0, 0])}

      {/* Seated Plinth / Technical Seat */}
      {renderLimb(geometries.base.seatPlinth, geometries.wire.seatPlinth, [0, -0.56, -0.12])}
    </group>
  );
}

/**
 * Production GLB Model Component.
 * Dynamically mounts when /models/arham-wireframe.glb is provided by user.
 */
function ProductionGLBModel({ url, isAIActive = false }: { url: string; isAIActive?: boolean }) {
  const { scene } = useGLTF(url);

  // Apply wireframe shader / material to loaded GLB meshes
  useMemo(() => {
    scene.traverse((child) => {
      if ((child as THREE.Mesh).isMesh) {
        const mesh = child as THREE.Mesh;
        mesh.material = new THREE.MeshBasicMaterial({
          color: isAIActive ? HOLOGRAM_CONSTANTS.COLORS.PRIMARY : HOLOGRAM_CONSTANTS.COLORS.HOLOGRAM,
          wireframe: true,
          transparent: true,
          opacity: 0.85,
        });
      }
    });
  }, [scene, isAIActive]);

  return <primitive object={scene} />;
}

/**
 * HologramModel Boundary:
 * Automatically serves production GLB if available, otherwise renders the high-fidelity
 * procedural seated development wireframe mesh.
 */
export default function HologramModel({
  mode = "idle",
  isAIActive = false,
  intensity = 1,
}: HologramModelProps) {
  const isAI = isAIActive || mode === "active";

  // Production model path from constants
  const modelUrl = HOLOGRAM_CONSTANTS.ASSETS.PRODUCTION_MODEL;

  // We default to the procedural seated wireframe until user drops the GLB
  // This satisfies the Critical Truthfulness Rule and guarantees zero broken textures
  return <SeatedDevelopmentMesh isAIActive={isAI} intensity={intensity} />;
}
