"use client";

import { useEffect, useState } from "react";
import { getCorridorRoutePath } from "./journey.utils";

interface JourneyRouteProps {
  activeMilestoneId: string;
}

export default function JourneyRoute({ activeMilestoneId }: JourneyRouteProps) {
  const pathD = getCorridorRoutePath();
  const [prefersReducedMotion, setPrefersReducedMotion] = useState(false);

  // Listen for reduced motion user preference
  useEffect(() => {
    if (typeof window === "undefined" || typeof window.matchMedia !== "function") return;
    const mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
    setPrefersReducedMotion(mediaQuery.matches);

    const handleChange = (e: MediaQueryListEvent) => {
      setPrefersReducedMotion(e.matches);
    };

    if (mediaQuery.addEventListener) {
      mediaQuery.addEventListener("change", handleChange);
      return () => mediaQuery.removeEventListener("change", handleChange);
    }
  }, []);

  return (
    <g className="journey-route-layer" aria-hidden="true">
      <defs>
        {/* Glow corridor filter */}
        <filter id="corridor-glow" x="-20%" y="-20%" width="140%" height="140%">
          <feGaussianBlur stdDeviation="3" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>

        {/* Linear gradient along Java corridor */}
        <linearGradient id="journey-corridor-gradient" x1="0%" y1="0%" x2="100%" y2="50%">
          <stop offset="0%" stopColor="#06b6d4" stopOpacity="0.9" />
          <stop offset="50%" stopColor="#26b8ff" stopOpacity="0.95" />
          <stop offset="100%" stopColor="#38bdf8" stopOpacity="0.9" />
        </linearGradient>
      </defs>

      {/* Background guide glow corridor */}
      <path
        d={pathD}
        fill="none"
        stroke="#06b6d4"
        strokeWidth="6"
        strokeOpacity="0.12"
        strokeLinecap="round"
        filter="url(#corridor-glow)"
      />

      {/* Subtle dashed foundation track */}
      <path
        d={pathD}
        fill="none"
        stroke="rgba(148, 163, 184, 0.3)"
        strokeWidth="1.5"
        strokeDasharray="4 4"
        strokeLinecap="round"
      />

      {/* Active illuminated corridor route */}
      <path
        d={pathD}
        fill="none"
        stroke="url(#journey-corridor-gradient)"
        strokeWidth="2.5"
        strokeLinecap="round"
        className={
          prefersReducedMotion
            ? ""
            : "transition-all duration-1000 ease-out animate-[dash_1.5s_ease-out_forwards]"
        }
        style={{
          strokeDasharray: prefersReducedMotion ? "none" : "600",
          strokeDashoffset: prefersReducedMotion ? "0" : undefined,
        }}
        data-testid="journey-corridor-path"
      />

      {/* Decorative waypoint indicator rings along the corridor */}
      <circle cx="478.3" cy="203.2" r="3.5" fill="#06b6d4" fillOpacity="0.4" />
      <circle cx="445.3" cy="178.2" r="3.5" fill="#26b8ff" fillOpacity="0.4" />
      <circle cx="603.4" cy="239.2" r="3.5" fill="#38bdf8" fillOpacity="0.4" />
    </g>
  );
}
