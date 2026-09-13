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
              <div className="bg-primary/15 border border-primary/30 text-themeText-primary text-sm p-3 rounded-2xl rounded-tr-none max-w-[85%] break-words">
                {msg.question}
              </div>
              <div className="w-7 h-7 rounded-full bg-primary/20 border border-primary/40 flex items-center justify-center flex-shrink-0 text-primary">
                <User className="w-3.5 h-3.5" />
              </div>
            </div>
          )}

          {/* Assistant Response */}
          {msg.response ? (
            <div className="flex items-start gap-2.5 justify-start">
              <div className="w-7 h-7 rounded-full bg-surface border border-border flex items-center justify-center flex-shrink-0 text-primary">
                <Bot className="w-3.5 h-3.5" />
              </div>
              <div className="bg-surface/60 border border-border/80 p-3.5 rounded-2xl rounded-tl-none max-w-[92%] flex-1">
                <AskArhamAnswer
                  response={msg.response}
                  onActionTriggered={onActionTriggered}
                />
              </div>
            </div>
          ) : msg.errorMessage ? (
            <div className="flex items-start gap-2.5 justify-start">
              <div className="w-7 h-7 rounded-full bg-red-500/10 border border-red-500/30 flex items-center justify-center flex-shrink-0 text-red-400">
                <AlertTriangle className="w-3.5 h-3.5" />
              </div>
              <div className="bg-red-500/10 border border-red-500/20 text-red-300 text-sm p-3 rounded-2xl rounded-tl-none max-w-[90%]">
                {msg.errorMessage}
              </div>
            </div>
          ) : msg.status && msg.status !== "idle" && msg.status !== "ready" ? (
            <div className="flex items-start gap-2.5 justify-start">
              <div className="w-7 h-7 rounded-full bg-surface border border-border flex items-center justify-center flex-shrink-0 text-primary animate-pulse">
                <Bot className="w-3.5 h-3.5" />
              </div>
              <div className="bg-surface/40 border border-border/60 text-themeText-muted text-xs font-mono p-3 rounded-2xl rounded-tl-none flex items-center gap-2">
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
