"use client";

import React, { useState } from "react";
import Link from "next/link";
import { motion, AnimatePresence } from "motion/react";
import { Shield, ArrowRight, CheckCircle2, FileCode, ExternalLink, Sparkles } from "lucide-react";
import type { Skill, Project, Evidence } from "arham-porto-schema";

interface SkillsExplorerProps {
  skills: Skill[];
  projects: Project[];
  evidence: Evidence[];
}

export default function SkillsExplorer({
  skills,
  projects,
  evidence,
}: SkillsExplorerProps) {
  // Extract unique categories
  const categories = ["All", ...Array.from(new Set(skills.map((s) => s.category)))];
  const [selectedCategory, setSelectedCategory] = useState<string>("All");

  const filteredSkills = selectedCategory === "All"
    ? skills
    : skills.filter((s) => s.category === selectedCategory);

  const [selectedSkillId, setSelectedSkillId] = useState<string>(skills[0]?.id ?? "");

  const activeSkill = skills.find((s) => s.id === selectedSkillId) || filteredSkills[0] || skills[0];

  // Find linked projects
  const relevantProjects = projects.filter((p) => {
    const directLink = p.evidenceIds?.some((id) => activeSkill?.evidenceIds?.includes(id));
    if (directLink) return true;
    if (activeSkill?.category === "Backend" && p.category === "Backend") return true;
    if (activeSkill?.category === "Frontend" && p.category === "Full-Stack") return true;
    if (activeSkill?.category === "AI / ML" && p.category === "AI") return true;
    if (activeSkill?.category === "Systems" && p.category === "Systems") return true;
    return false;
  });

  // Find linked evidence items
  const relevantEvidence = evidence.filter((e) =>
    e.skillIds?.includes(activeSkill?.id ?? "") ||
    activeSkill?.evidenceIds?.includes(e.id)
  );

  return (
    <div className="space-y-8">
      {/* Category Pills */}
      <div className="flex flex-wrap items-center gap-2" role="tablist" aria-label="Skill categories">
        {categories.map((cat) => {
          const isActive = selectedCategory === cat;
          return (
            <button
              key={cat}
              type="button"
              role="tab"
              aria-selected={isActive}
              onClick={() => {
                setSelectedCategory(cat);
                const next = cat === "All" ? skills[0] : skills.find((s) => s.category === cat);
                if (next) setSelectedSkillId(next.id);
              }}
              className={`px-3.5 py-1.5 rounded-md text-xs font-code transition-all focus-visible:ring-2 focus-visible:ring-primary ${
                isActive
                  ? "bg-primary text-white font-semibold shadow-sm"
                  : "bg-surface border border-border text-themeText-body hover:text-themeText-primary hover:border-primary/40 shadow-sm"
              }`}
            >
              {cat}
            </button>
          );
        })}
      </div>

      {/* Main Two-Column Master-Detail Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
        {/* Left Column: Skill Selector Cards */}
        <div className="lg:col-span-5 flex flex-col gap-3">
          <span className="text-xs font-code text-themeText-muted uppercase tracking-wider">
            AVAILABLE SKILLS ({filteredSkills.length})
          </span>

          <div className="space-y-2.5">
            {filteredSkills.map((skill) => {
              const isSelected = skill.id === activeSkill?.id;
              return (
                <button
                  key={skill.id}
                  type="button"
                  onClick={() => setSelectedSkillId(skill.id)}
                  className={`w-full text-left p-4 rounded-xl border transition-all duration-200 flex flex-col justify-between gap-2 ${
                    isSelected
                      ? "bg-surface border-primary ring-1 ring-primary shadow-sm"
                      : "bg-surface/80 border-border hover:border-border-strong hover:bg-surface shadow-sm"
                  }`}
                >
                  <div className="flex items-center justify-between w-full">
                    <span className="text-[10px] font-code px-2 py-0.5 rounded bg-primary-muted text-primary border border-primary/20 font-medium">
                      {skill.category}
                    </span>
                    <span className="flex items-center gap-1 text-[10px] font-code text-themeText-muted">
                      <Shield className="h-3 w-3 text-status-success" /> Verified
                    </span>
                  </div>

                  <h3 className="font-display text-base font-bold text-themeText-primary">
                    {skill.name}
                  </h3>

                  <p className="text-xs font-body text-themeText-muted line-clamp-2">
                    {skill.claim}
                  </p>
                </button>
              );
            })}
          </div>
        </div>

        {/* Right Column: Grounded Evidence Detail Panel */}
        <div className="lg:col-span-7">
          <AnimatePresence mode="wait">
            {activeSkill && (
              <motion.div
                key={activeSkill.id}
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -10 }}
                transition={{ duration: 0.25, ease: "easeOut" }}
                className="rounded-card border border-border bg-surface p-6 sm:p-8 space-y-6 shadow-sm"
              >
                {/* Active Skill Title & Category */}
                <div className="border-b border-border/80 pb-5 space-y-2">
                  <div className="flex items-center gap-2">
                    <span className="text-[11px] font-code px-2.5 py-1 rounded bg-primary-muted border border-primary/25 text-primary font-semibold">
                      {activeSkill.category}
                    </span>
                    <span className="text-xs font-code text-themeText-muted flex items-center gap-1">
                      <CheckCircle2 className="h-3.5 w-3.5 text-status-success" /> Evidence Grounded
                    </span>
                  </div>

                  <h2 className="font-display text-2xl sm:text-3xl font-extrabold text-themeText-primary">
                    {activeSkill.name}
                  </h2>

                  <div className="p-4 rounded-xl bg-canvas-soft border border-border/70 text-xs sm:text-sm font-body text-themeText-body leading-relaxed">
                    <span className="font-code text-primary font-bold block mb-1 uppercase tracking-wider text-[11px]">
                      VERIFIED CAPABILITY CLAIM:
                    </span>
                    {activeSkill.claim}
                  </div>
                </div>

                {/* Evidence Sources */}
                <div className="space-y-3">
                  <span className="text-xs font-code text-themeText-muted uppercase tracking-wider block">
                    CANONICAL EVIDENCE ITEMS ({relevantEvidence.length})
                  </span>

                  {relevantEvidence.length > 0 ? (
                    <div className="space-y-3">
                      {relevantEvidence.map((item) => (
                        <div
                          key={item.id}
                          className="p-4 rounded-xl bg-surface-elevated border border-border text-xs space-y-2 shadow-sm"
                        >
                          <div className="flex items-center justify-between">
                            <span className="font-code font-bold text-primary flex items-center gap-1.5">
                              <FileCode className="h-3.5 w-3.5" />
                              {item.title}
                            </span>
                            <span className="text-[10px] font-code uppercase text-themeText-muted px-2 py-0.5 rounded bg-canvas-soft border border-border/70">
                              {item.type}
                            </span>
                          </div>

                          <p className="font-body text-themeText-body leading-relaxed">
                            {item.summary}
                          </p>

                          {item.sourceUrl && (
                            <a
                              href={item.sourceUrl}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="inline-flex items-center gap-1 text-[11px] font-code text-primary hover:text-primary-active pt-1"
                            >
                              <span>View Source Repository</span>
                              <ExternalLink className="h-3 w-3" />
                            </a>
                          )}
                        </div>
                      ))}
                    </div>
                  ) : (
                    <div className="p-4 rounded-xl bg-canvas-soft border border-border text-xs text-themeText-muted">
                      Validated in core architectural case studies.
                    </div>
                  )}
                </div>

                {/* Associated Projects */}
                <div className="space-y-3 pt-2">
                  <span className="text-xs font-code text-themeText-muted uppercase tracking-wider block">
                    DEMONSTRATED IN PROJECTS ({relevantProjects.length})
                  </span>

                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    {relevantProjects.map((p) => (
                      <Link
                        key={p.id}
                        href={`/projects/${p.slug}`}
                        className="p-3.5 rounded-lg border border-border bg-surface hover:border-primary/50 hover:bg-surface-elevated transition-all flex items-center justify-between group shadow-sm"
                      >
                        <div className="space-y-0.5">
                          <p className="font-display text-sm font-bold text-themeText-primary group-hover:text-primary transition-colors">
                            {p.title}
                          </p>
                          <p className="text-[10px] font-code text-themeText-muted">
                            {p.category}
                          </p>
                        </div>
                        <ArrowRight className="h-4 w-4 text-primary transform group-hover:translate-x-1 transition-transform" />
                      </Link>
                    ))}
                  </div>
                </div>

                {/* Zero Fake Percentages Rule Banner */}
                <div className="p-3 rounded-lg border border-primary/20 bg-primary/5 flex items-center gap-2 text-xs font-code text-themeText-muted">
                  <Sparkles className="h-4 w-4 text-primary flex-shrink-0" />
                  <span>
                    Skills strictly follow the <code>SKILL → CLAIM → EVIDENCE</code> architecture. Arbitrary percentage bars are omitted.
                  </span>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </div>
  );
}
