import type { JourneyStop } from "arham-porto-schema";
import type { PublicMilestone, GeographicLocation } from "./journey.types";
import { GEOGRAPHIC_LOCATIONS } from "./data/indonesia-outline";

/**
 * Maps a journey stop ID or category to its corresponding unique geographic location ID.
 * Exactly 3 unique locations exist: Karanganyar, Salatiga, Malang.
 */
export function getLocationIdForStop(stop: JourneyStop): "karanganyar" | "salatiga" | "malang" {
  if (stop.id === "origin-karanganyar" || stop.category === "birthplace") {
    return "karanganyar";
  }
  if (stop.id === "ma-assurkati-salatiga" || stop.category === "sma") {
    return "salatiga";
  }
  return "malang";
}

/**
 * Filters canonical journey stops to only those explicitly approved for public rendering.
 * Early education (TK, SD, SMP) with public=false are strictly excluded.
 */
export function filterPublicMilestones(stops: JourneyStop[] = []): PublicMilestone[] {
  if (!Array.isArray(stops)) return [];
  return stops
    .filter((stop) => stop && stop.public === true)
    .map((stop) => ({
      id: stop.id,
      category: stop.category,
      title: stop.title || stop.institution || "Milestone",
      institution: stop.institution,
      city: stop.city || "Indonesia",
      region: stop.region || "Indonesia",
      period: stop.period,
      description: stop.description,
      locationId: getLocationIdForStop(stop),
    }));
}

/**
 * Returns the geographic location metadata for a given location ID.
 */
export function getGeographicLocation(id: "karanganyar" | "salatiga" | "malang"): GeographicLocation {
  return GEOGRAPHIC_LOCATIONS[id];
}

/**
 * Constructs the SVG route path connecting the 3 unique locations:
 * Karanganyar (Origin) -> Salatiga (Secondary Education) -> Malang (University & Current Base).
 * No redundant segment exists between University and Current Base.
 */
export function getCorridorRoutePath(): string {
  const k = GEOGRAPHIC_LOCATIONS.karanganyar;
  const s = GEOGRAPHIC_LOCATIONS.salatiga;
  const m = GEOGRAPHIC_LOCATIONS.malang;

  // Segment 1: Karanganyar -> Salatiga (smooth curve across Central Java)
  // Segment 2: Salatiga -> Malang (smooth curve into East Java)
  return `M ${k.x} ${k.y} Q 460 185 ${s.x} ${s.y} Q 520 180 ${m.x} ${m.y}`;
}
