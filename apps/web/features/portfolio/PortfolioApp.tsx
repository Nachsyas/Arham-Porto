"use client";

import { useState } from "react";
import Navbar from "@/components/Navbar";
import QuickReviewDrawer from "@/components/QuickReviewDrawer";
import HeroSection from "@/features/hero/HeroSection";
import ProjectsSection from "@/features/projects/ProjectsSection";
import SkillsSection from "@/features/skills/SkillsSection";
import ExperienceSection from "@/features/experience/ExperienceSection";
import EducationSection from "@/features/education/EducationSection";
import ContactSection from "@/features/contact/ContactSection";
import Footer from "@/components/Footer";
import type {
  Profile,
  Project,
  Skill,
  ExperienceItem,
  EducationItem,
} from "arham-porto-schema";

interface PortfolioAppProps {
  profile: Profile;
  projects: Project[];
  skills: Skill[];
  experiences: ExperienceItem[];
  education: EducationItem[];
}

export default function PortfolioApp({
  profile,
  projects,
  skills,
  experiences,
  education,
}: PortfolioAppProps) {
  const [quickReviewOpen, setQuickReviewOpen] = useState(false);

  return (
    <div className="min-h-screen bg-canvas text-themeText-body selection:bg-primary/20 selection:text-primary">
      {/* Global Navigation */}
      <Navbar
        onOpenQuickReview={() => setQuickReviewOpen(true)}
        hasExperience={experiences.length > 0}
      />

      {/* Main Reviewer Experience */}
      <main id="main-content" className="flex flex-col">
        {/* 1. Hero Identity with Seated Hologram Stage Placeholder */}
        <HeroSection
          profile={profile}
          onOpenQuickReview={() => setQuickReviewOpen(true)}
        />

        {/* 2. Selected Work Showcase with Domain Filter */}
        <ProjectsSection projects={projects} />

        {/* 3. Skill & Evidence Explorer (No percentages) */}
        <SkillsSection skills={skills} projects={projects} />

        {/* 4. Professional Experience (Rendered when verified records exist) */}
        {experiences.length > 0 && <ExperienceSection experiences={experiences} />}

        {/* 5. Academic Background (Rendered when verified records exist) */}
        {education.length > 0 && <EducationSection education={education} />}

        {/* 6. Contact Section (Verified Channels Only) */}
        <ContactSection profile={profile} />
      </main>

      {/* Footer */}
      <Footer profile={profile} />

      {/* 60-Second Quick Review Drawer */}
      <QuickReviewDrawer
        isOpen={quickReviewOpen}
        onClose={() => setQuickReviewOpen(false)}
        profile={profile}
        projects={projects}
      />
    </div>
  );
}
