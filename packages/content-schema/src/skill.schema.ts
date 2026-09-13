import { z } from "zod";

export const SkillSchema = z.object({
  id: z.string().min(1),
  name: z.string().min(1),
  category: z.enum([
    "Backend",
    "Frontend",
    "AI / ML",
    "Systems",
    "DevOps",
    "Database",
  ]),
  claim: z.string().nullable().default(null),
  evidenceIds: z.array(z.string()).default([]),
  todo: z.array(z.string()).optional(),
});

export const SkillsContainerSchema = z.object({
  skills: z.array(SkillSchema),
});

export type Skill = z.infer<typeof SkillSchema>;
export type SkillsContainer = z.infer<typeof SkillsContainerSchema>;
