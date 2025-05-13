import { ReactNode, useState } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import {
  HomeIcon,
  BuildingOfficeIcon,
  UsersIcon,
  DocumentTextIcon,
  Cog6ToothIcon,
  Bars3Icon,
  XMarkIcon,
} from '@heroicons/react/24/outline';
import Navbar from './Navbar';

interface DashboardLayoutProps {
  children: ReactNode;
}

const navigation = [
  { name: 'Dashboard', href: '/dashboard', icon: HomeIcon },
  { name: 'Empresas', href: '/empresas', icon: BuildingOfficeIcon },
  { name: 'Usuarios', href: '/usuarios', icon: UsersIcon },
  { name: 'Documentos', href: '/documentos', icon: DocumentTextIcon },
  { name: 'Configuración', href: '/configuracion', icon: Cog6ToothIcon },
];

export default function DashboardLayout({ children }: DashboardLayoutProps) {
  const router = useRouter();
  const pathname = usePathname();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      {/* Navbar superior */}
      <nav className="fixed top-0 left-0 right-0 h-16 bg-white shadow z-30 flex items-center px-4 justify-between">
        <div className="flex items-center gap-2">
          {/* Botón hamburguesa en móvil */}
          <button
            className="lg:hidden p-2 rounded-md hover:bg-gray-100 focus:outline-none"
            onClick={() => setSidebarOpen(true)}
            aria-label="Abrir menú"
          >
            <Bars3Icon className="h-6 w-6 text-gray-700" />
          </button>
          <span className="font-bold text-xl text-indigo-700 tracking-tight">ERP Dashboard</span>
        </div>
        {/* Usuario/logout aquí si lo deseas */}
      </nav>

      {/* Sidebar oscura, colapsable */}
      {/* Overlay para móvil */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-40 z-40 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}
      <aside
        className={`fixed top-0 left-0 h-full w-64 bg-gray-900 text-white z-50 transform transition-transform duration-200 ease-in-out
          ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
          lg:translate-x-0 lg:static lg:block`}
      >
        <div className="flex items-center justify-between h-16 px-4 border-b border-gray-800">
          <span className="font-bold text-lg text-indigo-400">Menú</span>
          <button
            className="lg:hidden p-2 rounded-md hover:bg-gray-800 focus:outline-none"
            onClick={() => setSidebarOpen(false)}
            aria-label="Cerrar menú"
          >
            <XMarkIcon className="h-6 w-6 text-gray-300" />
          </button>
        </div>
        <nav className="mt-4 space-y-1 px-2">
          {navigation.map((item) => {
            const isActive = pathname.startsWith(item.href);
            return (
              <Link
                key={item.name}
                href={item.href}
                className={`group flex items-center px-3 py-2 text-base font-medium rounded-md transition-colors
                  ${isActive ? 'bg-indigo-700 text-white' : 'text-gray-300 hover:bg-gray-800 hover:text-white'}`}
                onClick={() => setSidebarOpen(false)}
              >
                <item.icon
                  className={`mr-3 h-5 w-5 flex-shrink-0 ${isActive ? 'text-white' : 'text-indigo-300 group-hover:text-white'}`}
                  aria-hidden="true"
                />
                {item.name}
              </Link>
            );
          })}
        </nav>
      </aside>

      {/* Main content */}
      <div className="flex-1 flex flex-col lg:pl-64 pt-16 transition-all">
        <main className="flex-1">
          <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-6">
            {children}
          </div>
        </main>
      </div>
    </div>
  );
} 