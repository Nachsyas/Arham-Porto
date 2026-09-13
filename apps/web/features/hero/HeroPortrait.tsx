"use client";

import React, { useState } from "react";
import Image from "next/image";
import { motion } from "motion/react";
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

  return (
    <motion.div
      initial={{ opacity: 0, y: 12 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.6, ease: "easeOut" }}
      className="relative w-full max-w-[340px] sm:max-w-[380px] lg:max-w-[420px] mx-auto group"
    >
      {/* Outer Technical Frame Container */}
      <div className="relative aspect-[3/4] rounded-2xl border border-border bg-gradient-to-b from-surfaceElevated/70 to-surface/90 p-2.5 sm:p-3 shadow-2xl overflow-hidden hover:border-primary/40 transition-colors duration-500">
        {/* Subtle Minimal Corner Markers */}
        <span className="absolute top-2 left-3 text-[9px] font-code text-primary/40 pointer-events-none z-20">
          ┌ PROFILE
        </span>
        <span className="absolute top-2 right-3 text-[9px] font-code text-primary/40 pointer-events-none z-20">
          ┐
        </span>

        {/* Inner Photographic Frame */}
        <div className="relative w-full h-full rounded-xl overflow-hidden bg-surfaceStrong border border-border/50">
          {!hasError ? (
            <>
              <Image
                src="/images/profile/nachsyas-arham.jpg"
                alt={`Portrait of ${fullName}`}
                fill
                priority
                sizes="(max-width: 640px) 90vw, (max-width: 1024px) 45vw, 420px"
                className="object-cover object-top transition-transform duration-700 ease-out group-hover:scale-[1.02]"
                onError={() => setHasError(true)}
              />
              {/* Very Subtle Grounding Gradient to integrate seamlessly with canvas dark theme */}
              <div className="absolute inset-0 pointer-events-none bg-gradient-to-t from-canvas/80 via-transparent to-transparent opacity-60" />
            </>
          ) : (
            /* Graceful Fallback if image fails to load */
            <div
              className="w-full h-full flex flex-col items-center justify-center p-6 text-center space-y-4 bg-surface"
              role="img"
              aria-label={`Visual profile placeholder for ${fullName}`}
            >
              <div className="flex h-16 w-16 items-center justify-center rounded-2xl bg-canvas border border-primary/40 text-primary">
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

          {/* Minimalist Profile Identity Overlay Strip */}
          <div className="absolute bottom-0 inset-x-0 p-3 bg-gradient-to-t from-canvas/95 via-canvas/80 to-transparent border-t border-border/40 backdrop-blur-xs flex items-center justify-between z-10">
            <span className="text-[11px] font-code text-themeText-primary font-medium tracking-wide">
              NACHSYAS ARHAM
            </span>
            <span className="text-[10px] font-code text-primary font-medium tracking-wider uppercase">
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
