"use client";

import { motion } from "motion/react";
import React from "react";

interface PageHeaderProps {
  eyebrow: string;
  title: string;
  description: string;
  badge?: string;
  children?: React.ReactNode;
}

export default function PageHeader({
  eyebrow,
  title,
  description,
  badge,
  children,
}: PageHeaderProps) {
  return (
    <header className="pt-10 pb-8 sm:pt-14 sm:pb-10 border-b border-border/70 mb-8 sm:mb-12">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 space-y-3 sm:space-y-4">
        {/* Eyebrow / Technical Tag */}
        <motion.div
          initial={{ opacity: 0, y: 8 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.35, ease: "easeOut" }}
          className="flex items-center gap-2"
        >
          <span className="h-1.5 w-1.5 rounded-full bg-primary" />
          <span className="text-xs font-code font-bold tracking-wider text-primary uppercase">
            {eyebrow}
          </span>
          {badge && (
            <span className="text-[10px] font-code px-2 py-0.5 rounded-full bg-primary-muted text-primary border border-primary/25">
              {badge}
            </span>
          )}
        </motion.div>

        {/* H1 Heading */}
        <motion.h1
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.08, ease: "easeOut" }}
          className="font-display text-3xl sm:text-5xl font-extrabold tracking-tight text-themeText-primary"
        >
          {title}
        </motion.h1>

        {/* Description Paragraph */}
        <motion.p
          initial={{ opacity: 0, y: 12 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4, delay: 0.16, ease: "easeOut" }}
          className="font-body text-base sm:text-lg text-themeText-body max-w-3xl leading-relaxed"
        >
          {description}
        </motion.p>

        {children && (
          <motion.div
            initial={{ opacity: 0, y: 8 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.4, delay: 0.24, ease: "easeOut" }}
            className="pt-2"
          >
            {children}
          </motion.div>
        )}
      </div>
    </header>
  );
}
