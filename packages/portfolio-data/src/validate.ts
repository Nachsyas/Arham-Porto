import {
  getProfile,
  getProjects,
  getSkills,
  getEvidence,
  getExperience,
  getEducation,
  getJourney,
  getAllowlist,
} from "./loader";

console.log("==================================================");
console.log("ARHAM PORTO — CANONICAL DATA RUNTIME VALIDATION");
console.log("==================================================");

let hasErrors = false;

function check(label: string, fn: () => unknown) {
  try {
    const data = fn();
    const count = Array.isArray(data) ? data.length : 1;
    console.log(`[PASS] ${label} (${count} item(s) validated)`);
  } catch (err) {
    console.error(`[FAIL] ${label}:`, err);
    hasErrors = true;
  }
}

check("Profile Schema", () => getProfile());
check("Projects Schema", () => getProjects());
check("Skills Schema", () => getSkills());
check("Evidence Schema", () => getEvidence());
check("Experience Schema", () => getExperience());
check("Education Schema", () => getEducation());
check("Journey Schema", () => getJourney());
check("GitHub Allowlist Schema", () => getAllowlist());

console.log("==================================================");
if (hasErrors) {
  console.error("FAILED: One or more canonical data files failed schema validation.");
  process.exit(1);
} else {
  console.log("SUCCESS: All canonical portfolio data passed Zod validation.");
  process.exit(0);
}
