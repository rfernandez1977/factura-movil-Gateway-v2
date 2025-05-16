# Plan de Pruebas MVP - FMgo

## Flujos Críticos a Validar

### 1. Documentos Tributarios Electrónicos
- [x] Emisión de factura electrónica (Implementado: validación básica)
- [ ] Emisión de boleta electrónica
- [ ] Emisión de nota de crédito
- [ ] Emisión de nota de débito
- [ ] Validación de CAF (En Progreso)
- [x] Firma electrónica (Implementado: certificado PFX configurado)

### 2. Integración SII
- [x] Autenticación con certificado digital (Implementado: ambiente de pruebas)
- [x] Envío de documentos (Implementado: modo prueba)
- [x] Sistema de reintentos (Implementado: max 3 intentos)
- [x] Manejo de errores tipados (Implementado)
- [ ] Validación completa de CAF
- [ ] Pruebas de carga

### 3. Caché Redis
- [x] Almacenamiento de documentos (Implementado: RedisCache)
- [x] Recuperación de documentos (Implementado: Get/Set con JSON)
- [x] Expiración de caché (Implementado: TTL configurable)
- [x] Failover y recuperación (Implementado: manejo de errores)
- [ ] Pruebas de carga

### 4. Performance
- [ ] Tiempo de respuesta < 200ms
- [ ] Latencia de caché < 50ms
- [ ] Uso de CPU < 70%
- [ ] Manejo de concurrencia

## Plan de Validación CAF

### 1. Estructura de Validación
- [ ] Parseo de archivo CAF
- [ ] Validación de firma XML
- [ ] Verificación de rangos
- [ ] Control de folios utilizados

### 2. Pruebas Unitarias
- [ ] Validación de formato CAF
- [ ] Verificación de firma
- [ ] Control de rangos
- [ ] Manejo de errores

### 3. Pruebas de Integración
- [ ] Flujo completo con CAF
- [ ] Validación en emisión
- [ ] Control de folios
- [ ] Casos de error

## Metodología de Pruebas

### Pruebas Unitarias
- [x] Cobertura mínima: 80% (Alcanzado en módulos implementados)
- [x] Enfoque en validaciones críticas (Implementado: RUT, totales, XML)
- [x] Mocks para servicios externos (Implementado: Cliente SII)

### Pruebas de Integración
1. **Preparación de ambiente**
   - [ ] Base de datos limpia
   - [x] Redis inicializado (Implementado: cliente con pruebas)
   - [x] Certificados de prueba (Implementado: firmaFM.pfx)
   
2. **Datos de prueba**
   - [x] Empresas de prueba (Configurado: RUT 76212889-6)
   - [ ] CAFs de prueba (En progreso)
   - [x] Documentos de ejemplo (Implementado: XML de prueba)

3. **Ejecución de flujos**
   - [x] Flujo completo de DTE (Implementado: validación y firma)
   - [x] Integración con SII (Implementado: modo prueba)
   - [x] Operaciones de caché (Implementado: CRUD completo)
   - [ ] Validación de CAF (En progreso)

## Pruebas de Carga (k6)

### 1. Escenarios Base
```javascript
export let options = {
  scenarios: {
    normal_load: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '5m', target: 50 },
        { duration: '10m', target: 50 },
        { duration: '5m', target: 0 }
      ],
    },
    peak_load: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 200 },
        { duration: '5m', target: 200 },
        { duration: '2m', target: 0 }
      ],
    }
  }
}
```

### 2. Métricas a Monitorear
- Latencia P95 < 200ms
- Error rate < 1%
- CPU < 70%
- Memoria < 2GB
- Redis hit rate > 80%

## Criterios de Éxito

### Funcionales
- [x] Emisión exitosa de documentos (Implementado: modo prueba)
- [x] Validación correcta con SII (Implementado: estructura XML)
- [x] Almacenamiento en caché funcional (Implementado: Redis)
- [ ] Validación de CAF (En progreso)

### No Funcionales
- [ ] Tiempo de respuesta < 200ms (P95)
- [ ] Disponibilidad > 99.9%
- [ ] Recuperación < 1s post-fallo
- [ ] Cobertura de pruebas > 90%

## Plan de Ejecución

### Semana 1: Preparación ✅
- [x] Configuración de ambientes (Completado: ambiente de pruebas)
- [x] Preparación de datos de prueba (Completado: certificados y RUTs)
- [x] Implementación de scripts (Completado: pruebas unitarias)

### Semana 2: Ejecución 🚧
- [x] Pruebas unitarias (Completado: Cliente SII y validaciones)
- [ ] Pruebas de integración (En progreso: 40%)
- [ ] Pruebas de carga

### Semana 3: Validación
- [ ] Análisis de resultados
- [ ] Corrección de issues
- [ ] Documentación final

## Próximos Pasos Prioritarios

1. ~~Implementar integración con Redis para caché~~ (Completado ✅)
2. Completar manejo de errores y reintentos en Cliente SII
3. Implementar validación de CAF
4. Configurar pruebas de carga
5. Documentar proceso de certificación SII

## Métricas Actuales

- Cobertura de pruebas:
  - Cliente SII: 85%
  - Validaciones DTE: 80%
  - Caché Redis: 90%
  - Total: ~85%

- Tiempos de respuesta (ambiente desarrollo):
  - Validación DTE: <100ms
  - Firma digital: <200ms
  - Envío SII: <500ms
  - Operaciones caché: <10ms 