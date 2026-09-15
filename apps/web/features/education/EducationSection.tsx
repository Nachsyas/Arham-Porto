import { GraduationCap, Clock } from "lucide-react";
import type { EducationItem } from "arham-porto-schema";

interface EducationSectionProps {
  education: EducationItem[];
}

export default function EducationSection({ education }: EducationSectionProps) {
  return (
    <section id="education" aria-label="Academic Education" className="py-20 border-b border-border px-4 sm:px-6 lg:px-8 bg-canvas-soft/40">
      <div className="mx-auto max-w-7xl">
        {/* Section Header */}
        <div className="mb-12">
          <span className="text-xs font-code text-primary uppercase tracking-wider block mb-2">
            04 // ACADEMIC FOUNDATION
          </span>
          <h2 className="font-display text-3xl sm:text-4xl font-extrabold text-themeText-primary tracking-tight">
            Education & Background
          </h2>
          <p className="mt-2 text-sm sm:text-base font-body text-themeText-muted max-w-2xl">
            Formal computer science and engineering coursework foundation.
          </p>
        </div>

        {/* Dynamic Content or Graceful Empty State */}
        {education.length > 0 ? (
          <div className="space-y-6">
            {education.map((item) => (
              <div
                key={item.id}
                className="rounded-card border border-border bg-surface p-6 flex flex-col md:flex-row md:items-start justify-between gap-4 shadow-sm"
              >
                <div>
                  <h3 className="font-display text-xl font-bold text-themeText-primary">
                    {item.institution ?? "University"}
                  </h3>
                  <p className="text-sm font-medium text-primary mt-1">
                    {item.degree} {item.field && `in ${item.field}`}
                  </p>
                  {item.description && (
                    <p className="mt-3 text-sm text-themeText-body leading-relaxed max-w-3xl">
                      {item.description}
                    </p>
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
          <div className="rounded-card border border-border bg-surface p-8 sm:p-12 text-center max-w-2xl mx-auto shadow-sm">
            <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-surface-elevated border border-border mx-auto mb-4 text-primary">
              <GraduationCap className="h-6 w-6" />
            </div>
            <h3 className="font-display text-lg font-semibold text-themeText-primary mb-2">
              Academic Background Verification
            </h3>
            <p className="text-xs sm:text-sm font-body text-themeText-muted leading-relaxed">
              Formal degree credentials, thesis topics, and academic institution records are currently undergoing verification before public listing.
            </p>
            <span className="inline-block mt-4 text-[10px] font-code px-3 py-1 rounded bg-surface-elevated text-primary border border-primary/20">
              STATUS: PENDING USER CONFIRMATION
            </span>
          </div>
        )}
      </div>
    </section>
  );
}
