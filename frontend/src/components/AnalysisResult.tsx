import type { AnalyzeResponse } from '../types/analysis';
import HeadingSummary from './HeadingSummary';
import LinkSummary from './LinkSummary';
import './AnalysisResult.css';

interface AnalysisResultProps {
  result: AnalyzeResponse;
}

export default function AnalysisResult({ result }: AnalysisResultProps) {
  return (
    <section className="result-card" aria-label="Analysis result">
      <div>
        <h2>Analysis Result</h2>
        <dl className="result-card__details">
          <div><dt>URL:</dt> <dd>{result.url}</dd></div>
          <div><dt>HTML Version:</dt> <dd>{result.htmlVersion}</dd></div>
          <div><dt>Title:</dt> <dd>{result.title || <em>No title</em>}</dd></div>
          <div>
            <dt>Login Form:</dt>{' '}
            <dd>
              <span className={`result-card__badge result-card__badge--${result.hasLoginForm ? 'yes' : 'no'}`}>
                {result.hasLoginForm ? 'Yes' : 'No'}
              </span>
            </dd>
          </div>
        </dl>
      </div>

      <HeadingSummary headings={result.headings} />
      <LinkSummary links={result.links} />
    </section>
  );
}
