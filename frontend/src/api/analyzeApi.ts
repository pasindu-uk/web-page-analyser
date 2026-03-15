import type { AnalyzeResponse, ErrorResponse } from '../types/analysis';

export class ApiError extends Error {
  statusCode: number;

  constructor(statusCode: number, message: string) {
    super(message);
    this.statusCode = statusCode;
  }
}

export async function analyzeUrl(url: string): Promise<AnalyzeResponse> {
  const response = await fetch('/api/analyze', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url }),
  });

  if (!response.ok) {
    const errorBody: ErrorResponse = await response.json().catch(() => ({
      statusCode: response.status,
      message: response.statusText,
    }));
    throw new ApiError(errorBody.statusCode, errorBody.message);
  }

  return response.json();
}
