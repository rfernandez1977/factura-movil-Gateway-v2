import { Inter } from 'next/font/google';
import './globals.css';
import ClientProviders from './ClientProviders';

const inter = Inter({ subsets: ['latin'] });

export const metadata = {
  title: 'Dashboard ERP',
  description: 'Sistema de gestión empresarial',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="es">
      <body className={inter.className}>
        <ClientProviders>
          {children}
        </ClientProviders>
      </body>
    </html>
  );
}
