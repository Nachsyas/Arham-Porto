import { Briefcase, Clock } from "lucide-react";
import type { ExperienceItem } from "arham-porto-schema";

interface ExperienceSectionProps {
  experiences: ExperienceItem[];
}

export default function ExperienceSection({ experiences }: ExperienceSectionProps) {
  return (
    <section id="experience" aria-label="Professional Experience" className="py-20 border-b border-border px-4 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-7xl">
        {/* Section Header */}
        <div className="mb-12">
          <span className="text-xs font-code text-primary uppercase tracking-wider block mb-2">
            03 // TRACK RECORD
          </span>
          <h2 className="font-display text-3xl sm:text-4xl font-extrabold text-themeText-primary tracking-tight">
            Professional Experience
          </h2>
          <p className="mt-2 text-sm sm:text-base font-body text-themeText-muted max-w-2xl">
            Verified engineering roles, responsibilities, and key architectural contributions.
          </p>
        </div>

        {/* Dynamic Content or Graceful Empty State */}
        {experiences.length > 0 ? (
          <div className="space-y-6">
            {experiences.map((item) => (
              <div
                key={item.id}
                className="rounded-card border border-border bg-surface p-6 flex flex-col md:flex-row md:items-start justify-between gap-4"
              >
                <div>
                  <h3 className="font-display text-xl font-bold text-themeText-primary">
                    {item.role ?? "Software Engineer"}
                  </h3>
                  <p className="text-sm font-medium text-primary mt-1">
                    {item.company}
                  </p>
                  {item.description && (
                    <p className="mt-3 text-sm text-themeText-body leading-relaxed max-w-3xl">
                      {item.description}
                    </p>
                  )}
                  {item.skills.length > 0 && (
                    <div className="mt-4 flex flex-wrap gap-1.5">
                      {item.skills.map((s) => (
                        <span
                          key={s}
                          className="px-2 py-0.5 rounded text-[10px] font-code bg-surfaceElevated border border-border text-themeText-muted"
                        >
                          {s}
                        </span>
                      ))}
                    </div>
                  )}
                </div>
                {item.period && (
                  <div className="flex items-center gap-1 text-xs font-code text-themeText-muted whitespace-nowrap">
                    <Clock className="h-3.5 w-3.5" />
                    <span>{item.period}</span>
                  </div>
                )}
              </div>
            ))}
          </div>
        ) : (
          /* Dignified Graceful Empty State */
          <div className="rounded-card border border-border bg-surface p-8 sm:p-12 text-center max-w-2xl mx-auto">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-surfaceElevated border border-border mx-auto mb-4 text-primary">
              <Briefcase className="h-6 w-6" />
            </div>
            <h3 className="font-display text-lg font-semibold text-themeText-primary mb-2">
              Experience Records in Verification
            </h3>
            <p className="text-xs sm:text-sm font-body text-themeText-muted leading-relaxed">
              Official professional employment records and enterprise contributions are currently undergoing verification against primary source documentation.
            </p>
            <span className="inline-block mt-4 text-[10px] font-code px-3 py-1 rounded bg-surfaceElevated text-primary border border-primary/20">
              STATUS: PENDING USER CONFIRMATION
            </span>
          </div>
        )}
      </div>
    </section>
  );
}
