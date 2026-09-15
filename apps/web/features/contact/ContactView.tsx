"use client";

import React from "react";
import { Sparkles, ExternalLink, Bot, ShieldCheck, Terminal } from "lucide-react";
import { GithubIcon } from "@/components/icons/GithubIcon";
import { usePortfolioUI } from "@/context/PortfolioUIContext";
import type { Profile } from "arham-porto-schema";

interface ContactViewProps {
  profile: Profile;
}

export default function ContactView({ profile }: ContactViewProps) {
  const { openAskArham, openQuickReview } = usePortfolioUI();

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      {/* Decorative Cyan Orbit Background Element */}
      <div className="relative py-8">
        <div
          className="absolute inset-0 -top-12 flex items-center justify-center pointer-events-none opacity-40"
          aria-hidden="true"
        >
          <div className="w-96 h-96 rounded-full border border-primary/20 animate-[spin_60s_linear_infinite]" />
          <div className="absolute w-72 h-72 rounded-full border border-dashed border-primary/25 animate-[spin_40s_linear_infinite_reverse]" />
          <div className="absolute w-48 h-48 rounded-full bg-primary/5 blur-2xl" />
        </div>

        {/* Channels Grid */}
        <div className="relative z-10 grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Card 1: Verified GitHub Channel */}
          <div className="rounded-card border border-border bg-surface p-6 sm:p-7 shadow-sm hover:border-primary/50 transition-all flex flex-col justify-between">
            <div>
              <div className="flex items-center justify-between gap-2 mb-4">
                <div className="h-10 w-10 rounded-lg bg-surface-elevated border border-border flex items-center justify-center text-themeText-primary">
                  <GithubIcon className="h-5 w-5" />
                </div>
                <span className="px-2.5 py-0.5 rounded text-[10px] font-code font-bold uppercase bg-primary-muted text-primary border border-primary/20">
                  Verified Identity
                </span>
              </div>

              <h2 className="font-display text-lg sm:text-xl font-bold text-themeText-primary tracking-tight">
                GitHub Repositories & Code
              </h2>
              <p className="mt-2 text-xs sm:text-sm font-body text-themeText-body leading-relaxed">
                Review verified commit history, Go backend microservices, Next.js web applications, and architectural source code.
              </p>
            </div>

            <div className="mt-6 pt-4 border-t border-border/80">
              {profile.github ? (
                <a
                  href={profile.github}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center justify-center gap-2 w-full px-4 py-2.5 rounded-md text-xs font-code font-semibold bg-surface-elevated border border-border text-themeText-primary hover:bg-surface-strong hover:border-primary transition-all shadow-sm focus-visible:ring-2 focus-visible:ring-primary"
                >
                  <GithubIcon className="h-4 w-4" />
                  <span>github.com/Nachsyas</span>
                  <ExternalLink className="h-3.5 w-3.5 opacity-70" />
                </a>
              ) : (
                <span className="text-xs font-code text-themeText-muted">Channel pending verification</span>
              )}
            </div>
          </div>

          {/* Card 2: Interactive AI Assistant (Ask Arham) */}
          <div className="rounded-card border border-primary/30 bg-surface p-6 sm:p-7 shadow-sm hover:border-primary transition-all flex flex-col justify-between relative overflow-hidden">
            <div className="absolute top-0 right-0 h-24 w-24 bg-primary/5 rounded-bl-full pointer-events-none" />

            <div>
              <div className="flex items-center justify-between gap-2 mb-4">
                <div className="h-10 w-10 rounded-lg bg-primary/10 border border-primary/30 flex items-center justify-center text-primary">
                  <Bot className="h-5 w-5 text-primary" />
                </div>
                <span className="px-2.5 py-0.5 rounded text-[10px] font-code font-bold uppercase bg-primary text-white">
                  Grounded AI
                </span>
              </div>

              <h2 className="font-display text-lg sm:text-xl font-bold text-themeText-primary tracking-tight flex items-center gap-2">
                <span>Ask Arham AI</span>
                <Sparkles className="h-4 w-4 text-primary" />
              </h2>
              <p className="mt-2 text-xs sm:text-sm font-body text-themeText-body leading-relaxed">
                Query Nachsyas Arham&apos;s background directly. Powered by hybrid vector search over verified repositories with strict citation anchoring.
              </p>
            </div>

            <div className="mt-6 pt-4 border-t border-border/80">
              <button
                type="button"
                onClick={openAskArham}
                className="inline-flex items-center justify-center gap-2 w-full px-4 py-2.5 rounded-md text-xs font-code font-semibold bg-primary text-white hover:bg-primary-active transition-all shadow-sm focus-visible:ring-2 focus-visible:ring-primary cursor-pointer"
              >
                <Bot className="h-4 w-4" />
                <span>Launch Ask Arham Assistant</span>
              </button>
            </div>
          </div>
        </div>

        {/* Quick Review Card Strip */}
        <div className="mt-6 rounded-card border border-border bg-surface-elevated p-5 sm:p-6 flex flex-col sm:flex-row items-center justify-between gap-4 shadow-sm">
          <div className="flex items-start gap-3">
            <div className="h-8 w-8 rounded-md bg-surface border border-border flex items-center justify-center text-primary flex-shrink-0 mt-0.5">
              <Sparkles className="h-4 w-4 text-primary" />
            </div>
            <div>
              <h3 className="font-display text-sm sm:text-base font-bold text-themeText-primary">
                Short on time? Use the 60-Second Quick Review
              </h3>
              <p className="text-xs font-body text-themeText-muted mt-0.5">
                Inspect architecture highlights, core competencies, and featured system metrics in a high-density reviewer drawer.
              </p>
            </div>
          </div>

          <button
            type="button"
            onClick={openQuickReview}
            className="inline-flex items-center justify-center gap-2 px-4 py-2 rounded-md text-xs font-code font-medium bg-surface border border-border text-primary hover:bg-surface-strong hover:border-primary transition-all whitespace-nowrap shadow-sm focus-visible:ring-2 focus-visible:ring-primary w-full sm:w-auto"
          >
            <Terminal className="h-3.5 w-3.5" />
            <span>Open Quick Review</span>
          </button>
        </div>

        {/* Security & Evidence Compliance Box */}
        <div className="mt-8 rounded-lg border border-border/80 bg-surface/60 p-4 sm:p-5 flex items-start gap-3 text-xs font-code text-themeText-muted">
          <ShieldCheck className="h-4 w-4 text-primary flex-shrink-0 mt-0.5" />
          <div className="space-y-1">
            <p className="font-semibold text-themeText-primary">
              Strict Zero-Trust & Evidence Integrity Standard
            </p>
            <p className="font-body text-themeText-muted leading-relaxed">
              Every technical claim, project metric, and architectural role presented on this site is anchored to verifiable repositories and source documentation. Non-public credentials or personal identifiers are strictly withheld.
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
