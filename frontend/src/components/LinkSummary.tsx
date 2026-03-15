import { LinkSummary as LinkSummaryType } from '../types/analysis';

interface LinkSummaryProps {
  links: LinkSummaryType;
}

export default function LinkSummary({ links }: LinkSummaryProps) {
  const entries = [
    { label: 'Internal', count: links.internal, color: '#16a34a' },
    { label: 'External', count: links.external, color: '#2563eb' },
    { label: 'Inaccessible', count: links.inaccessible, color: '#dc2626' },
  ];

  return (
    <div>
      <h3 style={{ marginBottom: '8px' }}>Links</h3>
      <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
        {entries.map(({ label, count, color }) => (
          <div
            key={label}
            style={{
              padding: '8px 16px',
              backgroundColor: '#f3f4f6',
              borderRadius: '6px',
              textAlign: 'center',
              minWidth: '80px',
            }}
          >
            <div style={{ fontWeight: 600, fontSize: '18px', color }}>{count}</div>
            <div style={{ fontSize: '13px', color: '#6b7280' }}>{label}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
