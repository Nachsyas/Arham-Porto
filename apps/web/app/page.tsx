import { getProfile, getProjects } from "arham-porto-data";

export default function SmokePage() {
  // Server-side validated canonical data loading (zero database dependency)
  const profile = getProfile();
  const projects = getProjects();

  return (
    <div className="flex min-h-screen flex-col items-center justify-center p-8 bg-canvas text-themeText-primary">
      <div className="max-w-xl w-full border border-border bg-surface p-8 rounded-card text-center shadow-lg">
        <span className="inline-block px-3 py-1 mb-4 text-xs font-code tracking-wider text-primary bg-primary/10 border border-primary/20 rounded-full">
          PHASE 0 — FOUNDATION READY
        </span>
        <h1 className="text-3xl font-display font-bold text-themeText-primary mb-2">
          {profile.fullName}
        </h1>
        <p className="text-lg font-body text-primary mb-6">
          {profile.role}
        </p>

        <div className="pt-6 border-t border-border text-sm text-themeText-muted font-body space-y-2 text-left">
          <div className="flex justify-between py-1">
            <span>Project:</span>
            <span className="text-themeText-primary font-medium">{profile.projectName}</span>
          </div>
          <div className="flex justify-between py-1">
            <span>AI Feature:</span>
            <span className="text-themeText-primary font-medium">{profile.aiFeature}</span>
          </div>
          <div className="flex justify-between py-1">
            <span>Validated Canonical Projects:</span>
            <span className="text-primary font-code font-bold">{projects.length} loaded</span>
          </div>
          <div className="flex justify-between py-1">
            <span>Runtime Schema:</span>
            <span className="text-status-success font-medium">Zod Verified</span>
          </div>
        </div>

        <p className="mt-8 text-xs text-themeText-mutedSoft italic font-body">
          Foundation verification smoke test. Full portfolio interface scheduled for Phase 1.
        </p>
      </div>
    </div>
  );
}
