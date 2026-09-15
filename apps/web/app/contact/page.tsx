import { getProfile } from "arham-porto-data";
import PageHeader from "@/components/PageHeader";
import ContactView from "@/features/contact/ContactView";
import InterPageNav from "@/components/InterPageNav";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Contact & Verified Channels",
  description:
    "Get in touch with Nachsyas Arham Mumtaz Nashohi. Access verified GitHub repositories, ask the grounded AI reviewer, or launch the 60-second quick review.",
  alternates: {
    canonical: "/contact",
  },
};

export default function ContactPage() {
  const profile = getProfile();

  return (
    <div className="pb-24">
      <PageHeader
        eyebrow="CONTACT // VERIFIED CHANNELS"
        title="Get in Touch & Technical Inquiries"
        description="Connect through verified channels, inspect public code repositories, or query the grounded AI reviewer for architectural citations."
        badge="Direct Channels"
      />

      <ContactView profile={profile} />

      <InterPageNav
        label="PORTFOLIO OVERVIEW"
        nextRoute="/"
        nextTitle="Return to Overview & Recruiter Snapshot"
        description="Jump back to the cinematic hero introduction, core metrics, and architecture summaries."
      />
    </div>
  );
}
