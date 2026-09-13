import { z } from "zod";

export const EvidenceSchema = z.object({
  id: z.string().min(1),
  type: z.enum(["project", "github", "cv", "experience", "education"]),
  title: z.string().min(1),
  skillIds: z.array(z.string()).default([]),
  sourceUrl: z.string().nullable().default(null),
  sourcePath: z.string().nullable().default(null),
  summary: z.string().min(1),
  verified: z.boolean().default(false),
  todo: z.array(z.string()).optional(),
});

export const EvidenceContainerSchema = z.object({
  items: z.array(EvidenceSchema).default([]),
});

export type Evidence = z.infer<typeof EvidenceSchema>;
export type EvidenceContainer = z.infer<typeof EvidenceContainerSchema>;
