import React, { useState } from "react";
import { ChevronDown, ChevronUp, FileCode, ShieldCheck } from "lucide-react";
import type { PublicEvidenceItem } from "./ask-arham.types";

interface EvidenceListProps {
  evidence: PublicEvidenceItem[];
}

export default function EvidenceList({ evidence }: EvidenceListProps) {
  const [expanded, setExpanded] = useState(false);

  if (!evidence || evidence.length === 0) {
    return null;
  }

  return (
    <div className="mt-3 border border-border rounded-lg bg-surface overflow-hidden shadow-sm">
      <button
        type="button"
        onClick={() => setExpanded(!expanded)}
        className="w-full flex items-center justify-between px-3 py-2 text-xs font-mono text-themeText-muted hover:text-themeText-primary hover:bg-surface-elevated transition-colors"
        aria-expanded={expanded}
      >
        <span className="flex items-center gap-1.5">
          <ShieldCheck className="w-3.5 h-3.5 text-primary" />
          <span>Verified Evidence ({evidence.length} sources examined)</span>
        </span>
        {expanded ? (
          <ChevronUp className="w-3.5 h-3.5" />
        ) : (
          <ChevronDown className="w-3.5 h-3.5" />
        )}
      </button>

      {expanded && (
        <div className="p-3 border-t border-border flex flex-col gap-2.5 bg-canvas-soft">
          {evidence.map((item) => (
            <div
              key={item.id}
              className="p-2.5 rounded border border-border bg-surface text-xs flex flex-col gap-1 shadow-sm"
            >
              <div className="flex items-center justify-between font-mono text-themeText-muted">
                <span className="flex items-center gap-1 font-semibold text-primary">
                  <FileCode className="w-3 h-3" />
                  [{item.id}] {item.title}
                </span>
                <span className="text-[10px] uppercase tracking-wider text-themeText-muted/80">
                  {item.kind === "github" && item.repository ? item.repository : "Canonical Portfolio"}
                </span>
              </div>
              <pre className="mt-1 p-2 rounded bg-canvas-soft text-[11px] font-mono text-themeText-body whitespace-pre-wrap break-words border border-border max-h-36 overflow-y-auto">
                <code>{item.excerpt}</code>
              </pre>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
