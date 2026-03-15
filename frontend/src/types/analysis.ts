export interface AnalyzeRequest {
  url: string;
}

export interface HeadingCount {
  h1: number;
  h2: number;
  h3: number;
  h4: number;
  h5: number;
  h6: number;
}

export interface LinkSummary {
  internal: number;
  external: number;
  inaccessible: number;
}

export interface AnalyzeResponse {
  url: string;
  htmlVersion: string;
  title: string;
  headings: HeadingCount;
  links: LinkSummary;
  hasLoginForm: boolean;
}

export interface ErrorResponse {
  statusCode: number;
  message: string;
}
