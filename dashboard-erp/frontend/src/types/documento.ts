export interface Item {
    id?: string;
    descripcion: string;
    nombre?: string;
    codigo?: string;
    cantidad: number;
    precio_unitario: number;
    descuento: number;
    exento?: boolean;
    porcentaje_iva?: number;
    monto_iva?: number;
    monto_item: number;
    subtotal?: number;
    impuestos_adicionales?: ImpuestoAdicional[];
}

export interface ImpuestoAdicional {
    codigo: string;
    nombre: string;
    porcentaje: number;
    monto: number;
    base_imponible?: number;
}

export interface EventoHistorial {
    fecha: string;
    estado: string;
    detalle?: string;
    usuario?: string;
}

export type EstadoDocumento =
    | 'PENDIENTE'
    | 'ENVIADO'
    | 'ACEPTADO'
    | 'RECHAZADO'
    | 'ERROR'
    | 'ANULADO';

export interface DocumentoTributario {
    id: string;
    tipo: string;
    folio: number;
    rut_emisor: string;
    rut_receptor: string;
    razon_social_emisor?: string;
    razon_social_receptor?: string;
    giro_emisor?: string;
    giro_receptor?: string;
    fecha_emision: string;
    fecha_vencimiento?: string;
    monto_total: number;
    monto_neto: number;
    monto_iva: number;
    monto_exento?: number;
    estado: EstadoDocumento;
    estado_sii?: string;
    glosa_sii?: string;
    track_id?: string;
    items?: Item[];
    xml?: string;
    pdf?: string;
    historial?: EventoHistorial[];
    created_at: string;
    updated_at: string;
}

export interface DocumentoRequest {
    tipo: string;
    rut_emisor: string;
    rut_receptor: string;
    razon_social_emisor?: string;
    razon_social_receptor?: string;
    giro_emisor?: string;
    giro_receptor?: string;
    fecha_emision: string;
    items: Item[];
}

export interface DocumentoEstadisticas {
    pendientes: number;
    enviados: number;
    aceptados: number;
    rechazados: number;
    error: number;
    anulados: number;
    total: number;
} 