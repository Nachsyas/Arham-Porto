"use client";

import { Star, GitFork, ExternalLink, Calendar } from "lucide-react";
import { GithubIcon } from "@/components/icons/GithubIcon";
import type { GithubProjectItem } from "./github-projects.types";

interface GithubProjectCardProps {
  project: GithubProjectItem;
}

function formatSyncDate(dateStr: string): string {
  try {
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return dateStr;
    return d.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  } catch {
    return dateStr;
  }
}

export default function GithubProjectCard({ project }: GithubProjectCardProps) {
  return (
    <article
      aria-label={`GitHub repository: ${project.name}`}
      data-testid="github-project-card"
      className="group relative flex flex-col justify-between rounded-card border border-border bg-surface p-6 sm:p-7 hover:border-primary/50 transition-all duration-300 hover:shadow-md transform hover:-translate-y-1"
    >
      {/* Top Meta */}
      <div>
        <div className="flex items-center justify-between gap-3 mb-3">
          <div className="flex items-center gap-2">
            {project.language ? (
              <span className="text-[11px] font-code px-2.5 py-0.5 rounded bg-primary-muted border border-primary/20 text-primary font-semibold tracking-wide">
                {project.language}
              </span>
            ) : (
              <span className="text-[11px] font-code px-2.5 py-0.5 rounded bg-canvas-soft border border-border text-themeText-muted font-medium">
                Repository
              </span>
            )}
            <span className="text-[10px] font-code text-themeText-muted flex items-center gap-1 bg-canvas-soft px-2 py-0.5 rounded border border-border/60">
              <GithubIcon className="h-3 w-3 text-themeText-muted" /> GitHub Sync
            </span>
          </div>

          {/* Metrics */}
          <div className="flex items-center gap-2.5 text-[11px] font-code text-themeText-muted">
            {project.stars > 0 && (
              <span className="flex items-center gap-1" title={`${project.stars} stars`}>
                <Star className="h-3 w-3 text-status-warning fill-status-warning" />
                <span>{project.stars}</span>
              </span>
            )}
            {project.forks > 0 && (
              <span className="flex items-center gap-1" title={`${project.forks} forks`}>
                <GitFork className="h-3 w-3 text-themeText-muted" />
                <span>{project.forks}</span>
              </span>
            )}
          </div>
        </div>

        <h3 className="font-display text-xl sm:text-2xl font-bold text-themeText-primary group-hover:text-primary transition-colors">
          {project.name}
        </h3>

        <p className="mt-2.5 text-xs sm:text-sm font-body text-themeText-body leading-relaxed line-clamp-3">
          {project.description || "Open source repository synchronized from verified GitHub profile."}
        </p>

        {/* Topics */}
        {project.topics && project.topics.length > 0 && (
          <div className="mt-4 flex flex-wrap gap-1.5" aria-label="Repository topics">
            {project.topics.map((topic) => (
              <span
                key={topic}
                className="px-2.5 py-0.5 rounded text-[10px] font-code bg-canvas-soft border border-border text-themeText-body group-hover:border-primary/20 transition-colors"
              >
                #{topic}
              </span>
            ))}
          </div>
        )}
      </div>

      {/* Footer Info & Action */}
      <div className="mt-6 pt-4 border-t border-border/80 flex items-center justify-between gap-4">
        {project.synced_at && (
          <span className="flex items-center gap-1.5 text-[11px] font-code text-themeText-muted" title="Last synchronized timestamp">
            <Calendar className="h-3 w-3 text-themeText-muted/70" />
            <span>Synced: {formatSyncDate(project.synced_at)}</span>
          </span>
        )}

        <a
          href={project.html_url}
          target="_blank"
          rel="noopener noreferrer"
          aria-label={`View ${project.name} repository on GitHub`}
          className="inline-flex items-center gap-1.5 text-xs font-code font-semibold text-primary hover:text-primary-active transition-colors focus-visible:ring-2 focus-visible:ring-primary rounded py-1 ml-auto"
        >
          <span>View Repository</span>
          <ExternalLink className="h-3.5 w-3.5 transform group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-transform" />
        </a>
      </div>
    </article>
  );
}
