import { AnalyzeResponse } from '../types/analysis';
import HeadingSummary from './HeadingSummary';
import LinkSummary from './LinkSummary';

interface AnalysisResultProps {
  result: AnalyzeResponse;
}

export default function AnalysisResult({ result }: AnalysisResultProps) {
  return (
    <div
      style={{
        border: '1px solid #e5e7eb',
        borderRadius: '8px',
        padding: '24px',
        display: 'flex',
        flexDirection: 'column',
        gap: '20px',
      }}
    >
      <div>
        <h2 style={{ marginBottom: '12px' }}>Analysis Result</h2>
        <div style={{ display: 'flex', flexDirection: 'column', gap: '6px', fontSize: '15px' }}>
          <div><strong>URL:</strong> {result.url}</div>
          <div><strong>HTML Version:</strong> {result.htmlVersion}</div>
          <div><strong>Title:</strong> {result.title || <em>No title</em>}</div>
          <div>
            <strong>Login Form:</strong>{' '}
            <span
              style={{
                padding: '2px 8px',
                borderRadius: '4px',
                fontSize: '13px',
                fontWeight: 600,
                backgroundColor: result.hasLoginForm ? '#dcfce7' : '#f3f4f6',
                color: result.hasLoginForm ? '#16a34a' : '#6b7280',
              }}
            >
              {result.hasLoginForm ? 'Yes' : 'No'}
            </span>
          </div>
        </div>
      </div>

      <HeadingSummary headings={result.headings} />
      <LinkSummary links={result.links} />
    </div>
  );
}
