-- Esquema para la tabla de DTE
CREATE TABLE IF NOT EXISTS dte (
    id VARCHAR(36) PRIMARY KEY,
    version VARCHAR(10) NOT NULL,
    tipo_dte INTEGER NOT NULL,
    folio INTEGER NOT NULL,
    fecha_emision TIMESTAMP NOT NULL,
    rut_emisor VARCHAR(12) NOT NULL,
    razon_social_emisor VARCHAR(100) NOT NULL,
    giro_emisor VARCHAR(100) NOT NULL,
    direccion_emisor VARCHAR(100) NOT NULL,
    comuna_emisor VARCHAR(50) NOT NULL,
    ciudad_emisor VARCHAR(50) NOT NULL,
    correo_emisor VARCHAR(100) NOT NULL,
    rut_receptor VARCHAR(12) NOT NULL,
    razon_social_receptor VARCHAR(100) NOT NULL,
    giro_receptor VARCHAR(100) NOT NULL,
    direccion_receptor VARCHAR(100) NOT NULL,
    comuna_receptor VARCHAR(50) NOT NULL,
    ciudad_receptor VARCHAR(50) NOT NULL,
    monto_neto DECIMAL(18,2) NOT NULL,
    tasa_iva DECIMAL(5,2) NOT NULL,
    monto_iva DECIMAL(18,2) NOT NULL,
    monto_total DECIMAL(18,2) NOT NULL,
    track_id VARCHAR(50),
    estado VARCHAR(20) NOT NULL,
    fecha_creacion TIMESTAMP NOT NULL,
    fecha_actualizacion TIMESTAMP,
    UNIQUE(tipo_dte, folio)
);

-- Esquema para la tabla de detalles de DTE
CREATE TABLE IF NOT EXISTS dte_detalle (
    id SERIAL PRIMARY KEY,
    dte_id VARCHAR(36) NOT NULL,
    numero_linea INTEGER NOT NULL,
    nombre_item VARCHAR(100) NOT NULL,
    cantidad DECIMAL(18,6) NOT NULL,
    precio_unitario DECIMAL(18,2) NOT NULL,
    monto_item DECIMAL(18,2) NOT NULL,
    FOREIGN KEY (dte_id) REFERENCES dte(id) ON DELETE CASCADE,
    UNIQUE(dte_id, numero_linea)
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_dte_rut_emisor ON dte(rut_emisor);
CREATE INDEX IF NOT EXISTS idx_dte_fecha_emision ON dte(fecha_emision);
CREATE INDEX IF NOT EXISTS idx_dte_estado ON dte(estado);
CREATE INDEX IF NOT EXISTS idx_dte_track_id ON dte(track_id);
CREATE INDEX IF NOT EXISTS idx_dte_detalle_dte_id ON dte_detalle(dte_id); 