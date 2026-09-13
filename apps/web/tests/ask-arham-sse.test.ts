import { describe, it, expect, vi } from "vitest";
import { parseSSEEvents, streamAskQuestion, APIError } from "../features/ask-arham/ask-arham.api";
import type { GroundedResponse, PublicEvidenceItem } from "../features/ask-arham/ask-arham.types";

describe("Ask Arham AI SSE Stream Parser", () => {
  it("parses single event with standard LF line endings", () => {
    const raw = 'event: status\ndata: {"status":"retrieving"}\n\n';
    const { events, remaining } = parseSSEEvents(raw);

    expect(remaining).toBe("");
    expect(events).toHaveLength(1);
    expect(events[0]).toEqual({
      event: "status",
      data: '{"status":"retrieving"}',
    });
  });

  it("handles CRLF line endings cleanly", () => {
    const raw = 'event: status\r\ndata: {"status":"generating"}\r\n\r\n';
    const { events, remaining } = parseSSEEvents(raw);

    expect(remaining).toBe("");
    expect(events).toHaveLength(1);
    expect(events[0]).toEqual({
      event: "status",
      data: '{"status":"generating"}',
    });
  });

  it("handles multiple events in a single network chunk", () => {
    const raw =
      'event: status\ndata: {"status":"retrieving"}\n\n' +
      'event: status\ndata: {"status":"generating"}\n\n' +
      'event: done\ndata: {}\n\n';
    const { events, remaining } = parseSSEEvents(raw);

    expect(remaining).toBe("");
    expect(events).toHaveLength(3);
    expect(events[0].event).toBe("status");
    expect(events[1].event).toBe("status");
    expect(events[2].event).toBe("done");
  });

  it("retains incomplete events across chunk boundaries", () => {
    const chunk1 = 'event: status\ndata: {"status":';
    const { events: events1, remaining: rem1 } = parseSSEEvents(chunk1);
    expect(events1).toHaveLength(0);
    expect(rem1).toBe('event: status\ndata: {"status":');

    const chunk2 = rem1 + '"retrieving"}\n\n';
    const { events: events2, remaining: rem2 } = parseSSEEvents(chunk2);
    expect(rem2).toBe("");
    expect(events2).toHaveLength(1);
    expect(events2[0]).toEqual({
      event: "status",
      data: '{"status":"retrieving"}',
    });
  });

  it("handles split across arbitrary byte boundaries (character by character)", () => {
    const fullMessage =
      'event: result\ndata: {"status":"supported","answer":"Grounded answer."}\n\n';

    let buffer = "";
    const collectedEvents: Array<{ event: string; data: string }> = [];

    for (let i = 0; i < fullMessage.length; i++) {
      buffer += fullMessage[i];
      const parsed = parseSSEEvents(buffer);
      buffer = parsed.remaining;
      collectedEvents.push(...parsed.events);
    }

    expect(buffer).toBe("");
    expect(collectedEvents).toHaveLength(1);
    expect(collectedEvents[0].event).toBe("result");
    expect(JSON.parse(collectedEvents[0].data)).toEqual({
      status: "supported",
      answer: "Grounded answer.",
    });
  });

  it("parses multi-line SSE data payloads correctly", () => {
    const raw = "event: message\ndata: line one\ndata: line two\n\n";
    const { events, remaining } = parseSSEEvents(raw);

    expect(remaining).toBe("");
    expect(events).toHaveLength(1);
    expect(events[0]).toEqual({
      event: "message",
      data: "line one\nline two",
    });
  });
});

describe("streamAskQuestion Client Workflow", () => {
  it("handles successful SSE stream lifecycle", async () => {
    const mockEvidence: PublicEvidenceItem[] = [
      {
        id: "E1",
        title: "EduTrace - Clean Architecture",
        kind: "github",
        repository: "Nachsyas/EduTrace",
        excerpt: "Clean Architecture in Go",
      },
    ];

    const mockResult: GroundedResponse = {
      status: "supported",
      answer: "Arham follows Clean Architecture in Go.",
      segments: [
        {
          text: "Arham follows Clean Architecture in Go.",
          evidence_ids: ["E1"],
        },
      ],
      evidence: mockEvidence,
      sources: [
        {
          id: "E1",
          kind: "github",
          label: "EduTrace: Clean Architecture",
          url: "https://github.com/Nachsyas/EduTrace",
          repository: "Nachsyas/EduTrace",
        },
      ],
      actions: [
        {
          id: "view-project-edutrace",
          label: "View EduTrace Case Study",
        },
      ],
    };

    const sseChunks = [
      'event: status\ndata: {"status":"retrieving"}\n\n',
      `event: evidence\ndata: ${JSON.stringify({ evidence: mockEvidence })}\n\n`,
      'event: status\ndata: {"status":"generating"}\n\n',
      `event: result\ndata: ${JSON.stringify(mockResult)}\n\n`,
      'event: done\ndata: {}\n\n',
    ];

    // Create a mock ReadableStream
    const stream = new ReadableStream({
      start(controller) {
        for (const chunk of sseChunks) {
          controller.enqueue(new TextEncoder().encode(chunk));
        }
        controller.close();
      },
    });

    const mockFetch = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      body: stream,
    });
    global.fetch = mockFetch;

    const statuses: string[] = [];
    let receivedEvidence: PublicEvidenceItem[] = [];
    let finalResult: GroundedResponse | null = null;

    await streamAskQuestion("Tell me about Go architecture", {
      onStatus: (s) => statuses.push(s),
      onEvidence: (ev) => {
        receivedEvidence = ev;
      },
      onResult: (res) => {
        finalResult = res;
      },
    });

    expect(statuses).toEqual(["retrieving", "generating"]);
    expect(receivedEvidence).toHaveLength(1);
    expect(receivedEvidence[0].id).toBe("E1");
    expect(finalResult).not.toBeNull();
    expect((finalResult as unknown as GroundedResponse).status).toBe("supported");
    expect((finalResult as unknown as GroundedResponse).segments).toHaveLength(1);
  });

  it("handles HTTP preflight errors before SSE headers (503 AI disabled)", async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 503,
      json: async () => ({
        error: {
          code: "service_unavailable",
          message: "Ask Arham AI is currently unavailable.",
        },
      }),
    });

    await expect(
      streamAskQuestion("What is EduTrace?", {})
    ).rejects.toThrow(APIError);

    try {
      await streamAskQuestion("What is EduTrace?", {});
    } catch (err) {
      expect(err).toBeInstanceOf(APIError);
      const apiErr = err as APIError;
      expect(apiErr.code).toBe("service_unavailable");
      expect(apiErr.status).toBe(503);
    }
  });

  it("handles HTTP 429 rate limited error", async () => {
    global.fetch = vi.fn().mockResolvedValue({
      ok: false,
      status: 429,
      json: async () => ({
        error: {
          code: "rate_limit_exceeded",
          message: "AI rate limit exceeded. Please wait a minute.",
        },
      }),
    });

    try {
      await streamAskQuestion("What is EduTrace?", {});
    } catch (err) {
      expect(err).toBeInstanceOf(APIError);
      const apiErr = err as APIError;
      expect(apiErr.code).toBe("rate_limit_exceeded");
      expect(apiErr.status).toBe(429);
    }
  });

  it("handles in-stream error event", async () => {
    const sseChunks = [
      'event: status\ndata: {"status":"retrieving"}\n\n',
      'event: error\ndata: {"code":"service_unavailable","message":"Retrieval failed","request_id":"req-123"}\n\n',
      'event: done\ndata: {}\n\n',
    ];

    const stream = new ReadableStream({
      start(controller) {
        for (const chunk of sseChunks) {
          controller.enqueue(new TextEncoder().encode(chunk));
        }
        controller.close();
      },
    });

    global.fetch = vi.fn().mockResolvedValue({
      ok: true,
      status: 200,
      body: stream,
    });

    let reportedError: { code: string; message: string } | null = null;

    await streamAskQuestion("Tell me about Go", {
      onError: (err) => {
        reportedError = err;
      },
    });

    expect(reportedError).not.toBeNull();
    expect((reportedError as unknown as { code: string; message: string }).code).toBe("service_unavailable");
    expect((reportedError as unknown as { code: string; message: string }).message).toBe("Retrieval failed");
  });

  it("handles user client abort signal gracefully without re-throwing", async () => {
    const controller = new AbortController();

    const stream = new ReadableStream({
      start(ctrl) {
        ctrl.enqueue(new TextEncoder().encode('event: status\ndata: {"status":"retrieving"}\n\n'));
      },
      cancel() {},
    });

    global.fetch = vi.fn().mockImplementation(async (_url, opts) => {
      // Simulate abort triggered immediately
      controller.abort();
      if (opts?.signal?.aborted) {
        const error = new Error("The operation was aborted");
        error.name = "AbortError";
        throw error;
      }
      return {
        ok: true,
        status: 200,
        body: stream,
      };
    });

    // streamAskQuestion should return cleanly on signal.aborted without throwing
    await expect(
      streamAskQuestion("Cancelled question", {}, controller.signal)
    ).resolves.toBeUndefined();
  });
});
