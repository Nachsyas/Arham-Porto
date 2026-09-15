import { getSkills, getProjects, getEvidence } from "arham-porto-data";
import PageHeader from "@/components/PageHeader";
import SkillsExplorer from "@/features/skills/SkillsExplorer";
import InterPageNav from "@/components/InterPageNav";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Skills & Technical Evidence",
  description:
    "Explore verified engineering competencies of Nachsyas Arham Mumtaz Nashohi. Grounded evidence chain connecting skills directly to code artifacts.",
  alternates: {
    canonical: "/skills",
  },
};

export default function SkillsPage() {
  const skills = getSkills();
  const projects = getProjects();
  const evidence = getEvidence();

  return (
    <div className="pb-24">
      <PageHeader
        eyebrow="SKILLS // VERIFIED CAPABILITIES"
        title="Capabilities Backed by Verified Evidence"
        description="Structured engineering capabilities mapped to repository evidence and architecture decisions. Evaluated on demonstrated evidence rather than arbitrary percentage ratings."
        badge={`${skills.length} Validated Disciplines`}
      />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <SkillsExplorer
          skills={skills}
          projects={projects}
          evidence={evidence}
        />
      </div>

      <InterPageNav
        label="CONTINUE JOURNEY"
        nextRoute="/journey"
        nextTitle="Geographic & Academic Evolution"
        description="Discover the formative stops, secondary schooling in Central Java, and university computer science education in East Java."
      />
    </div>
  );
}
