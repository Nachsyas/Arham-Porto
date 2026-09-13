"use client";

import { motion } from "motion/react";
import { ArrowDown, Sparkles } from "lucide-react";
import { GithubIcon } from "@/components/icons/GithubIcon";
import { HologramStage } from "@/features/hologram";
import type { Profile } from "arham-porto-schema";

interface HeroSectionProps {
  profile: Profile;
  onOpenQuickReview: () => void;
}

export default function HeroSection({ profile, onOpenQuickReview }: HeroSectionProps) {
  const focusTags = ["Full-Stack", "Backend", "AI / ML", "System Design"];

  return (
    <section
      id="overview"
      aria-label="Overview & Hero Identity"
      className="relative min-h-[calc(100vh-4rem)] flex items-center justify-center border-b border-border px-4 py-16 sm:px-6 lg:px-8 overflow-hidden"
    >
      {/* Subtle Background Ambience */}
      <div className="absolute inset-0 pointer-events-none bg-[radial-gradient(ellipse_80%_60%_at_50%_-10%,rgba(38,184,255,0.08),rgba(2,6,11,0))]" />

      <div className="relative mx-auto max-w-7xl w-full grid grid-cols-1 lg:grid-cols-12 gap-12 items-center">
        {/* Left Column: Reviewer Identity & Copy */}
        <div className="lg:col-span-7 flex flex-col space-y-6 text-left">
          {/* Status & Role Pill */}
          <div className="inline-flex items-center gap-2 self-start px-3 py-1.5 rounded-full border border-primary/30 bg-surfaceElevated text-xs font-code text-primary">
            <span className="flex h-2 w-2 rounded-full bg-primary" />
            <span>{profile.role}</span>
            <span className="text-themeText-muted">|</span>
            <span className="text-themeText-muted">Evidence-Driven Profile</span>
          </div>

          {/* Primary Name Display */}
          <h1 className="font-display text-4xl sm:text-6xl lg:text-7xl font-extrabold tracking-tight text-themeText-primary leading-[1.08]">
            NACHSYAS ARHAM
            <br />
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-themeText-primary via-primary to-hologram">
              MUMTAZ NASHOHI
            </span>
          </h1>

          {/* Professional Positioning */}
          {profile.positioning && (
            <p className="font-body text-base sm:text-lg text-themeText-body max-w-2xl leading-relaxed">
              {profile.positioning}
            </p>
          )}

          {/* Focus Tags */}
          <div className="flex flex-wrap gap-2 pt-2">
            {focusTags.map((tag) => (
              <span
                key={tag}
                className="px-3 py-1 rounded-md text-xs font-code bg-surface border border-border text-themeText-muted hover:border-primary/40 hover:text-primary transition-colors"
              >
                #{tag}
              </span>
            ))}
          </div>

          {/* CTAs & Actions */}
          <div className="flex flex-wrap items-center gap-4 pt-4">
            <a
              href="#projects"
              className="inline-flex items-center gap-2 px-6 py-3 rounded-md text-sm font-semibold font-body bg-primary text-canvas hover:bg-primaryActive transition-all shadow-[0_0_15px_rgba(38,184,255,0.25)] focus-visible:ring-2 focus-visible:ring-primary"
            >
              Explore My Work
              <ArrowDown className="h-4 w-4" />
            </a>

            <button
              type="button"
              onClick={onOpenQuickReview}
              className="inline-flex items-center gap-2 px-5 py-3 rounded-md text-sm font-medium font-body bg-surfaceElevated border border-border text-themeText-primary hover:border-primary hover:text-primary transition-all focus-visible:ring-2 focus-visible:ring-primary"
            >
              <Sparkles className="h-4 w-4 text-primary" />
              60-Second Quick Review
            </button>

            {profile.github && (
              <a
                href={profile.github}
                target="_blank"
                rel="noopener noreferrer"
                aria-label="GitHub Profile"
                className="flex h-11 w-11 items-center justify-center rounded-md border border-border bg-surface text-themeText-muted hover:text-themeText-primary hover:border-primary transition-colors focus-visible:ring-2 focus-visible:ring-primary"
              >
                <GithubIcon className="h-5 w-5" />
              </a>
            )}
          </div>
        </div>

        {/* Right Column: 3D Hologram Stage */}
        <div className="lg:col-span-5 flex items-center justify-center">
          <HologramStage />
        </div>
      </div>
    </section>
  );
}
