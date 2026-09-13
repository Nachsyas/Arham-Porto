import { z } from "zod";

export const ProjectSchema = z.object({
  id: z.string().min(1),
  title: z.string().min(1),
  slug: z.string().min(1),
  summary: z.string().nullable().default(null),
  problem: z.string().nullable().default(null),
  solution: z.string().nullable().default(null),
  role: z.array(z.string()).default([]),
  contributions: z.array(z.string()).default([]),
  technologies: z.array(z.string()).default([]),
  githubUrl: z.string().nullable().default(null),
  demoUrl: z.string().nullable().default(null),
  image: z.string().nullable().default(null),
  featured: z.boolean().default(false),
  category: z.enum(["AI", "Full-Stack", "Backend", "Systems"]).nullable().default(null),
  evidenceIds: z.array(z.string()).default([]),
  todo: z.array(z.string()).optional(),
});

export const ProjectsContainerSchema = z.object({
  projects: z.array(ProjectSchema),
});

export type Project = z.infer<typeof ProjectSchema>;
export type ProjectsContainer = z.infer<typeof ProjectsContainerSchema>;
