import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Métricas personalizadas
const errors = new Rate('errors');

// Configuración de la prueba
export const options = {
    stages: [
        { duration: '5m', target: 50 },   // Ramp-up a 50 RPS
        { duration: '10m', target: 100 }, // Ramp-up a 100 RPS
        { duration: '30m', target: 100 }, // Mantener 100 RPS
        { duration: '5m', target: 0 },    // Ramp-down a 0
    ],
    thresholds: {
        'http_req_duration': ['p(95)<200'], // 95% de requests bajo 200ms
        'errors': ['rate<0.01'],            // Error rate menor al 1%
    },
};

// Datos de prueba
const testData = {
    emisor: {
        rut: '76123456-7',
        razon_social: 'EMPRESA DE PRUEBA SPA',
        giro: 'SERVICIOS INFORMATICOS',
        direccion: 'CALLE EJEMPLO 123',
        comuna: 'SANTIAGO'
    },
    receptor: {
        rut: '77654321-8',
        razon_social: 'CLIENTE DE PRUEBA LTDA',
        giro: 'COMERCIO',
        direccion: 'AV CLIENTE 456',
        comuna: 'PROVIDENCIA'
    }
};

// Función principal
export default function () {
    // Headers comunes
    const headers = {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer test-token'
    };

    // Grupo de pruebas: Emisión de DTE
    const dtePayload = {
        tipo_dte: '33',
        emisor: testData.emisor,
        receptor: testData.receptor,
        detalles: [
            {
                cantidad: 1,
                descripcion: 'Servicio Profesional',
                precio_unitario: 100000,
                monto_total: 100000
            }
        ],
        totales: {
            monto_neto: 100000,
            tasa_iva: 19,
            iva: 19000,
            total: 119000
        }
    };

    // Emitir DTE
    const dteResponse = http.post('http://localhost:8080/api/v1/dte',
        JSON.stringify(dtePayload),
        { headers: headers }
    );

    // Verificar respuesta
    check(dteResponse, {
        'status is 200': (r) => r.status === 200,
        'response has id': (r) => JSON.parse(r.body).id !== undefined,
    }) || errors.add(1);

    if (dteResponse.status === 200) {
        const dteId = JSON.parse(dteResponse.body).id;

        // Consultar estado del DTE
        const statusResponse = http.get(
            `http://localhost:8080/api/v1/dte/${dteId}/estado`,
            { headers: headers }
        );

        // Verificar respuesta de estado
        check(statusResponse, {
            'status check is 200': (r) => r.status === 200,
            'has valid state': (r) => {
                const body = JSON.parse(r.body);
                return ['PENDIENTE', 'PROCESANDO', 'ACEPTADO'].includes(body.estado);
            },
        }) || errors.add(1);
    }

    // Pausa entre iteraciones (distribución normal entre 1s y 2s)
    sleep(Math.random() + 1);
} 