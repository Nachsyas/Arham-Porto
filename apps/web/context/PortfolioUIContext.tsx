"use client";

import React, { createContext, useContext, useState, useCallback, useMemo } from "react";

interface PortfolioUIContextValue {
  isAskArhamOpen: boolean;
  openAskArham: () => void;
  closeAskArham: () => void;
  isQuickReviewOpen: boolean;
  openQuickReview: () => void;
  closeQuickReview: () => void;
}

const PortfolioUIContext = createContext<PortfolioUIContextValue | undefined>(undefined);

export function PortfolioUIProvider({ children }: { children: React.ReactNode }) {
  const [isAskArhamOpen, setIsAskArhamOpen] = useState(false);
  const [isQuickReviewOpen, setIsQuickReviewOpen] = useState(false);

  const openAskArham = useCallback(() => setIsAskArhamOpen(true), []);
  const closeAskArham = useCallback(() => setIsAskArhamOpen(false), []);

  const openQuickReview = useCallback(() => setIsQuickReviewOpen(true), []);
  const closeQuickReview = useCallback(() => setIsQuickReviewOpen(false), []);

  const value = useMemo(
    () => ({
      isAskArhamOpen,
      openAskArham,
      closeAskArham,
      isQuickReviewOpen,
      openQuickReview,
      closeQuickReview,
    }),
    [isAskArhamOpen, openAskArham, closeAskArham, isQuickReviewOpen, openQuickReview, closeQuickReview]
  );

  return (
    <PortfolioUIContext.Provider value={value}>
      {children}
    </PortfolioUIContext.Provider>
  );
}

export function usePortfolioUI(): PortfolioUIContextValue {
  const context = useContext(PortfolioUIContext);
  if (!context) {
    throw new Error("usePortfolioUI must be used within a PortfolioUIProvider");
  }
  return context;
}

export function useOptionalPortfolioUI(): PortfolioUIContextValue | null {
  return useContext(PortfolioUIContext) ?? null;
}
