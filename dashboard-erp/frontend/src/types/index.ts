export interface Empresa {
  id: string;
  rut: string;
  razon_social: string;
  giro: string;
  direccion: string;
  comuna: string;
  ciudad: string;
  email: string;
  resolucion_numero: string;
  resolucion_fecha: string;
  firma_rut: string;
  firma_nombre: string;
  firma_expiracion: string;
  certificado_id: string;
  nombre: string;
  clave_sii: string;
  created_at: string;
  updated_at: string;
}

export interface Usuario {
  id: string;
  email: string;
  nombre: string;
  apellido: string;
  empresa_id: string;
  rol_id: string;
  activo: boolean;
  created_at: string;
  updated_at: string;
}

export interface Rol {
  id: string;
  nombre: string;
  descripcion: string;
  empresa_id: string;
  created_at: string;
  updated_at: string;
}

export interface Permiso {
  id: string;
  nombre: string;
  descripcion: string;
  created_at: string;
  updated_at: string;
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
} 