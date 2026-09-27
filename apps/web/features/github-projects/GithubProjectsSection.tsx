"use client";

import { useEffect, useState } from "react";
import { AlertCircle, RefreshCw } from "lucide-react";
import GithubProjectCard from "./GithubProjectCard";
import { getGithubProjects } from "./github-projects.api";
import type { GithubProjectItem } from "./github-projects.types";
import { GithubIcon } from "@/components/icons/GithubIcon";

interface GithubProjectsSectionProps {
  initialProjects?: GithubProjectItem[];
}

export default function GithubProjectsSection({ initialProjects }: GithubProjectsSectionProps) {
  const [projects, setProjects] = useState<GithubProjectItem[]>(initialProjects || []);
  const [isLoading, setIsLoading] = useState<boolean>(!initialProjects);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (initialProjects && initialProjects.length > 0) {
      return;
    }

    const controller = new AbortController();

    async function loadProjects() {
      setIsLoading(true);
      setError(null);
      try {
        const data = await getGithubProjects(1, 20, controller.signal);
        setProjects(data);
      } catch (err: unknown) {
        if (controller.signal.aborted) return;
        // Non-blocking error handling
        setError("GitHub projects are temporarily unavailable.");
      } finally {
        if (!controller.signal.aborted) {
          setIsLoading(false);
        }
      }
    }

    loadProjects();

    return () => controller.abort();
  }, [initialProjects]);

  return (
    <section
      id="github-repos"
      aria-label="GitHub Repositories and Open Source Work"
      data-testid="github-projects-section"
      className="mt-16 pt-12 border-t border-border px-4 sm:px-6 lg:px-8"
    >
      {/* Header */}
      <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4 mb-8">
        <div>
          <span className="text-xs font-code text-primary uppercase tracking-wider block mb-2">
            02 // OPEN SOURCE & REPOSITORIES
          </span>
          <h2 className="font-display text-2xl sm:text-3xl font-extrabold text-themeText-primary tracking-tight flex items-center gap-2.5">
            <GithubIcon className="h-6 w-6 text-themeText-primary" />
            <span>GitHub Repositories</span>
          </h2>
          <p className="mt-2 text-xs sm:text-sm font-body text-themeText-muted max-w-2xl">
            Live repositories automatically synchronized from Nachsyas GitHub profile as verifiable engineering evidence.
          </p>
        </div>

        <a
          href="https://github.com/Nachsyas"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 text-xs font-code px-3.5 py-2 rounded-md bg-surface border border-border text-themeText-body hover:text-themeText-primary hover:border-primary/50 transition-all shadow-sm w-fit"
        >
          <GithubIcon className="h-4 w-4" />
          <span>github.com/Nachsyas</span>
        </a>
      </div>

      {/* Loading Skeleton */}
      {isLoading && (
        <div data-testid="github-projects-loading" className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {[1, 2, 3, 4].map((i) => (
            <div
              key={i}
              className="rounded-card border border-border bg-surface p-6 sm:p-7 animate-pulse flex flex-col justify-between h-48"
            >
              <div>
                <div className="flex items-center justify-between gap-3 mb-4">
                  <div className="h-4 w-20 bg-canvas-soft rounded" />
                  <div className="h-4 w-12 bg-canvas-soft rounded" />
                </div>
                <div className="h-6 w-48 bg-canvas-soft rounded mb-3" />
                <div className="h-3 w-full bg-canvas-soft rounded mb-2" />
                <div className="h-3 w-3/4 bg-canvas-soft rounded" />
              </div>
              <div className="flex items-center justify-between pt-4 border-t border-border/60">
                <div className="h-3 w-28 bg-canvas-soft rounded" />
                <div className="h-4 w-24 bg-canvas-soft rounded" />
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Error Fallback Banner */}
      {!isLoading && error && (
        <div
          data-testid="github-projects-error"
          role="alert"
          className="rounded-card border border-status-warning/30 bg-status-warning/5 p-6 text-center sm:text-left flex flex-col sm:flex-row items-center gap-4"
        >
          <div className="h-10 w-10 rounded-full bg-status-warning/10 border border-status-warning/20 flex items-center justify-center shrink-0">
            <AlertCircle className="h-5 w-5 text-status-warning" />
          </div>
          <div className="flex-1">
            <h3 className="text-sm font-semibold text-themeText-primary font-display">
              Synchronization Notice
            </h3>
            <p className="text-xs text-themeText-body mt-0.5 font-body">
              {error} Canonical projects and case studies remain fully accessible above.
            </p>
          </div>
          <a
            href="https://github.com/Nachsyas"
            target="_blank"
            rel="noopener noreferrer"
            className="text-xs font-code text-primary hover:underline flex items-center gap-1.5 shrink-0"
          >
            <span>Browse on GitHub</span>
          </a>
        </div>
      )}

      {/* Success Grid */}
      {!isLoading && !error && projects.length > 0 && (
        <div data-testid="github-projects-grid" className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {projects.map((project) => (
            <GithubProjectCard key={project.id} project={project} />
          ))}
        </div>
      )}

      {/* Empty State */}
      {!isLoading && !error && projects.length === 0 && (
        <div className="rounded-card border border-border bg-surface p-10 text-center text-themeText-muted">
          <p className="text-sm font-body">No synchronized repositories found at this moment.</p>
        </div>
      )}
    </section>
  );
}
