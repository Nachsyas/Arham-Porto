import { notFound } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, ExternalLink, Shield, CheckCircle2, Layers } from "lucide-react";
import { GithubIcon } from "@/components/icons/GithubIcon";
import { getProjects, getProfile } from "arham-porto-data";
import type { Metadata } from "next";

interface CaseStudyPageProps {
  params: Promise<{
    slug: string;
  }>;
}

export async function generateStaticParams() {
  const projects = getProjects();
  return projects.map((p) => ({
    slug: p.slug,
  }));
}

export async function generateMetadata({ params }: CaseStudyPageProps): Promise<Metadata> {
  const { slug } = await params;
  const projects = getProjects();
  const project = projects.find((p) => p.slug === slug);

  if (!project) {
    return {
      title: "Project Not Found | Arham Porto",
    };
  }

  return {
    title: `${project.title} — Engineering Case Study | Arham Porto`,
    description: project.summary ?? "Architectural case study by Nachsyas Arham Mumtaz Nashohi.",
  };
}

export default async function CaseStudyPage({ params }: CaseStudyPageProps) {
  const { slug } = await params;
  const projects = getProjects();
  const project = projects.find((p) => p.slug === slug);

  if (!project) {
    notFound();
  }

  return (
    <div className="min-h-screen bg-canvas text-themeText-body selection:bg-primary/20 selection:text-primary pb-24">
      {/* Top Breadcrumb Navigation */}
      <nav className="sticky top-0 z-30 border-b border-border bg-canvas/90 backdrop-blur-md px-4 sm:px-6 lg:px-8 py-4">
        <div className="mx-auto max-w-5xl flex items-center justify-between">
          <Link
            href="/#projects"
            className="inline-flex items-center gap-2 text-xs font-code text-themeText-muted hover:text-primary transition-colors focus-visible:ring-2 focus-visible:ring-primary rounded py-1 px-2"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to Selected Work
          </Link>
          <span className="text-xs font-code text-themeText-mutedSoft">
            CASE STUDY // {project.slug.toUpperCase()}
          </span>
        </div>
      </nav>

      <main className="mx-auto max-w-4xl px-4 sm:px-6 lg:px-8 pt-12 space-y-12">
        {/* Header Section */}
        <header className="space-y-4 border-b border-border pb-10">
          <div className="flex flex-wrap items-center gap-2.5">
            {project.category && (
              <span className="px-3 py-1 rounded text-xs font-code font-medium bg-primary/10 border border-primary/25 text-primary">
                {project.category}
              </span>
            )}
            {project.featured && (
              <span className="px-3 py-1 rounded text-xs font-code bg-surfaceElevated border border-border text-themeText-muted flex items-center gap-1.5">
                <CheckCircle2 className="h-3.5 w-3.5 text-status-success" /> Verified Case Study
              </span>
            )}
          </div>

          <h1 className="font-display text-3xl sm:text-5xl font-extrabold text-themeText-primary tracking-tight">
            {project.title}
          </h1>

          {project.summary && (
            <p className="font-body text-base sm:text-lg text-themeText-body leading-relaxed max-w-3xl">
              {project.summary}
            </p>
          )}

          {/* Quick External Links */}
          <div className="flex flex-wrap items-center gap-4 pt-2">
            {project.githubUrl && (
              <a
                href={project.githubUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 px-4 py-2 rounded-md text-xs font-code bg-surface border border-border text-themeText-primary hover:border-primary hover:text-primary transition-colors focus-visible:ring-2 focus-visible:ring-primary"
              >
                <GithubIcon className="h-4 w-4" /> View GitHub Repository <ExternalLink className="h-3 w-3 opacity-60" />
              </a>
            )}
            {project.demoUrl && (
              <a
                href={project.demoUrl}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 px-4 py-2 rounded-md text-xs font-code bg-primary text-canvas hover:bg-primaryActive transition-colors font-medium focus-visible:ring-2 focus-visible:ring-primary"
              >
                Live Demo <ExternalLink className="h-3 w-3" />
              </a>
            )}
          </div>
        </header>

        {/* Problem Statement (if verified) */}
        {project.problem && (
          <section aria-labelledby="problem-heading" className="rounded-card border border-border bg-surface p-6 sm:p-8 space-y-3">
            <h2 id="problem-heading" className="text-xs font-code text-primary uppercase tracking-wider">
              01 // The Problem & Context
            </h2>
            <p className="font-body text-sm sm:text-base text-themeText-body leading-relaxed">
              {project.problem}
            </p>
          </section>
        )}

        {/* Architectural Solution (if verified) */}
        {project.solution && (
          <section aria-labelledby="solution-heading" className="rounded-card border border-border bg-surface p-6 sm:p-8 space-y-3">
            <h2 id="solution-heading" className="text-xs font-code text-primary uppercase tracking-wider">
              02 // Architectural Solution
            </h2>
            <p className="font-body text-sm sm:text-base text-themeText-body leading-relaxed">
              {project.solution}
            </p>
          </section>
        )}

        {/* Verified Technologies */}
        {project.technologies.length > 0 && (
          <section aria-labelledby="tech-heading" className="space-y-4">
            <h2 id="tech-heading" className="text-xs font-code text-primary uppercase tracking-wider">
              03 // Technology Stack
            </h2>
            <div className="flex flex-wrap gap-2">
              {project.technologies.map((tech) => (
                <span
                  key={tech}
                  className="px-3.5 py-1.5 rounded-md text-xs font-code bg-surfaceElevated border border-border text-themeText-primary font-medium"
                >
                  {tech}
                </span>
              ))}
            </div>
          </section>
        )}

        {/* Verification & Attribution Note */}
        <section className="rounded-card border border-primary/20 bg-primary/5 p-6 space-y-2">
          <div className="flex items-center gap-2 text-primary font-medium text-sm">
            <Shield className="h-4 w-4" />
            Verification & Ownership Policy
          </div>
          <p className="text-xs text-themeText-muted leading-relaxed">
            Detailed team attribution, specific PR commit hashes, and benchmark performance metrics are linked directly to approved GitHub repositories. Speculative or unverified contributions are omitted.
          </p>
        </section>
      </main>
    </div>
  );
}
