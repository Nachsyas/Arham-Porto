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

  // Location state checks
  const isJakarta = location.id === "jakarta";
  const isMalang = location.id === "malang";
  const isUniversity = activeMilestoneId === "university-uin-malang";
  const isCurrentBase = activeMilestoneId === "current-base-malang";

  const handleClick = () => {
    if (isJakarta) {
      // If no Jakarta milestone is currently active, select the first: residence-jakarta
      // If a Jakarta milestone is already active, keep the current one (Rule 17)
      if (!location.milestoneIds.includes(activeMilestoneId)) {
        onSelectMilestone("residence-jakarta");
      }
    } else if (isMalang) {
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
    if (isJakarta) {
      if (activeMilestoneId === "residence-jakarta") {
        statusBadge = "Residence";
      } else if (activeMilestoneId === "mi-al-hamid-jakarta") {
        statusBadge = "MI Al-Hamid";
      } else if (activeMilestoneId === "mtsn30-jakarta") {
        statusBadge = "MTsN 30";
      }
    } else if (isMalang) {
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
        className={`group relative flex flex-col items-center focus:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 focus-visible:ring-offset-canvas rounded-full p-2 transition-all ${
          isActive ? "z-30 scale-110" : "z-10 hover:scale-105"
        }`}
      >
        {/* Radar Pulse Effect (when active, respects motion preferences) */}
        {isActive && (
          <span
            className="absolute h-8 w-8 rounded-full bg-primary/25 motion-safe:animate-ping motion-reduce:hidden"
            aria-hidden="true"
          />
        )}

        {/* Outer Halo */}
        <span
          className={`relative flex h-5 w-5 items-center justify-center rounded-full transition-all duration-300 ${
            isActive
              ? "bg-primary/20 ring-2 ring-primary shadow-sm"
              : "bg-surface ring-1 ring-border hover:ring-primary/60 hover:bg-primary-muted/40 shadow-sm"
          }`}
        >
          {/* Inner Solid Core */}
          <span
            className={`h-2 w-2 rounded-full transition-all duration-300 ${
              isActive ? "bg-primary scale-125" : "bg-themeText-muted group-hover:bg-primary"
            }`}
          />
        </span>

        {/* Location Label Tooltip / Badge */}
        <span
          className={`flex flex-col items-center px-2 py-0.5 rounded text-[10px] font-code font-medium whitespace-nowrap transition-all duration-200 pointer-events-none ${
            location.id === "salatiga" ? "order-first mb-1.5" : "order-last mt-1.5"
          } ${
            isActive
              ? "bg-surface border border-primary text-primary font-semibold shadow-sm backdrop-blur-sm"
              : "bg-surface/90 border border-border text-themeText-body group-hover:text-themeText-primary group-hover:border-border-strong shadow-sm backdrop-blur-sm"
          }`}
        >
          <span>{location.name}</span>
          {statusBadge && (
            <span className="text-[9px] text-primary font-mono -mt-0.5">
              [{statusBadge}]
            </span>
          )}
        </span>
      </button>
    </div>
  );
}
