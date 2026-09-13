import type { JourneyStop } from "arham-porto-schema";

export type LocationId = "karanganyar" | "jakarta" | "salatiga" | "malang";

export interface GeographicLocation {
  id: LocationId;
  name: string;
  province: string;
  x: number;
  y: number;
  milestoneIds: string[];
}

export type MilestoneCategory =
  | "birthplace"
  | "residence"
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
  locationId: LocationId;
}
