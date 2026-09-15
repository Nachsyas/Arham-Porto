"use client";

import Link from "next/link";
import { ArrowRight, Code2, Cpu, Database, Globe, Compass, Sparkles, CheckCircle2, Shield } from "lucide-react";
import HeroSection from "@/features/hero/HeroSection";
import type { Profile, Project, Skill, JourneyStop } from "arham-porto-schema";

interface HomeViewProps {
  profile: Profile;
  featuredProjects: Project[];
  skills: Skill[];
  journeyStops: JourneyStop[];
}

export default function HomeView({
  profile,
  featuredProjects,
  skills,
  journeyStops,
}: HomeViewProps) {
  // Top 4 skill categories with icons
  const skillCategories = [
    { name: "Backend Engineering", icon: Database, desc: "Robust Go services, Clean Architecture, REST APIs" },
    { name: "Frontend & Web Apps", icon: Globe, desc: "Next.js App Router, TypeScript, Responsive Systems" },
    { name: "AI & Retrieval Systems", icon: Cpu, desc: "Grounded RAG, LLM integration, Vector Search" },
    { name: "Smart Contracts & Web3", icon: Code2, desc: "EVM, Solidity, Soulbound Tokens (ERC-5192)" },
  ];

  return (
    <div className="flex flex-col space-y-20 sm:space-y-28 pb-20">
      {/* 1. Cinematic Hero Identity */}
      <HeroSection profile={profile} />

      {/* 2. Recruiter Snapshot */}
      <section aria-labelledby="snapshot-heading" className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full">
        <div className="rounded-card border border-border bg-surface p-6 sm:p-10 shadow-sm">
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-center">
            <div className="lg:col-span-8 space-y-3">
              <div className="flex items-center gap-2 text-xs font-code font-bold text-primary tracking-wider uppercase">
                <Shield className="h-3.5 w-3.5" />
                <span>01 // RECRUITER SNAPSHOT</span>
              </div>
              <h2 id="snapshot-heading" className="font-display text-2xl sm:text-3xl font-bold text-themeText-primary">
                Engineering Focus & Core Competencies
              </h2>
              <p className="font-body text-sm sm:text-base text-themeText-body leading-relaxed max-w-3xl">
                Specialized in building decoupled backend services with Go and high-performance frontend interfaces with TypeScript and Next.js. Every claimed technical capability in this portfolio is grounded in verifiable repository commits and case studies.
              </p>

              {/* Verified core stack pills */}
              <div className="flex flex-wrap gap-2 pt-2">
                {["Go (Golang)", "TypeScript", "Next.js", "PostgreSQL", "Clean Architecture", "REST APIs", "Docker"].map((tech) => (
                  <span
                    key={tech}
                    className="inline-flex items-center gap-1.5 px-3 py-1 rounded-md text-xs font-code bg-canvas-soft border border-border text-themeText-primary shadow-sm"
                  >
                    <CheckCircle2 className="h-3 w-3 text-status-success" />
                    {tech}
                  </span>
                ))}
              </div>
            </div>

            <div className="lg:col-span-4 flex flex-col justify-center gap-3 p-5 rounded-xl bg-surface-elevated border border-border/80">
              <span className="text-xs font-code text-themeText-muted uppercase tracking-wider">
                CURRENT STATUS
              </span>
              <p className="font-display text-base font-semibold text-themeText-primary">
                Undergraduate in Computer Science
              </p>
              <p className="text-xs font-code text-themeText-muted">
                UIN Maulana Malik Ibrahim Malang
              </p>
              <div className="pt-2 border-t border-border flex items-center justify-between">
                <span className="text-xs font-code text-primary font-medium">
                  Verified Candidate
                </span>
                <span className="h-2 w-2 rounded-full bg-status-success" />
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* 3. Featured Work Preview */}
      <section aria-labelledby="work-preview-heading" className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full space-y-8">
        <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4 border-b border-border/70 pb-4">
          <div className="space-y-1">
            <span className="text-xs font-code font-bold text-primary tracking-wider uppercase">
              02 // SELECTED WORK
            </span>
            <h2 id="work-preview-heading" className="font-display text-2xl sm:text-3xl font-bold text-themeText-primary">
              Selected Engineering Work
            </h2>
            <p className="text-sm font-body text-themeText-muted max-w-xl">
              Production-grade systems demonstrating clean architecture, API design, and distributed persistence.
            </p>
          </div>
          <Link
            href="/work"
            className="inline-flex items-center gap-2 text-xs font-code font-semibold text-primary hover:text-primary-active transition-colors self-start sm:self-auto"
          >
            <span>View all projects</span>
            <ArrowRight className="h-3.5 w-3.5" />
          </Link>
        </div>

        {/* 2–3 Featured Cards Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {featuredProjects.map((project) => (
            <article
              key={project.id}
              className="group flex flex-col justify-between rounded-card border border-border bg-surface p-6 shadow-sm hover:border-primary/50 hover:shadow-md transition-all duration-300"
            >
              <div className="space-y-3">
                <div className="flex items-center justify-between">
                  <span className="px-2.5 py-0.5 rounded text-[11px] font-code font-semibold bg-primary-muted text-primary border border-primary/20">
                    {project.category}
                  </span>
                  <span className="text-[10px] font-code text-themeText-muted">FEATURED</span>
                </div>

                <h3 className="font-display text-xl font-bold text-themeText-primary group-hover:text-primary transition-colors">
                  {project.title}
                </h3>

                <p className="text-xs font-body text-themeText-body line-clamp-3 leading-relaxed">
                  {project.summary}
                </p>

                <div className="flex flex-wrap gap-1.5 pt-2">
                  {project.technologies.slice(0, 4).map((tech) => (
                    <span
                      key={tech}
                      className="px-2 py-0.5 rounded text-[10px] font-code bg-canvas-soft border border-border/80 text-themeText-muted"
                    >
                      {tech}
                    </span>
                  ))}
                </div>
              </div>

              <div className="pt-5 mt-5 border-t border-border/60 flex items-center justify-between">
                <Link
                  href={`/projects/${project.slug}`}
                  className="text-xs font-code font-semibold text-primary hover:text-primary-active inline-flex items-center gap-1 transition-colors"
                >
                  <span>Read Case Study</span>
                  <ArrowRight className="h-3 w-3 transform group-hover:translate-x-1 transition-transform" />
                </Link>
                {project.githubUrl && (
                  <span className="text-[10px] font-code text-themeText-muted">GitHub Verified</span>
                )}
              </div>
            </article>
          ))}
        </div>
      </section>

      {/* 4. Skills & Capabilities Preview */}
      <section aria-labelledby="skills-preview-heading" className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full space-y-8">
        <div className="flex flex-col sm:flex-row sm:items-end justify-between gap-4 border-b border-border/70 pb-4">
          <div className="space-y-1">
            <span className="text-xs font-code font-bold text-primary tracking-wider uppercase">
              03 // TECHNICAL COMPETENCE
            </span>
            <h2 id="skills-preview-heading" className="font-display text-2xl sm:text-3xl font-bold text-themeText-primary">
              Core Capabilities
            </h2>
            <p className="text-sm font-body text-themeText-muted max-w-xl">
              Structured engineering disciplines backed by verified evidence across repositories.
            </p>
          </div>
          <Link
            href="/skills"
            className="inline-flex items-center gap-2 text-xs font-code font-semibold text-primary hover:text-primary-active transition-colors self-start sm:self-auto"
          >
            <span>Explore skill evidence</span>
            <ArrowRight className="h-3.5 w-3.5" />
          </Link>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
          {skillCategories.map((cat) => {
            const Icon = cat.icon;
            return (
              <div
                key={cat.name}
                className="p-5 rounded-card border border-border bg-surface shadow-sm hover:border-primary/40 hover:shadow-sm transition-all"
              >
                <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-primary-muted border border-primary/25 text-primary mb-3">
                  <Icon className="h-5 w-5" />
                </div>
                <h3 className="font-display font-bold text-base text-themeText-primary mb-1">
                  {cat.name}
                </h3>
                <p className="text-xs font-body text-themeText-body leading-relaxed">
                  {cat.desc}
                </p>
              </div>
            );
          })}
        </div>
      </section>

      {/* 5. Journey Teaser Preview */}
      <section aria-labelledby="journey-preview-heading" className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full space-y-8">
        <div className="rounded-card border border-border bg-surface p-6 sm:p-10 shadow-sm">
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-center">
            <div className="lg:col-span-8 space-y-4">
              <span className="text-xs font-code font-bold text-primary tracking-wider uppercase">
                04 // GEOGRAPHIC & ACADEMIC ROUTE
              </span>
              <h2 id="journey-preview-heading" className="font-display text-2xl sm:text-3xl font-bold text-themeText-primary">
                From Central Java to East Java
              </h2>
              <p className="text-sm font-body text-themeText-body leading-relaxed max-w-2xl">
                Trace the educational milestones from Karanganyar and Jakarta to boarding school in Salatiga and university studies in Malang. Discover the interactive route map and formative timeline.
              </p>
              <div className="pt-2">
                <Link
                  href="/journey"
                  className="inline-flex items-center gap-2 px-5 py-2.5 rounded-md text-xs font-code font-semibold bg-primary text-white hover:bg-primary-active transition-all shadow-sm"
                >
                  <span>Explore Interactive Journey</span>
                  <Compass className="h-4 w-4" />
                </Link>
              </div>
            </div>

            <div className="lg:col-span-4 p-5 rounded-xl bg-canvas-soft border border-border/80 space-y-3">
              <span className="text-xs font-code text-themeText-muted uppercase tracking-wider">
                CORRIDOR HIGHLIGHTS
              </span>
              <div className="space-y-2 text-xs font-code text-themeText-body">
                <div className="flex items-center justify-between pb-1 border-b border-border/60">
                  <span>Origin: Karanganyar</span>
                  <span className="text-primary font-semibold">Jawa Tengah</span>
                </div>
                <div className="flex items-center justify-between pb-1 border-b border-border/60">
                  <span>Secondary: Salatiga</span>
                  <span className="text-primary font-semibold">Boarding School</span>
                </div>
                <div className="flex items-center justify-between">
                  <span>University: Malang</span>
                  <span className="text-primary font-semibold">Computer Science</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* 6. Contact CTA Callout */}
      <section aria-labelledby="contact-preview-heading" className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 w-full">
        <div className="rounded-card border border-primary/30 bg-primary-muted/40 p-8 sm:p-12 text-center space-y-4">
          <span className="text-xs font-code font-bold text-primary tracking-wider uppercase">
            05 // GET IN TOUCH
          </span>
          <h2 id="contact-preview-heading" className="font-display text-3xl sm:text-4xl font-extrabold text-themeText-primary">
            Let&apos;s Build Intelligent Systems Together
          </h2>
          <p className="text-sm sm:text-base font-body text-themeText-body max-w-xl mx-auto leading-relaxed">
            Open for software engineering opportunities, distributed backend collaboration, and systems development.
          </p>
          <div className="pt-2 flex justify-center gap-4">
            <Link
              href="/contact"
              className="inline-flex items-center gap-2 px-6 py-3 rounded-md text-sm font-semibold font-body bg-primary text-white hover:bg-primary-active transition-all shadow-sm hover:shadow-md"
            >
              <span>Connect on Contact Page</span>
              <ArrowRight className="h-4 w-4" />
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
}
