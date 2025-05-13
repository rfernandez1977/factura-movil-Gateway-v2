import axios from 'axios';
import { ApiResponse, Empresa, Usuario } from '@/types';

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor para agregar el token de autenticación
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Interceptor para manejar errores
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Manejar error de autenticación
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export const empresasApi = {
  listar: async (): Promise<ApiResponse<Empresa[]>> => {
    const response = await api.get('/empresas');
    return response.data;
  },
  obtener: async (id: string): Promise<ApiResponse<Empresa>> => {
    const response = await api.get(`/empresas/${id}`);
    return response.data;
  },
  crear: async (empresa: Omit<Empresa, 'id' | 'created_at' | 'updated_at'>): Promise<ApiResponse<Empresa>> => {
    const response = await api.post('/empresas', empresa);
    return response.data;
  },
  actualizar: async (id: string, empresa: Partial<Empresa>): Promise<ApiResponse<Empresa>> => {
    const response = await api.put(`/empresas/${id}`, empresa);
    return response.data;
  },
  eliminar: async (id: string): Promise<ApiResponse<void>> => {
    const response = await api.delete(`/empresas/${id}`);
    return response.data;
  },
};

export const usuariosApi = {
  listar: async (empresaId: string): Promise<ApiResponse<Usuario[]>> => {
    const response = await api.get('/usuarios', { params: { empresa_id: empresaId } });
    return response.data;
  },
  obtener: async (id: string): Promise<ApiResponse<Usuario>> => {
    const response = await api.get(`/usuarios/${id}`);
    return response.data;
  },
  crear: async (usuario: Omit<Usuario, 'id' | 'created_at' | 'updated_at'>): Promise<ApiResponse<Usuario>> => {
    const response = await api.post('/usuarios', usuario);
    return response.data;
  },
  actualizar: async (id: string, usuario: Partial<Usuario>): Promise<ApiResponse<Usuario>> => {
    const response = await api.put(`/usuarios/${id}`, usuario);
    return response.data;
  },
  eliminar: async (id: string): Promise<ApiResponse<void>> => {
    const response = await api.delete(`/usuarios/${id}`);
    return response.data;
  },
}; 