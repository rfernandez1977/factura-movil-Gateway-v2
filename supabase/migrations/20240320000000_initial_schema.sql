-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create empresas table
CREATE TABLE empresas (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rut VARCHAR(20) NOT NULL UNIQUE,
    razon_social VARCHAR(255) NOT NULL,
    giro VARCHAR(255),
    direccion VARCHAR(255),
    comuna VARCHAR(100),
    ciudad VARCHAR(100),
    email VARCHAR(255),
    resolucion_numero VARCHAR(50),
    resolucion_fecha TIMESTAMP WITH TIME ZONE,
    firma_rut VARCHAR(20),
    firma_nombre VARCHAR(255),
    firma_expiracion TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create documentos table
CREATE TABLE documentos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tipo VARCHAR(50) NOT NULL,
    rut_emisor VARCHAR(20) NOT NULL,
    rut_receptor VARCHAR(20) NOT NULL,
    folio INTEGER NOT NULL,
    monto_total DECIMAL(15,2) NOT NULL,
    estado VARCHAR(50) NOT NULL,
    xml TEXT,
    pdf TEXT,
    firma TEXT,
    ted TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create certificados table
CREATE TABLE certificados (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rut VARCHAR(20) NOT NULL UNIQUE,
    certificado TEXT NOT NULL,
    llave_privada TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create sesiones table
CREATE TABLE sesiones (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rut VARCHAR(20) NOT NULL,
    token TEXT NOT NULL,
    expiracion TIMESTAMP WITH TIME ZONE NOT NULL,
    estado VARCHAR(50) NOT NULL DEFAULT 'ACTIVA',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create xml_files table
CREATE TABLE xml_files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tipo VARCHAR(50) NOT NULL,
    ruta TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create csfs table (Códigos de Seguridad del Folio)
CREATE TABLE csfs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    empresa_id UUID NOT NULL REFERENCES empresas(id),
    tipo_documento VARCHAR(50) NOT NULL,
    folio_inicial INTEGER NOT NULL,
    folio_final INTEGER NOT NULL,
    folio_actual INTEGER NOT NULL DEFAULT 0,
    estado VARCHAR(50) NOT NULL,
    fecha_resolucion TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_documentos_rut_emisor ON documentos(rut_emisor);
CREATE INDEX idx_documentos_rut_receptor ON documentos(rut_receptor);
CREATE INDEX idx_documentos_estado ON documentos(estado);
CREATE INDEX idx_sesiones_token ON sesiones(token);
CREATE INDEX idx_sesiones_estado ON sesiones(estado);
CREATE INDEX idx_csfs_empresa_id ON csfs(empresa_id);
CREATE INDEX idx_csfs_tipo_documento ON csfs(tipo_documento);

-- Create RLS policies
ALTER TABLE empresas ENABLE ROW LEVEL SECURITY;
ALTER TABLE documentos ENABLE ROW LEVEL SECURITY;
ALTER TABLE certificados ENABLE ROW LEVEL SECURITY;
ALTER TABLE sesiones ENABLE ROW LEVEL SECURITY;
ALTER TABLE xml_files ENABLE ROW LEVEL SECURITY;
ALTER TABLE csfs ENABLE ROW LEVEL SECURITY;

-- Create policies for authenticated users
CREATE POLICY "Users can view their own data" ON empresas
    FOR SELECT USING (auth.uid() = id);

CREATE POLICY "Users can view their own documents" ON documentos
    FOR SELECT USING (auth.uid() = rut_emisor);

CREATE POLICY "Users can view their own certificates" ON certificados
    FOR SELECT USING (auth.uid() = rut);

CREATE POLICY "Users can view their own sessions" ON sesiones
    FOR SELECT USING (auth.uid() = rut);

-- Create functions for updating timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updating timestamps
CREATE TRIGGER update_empresas_updated_at
    BEFORE UPDATE ON empresas
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_documentos_updated_at
    BEFORE UPDATE ON documentos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_certificados_updated_at
    BEFORE UPDATE ON certificados
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sesiones_updated_at
    BEFORE UPDATE ON sesiones
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_xml_files_updated_at
    BEFORE UPDATE ON xml_files
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_csfs_updated_at
    BEFORE UPDATE ON csfs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column(); 