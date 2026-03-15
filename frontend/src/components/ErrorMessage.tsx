interface ErrorMessageProps {
  statusCode: number;
  message: string;
}

export default function ErrorMessage({ statusCode, message }: ErrorMessageProps) {
  return (
    <div
      style={{
        padding: '16px',
        backgroundColor: '#fef2f2',
        border: '1px solid #fecaca',
        borderRadius: '8px',
        color: '#991b1b',
      }}
    >
      <strong>Error {statusCode}:</strong> {message}
    </div>
  );
}
