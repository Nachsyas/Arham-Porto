"use client";

import { Compass, MapPin } from "lucide-react";
import type { GeographicLocation } from "./journey.types";
import {
  MAP_DIMENSIONS,
  JAVA_ISLAND_PATH,
  MADURA_ISLAND_PATH,
  BALI_ISLAND_PATH,
  INDONESIA_ARCHIPELAGO_PATH,
} from "./data/indonesia-outline";
import JourneyMarker from "./JourneyMarker";
import JourneyRoute from "./JourneyRoute";

interface JourneyMapProps {
  locations: GeographicLocation[];
  activeMilestoneId: string;
  onSelectMilestone: (milestoneId: string) => void;
}

export default function JourneyMap({
  locations,
  activeMilestoneId,
  onSelectMilestone,
}: JourneyMapProps) {
  return (
    <div className="relative w-full rounded-card border border-border bg-canvas/90 overflow-hidden shadow-2xl">
      {/* Decorative Technical Top Bar */}
      <div className="flex items-center justify-between px-4 py-2.5 border-b border-border bg-surface/50 text-[10px] font-code text-themeText-muted">
        <div className="flex items-center gap-2">
          <span className="h-1.5 w-1.5 rounded-full bg-cyan-400 animate-pulse" />
          <span className="text-cyan-400 font-medium">JAVA // CORRIDOR VISUALIZATION</span>
        </div>
        <div className="flex items-center gap-3">
          <span className="hidden sm:inline text-themeText-muted/70">
            CENTRAL JAVA → EAST JAVA
          </span>
          <div className="flex items-center gap-1 text-primary/90">
            <Compass className="h-3 w-3" />
            <span>INDONESIA</span>
          </div>
        </div>
      </div>

      {/* SVG Canvas Container */}
      <div className="relative aspect-[800/360] w-full select-none bg-radial from-surfaceElevated/20 via-canvas to-canvas">
        {/* Subtle Background Coordinate Grid */}
        <div
          className="absolute inset-0 opacity-[0.07] pointer-events-none"
          style={{
            backgroundImage:
              "linear-gradient(to right, #38bdf8 1px, transparent 1px), linear-gradient(to bottom, #38bdf8 1px, transparent 1px)",
            backgroundSize: "40px 40px",
          }}
          aria-hidden="true"
        />

        {/* Vector SVG Map */}
        <svg
          viewBox={MAP_DIMENSIONS.viewBox}
          className="w-full h-full block"
          role="img"
          aria-label="Stylized geographic map showing academic journey across Central Java and East Java, Indonesia"
        >
          <defs>
            {/* Java island gradient */}
            <linearGradient id="java-island-fill" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" stopColor="#1e293b" stopOpacity="0.9" />
              <stop offset="50%" stopColor="#0f172a" stopOpacity="0.95" />
              <stop offset="100%" stopColor="#1e293b" stopOpacity="0.9" />
            </linearGradient>

            {/* Subtle glow filter */}
            <filter id="map-glow" x="-10%" y="-10%" width="120%" height="120%">
              <feDropShadow dx="0" dy="2" stdDeviation="4" floodColor="#06b6d4" floodOpacity="0.2" />
            </filter>
          </defs>

          {/* 1. Low-contrast Ambient Indonesian Archipelago Silhouette (Overview Context) */}
          <g id="indonesia-archipelago-silhouette" opacity="0.35" aria-hidden="true">
            <path
              d={INDONESIA_ARCHIPELAGO_PATH}
              fill="rgba(148, 163, 184, 0.08)"
              stroke="rgba(148, 163, 184, 0.25)"
              strokeWidth="0.75"
              strokeLinejoin="round"
            />
          </g>

          {/* 2. Stronger Contrast Java Island & Neighbors (Primary Region) */}
          <g id="java-region-focus" filter="url(#map-glow)" aria-hidden="true">
            {/* Madura Island */}
            <path
              d={MADURA_ISLAND_PATH}
              fill="url(#java-island-fill)"
              stroke="rgba(6, 182, 212, 0.35)"
              strokeWidth="1.2"
              strokeLinejoin="round"
            />

            {/* Bali Island */}
            <path
              d={BALI_ISLAND_PATH}
              fill="url(#java-island-fill)"
              stroke="rgba(148, 163, 184, 0.25)"
              strokeWidth="1"
              strokeLinejoin="round"
              opacity="0.6"
            />

            {/* Main Java Island */}
            <path
              d={JAVA_ISLAND_PATH}
              fill="url(#java-island-fill)"
              stroke="rgba(6, 182, 212, 0.55)"
              strokeWidth="1.5"
              strokeLinejoin="round"
            />
          </g>

          {/* 3. Primary Corridor Route (Karanganyar -> Salatiga -> Malang) */}
          <JourneyRoute activeMilestoneId={activeMilestoneId} />

          {/* 4. Subtle Regional Labels (Non-intrusive) */}
          <text
            x="420"
            y="235"
            fill="rgba(148, 163, 184, 0.35)"
            fontSize="11"
            fontFamily="monospace"
            letterSpacing="2"
            textAnchor="middle"
            aria-hidden="true"
          >
            JAWA TENGAH
          </text>
          <text
            x="640"
            y="245"
            fill="rgba(148, 163, 184, 0.35)"
            fontSize="11"
            fontFamily="monospace"
            letterSpacing="2"
            textAnchor="middle"
            aria-hidden="true"
          >
            JAWA TIMUR
          </text>
        </svg>

        {/* 5. Accessible Map Markers Overlay (Real HTML Buttons positioned over SVG) */}
        <div
          className="absolute inset-0 pointer-events-none"
          data-testid="journey-markers-container"
        >
          {locations.map((location) => (
            <JourneyMarker
              key={location.id}
              location={location}
              activeMilestoneId={activeMilestoneId}
              onSelectMilestone={onSelectMilestone}
            />
          ))}
        </div>
      </div>

      {/* Decorative Technical Bottom Bar */}
      <div className="flex items-center justify-between px-4 py-2 border-t border-border/80 bg-surface/30 text-[9px] font-code text-themeText-muted/70">
        <div>STYLIZED GEOGRAPHIC VISUALIZATION</div>
        <div className="hidden sm:block">SOURCE: NATURAL EARTH (PUBLIC DOMAIN / CC0)</div>
        <div className="flex items-center gap-1 text-cyan-400/80">
          <MapPin className="h-2.5 w-2.5" />
          <span>3 UNIQUE LOCATIONS</span>
        </div>
      </div>
    </div>
  );
}
