import {
  getProfile,
  getProjects,
  getSkills,
  getExperience,
  getEducation,
  getJourney,
} from "arham-porto-data";
import PortfolioApp from "@/features/portfolio/PortfolioApp";

export default function HomePage() {
  // Load server-validated canonical data
  const profile = getProfile();
  const projects = getProjects();
  const skills = getSkills();
  const experiences = getExperience();
  const education = getEducation();
  const journeyStops = getJourney();

  return (
    <PortfolioApp
      profile={profile}
      projects={projects}
      skills={skills}
      journeyStops={journeyStops}
      experiences={experiences}
      education={education}
    />
  );
}
