import React from "react";
import { Sparkles } from "lucide-react";

interface SuggestedQuestionsProps {
  onSelect: (question: string) => void;
  disabled?: boolean;
}

const SUGGESTED_QUESTIONS = [
  "What does Arham specialize in?",
  "Which projects demonstrate backend engineering?",
  "What evidence shows Go experience?",
  "How is EduTrace architected?",
  "What technologies are used in Smart Kitchen?",
  "Does Maritime AI demonstrate RAG experience?",
  "Summarize Arham's engineering background.",
];

export default function SuggestedQuestions({
  onSelect,
  disabled = false,
}: SuggestedQuestionsProps) {
  return (
    <div className="flex flex-col gap-2.5 my-3">
      <div className="flex items-center gap-1.5 text-xs font-mono uppercase tracking-wider text-themeText-muted">
        <Sparkles className="w-3.5 h-3.5 text-primary" />
        <span>Suggested Reviewer Questions</span>
      </div>
      <div className="flex flex-wrap gap-2">
        {SUGGESTED_QUESTIONS.map((q, idx) => (
          <button
            key={idx}
            type="button"
            disabled={disabled}
            onClick={() => onSelect(q)}
            className="text-left text-xs bg-surface-elevated hover:bg-surface border border-border hover:border-primary text-themeText-body hover:text-primary px-3 py-1.5 rounded-full transition-all duration-150 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-1 focus:ring-primary shadow-sm"
          >
            {q}
          </button>
        ))}
      </div>
    </div>
  );
}
