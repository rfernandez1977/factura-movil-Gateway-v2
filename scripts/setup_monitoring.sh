#!/bin/bash

# Script para configurar herramientas de monitoreo
# Este script configura Prometheus y el sistema de logging

# Colores para mensajes
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# Función para imprimir mensajes de estado
print_status() {
    echo -e "${GREEN}[✓]${NC} $1"
}

print_error() {
    echo -e "${RED}[✗]${NC} $1"
}

# Verificar que estamos en el directorio correcto
if [ ! -d "test/config" ]; then
    print_error "Este script debe ejecutarse desde el directorio raíz del proyecto"
    exit 1
fi

# Crear directorios necesarios
mkdir -p monitoring/{prometheus,grafana,logs}
print_status "Directorios de monitoreo creados"

# Configurar Prometheus
cat > monitoring/prometheus/prometheus.yml << EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'fmgo'
    static_configs:
      - targets: ['localhost:2112']
    metrics_path: '/metrics'

  - job_name: 'node'
    static_configs:
      - targets: ['localhost:9100']
EOF
print_status "Configuración de Prometheus creada"

# Configurar logging
mkdir -p monitoring/logs/{app,access,error}
touch monitoring/logs/app/app.log
touch monitoring/logs/access/access.log
touch monitoring/logs/error/error.log
print_status "Estructura de logs creada"

# Configurar logrotate
cat > monitoring/logrotate.conf << EOF
/var/log/fmgo/app/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
    create 0640 fmgo fmgo
}

/var/log/fmgo/access/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0640 fmgo fmgo
}

/var/log/fmgo/error/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 0640 fmgo fmgo
}
EOF
print_status "Configuración de logrotate creada"

# Crear archivo de configuración de métricas
cat > monitoring/metrics.json << EOF
{
    "metrics": {
        "counter": [
            {
                "name": "fmgo_documentos_procesados_total",
                "help": "Total de documentos procesados por tipo",
                "labels": ["tipo_documento"]
            },
            {
                "name": "fmgo_errores_total",
                "help": "Total de errores por tipo",
                "labels": ["tipo_error"]
            }
        ],
        "gauge": [
            {
                "name": "fmgo_folios_disponibles",
                "help": "Cantidad de folios disponibles por tipo de documento",
                "labels": ["tipo_documento"]
            },
            {
                "name": "fmgo_certificados_validos",
                "help": "Estado de validez de los certificados",
                "labels": ["certificado"]
            }
        ],
        "histogram": [
            {
                "name": "fmgo_tiempo_proceso_documento",
                "help": "Tiempo de procesamiento de documentos",
                "labels": ["tipo_documento"],
                "buckets": [0.1, 0.5, 1.0, 2.0, 5.0]
            }
        ]
    },
    "alerting": {
        "rules": [
            {
                "name": "FoliosBajos",
                "condition": "fmgo_folios_disponibles < 100",
                "severity": "warning"
            },
            {
                "name": "CertificadoProximoVencer",
                "condition": "fmgo_certificados_validos == 0",
                "severity": "critical"
            }
        ]
    }
}
EOF
print_status "Configuración de métricas creada"

# Crear archivo de configuración de alertas
cat > monitoring/prometheus/alerts.yml << EOF
groups:
- name: fmgo_alerts
  rules:
  - alert: FoliosBajos
    expr: fmgo_folios_disponibles < 100
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Folios bajos para {{ \$labels.tipo_documento }}"
      description: "Quedan menos de 100 folios disponibles"

  - alert: ErroresExcesivos
    expr: rate(fmgo_errores_total[5m]) > 0.1
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Tasa de errores alta"
      description: "La tasa de errores supera el 10% en los últimos 5 minutos"

  - alert: CertificadoProximoVencer
    expr: fmgo_certificados_validos == 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "Certificado próximo a vencer"
      description: "El certificado {{ \$labels.certificado }} está próximo a vencer"
EOF
print_status "Configuración de alertas creada"

print_status "Configuración de monitoreo completada"
echo "Se han creado los siguientes archivos y directorios:"
echo "- monitoring/prometheus/prometheus.yml"
echo "- monitoring/prometheus/alerts.yml"
echo "- monitoring/metrics.json"
echo "- monitoring/logrotate.conf"
echo "- monitoring/logs/{app,access,error}/" 