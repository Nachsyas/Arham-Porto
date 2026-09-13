import {
  getProfile,
  getProjects,
  getSkills,
  getExperience,
  getEducation,
} from "arham-porto-data";
import PortfolioApp from "@/features/portfolio/PortfolioApp";

export default function HomePage() {
  // Load server-validated canonical data
  const profile = getProfile();
  const projects = getProjects();
  const skills = getSkills();
  const experiences = getExperience();
  const education = getEducation();

  return (
    <PortfolioApp
      profile={profile}
      projects={projects}
      skills={skills}
      experiences={experiences}
      education={education}
    />
  );
}
