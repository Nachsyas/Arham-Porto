export interface GithubProjectItem {
  id: string;
  name: string;
  full_name: string;
  description: string;
  html_url: string;
  homepage: string;
  language: string;
  stars: number;
  forks: number;
  topics: string[];
  synced_at: string;
}

export interface GithubProjectsResponse {
  status: string;
  data: GithubProjectItem[];
}
