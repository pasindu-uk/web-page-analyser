import { useState } from 'react';
import './App.css';
import AnalyzeForm from './components/AnalyzeForm';
import AnalysisResult from './components/AnalysisResult';
import ErrorMessage from './components/ErrorMessage';
import { analyzeUrl, ApiError } from './api/analyzeApi';
import { AnalyzeResponse } from './types/analysis';

function App() {
  const [isLoading, setIsLoading] = useState(false);
  const [result, setResult] = useState<AnalyzeResponse | null>(null);
  const [error, setError] = useState<{ statusCode: number; message: string } | null>(null);

  async function handleAnalyze(url: string) {
    setIsLoading(true);
    setResult(null);
    setError(null);

    try {
      const data = await analyzeUrl(url);
      setResult(data);
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

  return (
    <div className="app">
      <header className="app-header">
        <h1>Web Page Analyzer</h1>
        <p>Enter a URL to analyze its structure, headings, links, and more.</p>
      </header>

      <main className="app-main">
        <AnalyzeForm onSubmit={handleAnalyze} isLoading={isLoading} />

        {error && <ErrorMessage statusCode={error.statusCode} message={error.message} />}
        {result && <AnalysisResult result={result} />}
      </main>
    </div>
  );
}

export default App;
