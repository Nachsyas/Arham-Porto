"use client";

import { useEffect } from "react";
import { motion, AnimatePresence } from "motion/react";
import { X, ExternalLink, CheckCircle2, Shield, ArrowRight } from "lucide-react";
import { GithubIcon } from "@/components/icons/GithubIcon";
import Link from "next/link";
import type { Profile, Project } from "arham-porto-schema";

interface QuickReviewDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  profile: Profile;
  projects: Project[];
}

export default function QuickReviewDrawer({
  isOpen,
  onClose,
  profile,
  projects,
}: QuickReviewDrawerProps) {
  // Close on Escape key press
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    if (isOpen) {
      document.body.style.overflow = "hidden";
      window.addEventListener("keydown", handleKeyDown);
    }
    return () => {
      document.body.style.overflow = "unset";
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [isOpen, onClose]);

  const coreStack = [
    "Go (Golang)",
    "TypeScript",
    "Next.js",
    "React",
    "PostgreSQL",
    "REST APIs",
    "Clean Architecture",
    "Docker",
  ];

  return (
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-50 flex justify-end">
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            onClick={onClose}
            className="fixed inset-0 bg-canvas/80 backdrop-blur-sm"
            aria-hidden="true"
          />

          {/* Drawer Panel */}
          <motion.div
            initial={{ x: "100%" }}
            animate={{ x: 0 }}
            exit={{ x: "100%" }}
            transition={{ type: "spring", damping: 28, stiffness: 280 }}
            role="dialog"
            aria-modal="true"
            aria-labelledby="quick-review-title"
            className="relative z-10 flex h-full w-full max-w-lg flex-col bg-surface border-l border-border shadow-2xl overflow-y-auto"
          >
            {/* Header */}
            <div className="flex items-center justify-between border-b border-border p-6 bg-surfaceElevated/50">
              <div className="flex items-center gap-2">
                <span className="flex h-2 w-2 rounded-full bg-primary animate-pulse" />
                <h2 id="quick-review-title" className="font-display font-bold text-lg text-themeText-primary">
                  60-Second Quick Review
                </h2>
              </div>
              <button
                type="button"
                onClick={onClose}
                aria-label="Close review drawer"
                className="flex h-10 w-10 items-center justify-center rounded-md border border-border text-themeText-muted hover:text-themeText-primary hover:border-primary transition-colors focus-visible:ring-2 focus-visible:ring-primary"
              >
                <X className="h-5 w-5" />
              </button>
            </div>

            {/* Content Body */}
            <div className="flex-1 space-y-6 p-6">
              {/* Identity Snapshot */}
              <div className="rounded-card border border-border bg-surfaceElevated p-5">
                <p className="text-xs font-code tracking-wider text-primary uppercase mb-1">
                  Candidate Profile
                </p>
                <h3 className="font-display text-2xl font-bold text-themeText-primary">
                  {profile.fullName}
                </h3>
                <p className="text-sm font-medium text-primary mt-0.5">
                  {profile.role}
                </p>
                {profile.positioning && (
                  <p className="text-xs text-themeText-muted mt-2 leading-relaxed">
                    {profile.positioning}
                  </p>
                )}
              </div>

              {/* Verified Core Stack */}
              <div>
                <h4 className="text-xs font-code text-themeText-muted uppercase tracking-wider mb-3">
                  Core Engineering Stack
                </h4>
                <div className="flex flex-wrap gap-2">
                  {coreStack.map((tech) => (
                    <span
                      key={tech}
                      className="inline-flex items-center gap-1.5 px-3 py-1 rounded-md text-xs font-code bg-canvas border border-border text-themeText-primary"
                    >
                      <CheckCircle2 className="h-3 w-3 text-status-success" />
                      {tech}
                    </span>
                  ))}
                </div>
              </div>

              {/* Selected Work Shortlist */}
              <div>
                <h4 className="text-xs font-code text-themeText-muted uppercase tracking-wider mb-3">
                  Selected Work Showcase
                </h4>
                <div className="space-y-3">
                  {projects.map((proj) => (
                    <div
                      key={proj.id}
                      className="rounded-md border border-border bg-surfaceElevated/40 p-3 hover:border-borderStrong transition-colors"
                    >
                      <div className="flex items-center justify-between mb-1">
                        <span className="font-display font-semibold text-sm text-themeText-primary">
                          {proj.title}
                        </span>
                        {proj.category && (
                          <span className="text-[10px] font-code px-2 py-0.5 rounded bg-primary/10 text-primary border border-primary/20">
                            {proj.category}
                          </span>
                        )}
                      </div>
                      {proj.summary && (
                        <p className="text-xs text-themeText-muted line-clamp-2 mb-2 leading-relaxed">
                          {proj.summary}
                        </p>
                      )}
                      <div className="flex items-center justify-between text-xs pt-2 border-t border-border/50">
                        <Link
                          href={`/projects/${proj.slug}`}
                          onClick={onClose}
                          className="text-primary hover:text-primaryActive inline-flex items-center gap-1 font-medium focus-visible:ring-1 focus-visible:ring-primary rounded"
                        >
                          Case Study <ArrowRight className="h-3 w-3" />
                        </Link>
                        {proj.githubUrl && (
                          <a
                            href={proj.githubUrl}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-themeText-muted hover:text-themeText-primary inline-flex items-center gap-1"
                          >
                            <GithubIcon className="h-3 w-3" /> Repo
                          </a>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              {/* Evidence Architecture Note */}
              <div className="rounded-md border border-primary/20 bg-primary/5 p-4 text-xs text-themeText-muted leading-relaxed">
                <div className="flex items-center gap-2 text-primary font-medium mb-1">
                  <Shield className="h-4 w-4" /> Evidence-Grounded Portfolio
                </div>
                Skills and achievements in this portfolio are validated against verified project repositories and case studies. Zero arbitrary skill percentage scores are claimed.
              </div>
            </div>

            {/* Footer Actions */}
            <div className="border-t border-border p-6 bg-surfaceElevated/50 flex flex-col gap-3">
              {profile.github && (
                <a
                  href={profile.github}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center justify-center gap-2 w-full py-2.5 px-4 rounded-md border border-border bg-surface text-themeText-primary hover:border-primary text-sm font-medium transition-colors focus-visible:ring-2 focus-visible:ring-primary"
                >
                  <GithubIcon className="h-4 w-4" /> Visit GitHub Profile
                </a>
              )}
              <button
                type="button"
                onClick={onClose}
                className="w-full py-2 px-4 text-center text-xs text-themeText-muted hover:text-themeText-primary transition-colors"
              >
                Close Quick Review
              </button>
            </div>
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  );
}
