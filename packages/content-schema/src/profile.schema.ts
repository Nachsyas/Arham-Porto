import { z } from "zod";

export const ProfileSchema = z.object({
  fullName: z.string().min(1),
  role: z.string().min(1),
  projectName: z.string().min(1),
  aiFeature: z.string().min(1),
  positioning: z.string().nullable().default(null),
  bio: z.string().nullable().default(null),
  email: z.string().nullable().default(null),
  linkedin: z.string().nullable().default(null),
  github: z.string().nullable().default(null),
  cvUrl: z.string().nullable().default(null),
  currentCity: z.string().nullable().default(null),
  availability: z.string().nullable().default(null),
  todo: z.array(z.string()).optional(),
});

export type Profile = z.infer<typeof ProfileSchema>;
