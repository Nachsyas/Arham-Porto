import type { MetadataRoute } from "next";

export default function sitemap(): MetadataRoute.Sitemap {
  const siteUrl = process.env.NEXT_PUBLIC_SITE_URL || "https://arham-porto.vercel.app";

  const routes = [
    "",
    "/projects/edutrace",
    "/projects/gdgoc-ecommerce",
    "/projects/maritime-ai-dashboard",
    "/projects/smart-kitchen",
  ];

  return routes.map((route) => ({
    url: `${siteUrl}${route}`,
    lastModified: new Date(),
    changeFrequency: "weekly" as const,
    priority: route === "" ? 1.0 : 0.8,
  }));
}
