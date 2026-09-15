import React, { useEffect, useRef, useState, useCallback } from "react";
import { X, Sparkles, RotateCcw, Bot } from "lucide-react";
import type { AskMessage, GroundedResponse } from "./ask-arham.types";
import { streamAskQuestion, APIError } from "./ask-arham.api";
import AskArhamMessages from "./AskArhamMessages";
import AskArhamInput from "./AskArhamInput";
import SuggestedQuestions from "./SuggestedQuestions";

interface AskArhamPanelProps {
  isOpen: boolean;
  onClose: () => void;
  returnFocusRef?: React.RefObject<HTMLElement | null>;
}

export default function AskArhamPanel({
  isOpen,
  onClose,
  returnFocusRef,
}: AskArhamPanelProps) {
  const [messages, setMessages] = useState<AskMessage[]>([]);
  const [currentStatus, setCurrentStatus] = useState<string>("idle");
  const [announcement, setAnnouncement] = useState<string>("");
  const [isUnavailable, setIsUnavailable] = useState(false);
  const panelRef = useRef<HTMLDivElement>(null);
  const abortControllerRef = useRef<AbortController | null>(null);

  // Focus management & Escape key handling (Correction 45)
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        e.preventDefault();
        handleClose();
      }
    };

    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [isOpen]);

  // Return focus on close (Correction 45)
  const handleClose = useCallback(() => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      abortControllerRef.current = null;
    }
    onClose();
    setTimeout(() => {
      returnFocusRef?.current?.focus();
    }, 50);
  }, [onClose, returnFocusRef]);

  const handleStopStream = () => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort();
      abortControllerRef.current = null;
      setCurrentStatus("idle");
      setAnnouncement("Generation stopped.");
    }
  };

  const handleClearSession = () => {
    handleStopStream();
    setMessages([]);
    setCurrentStatus("idle");
    setAnnouncement("Conversation cleared.");
    setIsUnavailable(false);
  };

  const handleSubmitQuestion = async (question: string) => {
    const userMsgId = `usr_${Date.now()}`;
    const assistantMsgId = `ast_${Date.now() + 1}`;

    const userMsg: AskMessage = {
      id: userMsgId,
      role: "user",
      question,
      createdAt: new Date(),
    };

    const assistantMsg: AskMessage = {
      id: assistantMsgId,
      role: "assistant",
      status: "retrieving",
      createdAt: new Date(),
    };

    setMessages((prev) => [...prev, userMsg, assistantMsg]);
    setCurrentStatus("retrieving");
    setAnnouncement("Retrieving evidence...");

    const controller = new AbortController();
    abortControllerRef.current = controller;

    try {
      await streamAskQuestion(
        question,
        {
          onStatus: (status) => {
            setCurrentStatus(status);
            if (status === "generating") {
              setAnnouncement("Generating response...");
            }
            setMessages((prev) =>
              prev.map((m) =>
                m.id === assistantMsgId ? { ...m, status: status as any } : m
              )
            );
          },
          onEvidence: (evidence) => {
            // Early evidence metadata received
          },
          onResult: (result: GroundedResponse) => {
            setCurrentStatus("ready");
            setAnnouncement("Response ready.");
            setMessages((prev) =>
              prev.map((m) =>
                m.id === assistantMsgId
                  ? { ...m, response: result, status: "ready" }
                  : m
              )
            );
          },
          onError: (err) => {
            setCurrentStatus("error");
            setAnnouncement("An error occurred.");
            setMessages((prev) =>
              prev.map((m) =>
                m.id === assistantMsgId
                  ? { ...m, status: "error", errorMessage: err.message }
                  : m
              )
            );
          },
        },
        controller.signal
      );
    } catch (err: unknown) {
      if (controller.signal.aborted) {
        return;
      }

      setCurrentStatus("error");
      setAnnouncement("Request failed.");

      let errMsg = "Ask Arham AI is temporarily unavailable.";
      if (err instanceof APIError) {
        if (err.status === 503) {
          setIsUnavailable(true);
          errMsg = err.message || "Ask Arham AI is currently unavailable.";
        } else if (err.status === 429) {
          errMsg = "Rate limit reached. Please wait a minute before asking another question.";
        } else if (err.message) {
          errMsg = err.message;
        }
      } else if (err instanceof Error && err.message) {
        errMsg = err.message;
      }

      setMessages((prev) =>
        prev.map((m) =>
          m.id === assistantMsgId
            ? { ...m, status: "error", errorMessage: errMsg }
            : m
        )
      );
    } finally {
      abortControllerRef.current = null;
    }
  };

  if (!isOpen) return null;

  const isLoading = currentStatus === "retrieving" || currentStatus === "generating";

  return (
    <div
      className="fixed inset-0 z-50 flex justify-end bg-slate-900/25 backdrop-blur-sm transition-opacity"
      role="dialog"
      aria-modal="true"
      aria-labelledby="ask-arham-title"
    >
      {/* Screen Reader Announcement Region (Correction 43) */}
      <div className="sr-only" aria-live="polite">
        {announcement}
      </div>

      {/* Slide-over Drawer / Mobile Sheet (Correction 46) */}
      <div
        ref={panelRef}
        className="relative w-full max-w-lg bg-surface border-l border-border flex flex-col shadow-2xl h-[100dvh] max-h-[100dvh] pb-[env(safe-area-inset-bottom)]"
      >
        {/* Header */}
        <div className="flex items-center justify-between px-4 py-3.5 border-b border-border bg-surface-elevated">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-primary/10 border border-primary/25 flex items-center justify-center text-primary">
              <Bot className="w-4 h-4" />
            </div>
            <div>
              <h2 id="ask-arham-title" className="text-sm font-bold text-themeText-primary font-display flex items-center gap-1.5">
                Ask Arham AI
                <span className="text-[10px] font-mono font-normal px-2 py-0.5 rounded-full bg-primary-muted text-primary border border-primary/25">
                  Reviewer Copilot
                </span>
              </h2>
              <p className="text-xs text-themeText-muted font-mono">
                Grounded in verified portfolio evidence
              </p>
            </div>
          </div>

          <div className="flex items-center gap-1">
            {messages.length > 0 && (
              <button
                type="button"
                onClick={handleClearSession}
                className="p-1.5 rounded-md text-themeText-muted hover:text-themeText-primary hover:bg-surface-strong transition-colors"
                title="Reset conversation"
                aria-label="Reset conversation"
              >
                <RotateCcw className="w-4 h-4" />
              </button>
            )}
            <button
              type="button"
              onClick={handleClose}
              className="p-1.5 rounded-md text-themeText-muted hover:text-themeText-primary hover:bg-surface-strong transition-colors"
              aria-label="Close Ask Arham panel"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        </div>

        {/* Unavailable Banner if 503 received (Correction 46, 58) */}
        {isUnavailable && (
          <div className="p-3 bg-amber-50 border-b border-amber-200 text-xs text-amber-900 flex items-center gap-2">
            <Sparkles className="w-4 h-4 flex-shrink-0 text-status-warning" />
            <span>
              Ask Arham AI is currently unavailable. You can continue reviewing all verified projects, skills, and the journey map directly on the portfolio.
            </span>
          </div>
        )}

        {/* Content Area */}
        <div className="flex-1 overflow-y-auto flex flex-col justify-between">
          {messages.length === 0 ? (
            <div className="p-6 flex flex-col justify-center my-auto">
              <div className="w-12 h-12 rounded-2xl bg-primary/10 border border-primary/30 flex items-center justify-center text-primary mb-3">
                <Sparkles className="w-6 h-6" />
              </div>
              <h3 className="text-lg font-bold text-themeText-primary font-display">
                Reviewer-Oriented AI Assistant
              </h3>
              <p className="mt-1 text-sm text-themeText-muted leading-relaxed">
                Ask targeted technical questions about Nachsyas Arham Mumtaz Nashohi&apos;s verified engineering experience, architecture decisions, and technology stack.
              </p>

              <SuggestedQuestions
                onSelect={handleSubmitQuestion}
                disabled={isLoading || isUnavailable}
              />
            </div>
          ) : (
            <AskArhamMessages
              messages={messages}
              currentStatus={currentStatus}
              onActionTriggered={handleClose}
            />
          )}
        </div>

        {/* Input Form */}
        <AskArhamInput
          onSubmit={handleSubmitQuestion}
          onAbort={handleStopStream}
          isLoading={isLoading}
          disabled={isUnavailable}
        />
      </div>
    </div>
  );
}
