import type { AnalyzeResponse, ErrorResponse } from '../types/analysis';

export class ApiError extends Error {
  statusCode: number;

  constructor(statusCode: number, message: string) {
    super(message);
    this.statusCode = statusCode;
  }
}

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const errorBody: ErrorResponse = await response.json().catch(() => ({
      statusCode: response.status,
      message: response.statusText,
    }));
    throw new ApiError(errorBody.statusCode, errorBody.message);
  }
  return response.json();
}

export async function analyzeUrl(url: string): Promise<AnalyzeResponse> {
  const response = await fetch('/api/analyze', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ url }),
  });
  return handleResponse<AnalyzeResponse>(response);
}

export async function fetchAnalyses(): Promise<AnalyzeResponse[]> {
  const response = await fetch('/api/analyses');
  return handleResponse<AnalyzeResponse[]>(response);
}

export async function clearCache(): Promise<void> {
  const response = await fetch('/api/cache', { method: 'DELETE' });
  if (!response.ok) throw new ApiError(response.status, 'Failed to clear cache');
}
