import type { Metadata } from "next";
import { Space_Grotesk, JetBrains_Mono, Inter } from "next/font/google";
import "./globals.css";

const spaceGrotesk = Space_Grotesk({
  subsets: ["latin"],
  variable: "--font-space-grotesk",
  display: "swap",
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-jetbrains-mono",
  display: "swap",
});

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-geist-sans",
  display: "swap",
});

const siteUrl = process.env.NEXT_PUBLIC_SITE_URL || "https://arham-porto.vercel.app";

export const metadata: Metadata = {
  metadataBase: new URL(siteUrl),
  title: {
    default: "Nachsyas Arham Mumtaz Nashohi | Software Engineer",
    template: "%s | Arham Porto",
  },
  description:
    "Personal engineering portfolio and technical case-study platform for Nachsyas Arham Mumtaz Nashohi. Evidence-based architectural reviews and Ask Arham AI copilot.",
  alternates: {
    canonical: "/",
  },
  openGraph: {
    title: "Nachsyas Arham Mumtaz Nashohi | Software Engineer",
    description:
      "Personal engineering portfolio and technical case-study platform for Nachsyas Arham Mumtaz Nashohi.",
    url: siteUrl,
    siteName: "Arham Porto",
    locale: "en_US",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Nachsyas Arham Mumtaz Nashohi | Software Engineer",
    description:
      "Personal engineering portfolio and technical case-study platform for Nachsyas Arham Mumtaz Nashohi.",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={`${spaceGrotesk.variable} ${jetbrainsMono.variable} ${inter.variable}`}>
      <body className="bg-canvas text-themeText-body min-h-screen antialiased selection:bg-primary/20 selection:text-primary">
        <a
          href="#main-content"
          className="sr-only focus:not-sr-only focus:fixed focus:top-4 focus:left-4 focus:z-50 focus:px-4 focus:py-2 focus:bg-primary focus:text-white focus:rounded-md focus:font-medium shadow-sm"
        >
          Skip to content
        </a>
        <main id="main-content">{children}</main>
      </body>
    </html>
  );
}
