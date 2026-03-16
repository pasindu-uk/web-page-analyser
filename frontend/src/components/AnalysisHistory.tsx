import type { AnalyzeResponse } from '../types/analysis';
import './AnalysisHistory.css';

interface AnalysisHistoryProps {
  history: AnalyzeResponse[];
  onSelect: (result: AnalyzeResponse) => void;
}

export default function AnalysisHistory({ history, onSelect }: AnalysisHistoryProps) {
  if (history.length === 0) {
    return <div className="history__empty">No analysis history yet.</div>;
  }

  return (
    <section aria-label="Analysis history">
      <h2>History</h2>
      <div className="history__list">
        {history.map((item) => (
          <button
            key={item.url}
            onClick={() => onSelect(item)}
            className="history__item"
          >
            <div className="history__item-info">
              <div className="history__item-title">{item.title || 'No title'}</div>
              <div className="history__item-url">{item.url}</div>
            </div>
            <div className="history__item-meta">
              <span className="history__item-version">{item.htmlVersion}</span>
              {item.hasLoginForm && (
                <span className="history__login-badge">Login</span>
              )}
            </div>
          </button>
        ))}
      </div>
    </section>
  );
}
