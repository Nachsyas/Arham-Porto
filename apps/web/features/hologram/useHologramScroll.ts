"use client";

import { useEffect, useState } from "react";
import { HOLOGRAM_CONSTANTS } from "./hologram.constants";

export interface HologramScrollState {
  targetYaw: number;
  targetPitch: number;
  scrollProgress: number;
  isReducedMotion: boolean;
  isMobile: boolean;
}

export function useHologramScroll(): HologramScrollState {
  const [scrollProgress, setScrollProgress] = useState(0);
  const [isReducedMotion, setIsReducedMotion] = useState(false);
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    // Detect prefers-reduced-motion
    const mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
    setIsReducedMotion(mediaQuery.matches);

    const handleMotionChange = (e: MediaQueryListEvent) => {
      setIsReducedMotion(e.matches);
    };
    mediaQuery.addEventListener("change", handleMotionChange);

    // Detect mobile viewport
    const checkMobile = () => {
      setIsMobile(window.innerWidth < 768);
    };
    checkMobile();
    window.addEventListener("resize", checkMobile);

    // Track scroll
    const handleScroll = () => {
      const scrollY = window.scrollY;
      const heroHeight = Math.max(window.innerHeight * 0.8, 600);
      // Normalized progress [0, 1] across Hero view
      const progress = Math.min(Math.max(scrollY / heroHeight, 0), 1);
      setScrollProgress(progress);
    };

    window.addEventListener("scroll", handleScroll, { passive: true });
    handleScroll();

    return () => {
      mediaQuery.removeEventListener("change", handleMotionChange);
      window.removeEventListener("resize", checkMobile);
      window.removeEventListener("scroll", handleScroll);
    };
  }, []);

  if (isReducedMotion) {
    return {
      targetYaw: 0,
      targetPitch: 0,
      scrollProgress: 0,
      isReducedMotion: true,
      isMobile,
    };
  }

  // Linear mapping: scroll 0 -> -MAX_YAW, scroll 1 -> +MAX_YAW
  const targetYaw =
    -HOLOGRAM_CONSTANTS.ROTATION.MAX_YAW_RAD +
    scrollProgress * (2 * HOLOGRAM_CONSTANTS.ROTATION.MAX_YAW_RAD);

  // Subtle pitch offset based on scroll
  const targetPitch =
    Math.sin(scrollProgress * Math.PI) * HOLOGRAM_CONSTANTS.ROTATION.MAX_PITCH_RAD;

  return {
    targetYaw,
    targetPitch,
    scrollProgress,
    isReducedMotion,
    isMobile,
  };
}
