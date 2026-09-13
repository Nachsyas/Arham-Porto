import type { JourneyStop } from "arham-porto-schema";
import type { PublicMilestone, GeographicLocation, LocationId } from "./journey.types";
import { GEOGRAPHIC_LOCATIONS } from "./data/indonesia-outline";

/**
 * Maps a journey stop ID or category to its corresponding unique geographic location ID.
 * Exactly 4 unique locations exist: Karanganyar, Jakarta, Salatiga, Malang.
 */
export function getLocationIdForStop(stop: JourneyStop): LocationId {
  if (stop.id === "origin-karanganyar" || stop.category === "birthplace") {
    return "karanganyar";
  }
  if (
    stop.id === "residence-jakarta" ||
    stop.id === "mi-al-hamid-jakarta" ||
    stop.id === "mtsn30-jakarta" ||
    stop.category === "residence" ||
    stop.city === "Jakarta" ||
    stop.city === "Jakarta Timur"
  ) {
    return "jakarta";
  }
  if (stop.id === "ma-assurkati-salatiga" || stop.category === "sma") {
    return "salatiga";
  }
  return "malang";
}

/**
 * Filters canonical journey stops to only those explicitly approved for public rendering.
 * Early education (TK) with public=false is strictly excluded.
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
export function getGeographicLocation(id: LocationId): GeographicLocation {
  return GEOGRAPHIC_LOCATIONS[id];
}

/**
 * Constructs the SVG route path connecting the 4 unique locations:
 * Karanganyar (Origin) -> Jakarta (Residence / MI / MTs) -> Salatiga (MA) -> Malang (University & Current Base).
 * Exactly THREE geographic segments. No Jakarta -> Jakarta and no Malang -> Malang segments.
 */
export function getCorridorRoutePath(): string {
  const k = GEOGRAPHIC_LOCATIONS.karanganyar;
  const j = GEOGRAPHIC_LOCATIONS.jakarta;
  const s = GEOGRAPHIC_LOCATIONS.salatiga;
  const m = GEOGRAPHIC_LOCATIONS.malang;

  // Segment 1: Karanganyar -> Jakarta (arching northwest across northern Java)
  // Segment 2: Jakarta -> Salatiga (arching southeast back to Central Java)
  // Segment 3: Salatiga -> Malang (smooth curve into East Java)
  return `M ${k.x} ${k.y} Q 320 100 ${j.x} ${j.y} Q 310 160 ${s.x} ${s.y} Q 520 185 ${m.x} ${m.y}`;
}
