# Configuración del Ambiente - FMgo

## Variables de Entorno

### Base de Datos
| Variable | Descripción | Valor por Defecto | Requerido |
|----------|-------------|-------------------|-----------|
| DB_HOST | Host de la base de datos | localhost | Sí |
| DB_PORT | Puerto de la base de datos | 5432 | Sí |
| DB_NAME | Nombre de la base de datos | fmgo_{env} | Sí |
| DB_USER | Usuario de la base de datos | postgres | Sí |
| DB_PASSWORD | Contraseña de la base de datos | - | Sí |

### Redis
| Variable | Descripción | Valor por Defecto | Requerido |
|----------|-------------|-------------------|-----------|
| REDIS_HOST | Host de Redis | localhost | Sí |
| REDIS_PORT | Puerto de Redis | 6379 | Sí |
| REDIS_PASSWORD | Contraseña de Redis | - | No |
| REDIS_DB | Número de base de datos Redis | 0 | No |

### MongoDB
| Variable | Descripción | Valor por Defecto | Requerido |
|----------|-------------|-------------------|-----------|
| MONGO_URI | URI de conexión a MongoDB | mongodb://localhost:27017 | Sí |
| MONGO_DB_NAME | Nombre de la base de datos MongoDB | fmgo_{env} | Sí |

### RabbitMQ
| Variable | Descripción | Valor por Defecto | Requerido |
|----------|-------------|-------------------|-----------|
| RABBITMQ_HOST | Host de RabbitMQ | localhost | Sí |
| RABBITMQ_PORT | Puerto de RabbitMQ | 5672 | Sí |
| RABBITMQ_USER | Usuario de RabbitMQ | guest | Sí |
| RABBITMQ_PASSWORD | Contraseña de RabbitMQ | guest | Sí |

### SII
| Variable | Descripción | Valor por Defecto | Requerido |
|----------|-------------|-------------------|-----------|
| SII_URL | URL del SII | https://palena.sii.cl | Sí |
| SII_CERT_PATH | Ruta al certificado SII | config/certs/sii | Sí |
| SII_KEY_PATH | Ruta a la llave privada SII | config/certs/sii | Sí |

### Firma Digital
| Variable | Descripción | Valor por Defecto | Requerido |
|----------|-------------|-------------------|-----------|
| SIGN_CERT_PATH | Ruta al certificado de firma | config/certs/firma | Sí |
| SIGN_KEY_PATH | Ruta a la llave privada de firma | config/certs/firma | Sí |
| CAF_STORAGE_PATH | Ruta al almacenamiento de CAF | config/caf | Sí |

### Ambiente
| Variable | Descripción | Valor por Defecto | Requerido |
|----------|-------------|-------------------|-----------|
| ENV | Ambiente de ejecución | development | Sí |
| LOG_LEVEL | Nivel de logging | debug | No |
| API_PORT | Puerto de la API | 8080 | Sí |

### Métricas y Monitoreo
| Variable | Descripción | Valor por Defecto | Requerido |
|----------|-------------|-------------------|-----------|
| PROMETHEUS_PORT | Puerto para métricas Prometheus | 9090 | No |
| METRICS_ENABLED | Habilitar métricas | true | No |

## Notas de Configuración

### Ambientes
- **Desarrollo**: Usar prefijo `dev_` para bases de datos
- **Pruebas**: Usar prefijo `test_` para bases de datos
- **Producción**: No usar prefijos en bases de datos

### Seguridad
- No compartir archivos .env entre ambientes
- Mantener respaldos seguros de certificados y llaves
- Rotar contraseñas periódicamente
- Usar valores seguros para puertos en producción

### Certificados
- Almacenar certificados fuera del control de versiones
- Mantener respaldos cifrados
- Documentar fechas de expiración
- Configurar alertas para renovación

### CAF
- Mantener respaldo de archivos CAF
- Documentar proceso de solicitud
- Configurar alertas de uso/disponibilidad
- Implementar rotación automática 