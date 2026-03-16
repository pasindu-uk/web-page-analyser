import { useState, type FormEvent } from 'react';
import './AnalyzeForm.css';

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
    <form onSubmit={handleSubmit} className="analyze-form" aria-label="URL analysis form">
      <div className="analyze-form__row">
        <label htmlFor="url-input" className="sr-only">URL to analyze</label>
        <input
          id="url-input"
          type="url"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
          placeholder="https://example.com"
          disabled={isLoading}
          className="analyze-form__input"
          aria-describedby={error ? 'url-error' : undefined}
        />
        <button type="submit" disabled={isLoading} className="analyze-form__button">
          {isLoading ? 'Analyzing...' : 'Analyze'}
        </button>
      </div>
      {error && <p id="url-error" className="analyze-form__error" role="alert">{error}</p>}
    </form>
  );
}
