import type { HeadingCount } from '../types/analysis';
import './HeadingSummary.css';

interface HeadingSummaryProps {
  headings: HeadingCount;
}

const HEADING_LEVELS = ['h1', 'h2', 'h3', 'h4', 'h5', 'h6'] as const;

export default function HeadingSummary({ headings }: HeadingSummaryProps) {
  return (
    <div>
      <h3>Headings</h3>
      <div className="heading-summary__grid">
        {HEADING_LEVELS.map((level) => (
          <div key={level} className="heading-summary__item">
            <div className="heading-summary__count">{headings[level]}</div>
            <div className="heading-summary__label">{level.toUpperCase()}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
