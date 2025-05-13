-- Agregar NOTA_VENTA como tipo de documento válido
ALTER TYPE tipo_documento ADD VALUE 'NOTA_VENTA';

-- Actualizar la tabla de documentos para incluir validaciones específicas de nota de venta
ALTER TABLE documentos
ADD COLUMN IF NOT EXISTS tipo_nota_venta VARCHAR(50),
ADD COLUMN IF NOT EXISTS referencia_documento VARCHAR(50);

-- Crear índice para búsquedas por tipo de nota de venta
CREATE INDEX IF NOT EXISTS idx_documentos_tipo_nota_venta ON documentos(tipo_nota_venta);

-- Agregar comentarios descriptivos
COMMENT ON COLUMN documentos.tipo_nota_venta IS 'Tipo específico de nota de venta (ej: venta al contado, venta a crédito)';
COMMENT ON COLUMN documentos.referencia_documento IS 'Documento de referencia para la nota de venta'; 