package models

import "time"

// Configuracion representa la configuración de una empresa
type Configuracion struct {
	ID          string                  `json:"id" bson:"_id,omitempty"`
	Clave       string                  `json:"clave" bson:"clave"`
	EmpresaID   string                  `json:"empresa_id" bson:"empresa_id"`
	RUT         string                  `json:"rut" bson:"rut"`
	Nombre      string                  `json:"nombre" bson:"nombre"`
	ConfigSII   ConfiguracionSIIEmpresa `json:"config_sii" bson:"config_sii"`
	ConfigEmail ConfiguracionEmail      `json:"config_email" bson:"config_email"`
	ConfigAuth  ConfiguracionAuth       `json:"config_auth" bson:"config_auth"`
	ConfigApp   ConfiguracionApp        `json:"config_app" bson:"config_app"`
	ConfigDocs  ConfiguracionDocs       `json:"config_docs" bson:"config_docs"`
	Activo      bool                    `json:"activo" bson:"activo"`
	CreatedAt   time.Time               `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at" bson:"updated_at"`
}

// ConfiguracionSIIEmpresa representa la configuración del SII para una empresa
type ConfiguracionSIIEmpresa struct {
	EmpresaID                 string    `json:"empresa_id" bson:"empresa_id"`
	RUT                       string    `json:"rut" bson:"rut"`
	AmbienteSII               string    `json:"ambiente_sii" bson:"ambiente_sii"` // Produccion, Certificacion
	URLProduccion             string    `json:"url_produccion" bson:"url_produccion"`
	URLCertificacion          string    `json:"url_certificacion" bson:"url_certificacion"`
	CertificadoDigital        string    `json:"certificado_digital" bson:"certificado_digital"`
	ClavePrivada              string    `json:"clave_privada,omitempty" bson:"clave_privada,omitempty"`
	ResolucionNumero          int       `json:"resolucion_numero" bson:"resolucion_numero"`
	ResolucionFecha           time.Time `json:"resolucion_fecha" bson:"resolucion_fecha"`
	CAFS                      []string  `json:"cafs,omitempty" bson:"cafs,omitempty"`
	ConsultaAutomaticaDTEs    bool      `json:"consulta_automatica_dtes" bson:"consulta_automatica_dtes"`
	IntervaloConsultaDTEs     int       `json:"intervalo_consulta_dtes" bson:"intervalo_consulta_dtes"`
	IVA                       float64   `json:"iva" bson:"iva"`
	ValidarDTEAutomatico      bool      `json:"validar_dte_automatico" bson:"validar_dte_automatico"`
	EnviarDTEAutomatico       bool      `json:"enviar_dte_automatico" bson:"enviar_dte_automatico"`
	TimbreDTEAutomatico       bool      `json:"timbre_dte_automatico" bson:"timbre_dte_automatico"`
	EnviarDTECorreo           bool      `json:"enviar_dte_correo" bson:"enviar_dte_correo"`
	NextNumeroFactura         int64     `json:"next_numero_factura" bson:"next_numero_factura"`
	NextNumeroBoleta          int64     `json:"next_numero_boleta" bson:"next_numero_boleta"`
	NextNumeroNotaCredito     int64     `json:"next_numero_nota_credito" bson:"next_numero_nota_credito"`
	NextNumeroNotaDebito      int64     `json:"next_numero_nota_debito" bson:"next_numero_nota_debito"`
	NextNumeroGuiaDespacho    int64     `json:"next_numero_guia_despacho" bson:"next_numero_guia_despacho"`
	NextNumeroFacturaExenta   int64     `json:"next_numero_factura_exenta" bson:"next_numero_factura_exenta"`
	NextNumeroFacturaCompra   int64     `json:"next_numero_factura_compra" bson:"next_numero_factura_compra"`
	UltimaConsultaDTEs        time.Time `json:"ultima_consulta_dtes" bson:"ultima_consulta_dtes"`
	UltimaActualizacionFolios time.Time `json:"ultima_actualizacion_folios" bson:"ultima_actualizacion_folios"`
}

// ConfiguracionEmail representa la configuración de correo
type ConfiguracionEmail struct {
	EmpresaID       string            `json:"empresa_id" bson:"empresa_id"`
	SMTPServer      string            `json:"smtp_server" bson:"smtp_server"`
	SMTPPort        int               `json:"smtp_port" bson:"smtp_port"`
	SMTPUser        string            `json:"smtp_user" bson:"smtp_user"`
	SMTPPassword    string            `json:"smtp_password,omitempty" bson:"smtp_password,omitempty"`
	FromEmail       string            `json:"from_email" bson:"from_email"`
	FromName        string            `json:"from_name" bson:"from_name"`
	ReplyTo         string            `json:"reply_to" bson:"reply_to"`
	EnviarCopiaA    []string          `json:"enviar_copia_a" bson:"enviar_copia_a"`
	FirmaHTML       string            `json:"firma_html" bson:"firma_html"`
	UsarTLS         bool              `json:"usar_tls" bson:"usar_tls"`
	UsarSSL         bool              `json:"usar_ssl" bson:"usar_ssl"`
	EnviarPDF       bool              `json:"enviar_pdf" bson:"enviar_pdf"`
	EnviarXML       bool              `json:"enviar_xml" bson:"enviar_xml"`
	Plantillas      map[string]string `json:"plantillas" bson:"plantillas"`
	EnvioAutomatico bool              `json:"envio_automatico" bson:"envio_automatico"`
	IncluyeLogotipo bool              `json:"incluye_logotipo" bson:"incluye_logotipo"`
	ColorPrimario   string            `json:"color_primario" bson:"color_primario"`
	ColorSecundario string            `json:"color_secundario" bson:"color_secundario"`
	LogotipoBase64  string            `json:"logotipo_base64" bson:"logotipo_base64"`
}

// ConfiguracionAuth representa la configuración de autenticación
type ConfiguracionAuth struct {
	EmpresaID             string `json:"empresa_id" bson:"empresa_id"`
	ProveedorAuth         string `json:"proveedor_auth" bson:"proveedor_auth"` // JWT, OAuth2, OIDC
	DuracionTokenMinutos  int    `json:"duracion_token_minutos" bson:"duracion_token_minutos"`
	RequiereConfirmacion  bool   `json:"requiere_confirmacion" bson:"requiere_confirmacion"`
	PoliticaClaves        string `json:"politica_claves" bson:"politica_claves"`
	MaxIntentosLogin      int    `json:"max_intentos_login" bson:"max_intentos_login"`
	TiempoBloqueoMinutos  int    `json:"tiempo_bloqueo_minutos" bson:"tiempo_bloqueo_minutos"`
	PermitirLoginMultiple bool   `json:"permitir_login_multiple" bson:"permitir_login_multiple"`
	RequiereMultiFactor   bool   `json:"requiere_multi_factor" bson:"requiere_multi_factor"`
	ClaveSecretaJWT       string `json:"clave_secreta_jwt,omitempty" bson:"clave_secreta_jwt,omitempty"`
	OAuth2URL             string `json:"oauth2_url" bson:"oauth2_url"`
	OAuth2ClientID        string `json:"oauth2_client_id" bson:"oauth2_client_id"`
	OAuth2ClientSecret    string `json:"oauth2_client_secret,omitempty" bson:"oauth2_client_secret,omitempty"`
	OAuth2Scope           string `json:"oauth2_scope" bson:"oauth2_scope"`
}

// ConfiguracionApp representa la configuración de la aplicación
type ConfiguracionApp struct {
	EmpresaID                string            `json:"empresa_id" bson:"empresa_id"`
	TemaApp                  string            `json:"tema_app" bson:"tema_app"`
	ColorPrimario            string            `json:"color_primario" bson:"color_primario"`
	ColorSecundario          string            `json:"color_secundario" bson:"color_secundario"`
	LogotipoBase64           string            `json:"logotipo_base64" bson:"logotipo_base64"`
	LogotipoURLPublica       string            `json:"logotipo_url_publica" bson:"logotipo_url_publica"`
	MostrarMensajes          bool              `json:"mostrar_mensajes" bson:"mostrar_mensajes"`
	ModoOscuro               bool              `json:"modo_oscuro" bson:"modo_oscuro"`
	URLBaseApp               string            `json:"url_base_app" bson:"url_base_app"`
	TiempoSesionMinutos      int               `json:"tiempo_sesion_minutos" bson:"tiempo_sesion_minutos"`
	ModulosActivos           []string          `json:"modulos_activos" bson:"modulos_activos"`
	ConfiguracionesAvanzadas map[string]string `json:"configuraciones_avanzadas" bson:"configuraciones_avanzadas"`
	Languaje                 string            `json:"languaje" bson:"languaje"`
	MostrarIdiomas           []string          `json:"mostrar_idiomas" bson:"mostrar_idiomas"`
}

// ConfiguracionDocs representa la configuración de documentos
type ConfiguracionDocs struct {
	EmpresaID                   string            `json:"empresa_id" bson:"empresa_id"`
	FormatoFolio                string            `json:"formato_folio" bson:"formato_folio"`
	MostrarIVA                  bool              `json:"mostrar_iva" bson:"mostrar_iva"`
	MostrarMontoExento          bool              `json:"mostrar_monto_exento" bson:"mostrar_monto_exento"`
	MostrarMontoAfecto          bool              `json:"mostrar_monto_afecto" bson:"mostrar_monto_afecto"`
	MostrarMontoTotal           bool              `json:"mostrar_monto_total" bson:"mostrar_monto_total"`
	MostrarDetalleImpuestos     bool              `json:"mostrar_detalle_impuestos" bson:"mostrar_detalle_impuestos"`
	MostrarCodigoItem           bool              `json:"mostrar_codigo_item" bson:"mostrar_codigo_item"`
	MostrarDescripcionDetallada bool              `json:"mostrar_descripcion_detallada" bson:"mostrar_descripcion_detallada"`
	MostrarDescuentos           bool              `json:"mostrar_descuentos" bson:"mostrar_descuentos"`
	MostrarRecargos             bool              `json:"mostrar_recargos" bson:"mostrar_recargos"`
	PlantillaFactura            string            `json:"plantilla_factura" bson:"plantilla_factura"`
	PlantillaBoleta             string            `json:"plantilla_boleta" bson:"plantilla_boleta"`
	PlantillaNotaCredito        string            `json:"plantilla_nota_credito" bson:"plantilla_nota_credito"`
	PlantillaNotaDebito         string            `json:"plantilla_nota_debito" bson:"plantilla_nota_debito"`
	PlantillaGuiaDespacho       string            `json:"plantilla_guia_despacho" bson:"plantilla_guia_despacho"`
	TextosAdicionales           map[string]string `json:"textos_adicionales" bson:"textos_adicionales"`
	MostrarTimbreSII            bool              `json:"mostrar_timbre_sii" bson:"mostrar_timbre_sii"`
	MostrarTotalesEspecificos   bool              `json:"mostrar_totales_especificos" bson:"mostrar_totales_especificos"`
	MostrarCantidadConDecimales bool              `json:"mostrar_cantidad_con_decimales" bson:"mostrar_cantidad_con_decimales"`
}
