import React, { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useSupabaseClient } from '@supabase/auth-helpers-react';
import { DocumentoTributario } from '@/types/documento';
import { formatCurrency, formatDate, formatRut } from '@/utils/formatters';
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import {
    Table,
    TableBody,
    TableCell,
    TableHead,
    TableHeader,
    TableRow
} from '@/components/ui/table';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Separator } from '@/components/ui/separator';
import { ArrowLeft, FileDown, Send, Trash, Printer, RotateCcw } from 'lucide-react';
import Link from 'next/link';
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert';

const DocumentoDetail = () => {
    const params = useParams();
    const router = useRouter();
    const documentoId = params.id as string;
    const supabase = useSupabaseClient();

    const [documento, setDocumento] = useState<DocumentoTributario | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchDocumento = async () => {
            try {
                setLoading(true);
                setError(null);

                const { data, error } = await supabase
                    .from('documentos')
                    .select('*')
                    .eq('id', documentoId)
                    .single();

                if (error) throw error;

                setDocumento(data);
            } catch (err: any) {
                console.error('Error al cargar el documento:', err);
                setError(err.message || 'Error al cargar el documento');
            } finally {
                setLoading(false);
            }
        };

        if (documentoId) {
            fetchDocumento();
        }
    }, [documentoId]);

    const handleSendToSII = async () => {
        try {
            const { error } = await supabase.functions.invoke('send-to-sii', {
                body: { documentoId }
            });

            if (error) throw error;

            // Recargar el documento para obtener el trackID y estado actualizado
            const { data, error: fetchError } = await supabase
                .from('documentos')
                .select('*')
                .eq('id', documentoId)
                .single();

            if (fetchError) throw fetchError;

            setDocumento(data);

        } catch (err: any) {
            console.error('Error al enviar al SII:', err);
            setError(err.message || 'Error al enviar al SII');
        }
    };

    const handleDelete = async () => {
        if (confirm('¿Está seguro de eliminar este documento?')) {
            try {
                const { error } = await supabase
                    .from('documentos')
                    .delete()
                    .eq('id', documentoId);

                if (error) throw error;

                router.push('/documentos');
            } catch (err: any) {
                console.error('Error al eliminar el documento:', err);
                setError(err.message || 'Error al eliminar el documento');
            }
        }
    };

    const tipoDocumentoLabel = (tipo: string) => {
        switch (tipo) {
            case '33': return 'Factura Electrónica';
            case '34': return 'Factura Exenta Electrónica';
            case '39': return 'Boleta Electrónica';
            case '41': return 'Boleta Exenta Electrónica';
            case '56': return 'Nota Débito Electrónica';
            case '61': return 'Nota Crédito Electrónica';
            case '52': return 'Guía Despacho Electrónica';
            default: return tipo;
        }
    };

    const estadoBadgeColor = (estado: string) => {
        switch (estado) {
            case 'PENDIENTE': return 'bg-yellow-100 text-yellow-800';
            case 'ENVIADO': return 'bg-blue-100 text-blue-800';
            case 'ACEPTADO': return 'bg-green-100 text-green-800';
            case 'RECHAZADO': return 'bg-red-100 text-red-800';
            case 'ERROR': return 'bg-red-100 text-red-800';
            case 'ANULADO': return 'bg-gray-100 text-gray-800';
            default: return 'bg-gray-100 text-gray-800';
        }
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[50vh]">
                <p>Cargando documento...</p>
            </div>
        );
    }

    if (error) {
        return (
            <Alert variant="destructive" className="mb-4">
                <AlertTitle>Error</AlertTitle>
                <AlertDescription>{error}</AlertDescription>
            </Alert>
        );
    }

    if (!documento) {
        return (
            <Alert className="mb-4">
                <AlertTitle>Documento no encontrado</AlertTitle>
                <AlertDescription>El documento solicitado no existe o no tiene acceso a él.</AlertDescription>
            </Alert>
        );
    }

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <Button variant="ghost" onClick={() => router.back()}>
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Volver
                </Button>

                <div className="flex space-x-2">
                    {documento.estado === 'PENDIENTE' && (
                        <Button onClick={handleSendToSII}>
                            <Send className="mr-2 h-4 w-4" />
                            Enviar al SII
                        </Button>
                    )}

                    <Button variant="outline" asChild>
                        <Link href={`/documentos/${documentoId}/pdf`} target="_blank">
                            <FileDown className="mr-2 h-4 w-4" />
                            Descargar PDF
                        </Link>
                    </Button>

                    <Button variant="outline">
                        <Printer className="mr-2 h-4 w-4" />
                        Imprimir
                    </Button>

                    {documento.estado !== 'ACEPTADO' && (
                        <Button variant="destructive" onClick={handleDelete}>
                            <Trash className="mr-2 h-4 w-4" />
                            Eliminar
                        </Button>
                    )}
                </div>
            </div>

            <Card>
                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <div>
                        <CardTitle className="text-2xl font-bold">
                            {tipoDocumentoLabel(documento.tipo)} #{documento.folio}
                        </CardTitle>
                        <p className="text-sm text-gray-500">
                            Emitido el {formatDate(documento.fecha_emision)}
                        </p>
                    </div>
                    <Badge className={estadoBadgeColor(documento.estado)}>
                        {documento.estado}
                    </Badge>
                </CardHeader>

                <CardContent>
                    <Tabs defaultValue="detalles">
                        <TabsList>
                            <TabsTrigger value="detalles">Detalles</TabsTrigger>
                            <TabsTrigger value="items">Ítems</TabsTrigger>
                            <TabsTrigger value="historial">Historial</TabsTrigger>
                            {documento.xml && <TabsTrigger value="xml">XML</TabsTrigger>}
                        </TabsList>

                        <TabsContent value="detalles" className="space-y-6 pt-4">
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                <div className="space-y-4">
                                    <h3 className="text-lg font-semibold">Emisor</h3>
                                    <div className="space-y-2">
                                        <p><span className="font-medium">RUT:</span> {formatRut(documento.rut_emisor)}</p>
                                        <p><span className="font-medium">Razón Social:</span> {documento.razon_social_emisor || 'No especificado'}</p>
                                        <p><span className="font-medium">Giro:</span> {documento.giro_emisor || 'No especificado'}</p>
                                    </div>
                                </div>

                                <div className="space-y-4">
                                    <h3 className="text-lg font-semibold">Receptor</h3>
                                    <div className="space-y-2">
                                        <p><span className="font-medium">RUT:</span> {formatRut(documento.rut_receptor)}</p>
                                        <p><span className="font-medium">Razón Social:</span> {documento.razon_social_receptor || 'No especificado'}</p>
                                        <p><span className="font-medium">Giro:</span> {documento.giro_receptor || 'No especificado'}</p>
                                    </div>
                                </div>
                            </div>

                            <Separator />

                            <div className="space-y-4">
                                <h3 className="text-lg font-semibold">Totales</h3>
                                <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                                    <div className="bg-gray-50 p-4 rounded-md">
                                        <p className="text-sm text-gray-500">Neto</p>
                                        <p className="text-lg font-semibold">{formatCurrency(documento.monto_neto)}</p>
                                    </div>
                                    <div className="bg-gray-50 p-4 rounded-md">
                                        <p className="text-sm text-gray-500">IVA</p>
                                        <p className="text-lg font-semibold">{formatCurrency(documento.monto_iva)}</p>
                                    </div>
                                    <div className="bg-gray-50 p-4 rounded-md">
                                        <p className="text-sm text-gray-500">Exento</p>
                                        <p className="text-lg font-semibold">{formatCurrency(documento.monto_exento || 0)}</p>
                                    </div>
                                    <div className="bg-gray-50 p-4 rounded-md">
                                        <p className="text-sm text-gray-500">Total</p>
                                        <p className="text-lg font-semibold">{formatCurrency(documento.monto_total)}</p>
                                    </div>
                                </div>
                            </div>

                            {documento.track_id && (
                                <>
                                    <Separator />
                                    <div className="space-y-4">
                                        <h3 className="text-lg font-semibold">Información SII</h3>
                                        <div className="space-y-2">
                                            <p><span className="font-medium">TrackID:</span> {documento.track_id}</p>
                                            <p><span className="font-medium">Estado:</span> {documento.estado_sii || documento.estado}</p>
                                            {documento.glosa_sii && <p><span className="font-medium">Glosa:</span> {documento.glosa_sii}</p>}
                                        </div>
                                    </div>
                                </>
                            )}
                        </TabsContent>

                        <TabsContent value="items" className="pt-4">
                            <Table>
                                <TableHeader>
                                    <TableRow>
                                        <TableHead>#</TableHead>
                                        <TableHead>Descripción</TableHead>
                                        <TableHead className="text-right">Cantidad</TableHead>
                                        <TableHead className="text-right">Precio Unit.</TableHead>
                                        <TableHead className="text-right">Descuento</TableHead>
                                        <TableHead className="text-right">Total</TableHead>
                                    </TableRow>
                                </TableHeader>
                                <TableBody>
                                    {documento.items && documento.items.length > 0 ? (
                                        documento.items.map((item, index) => (
                                            <TableRow key={index}>
                                                <TableCell>{index + 1}</TableCell>
                                                <TableCell>{item.nombre || item.descripcion}</TableCell>
                                                <TableCell className="text-right">{item.cantidad}</TableCell>
                                                <TableCell className="text-right">{formatCurrency(item.precio_unitario)}</TableCell>
                                                <TableCell className="text-right">{item.descuento || 0}%</TableCell>
                                                <TableCell className="text-right">{formatCurrency(item.monto_item)}</TableCell>
                                            </TableRow>
                                        ))
                                    ) : (
                                        <TableRow>
                                            <TableCell colSpan={6} className="text-center py-4">
                                                No hay ítems disponibles
                                            </TableCell>
                                        </TableRow>
                                    )}
                                </TableBody>
                            </Table>
                        </TabsContent>

                        <TabsContent value="historial" className="pt-4 space-y-4">
                            {documento.historial && documento.historial.length > 0 ? (
                                <div className="space-y-4">
                                    {documento.historial.map((evento, index) => (
                                        <div key={index} className="border-l-2 border-gray-200 pl-4 py-2">
                                            <p className="text-sm text-gray-500">{formatDate(evento.fecha)}</p>
                                            <p className="font-medium">{evento.estado}</p>
                                            {evento.detalle && <p className="text-sm">{evento.detalle}</p>}
                                        </div>
                                    ))}
                                </div>
                            ) : (
                                <p className="text-center py-4">No hay historial disponible</p>
                            )}
                        </TabsContent>

                        {documento.xml && (
                            <TabsContent value="xml" className="pt-4">
                                <div className="bg-gray-100 p-4 rounded-md overflow-auto max-h-96">
                                    <pre className="whitespace-pre-wrap text-xs">
                                        {documento.xml}
                                    </pre>
                                </div>
                            </TabsContent>
                        )}
                    </Tabs>
                </CardContent>

                <CardFooter className="flex justify-between border-t pt-6">
                    <Button variant="ghost" onClick={() => router.back()}>
                        <ArrowLeft className="mr-2 h-4 w-4" />
                        Volver
                    </Button>

                    {documento.estado === 'RECHAZADO' || documento.estado === 'ERROR' ? (
                        <Button variant="outline">
                            <RotateCcw className="mr-2 h-4 w-4" />
                            Reintentar
                        </Button>
                    ) : null}
                </CardFooter>
            </Card>
        </div>
    );
};

export default DocumentoDetail; 