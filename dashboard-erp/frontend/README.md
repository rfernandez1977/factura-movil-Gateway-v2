# Dashboard ERP Frontend

Este es el frontend del sistema ERP, construido con Next.js, TypeScript y Tailwind CSS.

## Requisitos

- Node.js 18.x o superior
- npm 9.x o superior

## Instalación

1. Clona el repositorio
2. Navega al directorio del frontend:
   ```bash
   cd dashboard-erp/frontend
   ```
3. Instala las dependencias:
   ```bash
   npm install
   ```

## Configuración

1. Crea un archivo `.env.local` en el directorio raíz del frontend con las siguientes variables:
   ```
   NEXT_PUBLIC_API_URL=http://localhost:8080
   NEXT_PUBLIC_SUPABASE_URL=your-supabase-url
   NEXT_PUBLIC_SUPABASE_ANON_KEY=your-supabase-anon-key
   ```

2. Reemplaza los valores de las variables de entorno con tus propias credenciales.

## Desarrollo

Para iniciar el servidor de desarrollo:

```bash
npm run dev
```

El servidor se iniciará en `http://localhost:3000`.

## Construcción

Para construir la aplicación para producción:

```bash
npm run build
```

## Producción

Para iniciar la aplicación en modo producción:

```bash
npm start
```

## Estructura del Proyecto

```
src/
  ├── app/              # Páginas de la aplicación (Next.js App Router)
  ├── components/       # Componentes reutilizables
  ├── lib/             # Utilidades y configuraciones
  ├── hooks/           # Custom hooks
  ├── types/           # Definiciones de tipos TypeScript
  ├── utils/           # Funciones utilitarias
  ├── services/        # Servicios de API
  ├── store/           # Estado global (si se usa)
  └── styles/          # Estilos globales
```

## Características

- Autenticación con Supabase
- Gestión de empresas
- Gestión de usuarios
- Gestión de documentos
- Interfaz responsiva con Tailwind CSS
- Formularios con validación usando react-hook-form y yup
- Manejo de estado con React Query
- Tipado fuerte con TypeScript
