import type { GithubProjectItem, GithubProjectsResponse } from "./github-projects.types";

/**
 * Fetches synchronized GitHub repositories from the backend API.
 * Uses NEXT_PUBLIC_API_BASE_URL with safe development fallbacks.
 */
export async function getGithubProjects(
  page = 1,
  limit = 20,
  signal?: AbortSignal
): Promise<GithubProjectItem[]> {
  const apiBase =
    process.env.NEXT_PUBLIC_API_BASE_URL ||
    (typeof window !== "undefined" ? "http://localhost:8080" : "http://localhost:8080");

  const url = `${apiBase.replace(/\/+$/, "")}/api/v1/github/projects?page=${page}&limit=${limit}`;

  const response = await fetch(url, {
    method: "GET",
    headers: {
      Accept: "application/json",
    },
    signal,
  });

  if (!response.ok) {
    let errMsg = `GitHub projects API failed with status ${response.status}`;
    try {
      const errJson = await response.json();
      if (errJson?.error?.message) {
        errMsg = errJson.error.message;
      }
    } catch {
      // Fallback to HTTP status message
    }
    throw new Error(errMsg);
  }

  const json: GithubProjectsResponse = await response.json();
  return json.data || [];
}
