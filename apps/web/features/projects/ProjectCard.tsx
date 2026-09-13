import Link from "next/link";
import { ArrowRight, Layers, CheckCircle2 } from "lucide-react";
import { GithubIcon } from "@/components/icons/GithubIcon";
import type { Project } from "arham-porto-schema";

interface ProjectCardProps {
  project: Project;
}

export default function ProjectCard({ project }: ProjectCardProps) {
  return (
    <article className="group relative flex flex-col justify-between rounded-card border border-border bg-surface p-6 hover:border-borderStrong transition-all duration-200 hover:shadow-[0_8px_30px_rgba(2,6,11,0.5)]">
      {/* Header */}
      <div>
        <div className="flex items-center justify-between gap-3 mb-3">
          {project.category && (
            <span className="text-[11px] font-code px-2.5 py-0.5 rounded bg-primary/10 border border-primary/20 text-primary font-medium">
              {project.category}
            </span>
          )}
          {project.featured && (
            <span className="text-[10px] font-code text-themeText-muted flex items-center gap-1">
              <CheckCircle2 className="h-3 w-3 text-status-success" /> Featured
            </span>
          )}
        </div>

        <h3 className="font-display text-xl font-bold text-themeText-primary group-hover:text-primary transition-colors">
          {project.title}
        </h3>

        {project.summary && (
          <p className="mt-2.5 text-sm font-body text-themeText-body leading-relaxed">
            {project.summary}
          </p>
        )}

        {/* Problem / Solution Snapshot */}
        {project.solution && (
          <div className="mt-4 rounded-md border border-border bg-surfaceElevated/50 p-3 text-xs text-themeText-muted">
            <span className="font-code text-primary font-medium block mb-1">Architecture Solution:</span>
            <p className="line-clamp-2 leading-relaxed">{project.solution}</p>
          </div>
        )}

        {/* Technology Badges */}
        {project.technologies.length > 0 && (
          <div className="mt-5 flex flex-wrap gap-1.5">
            {project.technologies.map((tech) => (
              <span
                key={tech}
                className="px-2 py-0.5 rounded text-[11px] font-code bg-canvas border border-border text-themeText-muted"
              >
                {tech}
              </span>
            ))}
          </div>
        )}
      </div>

      {/* Footer Actions */}
      <div className="mt-6 pt-4 border-t border-border/80 flex items-center justify-between">
        <Link
          href={`/projects/${project.slug}`}
          className="inline-flex items-center gap-1.5 text-sm font-medium text-primary hover:text-primaryActive transition-colors focus-visible:ring-2 focus-visible:ring-primary rounded py-1"
        >
          Case Study <ArrowRight className="h-4 w-4" />
        </Link>

        <div className="flex items-center gap-2">
          {project.githubUrl && (
            <a
              href={project.githubUrl}
              target="_blank"
              rel="noopener noreferrer"
              aria-label={`GitHub repository for ${project.title}`}
              className="flex h-9 w-9 items-center justify-center rounded border border-border text-themeText-muted hover:text-themeText-primary hover:border-primary transition-colors focus-visible:ring-2 focus-visible:ring-primary"
            >
              <GithubIcon className="h-4 w-4" />
            </a>
          )}
        </div>
      </div>
    </article>
  );
}
