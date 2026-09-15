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
          <stop offset="0%" stopColor="#0799E6" stopOpacity="0.95" />
          <stop offset="50%" stopColor="#007CC3" stopOpacity="1" />
          <stop offset="100%" stopColor="#0799E6" stopOpacity="0.95" />
        </linearGradient>
      </defs>

      {/* Background guide glow corridor */}
      <path
        d={pathD}
        fill="none"
        stroke="#0799E6"
        strokeWidth="5"
        strokeOpacity="0.08"
        strokeLinecap="round"
        filter="url(#corridor-glow)"
      />

      {/* Subtle dashed foundation track */}
      <path
        d={pathD}
        fill="none"
        stroke="rgba(190, 211, 223, 0.5)"
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
          strokeDasharray: prefersReducedMotion ? "none" : "1200",
          strokeDashoffset: prefersReducedMotion ? "0" : undefined,
        }}
        data-testid="journey-corridor-path"
      />

      {/* Decorative waypoint indicator rings along the corridor */}
      <circle cx="478.3" cy="203.2" r="3.5" fill="#0799E6" fillOpacity="0.6" />
      <circle cx="172.0" cy="72.0" r="3.5" fill="#0799E6" fillOpacity="0.6" />
      <circle cx="445.3" cy="178.2" r="3.5" fill="#007CC3" fillOpacity="0.6" />
      <circle cx="603.4" cy="239.2" r="3.5" fill="#0799E6" fillOpacity="0.6" />
    </g>
  );
}
