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

  // Responsive label helper: concise for desktop horizontal strip, full for mobile vertical list
  const getMilestoneLabels = (milestone: PublicMilestone) => {
    switch (milestone.id) {
      case "origin-karanganyar":
        return {
          short: "Origin",
          medium: "Origin",
          full: "Origin — Karanganyar",
          accessible: "Origin in Karanganyar, Central Java",
        };
      case "residence-jakarta":
        return {
          short: "Jakarta",
          medium: "Jakarta",
          full: "Residence — Jakarta",
          accessible: "Residence in Jakarta, DKI Jakarta",
        };
      case "mi-al-hamid-jakarta":
        return {
          short: "MI",
          medium: "MI Al-Hamid",
          full: "MI Terpadu Al-Hamid",
          accessible: "Primary education at Madrasah Ibtidaiyah Terpadu Al-Hamid, Jakarta Timur",
        };
      case "mtsn30-jakarta":
        return {
          short: "MTs",
          medium: "MTsN 30",
          full: "MTsN 30 Jakarta",
          accessible: "Lower secondary education at Madrasah Tsanawiyah Negeri 30 Jakarta Timur",
        };
      case "ma-assurkati-salatiga":
        return {
          short: "MA",
          medium: "MA As-Surkati",
          full: "MA Tahfizh As-Surkati",
          accessible: "Secondary education at Madrasah Aliyah Tahfizhul Qur'an As-Surkati, Salatiga",
        };
      case "university-uin-malang":
        return {
          short: "Uni",
          medium: "UIN Malang",
          full: "UIN Malang — CS",
          accessible: "Computer Science undergraduate studies at Universitas Islam Negeri Maulana Malik Ibrahim Malang",
        };
      case "current-base-malang":
        return {
          short: "Current",
          medium: "Current Base",
          full: "Current Base — Malang",
          accessible: "Current engineering base in Malang, East Java",
        };
      default:
        return {
          short: milestone.city,
          medium: milestone.title,
          full: milestone.title,
          accessible: milestone.title,
        };
    }
  };

  const progressPct =
    milestones.length > 1
      ? (Math.max(0, activeIndex) / (milestones.length - 1)) * 100
      : 0;

  return (
    <nav
      aria-label="Milestone Progress Timeline"
      className="w-full rounded-card border border-border bg-surface/80 p-4 sm:p-5 backdrop-blur-sm relative overflow-hidden shadow-sm"
    >
      <div className="relative w-full">
        {/* Desktop background connector line (horizontal) */}
        <div
          className="hidden sm:block absolute left-0 top-1/2 -translate-y-1/2 w-full h-0.5 bg-border/80 z-0"
          aria-hidden="true"
        />

        {/* Desktop dynamic progress fill (horizontal) */}
        <div
          className="hidden sm:block absolute left-0 top-1/2 -translate-y-1/2 h-0.5 bg-primary z-0 transition-all duration-300"
          style={{ width: `${progressPct}%` }}
          aria-hidden="true"
        />

        {/* Mobile background connector line (vertical) */}
        <div
          className="sm:hidden absolute left-[19px] top-4 bottom-4 w-0.5 bg-border/80 z-0"
          aria-hidden="true"
        />

        {/* Mobile dynamic progress fill (vertical) */}
        <div
          className="sm:hidden absolute left-[19px] top-4 w-0.5 bg-primary z-0 transition-all duration-300"
          style={{ height: `${progressPct}%` }}
          aria-hidden="true"
        />

        {/* Milestone Steps List: vertical on mobile, horizontal on desktop */}
        <ol className="relative flex flex-col sm:flex-row sm:items-center sm:justify-between w-full gap-2 sm:gap-0">
          {milestones.map((milestone, index) => {
            const isActive = milestone.id === activeMilestoneId;
            const isPassed = index < activeIndex;
            const labels = getMilestoneLabels(milestone);

            return (
              <li
                key={milestone.id}
                className="relative z-10 flex flex-col sm:items-center"
              >
                <button
                  type="button"
                  onClick={() => onSelectMilestone(milestone.id)}
                  aria-current={isActive ? "step" : undefined}
                  aria-label={`Milestone ${index + 1}: ${labels.accessible}. ${
                    isActive ? "Currently selected." : "Click to view details."
                  }`}
                  data-testid={`timeline-step-${milestone.id}`}
                  className={`group flex sm:flex-col items-center gap-3 sm:gap-1.5 w-full sm:w-auto p-1.5 sm:p-1 rounded-lg text-left sm:text-center focus:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-2 focus-visible:ring-offset-canvas transition-all ${
                    isActive
                      ? "bg-surface-elevated sm:bg-transparent"
                      : "hover:bg-surface/50 sm:hover:bg-transparent"
                  }`}
                >
                  {/* Step Circle Indicator */}
                  <span
                    className={`flex h-8 w-8 sm:h-9 sm:w-9 flex-shrink-0 items-center justify-center rounded-full font-code text-xs font-bold transition-all duration-300 ${
                      isActive
                        ? "bg-primary text-white ring-4 ring-primary/20 scale-110 shadow-sm"
                        : isPassed
                        ? "bg-primary-muted border-2 border-primary text-primary hover:bg-primary/20"
                        : "bg-surface border border-border text-themeText-muted group-hover:border-border-strong group-hover:text-themeText-primary"
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
                    className={`font-code transition-colors duration-200 ${
                      isActive
                        ? "text-primary font-semibold text-xs sm:text-xs"
                        : isPassed
                        ? "text-themeText-primary text-[11px] sm:text-[11px]"
                        : "text-themeText-muted text-[11px] sm:text-[11px] group-hover:text-themeText-body"
                    }`}
                  >
                    {/* Full descriptive label on mobile */}
                    <span className="inline sm:hidden">{labels.full}</span>
                    {/* Concise label on tablet & desktop */}
                    <span className="hidden sm:inline lg:hidden">{labels.short}</span>
                    <span className="hidden lg:inline">{labels.medium}</span>
                  </span>
                </button>
              </li>
            );
          })}
        </ol>
      </div>
    </nav>
  );
}
