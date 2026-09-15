"use client";

import { ArrowLeft, ArrowRight, MapPin, Calendar, Compass } from "lucide-react";
import type { PublicMilestone } from "./journey.types";

interface JourneyMilestoneCardProps {
  milestone: PublicMilestone;
  currentIndex: number;
  totalCount: number;
  onPrevious: () => void;
  onNext: () => void;
  hasPrevious: boolean;
  hasNext: boolean;
}

export default function JourneyMilestoneCard({
  milestone,
  currentIndex,
  totalCount,
  onPrevious,
  onNext,
  hasPrevious,
  hasNext,
}: JourneyMilestoneCardProps) {
  // Category presentation label
  const getCategoryBadge = (category: string) => {
    switch (category) {
      case "birthplace":
        return "ORIGIN";
      case "residence":
        return "RESIDENCE";
      case "sd":
        return "PRIMARY EDUCATION";
      case "smp":
        return "LOWER SECONDARY EDUCATION";
      case "sma":
        return "SECONDARY EDUCATION";
      case "university":
        return "COMPUTER SCIENCE";
      case "current":
        return "CURRENT BASE";
      default:
        return category.toUpperCase();
    }
  };

  return (
    <article
      aria-label={`Milestone details: ${milestone.title}`}
      className="flex flex-col justify-between rounded-card border border-border bg-surface p-6 sm:p-7 shadow-sm relative overflow-hidden transition-all duration-300"
    >
      {/* Decorative Technical Header Strip */}
      <div className="flex items-center justify-between border-b border-border/80 pb-3 mb-5">
        <div className="flex items-center gap-2">
          <span className="flex h-2 w-2 rounded-full bg-primary" />
          <span className="text-[11px] font-code font-bold text-primary tracking-wider">
            STEP 0{currentIndex + 1} // 0{totalCount}
          </span>
        </div>
        <div className="flex items-center gap-1 text-[11px] font-code text-themeText-muted">
          <Compass className="h-3 w-3 text-primary" />
          <span>JAVA // INDONESIA</span>
        </div>
      </div>

      {/* Main Milestone Content */}
      <div className="space-y-4">
        {/* Category Pill */}
        <div className="inline-flex items-center">
          <span className="px-2.5 py-1 rounded bg-surface-elevated border border-border text-[11px] font-code font-semibold text-primary tracking-wide shadow-sm">
            {getCategoryBadge(milestone.category)}
          </span>
        </div>

        {/* Milestone Title & Institution */}
        <div>
          <h3 className="font-display text-xl sm:text-2xl font-bold text-themeText-primary tracking-tight">
            {milestone.title}
          </h3>
          {milestone.institution && milestone.institution !== milestone.title && (
            <p className="text-sm font-body text-themeText-body mt-1">
              {milestone.institution}
            </p>
          )}
        </div>

        {/* Location & Period Meta */}
        <div className="flex flex-wrap items-center gap-y-2 gap-x-5 text-xs font-code text-themeText-muted pt-1">
          {/* Safe City & Region Location */}
          <div className="flex items-center gap-1.5 text-themeText-primary">
            <MapPin className="h-3.5 w-3.5 text-primary flex-shrink-0" />
            <span>
              {milestone.city}, {milestone.region}
            </span>
          </div>

          {/* Period (Rendered ONLY if non-null, strictly no birth year for origin) */}
          {milestone.period && (
            <div className="flex items-center gap-1.5 text-themeText-muted">
              <Calendar className="h-3.5 w-3.5 text-primary/80 flex-shrink-0" />
              <span>{milestone.period}</span>
            </div>
          )}
        </div>

        {/* Authentic Story Narrative */}
        {milestone.description && (
          <p className="text-xs sm:text-sm font-body text-themeText-body leading-relaxed pt-1 border-t border-border/50">
            {milestone.description}
          </p>
        )}
      </div>

      {/* Navigation Controls */}
      <div className="flex items-center justify-between pt-6 mt-6 border-t border-border">
        <button
          type="button"
          onClick={onPrevious}
          disabled={!hasPrevious}
          aria-label="Go to previous milestone"
          className={`inline-flex items-center gap-2 px-3.5 py-2 rounded-md text-xs font-code font-medium transition-all focus-visible:ring-2 focus-visible:ring-primary ${
            hasPrevious
              ? "bg-surface-elevated border border-border text-themeText-primary hover:bg-surface-strong hover:border-primary/50 cursor-pointer shadow-sm"
              : "bg-surface/50 border border-border/40 text-themeText-muted/40 cursor-not-allowed"
          }`}
        >
          <ArrowLeft className="h-3.5 w-3.5" />
          <span>Previous</span>
        </button>

        {/* Visual Progress Dots */}
        <div className="flex items-center gap-1.5" aria-hidden="true">
          {Array.from({ length: totalCount }).map((_, idx) => (
            <span
              key={idx}
              className={`h-1.5 rounded-full transition-all duration-300 ${
                idx === currentIndex
                  ? "w-6 bg-primary"
                  : "w-1.5 bg-border hover:bg-themeText-muted"
              }`}
            />
          ))}
        </div>

        <button
          type="button"
          onClick={onNext}
          disabled={!hasNext}
          aria-label="Go to next milestone"
          className={`inline-flex items-center gap-2 px-3.5 py-2 rounded-md text-xs font-code font-medium transition-all focus-visible:ring-2 focus-visible:ring-primary ${
            hasNext
              ? "bg-primary text-white font-semibold hover:bg-primary-active shadow-sm cursor-pointer"
              : "bg-surface/50 border border-border/40 text-themeText-muted/40 cursor-not-allowed"
          }`}
        >
          <span>Next</span>
          <ArrowRight className="h-3.5 w-3.5" />
        </button>
      </div>
    </article>
  );
}
