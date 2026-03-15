import { useState, FormEvent } from 'react';

interface AnalyzeFormProps {
  onSubmit: (url: string) => void;
  isLoading: boolean;
}

export default function AnalyzeForm({ onSubmit, isLoading }: AnalyzeFormProps) {
  const [url, setUrl] = useState('');
  const [error, setError] = useState('');

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');

    if (!url.trim()) {
      setError('URL is required');
      return;
    }

    if (!url.startsWith('http://') && !url.startsWith('https://')) {
      setError('URL must start with http:// or https://');
      return;
    }

    onSubmit(url.trim());
  }

  return (
    <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: '8px' }}>
      <div style={{ display: 'flex', gap: '8px' }}>
        <input
          type="text"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://example.com"
          disabled={isLoading}
          style={{
            flex: 1,
            padding: '10px 14px',
            fontSize: '16px',
            border: '1px solid #d1d5db',
            borderRadius: '6px',
            outline: 'none',
          }}
        />
        <button
          type="submit"
          disabled={isLoading}
          style={{
            padding: '10px 24px',
            fontSize: '16px',
            backgroundColor: isLoading ? '#9ca3af' : '#2563eb',
            color: '#fff',
            border: 'none',
            borderRadius: '6px',
            cursor: isLoading ? 'not-allowed' : 'pointer',
          }}
        >
          {isLoading ? 'Analyzing...' : 'Analyze'}
        </button>
      </div>
      {error && <p style={{ color: '#dc2626', fontSize: '14px' }}>{error}</p>}
    </form>
  );
}
