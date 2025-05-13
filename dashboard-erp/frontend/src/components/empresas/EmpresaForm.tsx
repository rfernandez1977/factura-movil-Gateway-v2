import { useForm } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { Empresa } from '@/types';

const schema = yup.object({
  rut: yup.string().required('El RUT es requerido'),
  razon_social: yup.string().required('La razón social es requerida'),
  giro: yup.string().required('El giro es requerido'),
  direccion: yup.string().required('La dirección es requerida'),
  comuna: yup.string().required('La comuna es requerida'),
  ciudad: yup.string().required('La ciudad es requerida'),
  email: yup.string().email('Email inválido').required('El email es requerido'),
  resolucion_numero: yup.string().required('El número de resolución es requerido'),
  resolucion_fecha: yup.string().required('La fecha de resolución es requerida'),
  firma_rut: yup.string().required('El RUT de la firma es requerido'),
  firma_nombre: yup.string().required('El nombre de la firma es requerido'),
  firma_expiracion: yup.string().required('La fecha de expiración de la firma es requerida'),
  clave_sii: yup.string().required('La clave SII es requerida'),
}).required();

interface EmpresaFormProps {
  empresa?: Empresa;
  onSubmit: (data: Omit<Empresa, 'id' | 'created_at' | 'updated_at'>) => void;
  onCancel: () => void;
}

export default function EmpresaForm({ empresa, onSubmit, onCancel }: EmpresaFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: yupResolver(schema),
    defaultValues: empresa,
  });

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <div>
          <label htmlFor="rut" className="block text-sm font-medium text-gray-700">
            RUT
          </label>
          <input
            type="text"
            id="rut"
            {...register('rut')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.rut && (
            <p className="mt-1 text-sm text-red-600">{errors.rut.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="razon_social" className="block text-sm font-medium text-gray-700">
            Razón Social
          </label>
          <input
            type="text"
            id="razon_social"
            {...register('razon_social')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.razon_social && (
            <p className="mt-1 text-sm text-red-600">{errors.razon_social.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="giro" className="block text-sm font-medium text-gray-700">
            Giro
          </label>
          <input
            type="text"
            id="giro"
            {...register('giro')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.giro && (
            <p className="mt-1 text-sm text-red-600">{errors.giro.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="email" className="block text-sm font-medium text-gray-700">
            Email
          </label>
          <input
            type="email"
            id="email"
            {...register('email')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.email && (
            <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="direccion" className="block text-sm font-medium text-gray-700">
            Dirección
          </label>
          <input
            type="text"
            id="direccion"
            {...register('direccion')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.direccion && (
            <p className="mt-1 text-sm text-red-600">{errors.direccion.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="comuna" className="block text-sm font-medium text-gray-700">
            Comuna
          </label>
          <input
            type="text"
            id="comuna"
            {...register('comuna')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.comuna && (
            <p className="mt-1 text-sm text-red-600">{errors.comuna.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="ciudad" className="block text-sm font-medium text-gray-700">
            Ciudad
          </label>
          <input
            type="text"
            id="ciudad"
            {...register('ciudad')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.ciudad && (
            <p className="mt-1 text-sm text-red-600">{errors.ciudad.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="resolucion_numero" className="block text-sm font-medium text-gray-700">
            Número de Resolución
          </label>
          <input
            type="text"
            id="resolucion_numero"
            {...register('resolucion_numero')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.resolucion_numero && (
            <p className="mt-1 text-sm text-red-600">{errors.resolucion_numero.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="resolucion_fecha" className="block text-sm font-medium text-gray-700">
            Fecha de Resolución
          </label>
          <input
            type="date"
            id="resolucion_fecha"
            {...register('resolucion_fecha')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.resolucion_fecha && (
            <p className="mt-1 text-sm text-red-600">{errors.resolucion_fecha.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="firma_rut" className="block text-sm font-medium text-gray-700">
            RUT de la Firma
          </label>
          <input
            type="text"
            id="firma_rut"
            {...register('firma_rut')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.firma_rut && (
            <p className="mt-1 text-sm text-red-600">{errors.firma_rut.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="firma_nombre" className="block text-sm font-medium text-gray-700">
            Nombre de la Firma
          </label>
          <input
            type="text"
            id="firma_nombre"
            {...register('firma_nombre')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.firma_nombre && (
            <p className="mt-1 text-sm text-red-600">{errors.firma_nombre.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="firma_expiracion" className="block text-sm font-medium text-gray-700">
            Fecha de Expiración de la Firma
          </label>
          <input
            type="date"
            id="firma_expiracion"
            {...register('firma_expiracion')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.firma_expiracion && (
            <p className="mt-1 text-sm text-red-600">{errors.firma_expiracion.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="clave_sii" className="block text-sm font-medium text-gray-700">
            Clave SII
          </label>
          <input
            type="password"
            id="clave_sii"
            {...register('clave_sii')}
            className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
          {errors.clave_sii && (
            <p className="mt-1 text-sm text-red-600">{errors.clave_sii.message}</p>
          )}
        </div>
      </div>

      <div className="flex justify-end space-x-3">
        <button
          type="button"
          onClick={onCancel}
          className="rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
        >
          Cancelar
        </button>
        <button
          type="submit"
          className="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
        >
          {empresa ? 'Actualizar' : 'Crear'}
        </button>
      </div>
    </form>
  );
} 