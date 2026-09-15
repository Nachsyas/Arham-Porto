import React, { useState, useRef, useEffect } from "react";
import { Send, Square, AlertCircle } from "lucide-react";

interface AskArhamInputProps {
  onSubmit: (question: string) => void;
  onAbort?: () => void;
  isLoading: boolean;
  disabled?: boolean;
}

export default function AskArhamInput({
  onSubmit,
  onAbort,
  isLoading,
  disabled = false,
}: AskArhamInputProps) {
  const [question, setQuestion] = useState("");
  const inputRef = useRef<HTMLTextAreaElement>(null);

  const charCount = question.length;
  const isValidLength = charCount >= 2 && charCount <= 1000;
  const isOverLimit = charCount > 1000;

  useEffect(() => {
    if (!isLoading && !disabled) {
      inputRef.current?.focus();
    }
  }, [isLoading, disabled]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!isValidLength || isLoading || disabled) return;
    onSubmit(question);
    setQuestion("");
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      handleSubmit(e);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-1.5 p-3 border-t border-border bg-surface-elevated">
      <div className="relative flex items-end gap-2 bg-canvas-soft border border-border focus-within:border-primary rounded-xl p-2.5 transition-colors">
        <textarea
          ref={inputRef}
          value={question}
          onChange={(e) => setQuestion(e.target.value)}
          onKeyDown={handleKeyDown}
          disabled={isLoading || disabled}
          placeholder="Ask about projects, engineering evidence, or technical background..."
          rows={2}
          className="w-full resize-none bg-transparent text-sm text-themeText-primary placeholder:text-themeText-muted/70 focus:outline-none disabled:opacity-50 disabled:cursor-not-allowed"
          aria-label="Ask a question about Nachsyas's portfolio"
        />

        {isLoading ? (
          <button
            type="button"
            onClick={onAbort}
            className="flex-shrink-0 p-2 rounded-lg bg-red-50 hover:bg-red-100 text-red-700 border border-red-200 transition-colors focus:outline-none focus:ring-1 focus:ring-red-400"
            aria-label="Stop generation"
            title="Stop response"
          >
            <Square className="w-4 h-4 fill-current" />
          </button>
        ) : (
          <button
            type="submit"
            disabled={!isValidLength || disabled}
            className="flex-shrink-0 p-2 rounded-lg bg-primary hover:bg-primary-active text-white transition-colors disabled:opacity-40 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-primary shadow-sm"
            aria-label="Submit question"
          >
            <Send className="w-4 h-4" />
          </button>
        )}
      </div>

      {/* Character counter & length warning */}
      <div className="flex items-center justify-between px-1 text-[11px] font-mono text-themeText-muted">
        <span>Press Enter to send, Shift+Enter for new line</span>
        <span className={isOverLimit ? "text-status-error flex items-center gap-1" : ""}>
          {isOverLimit && <AlertCircle className="w-3 h-3" />}
          {charCount}/1000
        </span>
      </div>
    </form>
  );
}
