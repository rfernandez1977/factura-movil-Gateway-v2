package metrics

import (
	"github.com/cursor/FMgo/domain"
	"github.com/cursor/FMgo/models"
	"github.com/prometheus/client_golang/prometheus"
)

// MetricasDocumento contiene las métricas para documentos tributarios
type MetricasDocumento struct {
	MontoNeto                      prometheus.Histogram
	MontoIVA                       prometheus.Histogram
	MontoTotal                     prometheus.Histogram
	MontoExento                    prometheus.Histogram
	ImpuestosAdicionales           *prometheus.HistogramVec
	ImpuestosAdicionalesBase       *prometheus.HistogramVec
	ImpuestosAdicionalesPorcentaje *prometheus.HistogramVec
}

// TributarioMetrics representa el servicio para recopilar métricas de documentos tributarios
type TributarioMetrics struct {
	metricas *MetricasDocumento
}

// NewTributarioMetrics crea una nueva instancia del servicio TributarioMetrics
func NewTributarioMetrics(reg prometheus.Registerer) *TributarioMetrics {
	metricas := &MetricasDocumento{
		MontoNeto: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "documento_tributario_monto_neto",
			Help:    "Monto neto de documentos tributarios",
			Buckets: prometheus.LinearBuckets(1000, 10000, 10),
		}),
		MontoIVA: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "documento_tributario_monto_iva",
			Help:    "Monto IVA de documentos tributarios",
			Buckets: prometheus.LinearBuckets(190, 1900, 10),
		}),
		MontoTotal: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "documento_tributario_monto_total",
			Help:    "Monto total de documentos tributarios",
			Buckets: prometheus.LinearBuckets(1190, 11900, 10),
		}),
		MontoExento: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "documento_tributario_monto_exento",
			Help:    "Monto exento de documentos tributarios",
			Buckets: prometheus.LinearBuckets(1000, 10000, 10),
		}),
		ImpuestosAdicionales: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "documento_tributario_impuestos_adicionales_monto",
				Help:    "Monto de impuestos adicionales",
				Buckets: prometheus.LinearBuckets(100, 1000, 10),
			},
			[]string{"codigo", "nombre"},
		),
		ImpuestosAdicionalesBase: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "documento_tributario_impuestos_adicionales_base",
				Help:    "Base imponible de impuestos adicionales",
				Buckets: prometheus.LinearBuckets(1000, 10000, 10),
			},
			[]string{"codigo", "nombre"},
		),
		ImpuestosAdicionalesPorcentaje: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "documento_tributario_impuestos_adicionales_porcentaje",
				Help:    "Porcentaje de impuestos adicionales",
				Buckets: prometheus.LinearBuckets(1, 50, 10),
			},
			[]string{"codigo", "nombre"},
		),
	}

	// Registrar métricas
	reg.MustRegister(
		metricas.MontoNeto,
		metricas.MontoIVA,
		metricas.MontoTotal,
		metricas.MontoExento,
		metricas.ImpuestosAdicionales,
		metricas.ImpuestosAdicionalesBase,
		metricas.ImpuestosAdicionalesPorcentaje,
	)

	return &TributarioMetrics{
		metricas: metricas,
	}
}

// registrarMetricasImpuestos registra las métricas de impuestos
func (m *TributarioMetrics) registrarMetricasImpuestos(doc interface{}) {
	var (
		montoNeto            float64
		montoIVA             float64
		montoTotal           float64
		montoExento          float64
		impuestosAdicionales []models.ImpuestoAdicional
	)

	switch d := doc.(type) {
	case *models.Factura:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Obtener impuestos adicionales de los items
		for _, item := range d.Items {
			for _, impAdicional := range item.ImpuestosAdicionales {
				impuestosAdicionales = append(impuestosAdicionales, models.ImpuestoAdicional{
					Codigo:        impAdicional.Codigo,
					Nombre:        impAdicional.Nombre,
					Porcentaje:    impAdicional.Porcentaje,
					Monto:         impAdicional.Monto,
					BaseImponible: impAdicional.BaseImponible,
				})
			}
		}
	case *models.Boleta:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Obtener impuestos adicionales de los items
		for _, item := range d.Items {
			for _, impAdicional := range item.ImpuestosAdicionales {
				impuestosAdicionales = append(impuestosAdicionales, models.ImpuestoAdicional{
					Codigo:        impAdicional.Codigo,
					Nombre:        impAdicional.Nombre,
					Porcentaje:    impAdicional.Porcentaje,
					Monto:         impAdicional.Monto,
					BaseImponible: impAdicional.BaseImponible,
				})
			}
		}
	case *models.NotaCredito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Obtener impuestos adicionales de los items
		for _, item := range d.Items {
			for _, impAdicional := range item.ImpuestosAdicionales {
				impuestosAdicionales = append(impuestosAdicionales, models.ImpuestoAdicional{
					Codigo:        impAdicional.Codigo,
					Nombre:        impAdicional.Nombre,
					Porcentaje:    impAdicional.Porcentaje,
					Monto:         impAdicional.Monto,
					BaseImponible: impAdicional.BaseImponible,
				})
			}
		}
	case *models.NotaDebito:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Obtener impuestos adicionales de los items
		for _, item := range d.Items {
			for _, impAdicional := range item.ImpuestosAdicionales {
				impuestosAdicionales = append(impuestosAdicionales, models.ImpuestoAdicional{
					Codigo:        impAdicional.Codigo,
					Nombre:        impAdicional.Nombre,
					Porcentaje:    impAdicional.Porcentaje,
					Monto:         impAdicional.Monto,
					BaseImponible: impAdicional.BaseImponible,
				})
			}
		}
	case *models.GuiaDespacho:
		montoNeto = d.MontoNeto
		montoIVA = d.MontoIVA
		montoTotal = d.MontoTotal
		montoExento = d.MontoExento
		// Obtener impuestos adicionales de los items
		for _, item := range d.Items {
			for _, impAdicional := range item.ImpuestosAdicionales {
				impuestosAdicionales = append(impuestosAdicionales, models.ImpuestoAdicional{
					Codigo:        impAdicional.Codigo,
					Nombre:        impAdicional.Nombre,
					Porcentaje:    impAdicional.Porcentaje,
					Monto:         impAdicional.Monto,
					BaseImponible: impAdicional.BaseImponible,
				})
			}
		}
	}

	// Registrar métricas de impuestos
	m.metricas.MontoNeto.Observe(montoNeto)
	m.metricas.MontoIVA.Observe(montoIVA)
	m.metricas.MontoExento.Observe(montoExento)
	m.metricas.MontoTotal.Observe(montoTotal)

	// Registrar métricas de impuestos adicionales
	for _, impuesto := range impuestosAdicionales {
		m.metricas.ImpuestosAdicionales.WithLabelValues(impuesto.Codigo, impuesto.Nombre).Observe(impuesto.Monto)
		m.metricas.ImpuestosAdicionalesBase.WithLabelValues(impuesto.Codigo, impuesto.Nombre).Observe(impuesto.BaseImponible)
		m.metricas.ImpuestosAdicionalesPorcentaje.WithLabelValues(impuesto.Codigo, impuesto.Nombre).Observe(impuesto.Porcentaje)
	}
}

// CollectMetrics recopila métricas de un documento tributario
func (m *TributarioMetrics) CollectMetrics(doc interface{}) error {
	// Usar la función interna para registrar métricas
	m.registrarMetricasImpuestos(doc)

	// Otras métricas específicas por tipo de documento
	switch d := doc.(type) {
	case *models.Factura:
		// Métricas específicas de factura
		if d.FormaPago != "" {
			// Registrar forma de pago (se podría agregar a las métricas)
		}
	case *models.Boleta:
		// Métricas específicas de boleta
		if d.MedioPago != "" {
			// Registrar medio de pago (se podría agregar a las métricas)
		}
	case *models.NotaCredito:
		// Métricas específicas de nota de crédito
		if d.TipoReferencia != "" {
			// Registrar tipo de referencia (se podría agregar a las métricas)
		}
	case *models.NotaDebito:
		// Métricas específicas de nota de débito
		if d.TipoReferencia != "" {
			// Registrar tipo de referencia (se podría agregar a las métricas)
		}
	case *models.GuiaDespacho:
		// Métricas específicas de guía de despacho
		if d.RutTransportista != "" {
			// Registrar información de transporte (se podría agregar a las métricas)
		}
	}

	return nil
}

// ObtenerImpuestosAdicionales obtiene los impuestos adicionales de los items de tipo domain.Item
func (s *TributarioMetrics) ObtenerImpuestosAdicionales(items []domain.Item) []models.ImpuestoAdicional {
	var impuestos []models.ImpuestoAdicional
	impuestosMap := make(map[string]models.ImpuestoAdicional)

	// En el caso de domain.Item, no tenemos el campo ImpuestosAdicionales
	// Por lo tanto, devolvemos un slice vacío
	return impuestos
}

// ObtenerImpuestosAdicionalesBoleta obtiene los impuestos adicionales de los items de tipo DetalleBoleta
func (s *TributarioMetrics) ObtenerImpuestosAdicionalesBoleta(items []*models.DetalleBoleta) []models.ImpuestoAdicional {
	var impuestos []models.ImpuestoAdicional
	impuestosMap := make(map[string]models.ImpuestoAdicional)

	// En el caso de DetalleBoleta, no tenemos el campo ImpuestosAdicionales
	// Por lo tanto, devolvemos un slice vacío
	return impuestos
}

// AgregarImpuestoAdicionalesItems agrega los impuestos adicionales a las métricas
func (s *TributarioMetrics) AgregarImpuestoAdicionalesItems(metric *map[string]int64, montos *map[string]float64, items []domain.Item) {
	/*
		// Este código causa errores porque domain.Item no tiene el campo ImpuestosAdicionales
		for _, item := range items {
			for _, impuesto := range item.ImpuestosAdicionales {
				key := fmt.Sprintf("impuesto_adicional_%s", impuesto.Codigo)
				if _, ok := (*metric)[key]; !ok {
					(*metric)[key] = 0
					(*montos)[key] = 0
				}
				(*metric)[key]++
				(*montos)[key] += impuesto.Monto
			}
		}
	*/
}

// AgregarImpuestoAdicionalesBoleta agrega los impuestos adicionales de una boleta a las métricas
func (s *TributarioMetrics) AgregarImpuestoAdicionalesBoleta(metric *map[string]int64, montos *map[string]float64, items []*models.DetalleBoleta) {
	/*
		// Este código causa errores porque DetalleBoleta no tiene el campo ImpuestosAdicionales
		for _, item := range items {
			for _, impuesto := range item.ImpuestosAdicionales {
				key := fmt.Sprintf("impuesto_adicional_%s", impuesto.Codigo)
				if _, ok := (*metric)[key]; !ok {
					(*metric)[key] = 0
					(*montos)[key] = 0
				}
				(*metric)[key]++
				(*montos)[key] += impuesto.Monto
			}
		}
	*/
}
