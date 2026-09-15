import React, { useEffect, useRef } from "react";
import { Bot, User, Loader2, AlertTriangle } from "lucide-react";
import type { AskMessage } from "./ask-arham.types";
import AskArhamAnswer from "./AskArhamAnswer";

interface AskArhamMessagesProps {
  messages: AskMessage[];
  currentStatus: string;
  onActionTriggered?: () => void;
}

export default function AskArhamMessages({
  messages,
  currentStatus,
  onActionTriggered,
}: AskArhamMessagesProps) {
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages, currentStatus]);

  if (messages.length === 0) {
    return null;
  }

  return (
    <div className="flex-1 overflow-y-auto p-4 flex flex-col gap-4">
      {messages.map((msg) => (
        <div key={msg.id} className="flex flex-col gap-2">
          {/* User Question */}
          {msg.question && (
            <div className="flex items-start gap-2.5 justify-end">
              <div className="bg-primary-muted border border-primary/25 text-themeText-primary text-sm p-3 rounded-2xl rounded-tr-none max-w-[85%] break-words shadow-sm">
                {msg.question}
              </div>
              <div className="w-7 h-7 rounded-full bg-primary-muted border border-primary/30 flex items-center justify-center flex-shrink-0 text-primary">
                <User className="w-3.5 h-3.5" />
              </div>
            </div>
          )}

          {/* Assistant Response */}
          {msg.response ? (
            <div className="flex items-start gap-2.5 justify-start">
              <div className="w-7 h-7 rounded-full bg-surface border border-border flex items-center justify-center flex-shrink-0 text-primary shadow-sm">
                <Bot className="w-3.5 h-3.5" />
              </div>
              <div className="bg-surface border border-border p-3.5 rounded-2xl rounded-tl-none max-w-[92%] flex-1 shadow-sm">
                <AskArhamAnswer
                  response={msg.response}
                  onActionTriggered={onActionTriggered}
                />
              </div>
            </div>
          ) : msg.errorMessage ? (
            <div className="flex items-start gap-2.5 justify-start">
              <div className="w-7 h-7 rounded-full bg-red-50 border border-red-200 flex items-center justify-center flex-shrink-0 text-status-error">
                <AlertTriangle className="w-3.5 h-3.5" />
              </div>
              <div className="bg-red-50 border border-red-200 text-red-800 text-sm p-3 rounded-2xl rounded-tl-none max-w-[90%] shadow-sm">
                {msg.errorMessage}
              </div>
            </div>
          ) : msg.status && msg.status !== "idle" && msg.status !== "ready" ? (
            <div className="flex items-start gap-2.5 justify-start">
              <div className="w-7 h-7 rounded-full bg-surface border border-border flex items-center justify-center flex-shrink-0 text-primary animate-pulse shadow-sm">
                <Bot className="w-3.5 h-3.5" />
              </div>
              <div className="bg-surface-elevated border border-border text-themeText-muted text-xs font-mono p-3 rounded-2xl rounded-tl-none flex items-center gap-2 shadow-sm">
                <Loader2 className="w-3.5 h-3.5 animate-spin text-primary" />
                <span>
                  {msg.status === "retrieving"
                    ? "Retrieving approved evidence from Phase 5 index..."
                    : "Analyzing verified evidence and composing grounded response..."}
                </span>
              </div>
            </div>
          ) : null}
        </div>
      ))}
      <div ref={bottomRef} />
    </div>
  );
}
