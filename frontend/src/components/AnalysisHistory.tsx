import type { AnalyzeResponse } from '../types/analysis';

interface AnalysisHistoryProps {
  history: AnalyzeResponse[];
  onSelect: (result: AnalyzeResponse) => void;
}

export default function AnalysisHistory({ history, onSelect }: AnalysisHistoryProps) {
  if (history.length === 0) {
    return (
      <div style={{ color: '#6b7280', fontSize: '14px', textAlign: 'center', padding: '16px' }}>
        No analysis history yet.
      </div>
    );
  }

  return (
    <div>
      <h2 style={{ fontSize: '20px', marginBottom: '12px' }}>History</h2>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
        {history.map((item, index) => (
          <button
            key={index}
            onClick={() => onSelect(item)}
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              padding: '12px 16px',
              backgroundColor: '#f9fafb',
              border: '1px solid #e5e7eb',
              borderRadius: '6px',
              cursor: 'pointer',
              textAlign: 'left',
              fontSize: '14px',
              transition: 'background-color 0.15s',
            }}
            onMouseEnter={(e) => (e.currentTarget.style.backgroundColor = '#f3f4f6')}
            onMouseLeave={(e) => (e.currentTarget.style.backgroundColor = '#f9fafb')}
          >
            <div style={{ flex: 1, minWidth: 0 }}>
              <div style={{ fontWeight: 600, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                {item.title || 'No title'}
              </div>
              <div style={{ color: '#6b7280', fontSize: '13px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                {item.url}
              </div>
            </div>
            <div style={{ display: 'flex', gap: '8px', marginLeft: '12px', flexShrink: 0, fontSize: '12px' }}>
              <span style={{ color: '#6b7280' }}>{item.htmlVersion}</span>
              {item.hasLoginForm && (
                <span style={{ padding: '1px 6px', backgroundColor: '#dcfce7', color: '#16a34a', borderRadius: '4px', fontWeight: 600 }}>
                  Login
                </span>
              )}
            </div>
          </button>
        ))}
      </div>
    </div>
  );
}
