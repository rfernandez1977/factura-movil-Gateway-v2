-- Crear tabla de documentos tributarios
CREATE TABLE documentos_tributarios (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    tipo VARCHAR(50) NOT NULL,
    folio BIGINT NOT NULL,
    fecha_emision DATETIME NOT NULL,
    rut_emisor VARCHAR(12) NOT NULL,
    razon_social_emisor VARCHAR(100) NOT NULL,
    rut_receptor VARCHAR(12) NOT NULL,
    razon_social_receptor VARCHAR(100) NOT NULL,
    monto_neto DECIMAL(18,2) NOT NULL,
    monto_iva DECIMAL(18,2) NOT NULL,
    monto_total DECIMAL(18,2) NOT NULL,
    estado VARCHAR(50) NOT NULL,
    track_id VARCHAR(100),
    xml TEXT,
    pdf TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    UNIQUE KEY uk_documento_tipo_folio (tipo, folio)
);

-- Crear tabla de facturas
CREATE TABLE facturas (
    id BIGINT PRIMARY KEY,
    tipo_factura VARCHAR(50) NOT NULL,
    forma_pago VARCHAR(50) NOT NULL,
    fecha_vencimiento DATETIME NOT NULL,
    FOREIGN KEY (id) REFERENCES documentos_tributarios(id)
);

-- Crear tabla de boletas
CREATE TABLE boletas (
    id BIGINT PRIMARY KEY,
    FOREIGN KEY (id) REFERENCES documentos_tributarios(id)
);

-- Crear tabla de guías de despacho
CREATE TABLE guias_despacho (
    id BIGINT PRIMARY KEY,
    indicador_traslado VARCHAR(50) NOT NULL,
    indicador_servicio VARCHAR(50) NOT NULL,
    indicador_ventas VARCHAR(50) NOT NULL,
    patente VARCHAR(10),
    rut_transportista VARCHAR(12),
    nombre_chofer VARCHAR(100),
    direccion_origen VARCHAR(200),
    comuna_origen VARCHAR(100),
    direccion_destino VARCHAR(200),
    comuna_destino VARCHAR(100),
    FOREIGN KEY (id) REFERENCES documentos_tributarios(id)
);

-- Crear tabla de notas de crédito
CREATE TABLE notas_credito (
    id BIGINT PRIMARY KEY,
    indicador_servicio VARCHAR(50) NOT NULL,
    indicador_ventas VARCHAR(50) NOT NULL,
    tipo_documento_referencia VARCHAR(50) NOT NULL,
    folio_referencia BIGINT NOT NULL,
    fecha_referencia DATETIME NOT NULL,
    razon_referencia VARCHAR(200) NOT NULL,
    FOREIGN KEY (id) REFERENCES documentos_tributarios(id)
);

-- Crear tabla de notas de débito
CREATE TABLE notas_debito (
    id BIGINT PRIMARY KEY,
    indicador_servicio VARCHAR(50) NOT NULL,
    indicador_ventas VARCHAR(50) NOT NULL,
    tipo_documento_referencia VARCHAR(50) NOT NULL,
    folio_referencia BIGINT NOT NULL,
    fecha_referencia DATETIME NOT NULL,
    razon_referencia VARCHAR(200) NOT NULL,
    FOREIGN KEY (id) REFERENCES documentos_tributarios(id)
);

-- Crear tabla de ítems
CREATE TABLE items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    documento_id BIGINT NOT NULL,
    codigo VARCHAR(50) NOT NULL,
    descripcion VARCHAR(200) NOT NULL,
    cantidad DECIMAL(18,2) NOT NULL,
    precio_unitario DECIMAL(18,2) NOT NULL,
    monto_total DECIMAL(18,2) NOT NULL,
    unidad_medida VARCHAR(50),
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (documento_id) REFERENCES documentos_tributarios(id)
);

-- Crear tabla de estados de documentos
CREATE TABLE estados_documentos (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    documento_id BIGINT NOT NULL,
    estado VARCHAR(50) NOT NULL,
    glosa VARCHAR(200),
    fecha DATETIME NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (documento_id) REFERENCES documentos_tributarios(id)
);

-- Crear tabla de errores de documentos
CREATE TABLE errores_documentos (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    documento_id BIGINT NOT NULL,
    codigo VARCHAR(50) NOT NULL,
    glosa VARCHAR(200) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    FOREIGN KEY (documento_id) REFERENCES documentos_tributarios(id)
);

-- Crear índices
CREATE INDEX idx_documentos_tributarios_rut_emisor ON documentos_tributarios(rut_emisor);
CREATE INDEX idx_documentos_tributarios_rut_receptor ON documentos_tributarios(rut_receptor);
CREATE INDEX idx_documentos_tributarios_fecha_emision ON documentos_tributarios(fecha_emision);
CREATE INDEX idx_documentos_tributarios_estado ON documentos_tributarios(estado);
CREATE INDEX idx_items_documento_id ON items(documento_id);
CREATE INDEX idx_estados_documentos_documento_id ON estados_documentos(documento_id);
CREATE INDEX idx_errores_documentos_documento_id ON errores_documentos(documento_id); 