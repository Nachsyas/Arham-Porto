import { z } from "zod";

export const RepositoryAllowlistEntrySchema = z.object({
  repository: z.string().min(1),
  enabled: z.boolean().default(true),
  indexingStatus: z.enum(["not_indexed", "indexed", "error"]).default("not_indexed"),
  lastIndexedAt: z.string().nullable().default(null),
  notes: z.string().nullable().default(null),
});

export const GitHubAllowlistSchema = z.object({
  repositories: z.array(RepositoryAllowlistEntrySchema),
});

export type RepositoryAllowlistEntry = z.infer<typeof RepositoryAllowlistEntrySchema>;
export type GitHubAllowlist = z.infer<typeof GitHubAllowlistSchema>;
