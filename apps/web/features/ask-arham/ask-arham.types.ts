export type AskStatus = "idle" | "retrieving" | "generating" | "ready" | "error";

export type AnswerStatus =
  | "supported"
  | "insufficient_evidence"
  | "privacy_refusal"
  | "scope_refusal";

export interface GroundedSegment {
  text: string;
  evidence_ids: string[];
}

export interface PublicEvidenceItem {
  id: string;
  kind: "github" | "portfolio";
  title: string;
  repository?: string;
  path?: string;
  excerpt: string;
  citation_id?: string;
}

export interface SourceCitation {
  id: string;
  kind: "github" | "portfolio";
  label: string;
  url?: string;
  repository?: string;
  path?: string;
  commit_sha?: string;
}

export interface SafeAction {
  id: string;
  label: string;
}

export interface GroundedResponse {
  status: AnswerStatus;
  answer: string;
  segments: GroundedSegment[];
  evidence: PublicEvidenceItem[];
  sources: SourceCitation[];
  actions: SafeAction[];
}

export interface AskMessage {
  id: string;
  role: "user" | "assistant";
  question?: string;
  response?: GroundedResponse;
  status?: AskStatus;
  errorMessage?: string;
  createdAt: Date;
}
