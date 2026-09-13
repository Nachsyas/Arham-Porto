"use client";

import { Check } from "lucide-react";
import type { PublicMilestone } from "./journey.types";

interface JourneyTimelineProps {
  milestones: PublicMilestone[];
  activeMilestoneId: string;
  onSelectMilestone: (id: string) => void;
}

export default function JourneyTimeline({
  milestones,
  activeMilestoneId,
  onSelectMilestone,
}: JourneyTimelineProps) {
  const activeIndex = milestones.findIndex((m) => m.id === activeMilestoneId);

  // Responsive label helper: short for mobile, descriptive for desktop
  const getMilestoneLabels = (milestone: PublicMilestone) => {
    switch (milestone.category) {
      case "birthplace":
        return { short: "Origin", full: "Origin — Karanganyar" };
      case "sma":
        return { short: "MA", full: "MA Tahfizh As-Surkati" };
      case "university":
        return { short: "University", full: "UIN Malang — CS" };
      case "current":
        return { short: "Current", full: "Current Base — Malang" };
      default:
        return { short: milestone.city, full: milestone.title };
    }
  };

  return (
    <nav
      aria-label="Milestone Progress Timeline"
      className="w-full rounded-card border border-border bg-surface/60 p-4 sm:p-5 backdrop-blur-xs"
    >
      <ol className="relative flex items-center justify-between w-full">
        {/* Background connector line */}
        <div
          className="absolute left-0 top-1/2 -translate-y-1/2 w-full h-0.5 bg-border/80 z-0"
          aria-hidden="true"
        />

        {/* Dynamic active progress fill */}
        <div
          className="absolute left-0 top-1/2 -translate-y-1/2 h-0.5 bg-cyan-400 z-0 transition-all duration-300 shadow-[0_0_8px_#22d3ee]"
          style={{
            width: `${(Math.max(0, activeIndex) / Math.max(1, milestones.length - 1)) * 100}%`,
          }}
          aria-hidden="true"
        />

        {/* Milestone Steps */}
        {milestones.map((milestone, index) => {
          const isActive = milestone.id === activeMilestoneId;
          const isPassed = index < activeIndex;
          const labels = getMilestoneLabels(milestone);

          return (
            <li
              key={milestone.id}
              className="relative z-10 flex flex-col items-center"
            >
              <button
                type="button"
                onClick={() => onSelectMilestone(milestone.id)}
                aria-current={isActive ? "step" : undefined}
                aria-label={`Milestone ${index + 1}: ${labels.full}. ${
                  isActive ? "Currently selected." : "Click to view details."
                }`}
                data-testid={`timeline-step-${milestone.id}`}
                className={`group flex flex-col items-center focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400 focus-visible:ring-offset-2 focus-visible:ring-offset-canvas rounded-lg p-1.5 transition-all`}
              >
                {/* Step Circle Indicator */}
                <span
                  className={`flex h-8 w-8 sm:h-9 sm:w-9 items-center justify-center rounded-full font-code text-xs font-bold transition-all duration-300 ${
                    isActive
                      ? "bg-cyan-400 text-canvas ring-4 ring-cyan-400/20 shadow-[0_0_12px_#22d3ee] scale-110"
                      : isPassed
                      ? "bg-surfaceElevated border-2 border-cyan-500/60 text-cyan-400 hover:border-cyan-400"
                      : "bg-surfaceStrong border border-border text-themeText-muted hover:border-borderStrong hover:text-themeText-primary"
                  }`}
                >
                  {isPassed ? (
                    <Check className="h-4 w-4 stroke-[2.5]" />
                  ) : (
                    <span>0{index + 1}</span>
                  )}
                </span>

                {/* Milestone Title Label */}
                <span
                  className={`mt-2 text-center font-code transition-colors duration-200 ${
                    isActive
                      ? "text-cyan-400 font-semibold text-xs sm:text-sm"
                      : isPassed
                      ? "text-themeText-primary text-[11px] sm:text-xs"
                      : "text-themeText-muted text-[11px] sm:text-xs group-hover:text-themeText-body"
                  }`}
                >
                  {/* Short label on mobile screens */}
                  <span className="inline sm:hidden">{labels.short}</span>
                  {/* Descriptive label on tablet & desktop */}
                  <span className="hidden sm:inline">{labels.full}</span>
                </span>
              </button>
            </li>
          );
        })}
      </ol>
    </nav>
  );
}
