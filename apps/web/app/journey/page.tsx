import { getJourney, getEducation } from "arham-porto-data";
import PageHeader from "@/components/PageHeader";
import JourneySection from "@/features/journey/JourneySection";
import InterPageNav from "@/components/InterPageNav";
import { GraduationCap, MapPin, Calendar, BookOpen } from "lucide-react";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Geographic & Academic Journey",
  description:
    "Interactive map tracing Nachsyas Arham Mumtaz Nashohi's academic and geographic journey across Indonesia — from Central Java to Jakarta, Salatiga, and Malang.",
  alternates: {
    canonical: "/journey",
  },
};

export default function JourneyPage() {
  const journeyStops = getJourney();
  const educationItems = getEducation();

  // Highlight public academic stops from journey data
  const academicStops = journeyStops.filter(
    (stop) =>
      stop.public &&
      (stop.category === "university" ||
        stop.category === "sma" ||
        stop.category === "smp" ||
        stop.category === "sd")
  );

  return (
    <div className="pb-24">
      <PageHeader
        eyebrow="JOURNEY // GEOGRAPHIC & ACADEMIC"
        title="Geographic & Academic Evolution"
        description="An interactive corridor tracing formative secondary education in Central Java through university computer science in East Java."
        badge="Java Corridor"
      />

      {/* Main Interactive Map & Milestone Story Card */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <JourneySection
          stops={journeyStops}
          hideHeader={true}
          className="pb-16 bg-transparent"
        />

        {/* Formative & Academic Timeline Overview */}
        <section aria-label="Academic Records" className="mt-8 border-t border-border pt-12">
          <div className="flex items-center gap-2 mb-6">
            <GraduationCap className="h-5 w-5 text-primary" />
            <h2 className="font-display text-xl sm:text-2xl font-bold text-themeText-primary tracking-tight">
              Academic & Formative Records
            </h2>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {academicStops.map((stop) => (
              <div
                key={stop.id}
                className="rounded-card border border-border bg-surface p-5 sm:p-6 shadow-sm hover:border-primary/40 transition-all flex flex-col justify-between"
              >
                <div>
                  <div className="flex items-center justify-between gap-2 mb-2">
                    <span className="px-2 py-0.5 rounded text-[10px] font-code font-semibold uppercase bg-surface-elevated border border-border text-primary">
                      {stop.category === "university"
                        ? "Undergraduate"
                        : stop.category === "sma"
                        ? "Secondary / Tahfizh"
                        : stop.category === "smp"
                        ? "Lower Secondary"
                        : "Primary"}
                    </span>
                    {stop.period && (
                      <span className="inline-flex items-center gap-1 text-[11px] font-code text-themeText-muted">
                        <Calendar className="h-3 w-3 text-primary/70" />
                        {stop.period}
                      </span>
                    )}
                  </div>

                  <h3 className="font-display text-base sm:text-lg font-bold text-themeText-primary tracking-tight">
                    {stop.institution || stop.title}
                  </h3>

                  {stop.title && stop.title !== stop.institution && (
                    <p className="text-xs font-body text-themeText-body mt-1 flex items-center gap-1.5">
                      <BookOpen className="h-3 w-3 text-primary" />
                      <span>{stop.title}</span>
                    </p>
                  )}
                </div>

                <div className="mt-4 pt-3 border-t border-border/60 flex items-center gap-1.5 text-xs font-code text-themeText-muted">
                  <MapPin className="h-3.5 w-3.5 text-primary flex-shrink-0" />
                  <span>
                    {stop.city}, {stop.region}
                  </span>
                </div>
              </div>
            ))}
          </div>

          {educationItems.length > 0 && (
            <div className="mt-6 p-4 rounded-md bg-surface-elevated border border-border text-xs text-themeText-muted">
              {educationItems.map((item) => (
                <div key={item.id} className="py-1">
                  <span className="font-medium text-themeText-primary">{item.institution}</span> — {item.degree} ({item.field})
                </div>
              ))}
            </div>
          )}
        </section>
      </div>

      <InterPageNav
        label="FINAL CHAPTER"
        nextRoute="/contact"
        nextTitle="Get in Touch & Ask Arham"
        description="Connect through verified channels or query the grounded AI reviewer directly with citations to code artifacts."
      />
    </div>
  );
}
