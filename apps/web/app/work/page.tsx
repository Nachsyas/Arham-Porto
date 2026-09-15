import { getProjects } from "arham-porto-data";
import PageHeader from "@/components/PageHeader";
import ProjectsSection from "@/features/projects/ProjectsSection";
import InterPageNav from "@/components/InterPageNav";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Selected Work & Systems",
  description:
    "Explore architectural case studies by Nachsyas Arham Mumtaz Nashohi across backend engineering, full-stack platforms, and AI retrieval systems.",
  alternates: {
    canonical: "/work",
  },
};

export default function WorkPage() {
  const projects = getProjects();

  return (
    <div className="pb-24">
      <PageHeader
        eyebrow="WORK // SELECTED PROJECTS"
        title="Selected Systems & Engineering Projects"
        description="Architectural case studies demonstrating decoupled Go backend services, Next.js user interfaces, and verifiable repository commits."
        badge={`${projects.length} Verified Systems`}
      />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <ProjectsSection projects={projects} hideHeader={true} />
      </div>

      <InterPageNav
        label="CONTINUE JOURNEY"
        nextRoute="/skills"
        nextTitle="Skills & Technical Evidence"
        description="Explore the technical capabilities, language proficiency, and architecture practices backed by verifiable repository evidence."
      />
    </div>
  );
}
