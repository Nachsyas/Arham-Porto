import Link from "next/link";
import { Shield } from "lucide-react";
import { GithubIcon } from "@/components/icons/GithubIcon";
import type { Profile } from "arham-porto-schema";

interface FooterProps {
  profile: Profile;
}

export default function Footer({ profile }: FooterProps) {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="border-t border-border bg-canvasSoft py-12 px-4 sm:px-6 lg:px-8 text-xs font-body text-themeText-muted">
      <div className="mx-auto max-w-7xl flex flex-col md:flex-row items-center justify-between gap-6">
        {/* Left: Identity & Copyright */}
        <div className="flex flex-col sm:flex-row items-center gap-3 text-center sm:text-left">
          <span className="font-display font-bold text-themeText-primary text-sm">
            {profile.fullName}
          </span>
          <span className="hidden sm:inline text-borderStrong">|</span>
          <span>© {currentYear} {profile.projectName}. All rights reserved.</span>
        </div>

        {/* Center: Subtle Professional Note */}
        <div className="text-center font-code text-[11px] text-themeText-mutedSoft">
          Evidence-backed engineering portfolio
        </div>

        {/* Right: Links */}
        <div className="flex items-center gap-4">
          {profile.github && (
            <a
              href={profile.github}
              target="_blank"
              rel="noopener noreferrer"
              aria-label="GitHub Profile"
              className="text-themeText-muted hover:text-themeText-primary transition-colors focus-visible:ring-2 focus-visible:ring-primary rounded p-1"
            >
              <GithubIcon className="h-4 w-4" />
            </a>
          )}
          <a
            href="#overview"
            className="hover:text-themeText-primary transition-colors focus-visible:ring-2 focus-visible:ring-primary rounded px-2 py-1"
          >
            Back to top ↑
          </a>
        </div>
      </div>
    </footer>
  );
}
