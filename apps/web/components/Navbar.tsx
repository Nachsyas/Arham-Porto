"use client";

import { useState } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Menu, X, Sparkles } from "lucide-react";
import { motion } from "motion/react";
import { useOptionalPortfolioUI } from "@/context/PortfolioUIContext";

interface NavbarProps {
  onOpenQuickReview?: () => void;
  hasExperience?: boolean;
}

const NAV_ITEMS = [
  { label: "Home", href: "/" },
  { label: "Work", href: "/work" },
  { label: "Skills", href: "/skills" },
  { label: "Journey", href: "/journey" },
  { label: "Contact", href: "/contact" },
];

export default function Navbar({ onOpenQuickReview }: NavbarProps) {
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const pathname = usePathname();
  const uiContext = useOptionalPortfolioUI();

  const handleQuickReview = () => {
    if (onOpenQuickReview) {
      onOpenQuickReview();
    } else if (uiContext) {
      uiContext.openQuickReview();
    }
  };

  const isActive = (href: string) => {
    if (!pathname) {
      return href === "/";
    }
    if (href === "/") {
      return pathname === "/";
    }
    if (href === "/work") {
      return pathname === "/work" || pathname.startsWith("/projects/");
    }
    return pathname.startsWith(href);
  };

  return (
    <header className="sticky top-0 z-40 w-full border-b border-border bg-surface/90 backdrop-blur-md">
      <div className="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
        {/* Brand / Identity */}
        <Link
          href="/"
          className="flex items-center gap-2.5 focus-visible:ring-2 focus-visible:ring-primary rounded p-1"
        >
          <span className="flex h-8 w-8 items-center justify-center rounded-md bg-surface-strong border border-primary/30 text-primary font-code font-bold text-sm shadow-sm">
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

        {/* Desktop Route Navigation */}
        <nav aria-label="Main Navigation" className="hidden md:flex items-center gap-1.5 bg-canvas-soft/80 p-1.5 rounded-full border border-border/80 shadow-sm">
          {NAV_ITEMS.map((item) => {
            const active = isActive(item.href);
            return (
              <Link
                key={item.label}
                href={item.href}
                className={`relative px-3.5 py-1.5 text-xs font-code font-medium rounded-full transition-colors focus-visible:ring-2 focus-visible:ring-primary ${
                  active
                    ? "text-primary font-semibold"
                    : "text-themeText-body hover:text-themeText-primary"
                }`}
              >
                {active && (
                  <motion.span
                    layoutId="navbar-active-pill"
                    className="absolute inset-0 rounded-full bg-surface border border-primary/30 shadow-sm"
                    transition={{ type: "spring", stiffness: 380, damping: 30 }}
                  />
                )}
                <span className="relative z-10">{item.label}</span>
              </Link>
            );
          })}
        </nav>

        {/* Desktop Actions */}
        <div className="hidden md:flex items-center gap-3">
          <button
            type="button"
            onClick={handleQuickReview}
            className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-md text-xs font-medium font-code bg-surface-elevated border border-border text-primary hover:bg-primary-muted hover:border-primary transition-all shadow-sm focus-visible:ring-2 focus-visible:ring-primary"
          >
            <Sparkles className="h-3.5 w-3.5 text-primary" />
            Quick Review
          </button>

          <Link
            href="/contact"
            className="inline-flex items-center gap-1.5 px-3.5 py-1.5 rounded-md text-xs font-medium font-body bg-primary text-white hover:bg-primary-active transition-colors shadow-sm focus-visible:ring-2 focus-visible:ring-primary"
          >
            Get in Touch
          </Link>
        </div>

        {/* Mobile Menu Toggle */}
        <div className="flex md:hidden items-center gap-2">
          <button
            type="button"
            onClick={handleQuickReview}
            className="inline-flex items-center gap-1 px-2.5 py-1 rounded text-xs font-code bg-surface-elevated border border-border text-primary shadow-sm"
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
        <div className="md:hidden border-b border-border bg-surface px-4 py-4 space-y-3 shadow-md">
          <nav aria-label="Mobile Navigation" className="flex flex-col space-y-1.5">
            {NAV_ITEMS.map((item) => {
              const active = isActive(item.href);
              return (
                <Link
                  key={item.label}
                  href={item.href}
                  onClick={() => setMobileMenuOpen(false)}
                  className={`rounded-md px-3 py-2 text-sm font-code transition-colors flex items-center justify-between ${
                    active
                      ? "bg-primary-muted text-primary font-semibold border border-primary/20"
                      : "text-themeText-body hover:bg-surface-elevated hover:text-primary"
                  }`}
                >
                  <span>{item.label}</span>
                  {active && <span className="h-1.5 w-1.5 rounded-full bg-primary" />}
                </Link>
              );
            })}
          </nav>
          <div className="pt-3 border-t border-border flex flex-col gap-2">
            <button
              type="button"
              onClick={() => {
                setMobileMenuOpen(false);
                handleQuickReview();
              }}
              className="flex w-full items-center justify-center gap-2 py-2 px-3 rounded-md text-xs font-code bg-surface-elevated border border-primary/30 text-primary shadow-sm"
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
