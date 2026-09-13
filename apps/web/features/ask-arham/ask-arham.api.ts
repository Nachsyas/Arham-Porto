import type { GroundedResponse, PublicEvidenceItem } from "./ask-arham.types";

export interface StreamCallbacks {
  onStatus?: (status: string) => void;
  onEvidence?: (evidence: PublicEvidenceItem[]) => void;
  onResult?: (result: GroundedResponse) => void;
  onError?: (err: { code: string; message: string }) => void;
}

export class APIError extends Error {
  code: string;
  status: number;

  constructor(code: string, message: string, status: number) {
    super(message);
    this.name = "APIError";
    this.code = code;
    this.status = status;
  }
}

/**
 * Parses raw SSE chunk data into individual events.
 * Exported for independent unit testing of fragmented network scenarios.
 */
export function parseSSEEvents(buffer: string): {
  events: Array<{ event: string; data: string }>;
  remaining: string;
} {
  const events: Array<{ event: string; data: string }> = [];
  // Normalize CRLF to LF
  const normalized = buffer.replace(/\r\n/g, "\n");
  const rawEvents = normalized.split("\n\n");

  // The last part after the last \n\n may be an incomplete event
  const remaining = rawEvents.pop() ?? "";

  for (const block of rawEvents) {
    const trimmedBlock = block.trim();
    if (!trimmedBlock) continue;

    let eventType = "message";
    const dataLines: string[] = [];

    const lines = trimmedBlock.split("\n");
    for (const line of lines) {
      if (line.startsWith("event:")) {
        eventType = line.slice(6).trim();
      } else if (line.startsWith("data:")) {
        dataLines.push(line.slice(5).trim());
      }
    }

    if (dataLines.length > 0) {
      events.push({
        event: eventType,
        data: dataLines.join("\n"),
      });
    }
  }

  return { events, remaining };
}

/**
 * Sends a question to Ask Arham AI and consumes the SSE stream.
 */
export async function streamAskQuestion(
  question: string,
  callbacks: StreamCallbacks,
  signal?: AbortSignal
): Promise<void> {
  const apiBase =
    process.env.NEXT_PUBLIC_API_BASE_URL ||
    (typeof window !== "undefined" ? "http://localhost:8080" : "http://localhost:8080");

  try {
    const response = await fetch(`${apiBase}/api/v1/ai/ask`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "text/event-stream",
      },
      body: JSON.stringify({ question }),
      signal,
    });

    // Preflight error handling (Correction 9)
    if (!response.ok) {
      let errCode = "unknown_error";
      let errMsg = "An error occurred while contacting Ask Arham AI.";

      try {
        const errJson = await response.json();
        if (errJson?.error) {
          errCode = errJson.error.code || errCode;
          errMsg = errJson.error.message || errMsg;
        }
      } catch {
        errMsg = `Request failed with status ${response.status}`;
      }

      throw new APIError(errCode, errMsg, response.status);
    }

    // Stream consumption
    if (!response.body) {
      throw new APIError("no_body", "No response body received", 500);
    }

    const reader = response.body.getReader();
    const decoder = new TextDecoder("utf-8");
    let buffer = "";

    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const parsed = parseSSEEvents(buffer);
      buffer = parsed.remaining;

      for (const ev of parsed.events) {
        if (ev.event === "status") {
          try {
            const data = JSON.parse(ev.data);
            callbacks.onStatus?.(data.status);
          } catch {
            callbacks.onStatus?.(ev.data);
          }
        } else if (ev.event === "evidence") {
          try {
            const data = JSON.parse(ev.data);
            callbacks.onEvidence?.(data.evidence || []);
          } catch (e) {
            console.error("Failed to parse evidence SSE event", e);
          }
        } else if (ev.event === "result") {
          try {
            const result = JSON.parse(ev.data) as GroundedResponse;
            callbacks.onResult?.(result);
          } catch (e) {
            console.error("Failed to parse result SSE event", e);
          }
        } else if (ev.event === "error") {
          try {
            const errData = JSON.parse(ev.data);
            callbacks.onError?.({
              code: errData.code || "service_unavailable",
              message: errData.message || "Service temporarily unavailable",
            });
          } catch {
            callbacks.onError?.({
              code: "service_unavailable",
              message: ev.data,
            });
          }
        } else if (ev.event === "done") {
          // Stream completed cleanly
          return;
        }
      }
    }
  } catch (err: unknown) {
    if (signal?.aborted || (err instanceof Error && err.name === "AbortError")) {
      // User aborted stream cleanly
      return;
    }
    throw err;
  }
}
