"use client";

import React, { useEffect, useState } from "react";
import { usePathname } from "next/navigation";
import { useOptionalNavigationProgress } from "@/context/NavigationProgressContext";

export default function RouteProgress() {
  const pathname = usePathname();
  const navProgress = useOptionalNavigationProgress();
  const [status, setStatus] = useState<"idle" | "navigating" | "completing">("idle");
  const [prefersReduced, setPrefersReduced] = useState(false);

  // Respect prefers-reduced-motion
  useEffect(() => {
    if (typeof window === "undefined") return;
    const mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
    setPrefersReduced(mediaQuery.matches);
    const listener = (e: MediaQueryListEvent) => setPrefersReduced(e.matches);
    mediaQuery.addEventListener("change", listener);
    return () => mediaQuery.removeEventListener("change", listener);
  }, []);

  // React to startProgress from context
  useEffect(() => {
    if (navProgress?.isNavigating) {
      setStatus("navigating");
    }
  }, [navProgress?.isNavigating]);

  // When pathname changes, complete and fade
  useEffect(() => {
    setStatus("completing");
    const finishTimer = setTimeout(() => {
      setStatus("idle");
    }, 350);

    return () => clearTimeout(finishTimer);
  }, [pathname]);

  if (prefersReduced || status === "idle") {
    return null;
  }

  return (
    <div
      aria-hidden="true"
      data-testid="route-progress-bar"
      className="fixed top-0 left-0 right-0 z-50 h-[2.5px] bg-primary pointer-events-none origin-left transition-all duration-300 motion-reduce:hidden"
      style={{
        width: status === "navigating" ? "75%" : "100%",
        opacity: status === "completing" ? 0 : 1,
        transition:
          status === "completing"
            ? "width 200ms ease-out, opacity 250ms ease-in"
            : "width 400ms cubic-bezier(0.1, 0.9, 0.2, 1)",
      }}
    />
  );
}
