import React from "react";
import { ArrowRight } from "lucide-react";
import { useRouter } from "next/navigation";
import type { SafeAction } from "./ask-arham.types";

interface SafeActionButtonsProps {
  actions: SafeAction[];
  onActionTriggered?: () => void;
}

// Certified client-side target map matching verified routes (Correction 29, 30)
const SAFE_TARGET_MAP: Record<string, string> = {
  "view-project-edutrace": "/projects/edutrace",
  "view-project-gdgoc-ecommerce": "/projects/gdgoc-ecommerce",
  "view-project-maritime-ai-dashboard": "/projects/maritime-ai-dashboard",
  "view-project-smart-kitchen": "/projects/smart-kitchen",
  "go-to-projects": "/#projects",
  "go-to-skills": "/#skills",
  "go-to-journey": "/#journey",
  "go-to-contact": "/#contact",
};

export default function SafeActionButtons({
  actions,
  onActionTriggered,
}: SafeActionButtonsProps) {
  const router = useRouter();

  if (!actions || actions.length === 0) {
    return null;
  }

  const handleActionClick = (actionId: string) => {
    const target = SAFE_TARGET_MAP[actionId];
    if (!target) {
      console.warn(`Untrusted or unmapped action ID: ${actionId}`);
      return;
    }

    onActionTriggered?.();

    if (target.startsWith("/#")) {
      const sectionId = target.slice(2);
      const element = document.getElementById(sectionId);
      if (element) {
        element.scrollIntoView({ behavior: "smooth" });
      } else {
        router.push(target);
      }
    } else {
      router.push(target);
    }
  };

  return (
    <div className="mt-3 flex flex-col gap-1.5">
      <div className="text-[11px] font-mono text-themeText-muted uppercase tracking-wider">
        Contextual Actions
      </div>
      <div className="flex flex-wrap gap-2">
        {actions.map((act) => {
          // Verify ID exists in certified client map
          if (!SAFE_TARGET_MAP[act.id]) return null;

          return (
            <button
              key={act.id}
              type="button"
              onClick={() => handleActionClick(act.id)}
              className="inline-flex items-center gap-1.5 text-xs px-3 py-1.5 rounded-lg border border-primary/30 bg-primary-muted hover:bg-primary/20 text-primary font-medium transition-colors focus:outline-none focus:ring-1 focus:ring-primary shadow-sm"
            >
              <span>{act.label}</span>
              <ArrowRight className="w-3 h-3" />
            </button>
          );
        })}
      </div>
    </div>
  );
}
