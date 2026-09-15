"use client";

import Link from "next/link";
import { ArrowRight } from "lucide-react";

interface InterPageNavProps {
  label: string;
  nextRoute: string;
  nextTitle: string;
  description?: string;
}

export default function InterPageNav({
  label = "NEXT",
  nextRoute,
  nextTitle,
  description,
}: InterPageNavProps) {
  return (
    <section className="mt-16 sm:mt-24 pt-10 border-t border-border/80">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <Link
          href={nextRoute}
          aria-label={`Navigate to next section: ${nextTitle}`}
          data-testid="inter-page-nav"
          className="group block p-6 sm:p-8 rounded-card border border-border bg-surface hover:border-primary/50 hover:shadow-md transition-all duration-300"
        >
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div className="space-y-1">
              <span className="text-xs font-code font-bold text-primary tracking-wider uppercase">
                {label} // CONTINUATION
              </span>
              <h3 className="font-display text-xl sm:text-2xl font-bold text-themeText-primary group-hover:text-primary transition-colors flex items-center gap-2">
                <span>{nextTitle}</span>
              </h3>
              {description && (
                <p className="font-body text-xs sm:text-sm text-themeText-muted max-w-xl">
                  {description}
                </p>
              )}
            </div>

            <div className="flex items-center gap-2 font-code text-xs font-semibold text-primary">
              <span>Explore Section</span>
              <span className="flex h-8 w-8 items-center justify-center rounded-full bg-primary-muted border border-primary/25 group-hover:bg-primary group-hover:text-white transition-all transform group-hover:translate-x-1.5 shadow-sm">
                <ArrowRight className="h-4 w-4" />
              </span>
            </div>
          </div>
        </Link>
      </div>
    </section>
  );
}
