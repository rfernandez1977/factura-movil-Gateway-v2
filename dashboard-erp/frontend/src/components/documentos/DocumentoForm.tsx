import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useSupabaseClient } from '@supabase/auth-helpers-react';
import { useForm, useFieldArray, Controller } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';

import { Card, CardContent, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage
} from '@/components/ui/form';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert';
import { Plus, Trash, ArrowLeft, Check } from 'lucide-react';

// Esquema de validación para el formulario
const itemSchema = z.object({
    descripcion: z.string().min(1, 'La descripción es requerida'),
    cantidad: z.coerce.number().min(0.01, 'La cantidad debe ser mayor a 0'),
    precio_unitario: z.coerce.number().min(1, 'El precio debe ser mayor a 0'),
    descuento: z.coerce.number().min(0).max(100, 'El descuento debe estar entre 0 y 100'),
});

const documentoSchema = z.object({
    tipo: z.string().min(1, 'El tipo de documento es requerido'),
    rut_emisor: z.string().min(1, 'El RUT del emisor es requerido'),
    rut_receptor: z.string().min(1, 'El RUT del receptor es requerido'),
    razon_social_emisor: z.string().optional(),
    razon_social_receptor: z.string().optional(),
    giro_emisor: z.string().optional(),
    giro_receptor: z.string().optional(),
    fecha_emision: z.string().optional(),
    items: z.array(itemSchema).min(1, 'Debe agregar al menos un ítem'),
});

type DocumentoFormValues = z.infer<typeof documentoSchema>;

type DocumentoFormProps = {
    documentoId?: string; // Opcional, para edición
};

const DocumentoForm = ({ documentoId }: DocumentoFormProps) => {
    const router = useRouter();
    const supabase = useSupabaseClient();
    const [empresas, setEmpresas] = useState<any[]>([]);
    const [receptores, setReceptores] = useState<any[]>([]);
    const [loading, setLoading] = useState(false);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [totales, setTotales] = useState({
        neto: 0,
        iva: 0,
        total: 0
    });

    // Inicializar formulario
    const form = useForm<DocumentoFormValues>({
        resolver: zodResolver(documentoSchema),
        defaultValues: {
            tipo: '33', // Factura Electrónica por defecto
            rut_emisor: '',
            rut_receptor: '',
            razon_social_emisor: '',
            razon_social_receptor: '',
            giro_emisor: '',
            giro_receptor: '',
            fecha_emision: new Date().toISOString().split('T')[0],
            items: [
                { descripcion: '', cantidad: 1, precio_unitario: 0, descuento: 0 }
            ]
        }
    });

    // Field array para manejar los ítems
    const { fields, append, remove } = useFieldArray({
        control: form.control,
        name: "items"
    });

    // Cargar datos existentes para edición
    useEffect(() => {
        const fetchDocumento = async () => {
            if (!documentoId) return;

            try {
                setLoading(true);
                const { data, error } = await supabase
                    .from('documentos')
                    .select('*')
                    .eq('id', documentoId)
                    .single();

                if (error) throw error;

                if (data) {
                    // Mapear datos al formulario
                    form.reset({
                        tipo: data.tipo,
                        rut_emisor: data.rut_emisor,
                        rut_receptor: data.rut_receptor,
                        razon_social_emisor: data.razon_social_emisor,
                        razon_social_receptor: data.razon_social_receptor,
                        giro_emisor: data.giro_emisor,
                        giro_receptor: data.giro_receptor,
                        fecha_emision: data.fecha_emision.split('T')[0],
                        items: data.items || [{ descripcion: '', cantidad: 1, precio_unitario: 0, descuento: 0 }]
                    });

                    calcularTotales(data.items);
                }
            } catch (err: any) {
                console.error('Error al cargar documento:', err);
                setError(err.message || 'Error al cargar documento');
            } finally {
                setLoading(false);
            }
        };

        const fetchEmpresas = async () => {
            try {
                const { data, error } = await supabase
                    .from('empresas')
                    .select('*')
                    .order('razon_social', { ascending: true });

                if (error) throw error;

                setEmpresas(data || []);
            } catch (err) {
                console.error('Error al cargar empresas:', err);
            }
        };

        const fetchReceptores = async () => {
            try {
                const { data, error } = await supabase
                    .from('clientes')
                    .select('*')
                    .order('razon_social', { ascending: true });

                if (error) throw error;

                setReceptores(data || []);
            } catch (err) {
                console.error('Error al cargar receptores:', err);
            }
        };

        fetchEmpresas();
        fetchReceptores();
        fetchDocumento();
    }, [documentoId]);

    // Calcular totales cuando cambian los ítems
    const calcularTotales = (items: any[]) => {
        if (!items || items.length === 0) {
            setTotales({ neto: 0, iva: 0, total: 0 });
            return;
        }

        let neto = 0;

        items.forEach(item => {
            if (!item.cantidad || !item.precio_unitario) return;

            const subtotal = item.cantidad * item.precio_unitario;
            const descuento = subtotal * (item.descuento || 0) / 100;
            neto += subtotal - descuento;
        });

        const iva = neto * 0.19; // 19% IVA en Chile
        const total = neto + iva;

        setTotales({
            neto: Math.round(neto),
            iva: Math.round(iva),
            total: Math.round(total)
        });
    };

    // Recalcular totales cuando cambian los ítems
    useEffect(() => {
        const subscription = form.watch((value, { name }) => {
            if (name?.includes('items')) {
                calcularTotales(value.items as any[]);
            }
        });

        return () => subscription.unsubscribe();
    }, [form.watch]);

    // Manejar selección de empresa
    const handleSelectEmpresa = (empresaId: string) => {
        const empresa = empresas.find(e => e.id === empresaId);
        if (empresa) {
            form.setValue('rut_emisor', empresa.rut);
            form.setValue('razon_social_emisor', empresa.razon_social);
            form.setValue('giro_emisor', empresa.giro);
        }
    };

    // Manejar selección de receptor
    const handleSelectReceptor = (receptorId: string) => {
        const receptor = receptores.find(r => r.id === receptorId);
        if (receptor) {
            form.setValue('rut_receptor', receptor.rut);
            form.setValue('razon_social_receptor', receptor.razon_social);
            form.setValue('giro_receptor', receptor.giro);
        }
    };

    // Enviar formulario
    const onSubmit = async (data: DocumentoFormValues) => {
        try {
            setIsSubmitting(true);
            setError(null);

            // Calcular montos
            calcularTotales(data.items);

            const documentoData = {
                ...data,
                monto_neto: totales.neto,
                monto_iva: totales.iva,
                monto_total: totales.total,
                estado: 'PENDIENTE',
                fecha_emision: new Date(data.fecha_emision || new Date()).toISOString()
            };

            let result;

            if (documentoId) {
                // Actualizar documento existente
                result = await supabase
                    .from('documentos')
                    .update(documentoData)
                    .eq('id', documentoId);
            } else {
                // Crear nuevo documento
                result = await supabase
                    .from('documentos')
                    .insert(documentoData)
                    .select();
            }

            if (result.error) throw result.error;

            // Redirigir a la vista de detalle o lista
            if (documentoId) {
                router.push(`/documentos/${documentoId}`);
            } else if (result.data && result.data.length > 0) {
                router.push(`/documentos/${result.data[0].id}`);
            } else {
                router.push('/documentos');
            }

        } catch (err: any) {
            console.error('Error al guardar documento:', err);
            setError(err.message || 'Error al guardar documento');
        } finally {
            setIsSubmitting(false);
        }
    };

    // Agregar nuevo ítem
    const addItem = () => {
        append({ descripcion: '', cantidad: 1, precio_unitario: 0, descuento: 0 });
    };

    if (loading) {
        return (
            <div className="flex items-center justify-center min-h-[50vh]">
                <p>Cargando formulario...</p>
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div className="flex items-center justify-between">
                <Button variant="ghost" onClick={() => router.back()}>
                    <ArrowLeft className="mr-2 h-4 w-4" />
                    Volver
                </Button>

                <h1 className="text-2xl font-bold">
                    {documentoId ? 'Editar Documento' : 'Crear Nuevo Documento'}
                </h1>
            </div>

            {error && (
                <Alert variant="destructive">
                    <AlertTitle>Error</AlertTitle>
                    <AlertDescription>{error}</AlertDescription>
                </Alert>
            )}

            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)}>
                    <Card>
                        <CardHeader>
                            <CardTitle>Información General</CardTitle>
                        </CardHeader>
                        <CardContent>
                            <Tabs defaultValue="basico" className="w-full">
                                <TabsList className="mb-4">
                                    <TabsTrigger value="basico">Datos Básicos</TabsTrigger>
                                    <TabsTrigger value="emisor">Emisor</TabsTrigger>
                                    <TabsTrigger value="receptor">Receptor</TabsTrigger>
                                    <TabsTrigger value="items">Ítems</TabsTrigger>
                                </TabsList>

                                <TabsContent value="basico" className="space-y-4">
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                        <FormField
                                            control={form.control}
                                            name="tipo"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Tipo de Documento</FormLabel>
                                                    <Select
                                                        onValueChange={field.onChange}
                                                        defaultValue={field.value}
                                                    >
                                                        <FormControl>
                                                            <SelectTrigger>
                                                                <SelectValue placeholder="Seleccione tipo de documento" />
                                                            </SelectTrigger>
                                                        </FormControl>
                                                        <SelectContent>
                                                            <SelectItem value="33">Factura Electrónica</SelectItem>
                                                            <SelectItem value="34">Factura Exenta Electrónica</SelectItem>
                                                            <SelectItem value="39">Boleta Electrónica</SelectItem>
                                                            <SelectItem value="56">Nota Débito Electrónica</SelectItem>
                                                            <SelectItem value="61">Nota Crédito Electrónica</SelectItem>
                                                            <SelectItem value="52">Guía Despacho Electrónica</SelectItem>
                                                        </SelectContent>
                                                    </Select>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />

                                        <FormField
                                            control={form.control}
                                            name="fecha_emision"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Fecha de Emisión</FormLabel>
                                                    <FormControl>
                                                        <Input type="date" {...field} />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                    </div>

                                    <div className="flex flex-col p-4 mt-4 bg-gray-50 rounded-md">
                                        <div className="grid grid-cols-3 gap-4">
                                            <div>
                                                <p className="text-sm font-medium text-gray-500">Neto</p>
                                                <p className="text-lg font-semibold">${totales.neto.toLocaleString()}</p>
                                            </div>
                                            <div>
                                                <p className="text-sm font-medium text-gray-500">IVA (19%)</p>
                                                <p className="text-lg font-semibold">${totales.iva.toLocaleString()}</p>
                                            </div>
                                            <div>
                                                <p className="text-sm font-medium text-gray-500">Total</p>
                                                <p className="text-lg font-semibold">${totales.total.toLocaleString()}</p>
                                            </div>
                                        </div>
                                    </div>
                                </TabsContent>

                                <TabsContent value="emisor" className="space-y-4">
                                    {empresas.length > 0 && (
                                        <div className="mb-4">
                                            <FormLabel>Seleccionar Empresa</FormLabel>
                                            <Select onValueChange={handleSelectEmpresa}>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Seleccione una empresa" />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    {empresas.map(empresa => (
                                                        <SelectItem key={empresa.id} value={empresa.id}>
                                                            {empresa.razon_social}
                                                        </SelectItem>
                                                    ))}
                                                </SelectContent>
                                            </Select>
                                        </div>
                                    )}

                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                        <FormField
                                            control={form.control}
                                            name="rut_emisor"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>RUT Emisor</FormLabel>
                                                    <FormControl>
                                                        <Input {...field} />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />

                                        <FormField
                                            control={form.control}
                                            name="razon_social_emisor"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Razón Social</FormLabel>
                                                    <FormControl>
                                                        <Input {...field} />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />

                                        <FormField
                                            control={form.control}
                                            name="giro_emisor"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Giro</FormLabel>
                                                    <FormControl>
                                                        <Input {...field} />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                    </div>
                                </TabsContent>

                                <TabsContent value="receptor" className="space-y-4">
                                    {receptores.length > 0 && (
                                        <div className="mb-4">
                                            <FormLabel>Seleccionar Cliente</FormLabel>
                                            <Select onValueChange={handleSelectReceptor}>
                                                <SelectTrigger>
                                                    <SelectValue placeholder="Seleccione un cliente" />
                                                </SelectTrigger>
                                                <SelectContent>
                                                    {receptores.map(receptor => (
                                                        <SelectItem key={receptor.id} value={receptor.id}>
                                                            {receptor.razon_social}
                                                        </SelectItem>
                                                    ))}
                                                </SelectContent>
                                            </Select>
                                        </div>
                                    )}

                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                                        <FormField
                                            control={form.control}
                                            name="rut_receptor"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>RUT Receptor</FormLabel>
                                                    <FormControl>
                                                        <Input {...field} />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />

                                        <FormField
                                            control={form.control}
                                            name="razon_social_receptor"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Razón Social</FormLabel>
                                                    <FormControl>
                                                        <Input {...field} />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />

                                        <FormField
                                            control={form.control}
                                            name="giro_receptor"
                                            render={({ field }) => (
                                                <FormItem>
                                                    <FormLabel>Giro</FormLabel>
                                                    <FormControl>
                                                        <Input {...field} />
                                                    </FormControl>
                                                    <FormMessage />
                                                </FormItem>
                                            )}
                                        />
                                    </div>
                                </TabsContent>

                                <TabsContent value="items" className="space-y-4">
                                    <div className="mb-4 flex justify-end">
                                        <Button type="button" onClick={addItem}>
                                            <Plus className="mr-2 h-4 w-4" />
                                            Agregar Ítem
                                        </Button>
                                    </div>

                                    <div className="border rounded-md">
                                        <Table>
                                            <TableHeader>
                                                <TableRow>
                                                    <TableHead>#</TableHead>
                                                    <TableHead>Descripción</TableHead>
                                                    <TableHead>Cantidad</TableHead>
                                                    <TableHead>Precio Unit.</TableHead>
                                                    <TableHead>Descuento (%)</TableHead>
                                                    <TableHead>Subtotal</TableHead>
                                                    <TableHead></TableHead>
                                                </TableRow>
                                            </TableHeader>
                                            <TableBody>
                                                {fields.map((field, index) => {
                                                    const cantidad = form.watch(`items.${index}.cantidad`) || 0;
                                                    const precio = form.watch(`items.${index}.precio_unitario`) || 0;
                                                    const descuento = form.watch(`items.${index}.descuento`) || 0;
                                                    const subtotal = cantidad * precio * (1 - descuento / 100);

                                                    return (
                                                        <TableRow key={field.id}>
                                                            <TableCell>{index + 1}</TableCell>
                                                            <TableCell>
                                                                <FormField
                                                                    control={form.control}
                                                                    name={`items.${index}.descripcion`}
                                                                    render={({ field }) => (
                                                                        <FormItem>
                                                                            <FormControl>
                                                                                <Input {...field} placeholder="Descripción" />
                                                                            </FormControl>
                                                                            <FormMessage />
                                                                        </FormItem>
                                                                    )}
                                                                />
                                                            </TableCell>
                                                            <TableCell>
                                                                <FormField
                                                                    control={form.control}
                                                                    name={`items.${index}.cantidad`}
                                                                    render={({ field }) => (
                                                                        <FormItem>
                                                                            <FormControl>
                                                                                <Input
                                                                                    {...field}
                                                                                    type="number"
                                                                                    step="0.01"
                                                                                    min="0.01"
                                                                                    className="w-20"
                                                                                />
                                                                            </FormControl>
                                                                            <FormMessage />
                                                                        </FormItem>
                                                                    )}
                                                                />
                                                            </TableCell>
                                                            <TableCell>
                                                                <FormField
                                                                    control={form.control}
                                                                    name={`items.${index}.precio_unitario`}
                                                                    render={({ field }) => (
                                                                        <FormItem>
                                                                            <FormControl>
                                                                                <Input
                                                                                    {...field}
                                                                                    type="number"
                                                                                    min="0"
                                                                                    className="w-24"
                                                                                />
                                                                            </FormControl>
                                                                            <FormMessage />
                                                                        </FormItem>
                                                                    )}
                                                                />
                                                            </TableCell>
                                                            <TableCell>
                                                                <FormField
                                                                    control={form.control}
                                                                    name={`items.${index}.descuento`}
                                                                    render={({ field }) => (
                                                                        <FormItem>
                                                                            <FormControl>
                                                                                <Input
                                                                                    {...field}
                                                                                    type="number"
                                                                                    min="0"
                                                                                    max="100"
                                                                                    className="w-20"
                                                                                />
                                                                            </FormControl>
                                                                            <FormMessage />
                                                                        </FormItem>
                                                                    )}
                                                                />
                                                            </TableCell>
                                                            <TableCell className="font-medium">
                                                                ${subtotal.toLocaleString(undefined, { minimumFractionDigits: 0, maximumFractionDigits: 0 })}
                                                            </TableCell>
                                                            <TableCell>
                                                                <Button
                                                                    type="button"
                                                                    variant="ghost"
                                                                    size="sm"
                                                                    onClick={() => remove(index)}
                                                                    disabled={fields.length === 1}
                                                                >
                                                                    <Trash className="h-4 w-4 text-red-500" />
                                                                </Button>
                                                            </TableCell>
                                                        </TableRow>
                                                    )
                                                })}
                                            </TableBody>
                                        </Table>
                                    </div>
                                </TabsContent>
                            </Tabs>
                        </CardContent>
                        <CardFooter className="flex justify-between pt-6">
                            <Button type="button" variant="outline" onClick={() => router.back()}>
                                Cancelar
                            </Button>
                            <Button type="submit" disabled={isSubmitting}>
                                {isSubmitting ? 'Guardando...' : (
                                    <>
                                        <Check className="mr-2 h-4 w-4" />
                                        {documentoId ? 'Actualizar Documento' : 'Crear Documento'}
                                    </>
                                )}
                            </Button>
                        </CardFooter>
                    </Card>
                </form>
            </Form>
        </div>
    );
};

export default DocumentoForm; 