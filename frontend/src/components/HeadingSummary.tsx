import { HeadingCount } from '../types/analysis';

interface HeadingSummaryProps {
  headings: HeadingCount;
}

export default function HeadingSummary({ headings }: HeadingSummaryProps) {
  const entries = [
    { label: 'H1', count: headings.h1 },
    { label: 'H2', count: headings.h2 },
    { label: 'H3', count: headings.h3 },
    { label: 'H4', count: headings.h4 },
    { label: 'H5', count: headings.h5 },
    { label: 'H6', count: headings.h6 },
  ];

  return (
    <div>
      <h3 style={{ marginBottom: '8px' }}>Headings</h3>
      <div style={{ display: 'flex', gap: '12px', flexWrap: 'wrap' }}>
        {entries.map(({ label, count }) => (
          <div
            key={label}
            style={{
              padding: '8px 16px',
              backgroundColor: '#f3f4f6',
              borderRadius: '6px',
              textAlign: 'center',
              minWidth: '60px',
            }}
          >
            <div style={{ fontWeight: 600, fontSize: '18px' }}>{count}</div>
            <div style={{ fontSize: '13px', color: '#6b7280' }}>{label}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
