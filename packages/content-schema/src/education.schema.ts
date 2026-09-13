import { z } from "zod";

export const EducationItemSchema = z.object({
  id: z.string().min(1),
  institution: z.string().nullable().default(null),
  degree: z.string().nullable().default(null),
  field: z.string().nullable().default(null),
  period: z.string().nullable().default(null),
  description: z.string().nullable().default(null),
  evidenceIds: z.array(z.string()).default([]),
  todo: z.array(z.string()).optional(),
});

export const EducationContainerSchema = z.object({
  education: z.array(EducationItemSchema).default([]),
  todo: z.array(z.string()).optional(),
});

export type EducationItem = z.infer<typeof EducationItemSchema>;
export type EducationContainer = z.infer<typeof EducationContainerSchema>;
