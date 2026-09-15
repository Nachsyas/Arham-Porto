import { MessageSquare, ExternalLink } from "lucide-react";
import { GithubIcon } from "@/components/icons/GithubIcon";
import type { Profile } from "arham-porto-schema";

interface ContactSectionProps {
  profile: Profile;
}

export default function ContactSection({ profile }: ContactSectionProps) {
  return (
    <section id="contact" aria-label="Contact & Collaboration" className="py-24 px-4 sm:px-6 lg:px-8">
      <div className="mx-auto max-w-4xl text-center">
        <span className="text-xs font-code text-primary uppercase tracking-wider block mb-2">
          05 // CONTACT
        </span>
        <h2 className="font-display text-3xl sm:text-5xl font-extrabold text-themeText-primary tracking-tight">
          Let&apos;s Build Together
        </h2>
        <p className="mt-4 text-base sm:text-lg font-body text-themeText-muted max-w-xl mx-auto leading-relaxed">
          Open to software engineering opportunities, technical collaborations, and distributed systems challenges.
        </p>

        {/* Contact Action Cards */}
        <div className="mt-10 flex flex-col sm:flex-row items-center justify-center gap-4">
          {profile.github && (
            <a
              href={profile.github}
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2.5 px-6 py-3.5 rounded-md text-sm font-semibold font-body bg-primary text-white hover:bg-primary-active transition-all shadow-sm hover:shadow-md focus-visible:ring-2 focus-visible:ring-primary w-full sm:w-auto justify-center"
            >
              <GithubIcon className="h-4 w-4" />
              Connect on GitHub
              <ExternalLink className="h-3.5 w-3.5 opacity-80" />
            </a>
          )}
        </div>

        {/* Technical Privacy Note */}
        <p className="mt-12 text-xs font-code text-themeText-mutedSoft">
          Privacy Policy: Direct contact routes are verified against primary channels. Zero unsolicited marketing data collected.
        </p>
      </div>
    </section>
  );
}
