/**
 * Formatea una fecha a formato DD/MM/YYYY
 */
export const formatDate = (dateString?: string): string => {
    if (!dateString) return '';

    const date = new Date(dateString);
    if (isNaN(date.getTime())) return '';

    return date.toLocaleDateString('es-CL', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
    });
};

/**
 * Formatea un número a moneda chilena
 */
export const formatCurrency = (amount?: number): string => {
    if (amount === undefined || amount === null) return '$0';

    return new Intl.NumberFormat('es-CL', {
        style: 'currency',
        currency: 'CLP',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0
    }).format(amount);
};

/**
 * Formatea un RUT chileno (12345678-9)
 */
export const formatRut = (rut?: string): string => {
    if (!rut) return '';

    // Eliminar puntos y guiones
    let valor = rut.replace(/\./g, '').replace(/-/g, '');

    // Obtener el dígito verificador
    const dv = valor.slice(-1);

    // Obtener el cuerpo del RUT
    const rutBody = valor.slice(0, -1);

    // Formatear con puntos y guión
    let rutFormateado = '';
    let i = rutBody.length;

    while (i > 0) {
        const inicio = Math.max(i - 3, 0);
        rutFormateado = rutBody.substring(inicio, i) + (rutFormateado ? '.' + rutFormateado : '');
        i = inicio;
    }

    return rutFormateado + '-' + dv;
};

/**
 * Valida un RUT chileno
 */
export const validarRut = (rut: string): boolean => {
    // Eliminar puntos y guiones
    const rutLimpio = rut.replace(/\./g, '').replace(/-/g, '');

    if (rutLimpio.length < 2) return false;

    // Obtener el dígito verificador
    const dv = rutLimpio.charAt(rutLimpio.length - 1);

    // Obtener el cuerpo del RUT
    const rutBody = rutLimpio.slice(0, -1);

    // Calcular dígito verificador
    let suma = 0;
    let multiplo = 2;

    // Para cada dígito del RUT
    for (let i = rutBody.length - 1; i >= 0; i--) {
        suma += Number(rutBody.charAt(i)) * multiplo;
        multiplo = multiplo < 7 ? multiplo + 1 : 2;
    }

    let dvEsperado: string | number = 11 - (suma % 11);

    // Si el dígito verificador es 11, corresponde a 0
    if (dvEsperado === 11) dvEsperado = 0;

    // Si el dígito verificador es 10, corresponde a K
    if (dvEsperado === 10) dvEsperado = 'K';

    // Comparar dígito verificador
    return dv.toUpperCase() === String(dvEsperado);
};

/**
 * Formatea un número de teléfono chileno
 */
export const formatPhone = (phone?: string): string => {
    if (!phone) return '';

    // Eliminar espacios y caracteres no numéricos
    const numeroLimpio = phone.replace(/\D/g, '');

    // Si empieza con +56 o 56, formatear con código de país
    if (numeroLimpio.startsWith('56')) {
        const numero = numeroLimpio.substring(2);
        if (numero.length === 9) {
            return `+56 ${numero.substring(0, 1)} ${numero.substring(1, 5)} ${numero.substring(5)}`;
        }
    }

    // Si tiene 9 dígitos, formatear como número chileno
    if (numeroLimpio.length === 9) {
        return `+56 ${numeroLimpio.substring(0, 1)} ${numeroLimpio.substring(1, 5)} ${numeroLimpio.substring(5)}`;
    }

    // Si tiene 8 dígitos, asumir que falta el 9 inicial
    if (numeroLimpio.length === 8) {
        return `+56 9 ${numeroLimpio.substring(0, 4)} ${numeroLimpio.substring(4)}`;
    }

    // En otro caso, devolver como está
    return phone;
}; 