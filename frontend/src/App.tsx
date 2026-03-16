import { useState, useEffect } from 'react';
import './App.css';
import AnalyzeForm from './components/AnalyzeForm';
import AnalysisResult from './components/AnalysisResult';
import AnalysisHistory from './components/AnalysisHistory';
import ErrorMessage from './components/ErrorMessage';
import { analyzeUrl, fetchAnalyses, ApiError } from './api/analyzeApi';
import type { AnalyzeResponse } from './types/analysis';

function App() {
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState<AnalyzeResponse | null>(null);
  const [error, setError] = useState<{ statusCode: number; message: string } | null>(null);
  const [history, setHistory] = useState<AnalyzeResponse[]>([]);

  useEffect(() => {
    loadHistory();
  }, []);

  async function loadHistory() {
    try {
      const data = await fetchAnalyses();
      setHistory(data);
    } catch {
      // History not available (e.g., no MySQL configured) — silently ignore
    }
  }

  async function handleAnalyze(url: string) {
    setIsLoading(true);
    setResult(null);
    setError(null);

    try {
      const data = await analyzeUrl(url);
      setResult(data);
      loadHistory(); // refresh history after new analysis
    } catch (err) {
      if (err instanceof ApiError) {
        setError({ statusCode: err.statusCode, message: err.message });
      } else {
        setError({ statusCode: 0, message: 'An unexpected error occurred. Is the backend running?' });
      }
    } finally {
      setIsLoading(false);
    }
  }

  function handleSelectHistory(item: AnalyzeResponse) {
    setResult(item);
    setError(null);
  }

  return (
    <div className="app">
      <header className="app-header">
        <h1>Web Page Analyzer</h1>
        <p>Enter a URL to analyze its structure, headings, links, and more.</p>
      </header>

      <main className="app-main">
        <AnalyzeForm onSubmit={handleAnalyze} isLoading={isLoading} />

        <div aria-live="polite">
          {error && <ErrorMessage statusCode={error.statusCode} message={error.message} />}
          {result && <AnalysisResult result={result} />}
        </div>

        <AnalysisHistory history={history} onSelect={handleSelectHistory} />
      </main>
    </div>
  );
}

export default App;
