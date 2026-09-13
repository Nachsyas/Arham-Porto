"use client";

import { motion } from "motion/react";
import { ArrowDown, Sparkles, Terminal, Layers } from "lucide-react";
import { GithubIcon } from "@/components/icons/GithubIcon";
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

        {/* Right Column: Hologram Stage Development Placeholder */}
        <div className="lg:col-span-5 flex items-center justify-center">
          <div
            className="relative w-full max-w-md aspect-square rounded-2xl border border-border bg-gradient-to-b from-surfaceElevated/60 to-surface/40 p-6 flex flex-col items-center justify-center shadow-xl overflow-hidden group"
            aria-label="3D Hologram Stage Placeholder"
          >
            {/* Technical Corner Accents */}
            <span className="absolute top-2 left-2 text-[10px] font-code text-primary/40">┌ SYS:STAGE</span>
            <span className="absolute bottom-2 right-2 text-[10px] font-code text-primary/40">PHASE_02 ┘</span>

            {/* Geometric Wireframe Visual Grid */}
            <div className="relative w-64 h-64 flex items-center justify-center">
              {/* Outer concentric projection ring */}
              <div className="absolute inset-0 rounded-full border border-dashed border-primary/25 animate-[spin_40s_linear_infinite]" />
              {/* Inner ring */}
              <div className="absolute inset-6 rounded-full border border-hologram/30 shadow-[0_0_20px_rgba(97,216,255,0.15)]" />

              {/* Wireframe geometric matrix silhouette */}
              <div className="relative z-10 flex flex-col items-center justify-center p-6 text-center space-y-3">
                <div className="flex h-16 w-16 items-center justify-center rounded-xl bg-canvas border border-primary/40 shadow-[0_0_20px_rgba(38,184,255,0.3)]">
                  <Layers className="h-8 w-8 text-primary" />
                </div>
                <div className="space-y-1">
                  <p className="text-xs font-code font-bold tracking-wider text-hologram">
                    3D WIREFRAME STAGE
                  </p>
                  <p className="text-[11px] font-body text-themeText-muted max-w-[200px] leading-snug">
                    Seated human wireframe scene scheduled for Phase 2 integration.
                  </p>
                </div>
              </div>

              {/* Holographic scanning horizontal line */}
              <div className="absolute left-4 right-4 h-px bg-gradient-to-r from-transparent via-primary/50 to-transparent animate-pulse" />
            </div>

            {/* Technical Metadata Footer */}
            <div className="w-full mt-4 pt-3 border-t border-border/60 flex items-center justify-between text-[10px] font-code text-themeText-muted">
              <span>R3F / Three.js Target</span>
              <span className="text-primary font-medium">Pose: Seated Mesh</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
