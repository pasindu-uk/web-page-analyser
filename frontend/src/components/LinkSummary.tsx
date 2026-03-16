import type { LinkSummary as LinkSummaryType } from '../types/analysis';
import './LinkSummary.css';

interface LinkSummaryProps {
  links: LinkSummaryType;
}

const LINK_TYPES = [
  { key: 'internal', label: 'Internal' },
  { key: 'external', label: 'External' },
  { key: 'inaccessible', label: 'Inaccessible' },
] as const;

export default function LinkSummary({ links }: LinkSummaryProps) {
  return (
    <div>
      <h3>Links</h3>
      <div className="link-summary__grid">
        {LINK_TYPES.map(({ key, label }) => (
          <div key={key} className="link-summary__item">
            <div className={`link-summary__count link-summary__count--${key}`}>
              {links[key]}
            </div>
            <div className="link-summary__label">{label}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
