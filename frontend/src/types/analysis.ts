export type AnalyzeRequest = {
  url: string;
};

export type HeadingCount = {
  h1: number;
  h2: number;
  h3: number;
  h4: number;
  h5: number;
  h6: number;
};

export type LinkSummary = {
  internal: number;
  external: number;
  inaccessible: number;
};

export type AnalyzeResponse = {
  url: string;
  htmlVersion: string;
  title: string;
  headings: HeadingCount;
  links: LinkSummary;
  hasLoginForm: boolean;
};

export type ErrorResponse = {
  statusCode: number;
  message: string;
};
