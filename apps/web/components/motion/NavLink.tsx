"use client";

import React from "react";
import Link, { LinkProps } from "next/link";
import { usePathname } from "next/navigation";
import { useOptionalNavigationProgress } from "@/context/NavigationProgressContext";

export function isModifiedClick(event: React.MouseEvent): boolean {
  return Boolean(
    event.metaKey ||
    event.ctrlKey ||
    event.shiftKey ||
    event.altKey ||
    event.button !== 0
  );
}

export interface NavLinkProps
  extends Omit<React.AnchorHTMLAttributes<HTMLAnchorElement>, "href">,
    Omit<LinkProps, "as"> {
  children: React.ReactNode;
}

export default function NavLink({
  href,
  onClick,
  target,
  download,
  children,
  ...props
}: NavLinkProps) {
  const pathname = usePathname();
  const navProgress = useOptionalNavigationProgress();

  const handleClick = (e: React.MouseEvent<HTMLAnchorElement>) => {
    if (onClick) {
      onClick(e);
    }
    if (e.defaultPrevented) return;

    // Preserve native browser actions for modifiers, new tabs, downloads
    if (isModifiedClick(e) || target === "_blank" || download) {
      return;
    }

    const hrefStr = typeof href === "string" ? href : href.pathname ?? "";
    if (
      !hrefStr ||
      hrefStr.startsWith("#") ||
      hrefStr.startsWith("http://") ||
      hrefStr.startsWith("https://") ||
      hrefStr.startsWith("//") ||
      hrefStr.startsWith("mailto:")
    ) {
      return;
    }

    // If navigating to the current exact route, do not trigger progress
    if (hrefStr === pathname) {
      return;
    }

    navProgress?.startProgress();
  };

  return (
    <Link
      href={href}
      onClick={handleClick}
      target={target}
      download={download}
      {...props}
    >
      {children}
    </Link>
  );
}
