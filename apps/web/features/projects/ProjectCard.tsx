"use client";

import NavLink from "@/components/motion/NavLink";
import { ArrowRight, CheckCircle2 } from "lucide-react";
import { GithubIcon } from "@/components/icons/GithubIcon";
import type { Project } from "arham-porto-schema";

interface ProjectCardProps {
  project: Project;
}

export default function ProjectCard({ project }: ProjectCardProps) {
  return (
    <article
      aria-label={`Project case study: ${project.title}`}
      data-testid="project-card"
      className="group relative flex flex-col justify-between rounded-card border border-border bg-surface p-6 sm:p-7 hover:border-primary/50 transition-all duration-300 hover:shadow-md transform hover:-translate-y-1"
    >
      {/* Top Meta */}
      <div>
        <div className="flex items-center justify-between gap-3 mb-3">
          {project.category && (
            <span className="text-[11px] font-code px-2.5 py-0.5 rounded bg-primary-muted border border-primary/20 text-primary font-semibold tracking-wide">
              {project.category}
            </span>
          )}
          {project.featured && (
            <span className="text-[10px] font-code text-themeText-muted flex items-center gap-1 bg-canvas-soft px-2 py-0.5 rounded border border-border/60">
              <CheckCircle2 className="h-3 w-3 text-status-success" /> Flagship
            </span>
          )}
        </div>

        <h3 className="font-display text-xl sm:text-2xl font-bold text-themeText-primary group-hover:text-primary transition-colors">
          {project.title}
        </h3>

        {project.summary && (
          <p className="mt-2.5 text-xs sm:text-sm font-body text-themeText-body leading-relaxed">
            {project.summary}
          </p>
        )}

        {/* Problem / Solution Snapshot */}
        {project.solution && (
          <div className="mt-4 rounded-md border border-border bg-canvas-soft p-3 text-xs text-themeText-muted transition-colors group-hover:border-border-strong">
            <span className="font-code text-primary font-semibold block mb-1">Architecture Solution:</span>
            <p className="line-clamp-2 leading-relaxed text-themeText-body">{project.solution}</p>
          </div>
        )}

        {/* Technology Badges */}
        {project.technologies.length > 0 && (
          <div className="mt-5 flex flex-wrap gap-1.5">
            {project.technologies.map((tech) => (
              <span
                key={tech}
                className="px-2.5 py-0.5 rounded text-[10px] font-code bg-canvas-soft border border-border text-themeText-body group-hover:border-primary/20 transition-colors"
              >
                {tech}
              </span>
            ))}
          </div>
        )}
      </div>

      {/* Footer Actions */}
      <div className="mt-6 pt-4 border-t border-border/80 flex items-center justify-between">
        <NavLink
          href={`/projects/${project.slug}`}
          className="inline-flex items-center gap-1.5 text-xs font-code font-semibold text-primary hover:text-primary-active transition-colors focus-visible:ring-2 focus-visible:ring-primary rounded py-1"
        >
          <span>Explore Case Study</span>
          <ArrowRight className="h-3.5 w-3.5 transform group-hover:translate-x-1.5 transition-transform" />
        </NavLink>

        <div className="flex items-center gap-2">
          {project.githubUrl && (
            <a
              href={project.githubUrl}
              target="_blank"
              rel="noopener noreferrer"
              aria-label={`GitHub repository for ${project.title}`}
              className="flex h-9 w-9 items-center justify-center rounded border border-border bg-surface text-themeText-muted hover:text-themeText-primary hover:border-primary transition-all shadow-sm focus-visible:ring-2 focus-visible:ring-primary hover:scale-105"
            >
              <GithubIcon className="h-4 w-4" />
            </a>
          )}
        </div>
      </div>
    </article>
  );
}
