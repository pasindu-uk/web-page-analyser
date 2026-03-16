import './ErrorMessage.css';

interface ErrorMessageProps {
  statusCode: number;
  message: string;
}

export default function ErrorMessage({ statusCode, message }: ErrorMessageProps) {
  return (
    <div className="error-message" role="alert">
      <strong>Error {statusCode}:</strong> {message}
    </div>
  );
}
