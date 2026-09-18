import { Shield, CheckCircle2, ArrowRight } from "lucide-react";
import type { Skill, Project } from "arham-porto-schema";
import NavLink from "@/components/motion/NavLink";

interface SkillsSectionProps {
  skills: Skill[];
  projects: Project[];
}

export default function SkillsSection({ skills, projects }: SkillsSectionProps) {
  // Map project references strictly by canonical evidence linkage
  const getRelevantProjects = (skill: Skill) => {
    return projects.filter((p) => {
      const hasDirectEvidence = p.evidenceIds?.some((id) => skill.evidenceIds?.includes(id));
      return Boolean(hasDirectEvidence);
    });
  };

  return (
    <section id="skills" aria-label="Technical Skills & Evidence Chain" className="py-20 border-b border-border px-4 sm:px-6 lg:px-8 bg-canvas-soft/40">
      <div className="mx-auto max-w-7xl">
        {/* Section Header */}
        <div className="mb-12">
          <span className="text-xs font-code text-primary uppercase tracking-wider block mb-2">
            02 // EVIDENCE CHAIN
          </span>
          <h2 className="font-display text-3xl sm:text-4xl font-extrabold text-themeText-primary tracking-tight">
            Skill & Evidence Explorer
          </h2>
          <p className="mt-2 text-sm sm:text-base font-body text-themeText-muted max-w-3xl">
            Grounded technical capabilities mapped to verified engineering projects. Evaluated on demonstrated evidence rather than arbitrary percentage ratings.
          </p>
        </div>

        {/* Skills Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {skills.map((skill) => {
            const relevantProjects = getRelevantProjects(skill);
            const isVerified = skill.evidenceIds && skill.evidenceIds.length > 0;

            return (
              <div
                key={skill.id}
                className="flex flex-col justify-between rounded-card border border-border bg-surface p-6 hover:border-border-strong hover:shadow-sm transition-all"
              >
                <div>
                  <div className="flex items-center justify-between gap-2 mb-3">
                    <span className="text-[11px] font-code px-2.5 py-0.5 rounded bg-primary-muted border border-primary/20 text-primary font-medium">
                      {skill.category}
                    </span>
                    {isVerified && (
                      <span className="text-[10px] font-code text-themeText-muted flex items-center gap-1">
                        <Shield className="h-3 w-3 text-primary" /> Verified
                      </span>
                    )}
                  </div>

                  <h3 className="font-display text-lg font-bold text-themeText-primary mb-2">
                    {skill.name}
                  </h3>

                  {skill.claim && (
                    <p className="text-xs font-body text-themeText-body leading-relaxed mb-4">
                      {skill.claim}
                    </p>
                  )}
                </div>

                {/* Evidence Links */}
                <div className="pt-4 border-t border-border/60">
                  <span className="text-[10px] font-code text-themeText-muted uppercase tracking-wider block mb-2">
                    Demonstrated in Projects:
                  </span>
                  {relevantProjects.length > 0 ? (
                    <div className="space-y-1.5">
                      {relevantProjects.map((p) => (
                        <NavLink
                          key={p.id}
                          href={`/projects/${p.slug}`}
                          className="flex items-center justify-between text-xs text-themeText-primary hover:text-primary transition-colors py-1.5 px-2.5 rounded bg-canvas-soft hover:bg-primary-muted/40 border border-border/50"
                        >
                          <span className="font-medium">{p.title}</span>
                          <ArrowRight className="h-3 w-3 text-primary" />
                        </NavLink>
                      ))}
                    </div>
                  ) : (
                    <p className="text-[11px] text-themeText-muted font-body">
                      Core technical focus
                    </p>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
