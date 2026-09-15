"use client";

import React, { useState, useRef } from "react";
import Image from "next/image";
import { motion, useMotionValue, useTransform, useReducedMotion } from "motion/react";
import { UserCheck } from "lucide-react";

interface HeroPortraitProps {
  fullName?: string;
  role?: string;
}

export default function HeroPortrait({
  fullName = "Nachsyas Arham Mumtaz Nashohi",
  role = "Software Engineer",
}: HeroPortraitProps) {
  const [hasError, setHasError] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);
  const shouldReduceMotion = useReducedMotion();

  // Subtle 3D tilt motion values
  const mouseX = useMotionValue(0.5);
  const mouseY = useMotionValue(0.5);

  // Restrained rotation: rotateX <= 1.5deg, rotateY <= 2deg
  const rotateX = useTransform(mouseY, [0, 1], [1.5, -1.5]);
  const rotateY = useTransform(mouseX, [0, 1], [-2, 2]);

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (shouldReduceMotion || !containerRef.current) return;
    const rect = containerRef.current.getBoundingClientRect();
    const x = (e.clientX - rect.left) / rect.width;
    const y = (e.clientY - rect.top) / rect.height;
    mouseX.set(x);
    mouseY.set(y);
  };

  const handleMouseLeave = () => {
    mouseX.set(0.5);
    mouseY.set(0.5);
  };

  return (
    <motion.div
      ref={containerRef}
      onMouseMove={handleMouseMove}
      onMouseLeave={handleMouseLeave}
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.6, ease: "easeOut" }}
      style={{
        perspective: 1000,
        rotateX: shouldReduceMotion ? 0 : rotateX,
        rotateY: shouldReduceMotion ? 0 : rotateY,
      }}
      className="relative w-full max-w-[340px] sm:max-w-[380px] lg:max-w-[420px] mx-auto group transition-transform duration-200"
    >
      {/* Outer Technical Frame Container - Exact canonical dimensions preserved */}
      <div className="relative aspect-[3/4] rounded-2xl border border-border bg-surface p-2.5 sm:p-3 shadow-md overflow-hidden hover:border-primary/40 transition-colors duration-500 flex flex-col justify-between">
        {/* Subtle Minimal Corner Markers */}
        <span className="absolute top-2 left-3 text-[9px] font-code text-primary/50 pointer-events-none z-20">
          ┌ PROFILE
        </span>
        <span className="absolute top-2 right-3 text-[9px] font-code text-primary/50 pointer-events-none z-20">
          ┐
        </span>

        {/* Inner Stage Container */}
        <div className="relative w-full h-full rounded-xl bg-canvas-soft border border-border/50 flex flex-col items-center justify-center p-3 sm:p-4 overflow-hidden">
          {/* Photograph Stage: sized to ~86% of available stage */}
          <div className="relative w-[88%] h-[86%] rounded-lg overflow-hidden border border-border/70 shadow-sm bg-surface-strong">
            {!hasError ? (
              <Image
                src="/images/profile/nachsyas-arham.jpg"
                alt={`Portrait of ${fullName}`}
                fill
                priority
                sizes="(max-width: 640px) 90vw, (max-width: 1024px) 45vw, 420px"
                className="object-cover object-top transition-transform duration-500 ease-out group-hover:scale-[1.012]"
                onError={() => setHasError(true)}
              />
            ) : (
              /* Graceful Fallback if image fails to load */
              <div
                className="w-full h-full flex flex-col items-center justify-center p-6 text-center space-y-4 bg-surface"
                role="img"
                aria-label={`Visual profile placeholder for ${fullName}`}
              >
                <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-canvas-soft border border-primary/40 text-primary">
                  <UserCheck className="h-8 w-8" />
                </div>
                <div>
                  <p className="font-display font-bold text-sm text-themeText-primary">
                    {fullName}
                  </p>
                  <p className="font-code text-xs text-primary">{role}</p>
                </div>
              </div>
            )}
          </div>

          {/* Minimalist Profile Identity Overlay Strip */}
          <div className="w-[88%] mt-2 px-2.5 py-1.5 bg-surface/95 border border-border/70 rounded-md backdrop-blur-sm flex items-center justify-between z-10 shadow-sm">
            <span className="text-[10px] font-code text-themeText-primary font-semibold tracking-wide">
              NACHSYAS ARHAM
            </span>
            <span className="text-[9px] font-code text-primary font-semibold tracking-wider uppercase">
              {role}
            </span>
          </div>
        </div>

        {/* Bottom Corner Markers */}
        <span className="absolute bottom-1 left-3 text-[9px] font-code text-themeText-mutedSoft pointer-events-none z-20">
          └
        </span>
        <span className="absolute bottom-1 right-3 text-[9px] font-code text-themeText-mutedSoft pointer-events-none z-20">
          ┘
        </span>
      </div>
    </motion.div>
  );
}
