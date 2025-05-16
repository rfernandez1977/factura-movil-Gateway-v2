package models

// Ambiente representa el ambiente de ejecución (Producción o Certificación)
type Ambiente string

const (
	Produccion    Ambiente = "PRODUCCION"
	Certificacion Ambiente = "CERTIFICACION"
)

// URLs para ambiente de certificación
var (
	URLSemillaCert     = "https://maullin.sii.cl/DTEWS/CrSeed.jws"
	URLTokenCert       = "https://maullin.sii.cl/DTEWS/GetTokenFromSeed.jws"
	URLEnvioDTECert    = "https://maullin.sii.cl/cgi_dte/UPL/DTEUpload"
	URLEstadoDTECert   = "https://maullin.sii.cl/DTEWS/QueryEstDte.jws"
	URLEstadoEnvioCert = "https://maullin.sii.cl/DTEWS/QueryEstUp.jws"
)

// URLs para ambiente de producción
var (
	URLSemillaProd     = "https://palena.sii.cl/DTEWS/CrSeed.jws"
	URLTokenProd       = "https://palena.sii.cl/DTEWS/GetTokenFromSeed.jws"
	URLEnvioDTEProd    = "https://palena.sii.cl/cgi_dte/UPL/DTEUpload"
	URLEstadoDTEProd   = "https://palena.sii.cl/DTEWS/QueryEstDte.jws"
	URLEstadoEnvioProd = "https://palena.sii.cl/DTEWS/QueryEstUp.jws"
)
