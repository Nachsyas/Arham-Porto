"use client";

import type { GeographicLocation } from "./journey.types";
import { MAP_DIMENSIONS } from "./data/indonesia-outline";

interface JourneyMarkerProps {
  location: GeographicLocation;
  activeMilestoneId: string;
  onSelectMilestone: (milestoneId: string) => void;
}

export default function JourneyMarker({
  location,
  activeMilestoneId,
  onSelectMilestone,
}: JourneyMarkerProps) {
  // Check if this geographic location corresponds to the active milestone
  const isActive = location.milestoneIds.includes(activeMilestoneId);

  // Position relative to the 800x360 map viewBox coordinate space
  const leftPct = (location.x / MAP_DIMENSIONS.width) * 100;
  const topPct = (location.y / MAP_DIMENSIONS.height) * 100;

  // For shared Malang location, detect whether University or Current Base is active
  const isMalang = location.id === "malang";
  const isUniversity = activeMilestoneId === "university-uin-malang";
  const isCurrentBase = activeMilestoneId === "current-base-malang";

  const handleClick = () => {
    if (isMalang) {
      // Toggle or cycle between University and Current Base if already on Malang,
      // otherwise default to University
      if (isUniversity) {
        onSelectMilestone("current-base-malang");
      } else if (isCurrentBase) {
        onSelectMilestone("university-uin-malang");
      } else {
        onSelectMilestone("university-uin-malang");
      }
    } else {
      onSelectMilestone(location.milestoneIds[0]);
    }
  };

  // Status subtitle for the pin
  let statusBadge = "";
  if (isActive) {
    if (isMalang) {
      statusBadge = isCurrentBase ? "Current Base" : "University";
    } else if (location.id === "karanganyar") {
      statusBadge = "Origin";
    } else if (location.id === "salatiga") {
      statusBadge = "Secondary";
    }
  }

  return (
    <div
      className="absolute pointer-events-auto transform -translate-x-1/2 -translate-y-1/2 transition-transform duration-200"
      style={{ left: `${leftPct}%`, top: `${topPct}%` }}
    >
      <button
        type="button"
        onClick={handleClick}
        aria-pressed={isActive}
        aria-current={isActive ? "location" : undefined}
        aria-label={`Geographic location: ${location.name}, ${location.province}. ${
          isActive
            ? `Active location (${statusBadge || "selected"}).`
            : "Click to view milestones."
        }`}
        data-testid={`journey-marker-${location.id}`}
        className={`group relative flex flex-col items-center focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400 focus-visible:ring-offset-2 focus-visible:ring-offset-canvas rounded-full p-2 transition-all ${
          isActive ? "z-30 scale-110" : "z-10 hover:scale-105"
        }`}
      >
        {/* Radar Pulse Effect (when active, respects motion preferences) */}
        {isActive && (
          <span
            className="absolute h-8 w-8 rounded-full bg-cyan-400/25 motion-safe:animate-ping motion-reduce:hidden"
            aria-hidden="true"
          />
        )}

        {/* Outer Glow Halo */}
        <span
          className={`relative flex h-5 w-5 items-center justify-center rounded-full transition-all duration-300 ${
            isActive
              ? "bg-cyan-500/30 ring-2 ring-cyan-400 shadow-[0_0_15px_rgba(6,182,212,0.8)]"
              : "bg-surfaceStrong/80 ring-1 ring-border hover:ring-cyan-400/60 hover:bg-cyan-950/40"
          }`}
        >
          {/* Inner Solid Core */}
          <span
            className={`h-2 w-2 rounded-full transition-all duration-300 ${
              isActive ? "bg-cyan-300 scale-125 shadow-[0_0_8px_#38bdf8]" : "bg-cyan-500/60 group-hover:bg-cyan-400"
            }`}
          />
        </span>

        {/* Location Label Tooltip / Badge */}
        <span
          className={`flex flex-col items-center px-2 py-0.5 rounded text-[10px] font-code font-medium whitespace-nowrap transition-all duration-200 pointer-events-none ${
            location.id === "salatiga" ? "order-first mb-1.5" : "order-last mt-1.5"
          } ${
            isActive
              ? "bg-surfaceElevated/95 border border-cyan-400/50 text-cyan-300 shadow-[0_0_10px_rgba(6,182,212,0.3)] backdrop-blur-sm"
              : "bg-surface/80 border border-border text-themeText-muted group-hover:text-themeText-primary group-hover:border-borderStrong backdrop-blur-xs"
          }`}
        >
          <span>{location.name}</span>
          {statusBadge && (
            <span className="text-[9px] text-cyan-400/90 font-mono -mt-0.5">
              [{statusBadge}]
            </span>
          )}
        </span>
      </button>
    </div>
  );
}
