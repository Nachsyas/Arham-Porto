import type { JourneyStop } from "arham-porto-schema";

export interface GeographicLocation {
  id: "karanganyar" | "salatiga" | "malang";
  name: string;
  province: string;
  x: number;
  y: number;
  milestoneIds: string[];
}

export type MilestoneCategory =
  | "birthplace"
  | "tk"
  | "sd"
  | "smp"
  | "sma"
  | "university"
  | "current";

export interface PublicMilestone {
  id: string;
  category: MilestoneCategory;
  title: string;
  institution: string | null;
  city: string;
  region: string;
  period: string | null;
  description: string | null;
  locationId: "karanganyar" | "salatiga" | "malang";
}
