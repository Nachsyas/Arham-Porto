"use client";

import { useState } from "react";
import Link from "next/link";
import { Menu, X, FileText, Sparkles } from "lucide-react";

interface NavbarProps {
  onOpenQuickReview: () => void;
}

export default function Navbar({ onOpenQuickReview }: NavbarProps) {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const navLinks = [
    { label: "Overview", href: "#overview" },
    { label: "Projects", href: "#projects" },
    { label: "Skills", href: "#skills" },
    { label: "Experience", href: "#experience" },
    { label: "Contact", href: "#contact" },
  ];

  return (
    <header className="sticky top-0 z-40 w-full border-b border-border bg-canvas/85 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        {/* Brand / Identity */}
        <Link
          href="#overview"
          className="flex items-center gap-2.5 focus-visible:ring-2 focus-visible:ring-primary rounded p-1"
        >
          <span className="flex h-8 w-8 items-center justify-center rounded-md bg-surfaceStrong border border-primary/40 text-primary font-code font-bold text-sm shadow-[0_0_12px_rgba(38,184,255,0.2)]">
            AP
          </span>
          <div className="flex flex-col">
            <span className="font-display text-sm font-bold tracking-tight text-themeText-primary">
              ARHAM PORTO
            </span>
            <span className="text-[10px] font-code text-themeText-muted -mt-0.5">
              Nachsyas Arham
            </span>
          </div>
        </Link>

        {/* Desktop Navigation */}
        <nav aria-label="Main Navigation" className="hidden md:flex items-center gap-6">
          {navLinks.map((link) => (
            <Link
              key={link.label}
              href={link.href}
              className="text-sm font-body text-themeText-body hover:text-primary transition-colors focus-visible:ring-2 focus-visible:ring-primary rounded px-2 py-1"
            >
              {link.label}
            </Link>
          ))}
        </nav>

        {/* Desktop Actions */}
        <div className="hidden md:flex items-center gap-3">
          <button
            type="button"
            onClick={onOpenQuickReview}
            className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-md text-xs font-medium font-code bg-surfaceElevated border border-primary/30 text-primary hover:bg-primary/10 hover:border-primary transition-all shadow-[0_0_10px_rgba(38,184,255,0.15)] focus-visible:ring-2 focus-visible:ring-primary"
          >
            <Sparkles className="h-3.5 w-3.5 text-primary" />
            Quick Review
          </button>

          <a
            href="#contact"
            className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-md text-xs font-medium font-body bg-primary text-canvas hover:bg-primaryActive transition-colors focus-visible:ring-2 focus-visible:ring-primary"
          >
            Get in Touch
          </a>
        </div>

        {/* Mobile Menu Toggle */}
        <div className="flex md:hidden items-center gap-2">
          <button
            type="button"
            onClick={onOpenQuickReview}
            className="inline-flex items-center gap-1 px-2.5 py-1 rounded text-xs font-code bg-surface border border-primary/30 text-primary"
          >
            <Sparkles className="h-3 w-3" /> Review
          </button>
          <button
            type="button"
            onClick={() => setMobileMenuOpen(!mobileMenuOpen)}
            aria-expanded={mobileMenuOpen}
            aria-label="Toggle navigation menu"
            className="flex h-10 w-10 items-center justify-center rounded-md border border-border text-themeText-muted hover:text-themeText-primary focus-visible:ring-2 focus-visible:ring-primary"
          >
            {mobileMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </button>
        </div>
      </div>

      {/* Mobile Dropdown Menu */}
      {mobileMenuOpen && (
        <div className="md:hidden border-b border-border bg-surface px-4 py-4 space-y-3">
          <nav aria-label="Mobile Navigation" className="flex flex-col space-y-2">
            {navLinks.map((link) => (
              <Link
                key={link.label}
                href={link.href}
                onClick={() => setMobileMenuOpen(false)}
                className="rounded-md px-3 py-2 text-sm font-body text-themeText-body hover:bg-surfaceElevated hover:text-primary transition-colors"
              >
                {link.label}
              </Link>
            ))}
          </nav>
          <div className="pt-3 border-t border-border flex flex-col gap-2">
            <button
              type="button"
              onClick={() => {
                setMobileMenuOpen(false);
                onOpenQuickReview();
              }}
              className="flex w-full items-center justify-center gap-2 py-2 px-3 rounded-md text-xs font-code bg-surfaceElevated border border-primary/40 text-primary"
            >
              <Sparkles className="h-3.5 w-3.5" />
              Open 60-Second Quick Review
            </button>
          </div>
        </div>
      )}
    </header>
  );
}
