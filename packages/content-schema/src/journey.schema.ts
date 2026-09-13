import { z } from "zod";

export const JourneyCategorySchema = z.enum([
  "birthplace",
  "tk",
  "sd",
  "smp",
  "sma",
  "university",
  "current",
]);

export const JourneyStopSchema = z.object({
  id: z.string().min(1),
  category: JourneyCategorySchema,
  title: z.string().nullable().default(null),
  institution: z.string().nullable().default(null),
  city: z.string().nullable().default(null),
  region: z.string().nullable().default(null),
  country: z.string().default("Indonesia"),
  coordinates: z.tuple([z.number(), z.number()]).nullable().default(null),
  period: z.string().nullable().default(null),
  description: z.string().nullable().default(null),
  image: z.string().nullable().default(null),
  public: z.boolean().default(false),
  todo: z.array(z.string()).optional(),
});

export const JourneyContainerSchema = z.object({
  stops: z.array(JourneyStopSchema),
});

export type JourneyCategory = z.infer<typeof JourneyCategorySchema>;
export type JourneyStop = z.infer<typeof JourneyStopSchema>;
export type JourneyContainer = z.infer<typeof JourneyContainerSchema>;
