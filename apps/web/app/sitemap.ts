import type { MetadataRoute } from "next";
import { getProjects } from "arham-porto-data";

export default function sitemap(): MetadataRoute.Sitemap {
  const siteUrl = process.env.NEXT_PUBLIC_SITE_URL || "https://arham-porto.vercel.app";

  const staticRoutes: { path: string; priority: number }[] = [
    { path: "", priority: 1.0 },
    { path: "/work", priority: 0.9 },
    { path: "/skills", priority: 0.9 },
    { path: "/journey", priority: 0.8 },
    { path: "/contact", priority: 0.8 },
  ];

  const projects = getProjects();
  const projectRoutes = projects.map((p) => ({
    path: `/projects/${p.slug}`,
    priority: 0.7,
  }));

  const allRoutes = [...staticRoutes, ...projectRoutes];

  return allRoutes.map(({ path, priority }) => ({
    url: `${siteUrl}${path}`,
    lastModified: new Date(),
    changeFrequency: "weekly" as const,
    priority,
  }));
}

