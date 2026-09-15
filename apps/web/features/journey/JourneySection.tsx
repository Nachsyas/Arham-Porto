"use client";

import { useState, useMemo, useEffect, useCallback } from "react";
import type { JourneyStop } from "arham-porto-schema";
import { filterPublicMilestones } from "./journey.utils";
import { GEOGRAPHIC_LOCATIONS } from "./data/indonesia-outline";
import JourneyMap from "./JourneyMap";
import JourneyMilestoneCard from "./JourneyMilestoneCard";
import JourneyTimeline from "./JourneyTimeline";

interface JourneySectionProps {
  stops?: JourneyStop[];
  hideHeader?: boolean;
  className?: string;
}

export default function JourneySection({ stops = [], hideHeader = false, className = "" }: JourneySectionProps) {
  // 1. Strict public filtering: non-public stops (TK, SD, SMP) are excluded from DOM
  const publicMilestones = useMemo(() => filterPublicMilestones(stops), [stops]);

  // 2. Initial state: strictly defaults to ORIGIN (chronological start)
  const [activeMilestoneId, setActiveMilestoneId] = useState<string>(
    publicMilestones[0]?.id || "origin-karanganyar"
  );

  // Derive active milestone and index
  const currentIndex = publicMilestones.findIndex((m) => m.id === activeMilestoneId);
  const safeIndex = currentIndex >= 0 ? currentIndex : 0;
  const currentMilestone = publicMilestones[safeIndex];

  // 4 unique geographic locations (Karanganyar, Jakarta, Salatiga, Malang)
  const uniqueLocations = useMemo(
    () => [
      GEOGRAPHIC_LOCATIONS.karanganyar,
      GEOGRAPHIC_LOCATIONS.jakarta,
      GEOGRAPHIC_LOCATIONS.salatiga,
      GEOGRAPHIC_LOCATIONS.malang,
    ],
    []
  );

  // Navigation handlers
  const handlePrevious = useCallback(() => {
    if (safeIndex > 0) {
      setActiveMilestoneId(publicMilestones[safeIndex - 1].id);
    }
  }, [safeIndex, publicMilestones]);

  const handleNext = useCallback(() => {
    if (safeIndex < publicMilestones.length - 1) {
      setActiveMilestoneId(publicMilestones[safeIndex + 1].id);
    }
  }, [safeIndex, publicMilestones]);

  // Optional keyboard navigation (Left/Right arrow) when section is focused
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Don't intercept if user is typing in an input/textarea
      if (
        e.target instanceof HTMLInputElement ||
        e.target instanceof HTMLTextAreaElement
      ) {
        return;
      }

      if (e.key === "ArrowLeft") {
        handlePrevious();
      } else if (e.key === "ArrowRight") {
        handleNext();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [handlePrevious, handleNext]);

  if (!currentMilestone) {
    return null;
  }

  return (
    <section
      id="journey"
      aria-label="Academic and Geographic Journey Map"
      className={className || "py-20 border-b border-border px-4 sm:px-6 lg:px-8 bg-canvas relative overflow-hidden"}
    >
      <div className="mx-auto max-w-7xl">
        {/* Section Header */}
        {!hideHeader && (
          <div className="flex flex-col md:flex-row md:items-end justify-between mb-10 gap-4">
            <div>
              <span className="text-xs font-code text-primary uppercase tracking-wider block mb-2">
                03 // GEOGRAPHIC & ACADEMIC EVOLUTION
              </span>
              <h2 className="font-display text-3xl sm:text-4xl font-extrabold text-themeText-primary tracking-tight">
                Interactive Journey Map
              </h2>
              <p className="mt-2 text-sm sm:text-base font-body text-themeText-muted max-w-2xl">
                From Karanganyar to Jakarta, Salatiga, and Malang — a journey through formative education and computer science.
              </p>
            </div>

            {/* Quick Counter Badge */}
            <div className="hidden sm:flex items-center gap-2 px-3 py-1.5 rounded-md bg-surface border border-border text-xs font-code text-themeText-muted">
              <span className="text-primary font-bold">0{safeIndex + 1}</span>
              <span>/</span>
              <span>0{publicMilestones.length} Milestones</span>
            </div>
          </div>
        )}

        {/* Responsive Desktop / Tablet / Mobile Composition */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 lg:gap-8 items-stretch mb-6">
          {/* Map Column (7 cols on desktop) */}
          <div className="lg:col-span-7 flex flex-col justify-center">
            <JourneyMap
              locations={uniqueLocations}
              activeMilestoneId={activeMilestoneId}
              onSelectMilestone={setActiveMilestoneId}
            />
          </div>

          {/* Milestone Story Card Column (5 cols on desktop) */}
          <div className="lg:col-span-5 flex flex-col">
            <JourneyMilestoneCard
              milestone={currentMilestone}
              currentIndex={safeIndex}
              totalCount={publicMilestones.length}
              onPrevious={handlePrevious}
              onNext={handleNext}
              hasPrevious={safeIndex > 0}
              hasNext={safeIndex < publicMilestones.length - 1}
            />
          </div>
        </div>

        {/* Interactive Progress Timeline (Full Width) */}
        <div>
          <JourneyTimeline
            milestones={publicMilestones}
            activeMilestoneId={activeMilestoneId}
            onSelectMilestone={setActiveMilestoneId}
          />
        </div>
      </div>
    </section>
  );
}
