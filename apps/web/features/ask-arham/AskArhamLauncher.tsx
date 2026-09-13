"use client";

import React, { useEffect, forwardRef } from "react";
import { Sparkles, Bot } from "lucide-react";

interface AskArhamLauncherProps {
  onClick: () => void;
  isOpen: boolean;
}

export const AskArhamLauncher = forwardRef<
  HTMLButtonElement,
  AskArhamLauncherProps
>(function AskArhamLauncher({ onClick, isOpen }, ref) {
  // Optional enhancement: Alt+A shortcut (Correction 44)
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.altKey && (e.key === "a" || e.key === "A")) {
        e.preventDefault();
        onClick();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [onClick]);

  return (
    <div className="fixed bottom-6 right-6 z-40 flex items-center">
      <button
        ref={ref}
        type="button"
        onClick={onClick}
        aria-haspopup="dialog"
        aria-expanded={isOpen}
        aria-label="Open Ask Arham AI Reviewer Copilot"
        className="group relative flex items-center gap-2 px-4 py-3 rounded-full bg-surface border border-primary/40 shadow-lg hover:border-primary hover:shadow-primary/20 text-themeText-primary transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 focus:ring-offset-canvas"
      >
        <span className="relative flex h-2.5 w-2.5">
          <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
          <span className="relative inline-flex rounded-full h-2.5 w-2.5 bg-primary"></span>
        </span>

        <Bot className="w-4 h-4 text-primary group-hover:rotate-6 transition-transform" />

        <span className="text-xs font-mono font-semibold tracking-tight">
          Ask Arham AI
        </span>

        <span className="hidden sm:inline-flex text-[10px] font-mono text-themeText-muted px-1.5 py-0.5 rounded border border-border/80 bg-canvas/60">
          Alt+A
        </span>
      </button>
    </div>
  );
});

export default AskArhamLauncher;
