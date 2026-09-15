import {
  getProfile,
  getProjects,
  getSkills,
  getJourney,
} from "arham-porto-data";
import HomeView from "@/features/home/HomeView";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Nachsyas Arham Mumtaz Nashohi | Software Engineer",
  description:
    "Personal engineering portfolio of Nachsyas Arham Mumtaz Nashohi. Clean architecture, backend systems in Go, modern Next.js interfaces, and grounded AI retrieval.",
  alternates: {
    canonical: "/",
  },
};

export default function HomePage() {
  const profile = getProfile();
  const allProjects = getProjects();
  const featuredProjects = allProjects.filter((p) => p.featured).slice(0, 3);
  const skills = getSkills();
  const journeyStops = getJourney();

  return (
    <HomeView
      profile={profile}
      featuredProjects={featuredProjects}
      skills={skills}
      journeyStops={journeyStops}
    />
  );
}
