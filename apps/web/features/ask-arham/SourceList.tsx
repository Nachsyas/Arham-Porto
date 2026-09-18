import React from "react";
import { ExternalLink, Bookmark } from "lucide-react";
import NavLink from "@/components/motion/NavLink";
import { GithubIcon } from "@/components/icons/GithubIcon";
import type { SourceCitation } from "./ask-arham.types";

interface SourceListProps {
  sources: SourceCitation[];
}

export default function SourceList({ sources }: SourceListProps) {
  if (!sources || sources.length === 0) {
    return null;
  }

  return (
    <div className="mt-3 flex flex-col gap-1.5">
      <div className="text-[11px] font-mono text-themeText-muted uppercase tracking-wider">
        Certified Sources
      </div>
      <div className="flex flex-wrap gap-2">
        {sources.map((src) => {
          const isExternal = src.kind === "github" && src.url;

          if (isExternal && src.url) {
            return (
              <a
                key={src.id}
                href={src.url}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-1 text-xs px-2.5 py-1 rounded border border-border bg-surface hover:border-primary text-themeText-body hover:text-primary transition-colors font-mono shadow-sm"
              >
                <GithubIcon className="w-3 h-3 text-themeText-muted" />
                <span>{src.label}</span>
                <ExternalLink className="w-3 h-3 text-themeText-muted ml-0.5" />
              </a>
            );
          }

          if (src.url) {
            return (
              <NavLink
                key={src.id}
                href={src.url}
                className="inline-flex items-center gap-1 text-xs px-2.5 py-1 rounded border border-border bg-surface hover:border-primary text-themeText-body hover:text-primary transition-colors font-mono shadow-sm"
              >
                <Bookmark className="w-3 h-3 text-primary" />
                <span>{src.label}</span>
              </NavLink>
            );
          }

          return (
            <span
              key={src.id}
              className="inline-flex items-center gap-1 text-xs px-2.5 py-1 rounded border border-border bg-surface-elevated text-themeText-muted font-mono"
            >
              <Bookmark className="w-3 h-3 text-themeText-muted" />
              <span>{src.label}</span>
            </span>
          );
        })}
      </div>
    </div>
  );
}
