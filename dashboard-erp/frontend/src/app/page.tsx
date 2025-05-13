export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <h1 className="text-4xl font-bold text-center mb-6">
        Dashboard ERP - Página de Inicio
      </h1>
      <p className="text-center mb-6">
        Sistema de gestión empresarial con Next.js, TypeScript y Tailwind CSS
      </p>
      <div className="flex gap-4">
        <a
          href="/login"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          Iniciar Sesión
        </a>
        <a
          href="/dashboard"
          className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded"
        >
          Dashboard
        </a>
      </div>
    </main>
  );
}
