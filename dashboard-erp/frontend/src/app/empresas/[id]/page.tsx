'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import DashboardLayout from '@/components/layout/DashboardLayout';
import EmpresaForm from '@/components/empresas/EmpresaForm';
import { Empresa } from '@/types';
import { empresasApi } from '@/services/api';

interface EmpresaPageProps {
  params: {
    id: string;
  };
}

export default function EmpresaPage({ params }: EmpresaPageProps) {
  const router = useRouter();
  const [empresa, setEmpresa] = useState<Empresa | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const cargarEmpresa = async () => {
      try {
        setLoading(true);
        const response = await empresasApi.obtener(params.id);
        if (response.data) {
          setEmpresa(response.data);
        }
      } catch (err) {
        setError('Error al cargar la empresa');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    if (params.id !== 'nueva') {
      cargarEmpresa();
    } else {
      setLoading(false);
    }
  }, [params.id]);

  const handleSubmit = async (data: Omit<Empresa, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      if (params.id === 'nueva') {
        await empresasApi.crear(data);
      } else {
        await empresasApi.actualizar(params.id, data);
      }
      router.push('/empresas');
    } catch (err) {
      setError('Error al guardar la empresa');
      console.error(err);
    }
  };

  const handleCancel = () => {
    router.push('/empresas');
  };

  if (loading) {
    return (
      <DashboardLayout>
        <div className="py-6">
          <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <p className="text-center text-gray-500">Cargando...</p>
          </div>
        </div>
      </DashboardLayout>
    );
  }

  if (error) {
    return (
      <DashboardLayout>
        <div className="py-6">
          <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <p className="text-center text-red-500">{error}</p>
          </div>
        </div>
      </DashboardLayout>
    );
  }

  return (
    <DashboardLayout>
      <div className="py-6">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <h1 className="text-2xl font-semibold text-gray-900">
            {params.id === 'nueva' ? 'Nueva Empresa' : 'Editar Empresa'}
          </h1>
          <div className="mt-6">
            <EmpresaForm
              empresa={empresa || undefined}
              onSubmit={handleSubmit}
              onCancel={handleCancel}
            />
          </div>
        </div>
      </div>
    </DashboardLayout>
  );
} 