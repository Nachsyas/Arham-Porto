"use client";

import React, {
  createContext,
  useContext,
  useState,
  useCallback,
  useMemo,
  useEffect,
} from "react";
import { usePathname } from "next/navigation";

interface NavigationProgressContextValue {
  isNavigating: boolean;
  startProgress: () => void;
  completeProgress: () => void;
}

const NavigationProgressContext = createContext<
  NavigationProgressContextValue | undefined
>(undefined);

export function NavigationProgressProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [isNavigating, setIsNavigating] = useState(false);
  const pathname = usePathname();

  const startProgress = useCallback(() => {
    setIsNavigating(true);
  }, []);

  const completeProgress = useCallback(() => {
    setIsNavigating(false);
  }, []);

  // When the route commits and pathname changes, complete navigation
  useEffect(() => {
    if (isNavigating) {
      setIsNavigating(false);
    }
  }, [pathname]);

  // Safety fallback: if navigation stalls or is aborted, reset after 4s
  useEffect(() => {
    if (!isNavigating) return;
    const timeout = setTimeout(() => {
      setIsNavigating(false);
    }, 4000);
    return () => clearTimeout(timeout);
  }, [isNavigating]);

  const value = useMemo(
    () => ({
      isNavigating,
      startProgress,
      completeProgress,
    }),
    [isNavigating, startProgress, completeProgress]
  );

  return (
    <NavigationProgressContext.Provider value={value}>
      {children}
    </NavigationProgressContext.Provider>
  );
}

export function useNavigationProgress(): NavigationProgressContextValue {
  const context = useContext(NavigationProgressContext);
  if (!context) {
    throw new Error(
      "useNavigationProgress must be used within a NavigationProgressProvider"
    );
  }
  return context;
}

export function useOptionalNavigationProgress(): NavigationProgressContextValue | null {
  return useContext(NavigationProgressContext) ?? null;
}
