"use client";

import { useState } from "react";
import ProjectCard from "./ProjectCard";
import type { Project } from "arham-porto-schema";

interface ProjectsSectionProps {
  projects: Project[];
  hideHeader?: boolean;
}

export default function ProjectsSection({ projects, hideHeader = false }: ProjectsSectionProps) {
  const [activeFilter, setActiveFilter] = useState<string>("All");

  const filterCategories = ["All", "AI", "Full-Stack", "Backend", "Systems"];

  const filteredProjects = activeFilter === "All"
    ? projects
    : projects.filter((p) => p.category === activeFilter);

  return (
    <section id="projects" aria-label="Selected Engineering Work" className={`${hideHeader ? "py-6" : "py-20 border-b border-border"} px-4 sm:px-6 lg:px-8`}>
      <div className="mx-auto max-w-7xl">
        {/* Section Header or Compact Filter Bar */}
        {!hideHeader ? (
          <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-12">
            <div>
              <span className="text-xs font-code text-primary uppercase tracking-wider block mb-2">
                01 // SHOWCASE
              </span>
              <h2 className="font-display text-3xl sm:text-4xl font-extrabold text-themeText-primary tracking-tight">
                Selected Engineering Work
              </h2>
              <p className="mt-2 text-sm sm:text-base font-body text-themeText-muted max-w-2xl">
                Architectural case studies demonstrating backend scalability, full-stack systems, and AI engineering.
              </p>
            </div>

            {/* Filter Pills */}
            <div className="flex flex-wrap items-center gap-2" role="tablist" aria-label="Filter projects by domain">
              {filterCategories.map((category) => {
                const isActive = activeFilter === category;
                return (
                  <button
                    key={category}
                    type="button"
                    role="tab"
                    aria-selected={isActive}
                    onClick={() => setActiveFilter(category)}
                    className={`px-3.5 py-1.5 rounded-md text-xs font-code transition-all focus-visible:ring-2 focus-visible:ring-primary ${
                      isActive
                        ? "bg-primary text-white font-semibold shadow-sm"
                        : "bg-surface border border-border text-themeText-body hover:text-themeText-primary hover:border-primary/40 shadow-sm"
                    }`}
                  >
                    {category}
                  </button>
                );
              })}
            </div>
          </div>
        ) : (
          <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 mb-8 pb-4 border-b border-border/80">
            <span className="text-xs font-code text-themeText-muted uppercase tracking-wider">
              FILTER BY DOMAIN
            </span>
            <div className="flex flex-wrap items-center gap-2" role="tablist" aria-label="Filter projects by domain">
              {filterCategories.map((category) => {
                const isActive = activeFilter === category;
                return (
                  <button
                    key={category}
                    type="button"
                    role="tab"
                    aria-selected={isActive}
                    onClick={() => setActiveFilter(category)}
                    className={`px-3.5 py-1.5 rounded-md text-xs font-code transition-all focus-visible:ring-2 focus-visible:ring-primary ${
                      isActive
                        ? "bg-primary text-white font-semibold shadow-sm"
                        : "bg-surface border border-border text-themeText-body hover:text-themeText-primary hover:border-primary/40 shadow-sm"
                    }`}
                  >
                    {category}
                  </button>
                );
              })}
            </div>
          </div>
        )}

        {/* Project Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {filteredProjects.map((project) => (
            <ProjectCard key={project.id} project={project} />
          ))}
        </div>

        {filteredProjects.length === 0 && (
          <div className="rounded-card border border-border bg-surface p-12 text-center text-themeText-muted">
            <p className="text-sm font-body">No projects found matching the selected category filter.</p>
          </div>
        )}
      </div>
    </section>
  );
}
