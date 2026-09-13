import {
  ProfileSchema,
  type Profile,
  ProjectsContainerSchema,
  type Project,
  SkillsContainerSchema,
  type Skill,
  EvidenceContainerSchema,
  type Evidence,
  ExperienceContainerSchema,
  type ExperienceItem,
  EducationContainerSchema,
  type EducationItem,
  JourneyContainerSchema,
  type JourneyStop,
  GitHubAllowlistSchema,
  type RepositoryAllowlistEntry,
} from "arham-porto-schema";

// Static JSON imports guaranteeing server-safe bundling without fragile filesystem paths
import rawProfile from "../../../data/profile/profile.json";
import rawProjects from "../../../data/projects/projects.json";
import rawSkills from "../../../data/skills/skills.json";
import rawEvidence from "../../../data/evidence/evidence.json";
import rawExperience from "../../../data/experience/experience.json";
import rawEducation from "../../../data/education/education.json";
import rawJourney from "../../../data/journey/journey.json";
import rawAllowlist from "../../../data/ai/allowlist.json";

function validateOrThrow<T>(name: string, schema: { safeParse: (data: unknown) => { success: boolean; data?: T; error?: unknown } }, data: unknown): T {
  const result = schema.safeParse(data);
  if (!result.success) {
    throw new Error(`[portfolio-data] Validation failed for ${name}: ${JSON.stringify(result.error, null, 2)}`);
  }
  return result.data as T;
}

export function getProfile(): Profile {
  return validateOrThrow<Profile>("Profile", ProfileSchema, rawProfile);
}

export function getProjects(): Project[] {
  const container = validateOrThrow<{ projects: Project[] }>("Projects", ProjectsContainerSchema, rawProjects);
  return container.projects;
}

export function getSkills(): Skill[] {
  const container = validateOrThrow<{ skills: Skill[] }>("Skills", SkillsContainerSchema, rawSkills);
  return container.skills;
}

export function getEvidence(): Evidence[] {
  const container = validateOrThrow<{ items: Evidence[] }>("Evidence", EvidenceContainerSchema, rawEvidence);
  return container.items;
}

export function getExperience(): ExperienceItem[] {
  const container = validateOrThrow<{ experiences: ExperienceItem[] }>("Experience", ExperienceContainerSchema, rawExperience);
  return container.experiences;
}

export function getEducation(): EducationItem[] {
  const container = validateOrThrow<{ education: EducationItem[] }>("Education", EducationContainerSchema, rawEducation);
  return container.education;
}

export function getJourney(): JourneyStop[] {
  const container = validateOrThrow<{ stops: JourneyStop[] }>("Journey", JourneyContainerSchema, rawJourney);
  return container.stops;
}

export function getAllowlist(): RepositoryAllowlistEntry[] {
  const container = validateOrThrow<{ repositories: RepositoryAllowlistEntry[] }>("GitHubAllowlist", GitHubAllowlistSchema, rawAllowlist);
  return container.repositories;
}
