package models

// TipoERP define los diferentes tipos de ERP soportados
type TipoERP string

const (
	// ERP_SAP representa el ERP SAP
	ERP_SAP TipoERP = "SAP"

	// ERP_ORACLE representa el ERP Oracle
	ERP_ORACLE TipoERP = "ORACLE"

	// ERP_DYNAMICS representa el ERP Microsoft Dynamics
	ERP_DYNAMICS TipoERP = "DYNAMICS"

	// ERP_NETSUITE representa el ERP Oracle NetSuite
	ERP_NETSUITE TipoERP = "NETSUITE"

	// ERP_SAGE representa el ERP Sage
	ERP_SAGE TipoERP = "SAGE"

	// ERP_SOFTLAND representa el ERP Softland
	ERP_SOFTLAND TipoERP = "SOFTLAND"

	// ERP_PERSONALIZADO representa un ERP con integración personalizada
	ERP_PERSONALIZADO TipoERP = "PERSONALIZADO"

	// ERP_LEGACY representa un ERP legado
	ERP_LEGACY TipoERP = "LEGACY"
)
