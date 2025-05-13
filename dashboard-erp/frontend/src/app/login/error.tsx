"use client";

export default function Error({ error, reset }: { error: Error; reset: () => void }) {
  return (
    <div style={{ padding: 32 }}>
      <h2>Ocurrió un error en Login</h2>
      <p>{error.message}</p>
      <button onClick={() => reset()}>Reintentar</button>
    </div>
  );
} 