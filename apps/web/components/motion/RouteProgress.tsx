"use client";

import React, { useEffect, useState, useRef } from "react";
import { usePathname } from "next/navigation";
import { useOptionalNavigationProgress } from "@/context/NavigationProgressContext";

export default function RouteProgress() {
  const pathname = usePathname();
  const navProgress = useOptionalNavigationProgress();
  const [status, setStatus] = useState<"idle" | "navigating" | "completing">("idle");
  const [scale, setScale] = useState(0);
  const [opacity, setOpacity] = useState(0);
  const [prefersReduced, setPrefersReduced] = useState(false);
  const isFirstMount = useRef(true);

  // Respect prefers-reduced-motion
  useEffect(() => {
    if (typeof window === "undefined" || typeof window.matchMedia !== "function") return;
    const mediaQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
    setPrefersReduced(mediaQuery.matches);
    const listener = (e: MediaQueryListEvent) => setPrefersReduced(e.matches);
    mediaQuery.addEventListener("change", listener);
    return () => mediaQuery.removeEventListener("change", listener);
  }, []);

  // React to startProgress from context (navigation start)
  useEffect(() => {
    if (navProgress?.isNavigating) {
      setStatus("navigating");
      setOpacity(1);
      setScale(0.05);

      const timer = setTimeout(() => {
        setScale(0.75);
      }, 20);

      return () => clearTimeout(timer);
    }
  }, [navProgress?.isNavigating]);

  // When pathname changes (route commit)
  useEffect(() => {
    if (isFirstMount.current) {
      isFirstMount.current = false;
      return;
    }

    setStatus("completing");
    setScale(1);
    setOpacity(1);

    const fadeTimer = setTimeout(() => {
      setOpacity(0);
    }, 180);

    const idleTimer = setTimeout(() => {
      setStatus("idle");
      setScale(0);
    }, 400);

    return () => {
      clearTimeout(fadeTimer);
      clearTimeout(idleTimer);
    };
  }, [pathname]);

  if (prefersReduced) {
    return null;
  }

  const getTransition = () => {
    if (status === "completing") {
      return "transform 180ms ease-out, opacity 200ms ease-in";
    }
    if (status === "navigating" && scale > 0.05) {
      return "transform 600ms cubic-bezier(0.1, 0.9, 0.2, 1)";
    }
    return "none";
  };

  return (
    <div
      aria-hidden="true"
      data-testid="route-progress-bar"
      className="fixed top-0 left-0 right-0 z-50 h-[2px] bg-primary pointer-events-none origin-left motion-reduce:hidden"
      style={{
        transform: `scaleX(${scale})`,
        opacity,
        transition: getTransition(),
      }}
    />
  );
}
