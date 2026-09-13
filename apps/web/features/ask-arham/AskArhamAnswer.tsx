import React from "react";
import { ShieldCheck, AlertCircle, Lock, Slash } from "lucide-react";
import type { GroundedResponse } from "./ask-arham.types";
import EvidenceList from "./EvidenceList";
import SourceList from "./SourceList";
import SafeActionButtons from "./SafeActionButtons";

interface AskArhamAnswerProps {
  response: GroundedResponse;
  onActionTriggered?: () => void;
}

export default function AskArhamAnswer({
  response,
  onActionTriggered,
}: AskArhamAnswerProps) {
  const { status, segments, evidence, sources, actions } = response;

  const renderStatusBadge = () => {
    switch (status) {
      case "supported":
        return (
          <span className="inline-flex items-center gap-1 text-[11px] font-mono px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/30">
            <ShieldCheck className="w-3 h-3" />
            <span>Verified Evidence</span>
          </span>
        );
      case "insufficient_evidence":
        return (
          <span className="inline-flex items-center gap-1 text-[11px] font-mono px-2 py-0.5 rounded-full bg-amber-500/10 text-amber-400 border border-amber-500/30">
            <AlertCircle className="w-3 h-3" />
            <span>Insufficient Evidence</span>
          </span>
        );
      case "privacy_refusal":
        return (
          <span className="inline-flex items-center gap-1 text-[11px] font-mono px-2 py-0.5 rounded-full bg-blue-500/10 text-blue-400 border border-blue-500/30">
            <Lock className="w-3 h-3" />
            <span>Privacy Guard</span>
          </span>
        );
      case "scope_refusal":
        return (
          <span className="inline-flex items-center gap-1 text-[11px] font-mono px-2 py-0.5 rounded-full bg-neutral-500/10 text-neutral-400 border border-neutral-500/30">
            <Slash className="w-3 h-3" />
            <span>Scope Bounded</span>
          </span>
        );
    }
  };

  return (
    <div className="flex flex-col gap-2 text-sm text-themeText-primary">
      {/* Status Badge */}
      <div className="flex items-center justify-between gap-2">
        {renderStatusBadge()}
      </div>

      {/* Answer Paragraphs with Claim-Level Evidence Indicators */}
      <div className="leading-relaxed text-themeText-body whitespace-pre-wrap break-words">
        {segments && segments.length > 0 ? (
          segments.map((seg, idx) => (
            <span key={idx} className="mr-1">
              {seg.text}
              {seg.evidence_ids && seg.evidence_ids.length > 0 && (
                <span className="inline-flex gap-1 ml-1 align-baseline">
                  {seg.evidence_ids.map((id) => (
                    <span
                      key={id}
                      className="text-[10px] font-mono px-1 py-0.2 rounded bg-primary/10 text-primary border border-primary/30"
                    >
                      [{id}]
                    </span>
                  ))}
                </span>
              )}
            </span>
          ))
        ) : (
          <span>{response.answer}</span>
        )}
      </div>

      {/* Evidence Accordion (No similarity scores) */}
      <EvidenceList evidence={evidence} />

      {/* Immutable Source Links */}
      <SourceList sources={sources} />

      {/* Verified Contextual Actions */}
      <SafeActionButtons
        actions={actions}
        onActionTriggered={onActionTriggered}
      />
    </div>
  );
}
