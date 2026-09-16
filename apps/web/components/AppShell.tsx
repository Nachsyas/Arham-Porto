"use client";

import React, { useRef } from "react";
import dynamic from "next/dynamic";
import Navbar from "@/components/Navbar";
import Footer from "@/components/Footer";
import QuickReviewDrawer from "@/components/QuickReviewDrawer";
import RouteProgress from "@/components/motion/RouteProgress";
import { AskArhamLauncher } from "@/features/ask-arham";
import { PortfolioUIProvider, usePortfolioUI } from "@/context/PortfolioUIContext";
import { NavigationProgressProvider } from "@/context/NavigationProgressContext";
import type { Profile, Project } from "arham-porto-schema";

const AskArhamPanel = dynamic(
  () => import("@/features/ask-arham/AskArhamPanel"),
  { ssr: false }
);

interface AppShellProps {
  children: React.ReactNode;
  profile: Profile;
  projects: Project[];
}

function AppShellInner({
  children,
  profile,
  projects,
}: {
  children: React.ReactNode;
  profile: Profile;
  projects: Project[];
}) {
  const {
    isAskArhamOpen,
    openAskArham,
    closeAskArham,
    isQuickReviewOpen,
    closeQuickReview,
  } = usePortfolioUI();

  const launcherRef = useRef<HTMLButtonElement>(null);

  return (
    <div className="min-h-screen flex flex-col bg-canvas text-themeText-body selection:bg-primary/20 selection:text-primary">
      {/* Top ~2px route transition indicator */}
      <RouteProgress />

      {/* Persistent Route-Aware Header */}
      <Navbar />

      {/* Dynamic Route Content */}
      <div className="flex-1 flex flex-col">{children}</div>

      {/* Global Footer */}
      <Footer profile={profile} />

      {/* Global 60-Second Quick Review Drawer */}
      <QuickReviewDrawer
        isOpen={isQuickReviewOpen}
        onClose={closeQuickReview}
        profile={profile}
        projects={projects}
      />

      {/* Persistent Ask Arham Launcher */}
      <AskArhamLauncher
        ref={launcherRef}
        isOpen={isAskArhamOpen}
        onClick={openAskArham}
      />

      {/* Persistent Ask Arham Panel */}
      <AskArhamPanel
        isOpen={isAskArhamOpen}
        onClose={closeAskArham}
        returnFocusRef={launcherRef}
      />
    </div>
  );
}

export default function AppShell({ children, profile, projects }: AppShellProps) {
  return (
    <PortfolioUIProvider>
      <NavigationProgressProvider>
        <AppShellInner profile={profile} projects={projects}>
          {children}
        </AppShellInner>
      </NavigationProgressProvider>
    </PortfolioUIProvider>
  );
}
