import { z } from "zod";

export const ExperienceItemSchema = z.object({
  id: z.string().min(1),
  role: z.string().nullable().default(null),
  company: z.string().nullable().default(null),
  period: z.string().nullable().default(null),
  description: z.string().nullable().default(null),
  skills: z.array(z.string()).default([]),
  evidenceIds: z.array(z.string()).default([]),
  todo: z.array(z.string()).optional(),
});

export const ExperienceContainerSchema = z.object({
  experiences: z.array(ExperienceItemSchema).default([]),
  todo: z.array(z.string()).optional(),
});

export type ExperienceItem = z.infer<typeof ExperienceItemSchema>;
export type ExperienceContainer = z.infer<typeof ExperienceContainerSchema>;
