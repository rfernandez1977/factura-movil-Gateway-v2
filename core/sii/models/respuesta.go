package models

import "encoding/xml"

// RespuestaSII es la estructura base para las respuestas del SII
type RespuestaSII struct {
	XMLName xml.Name `xml:"RESPUESTA"`
	Header  RespuestaHeader `xml:"RESP_HDR"`
	Body    RespuestaBody   `xml:"RESP_BODY"`
}

// RespuestaHeader contiene la información de estado de la respuesta
type RespuestaHeader struct {
	Estado string `xml:"ESTADO"`
	Glosa  string `xml:"GLOSA"`
}

// RespuestaBody contiene el cuerpo de la respuesta
type RespuestaBody struct {
	Semilla string `xml:"SEMILLA,omitempty"`
	Token   string `xml:"TOKEN,omitempty"`
	TrackID string `xml:"TRACKID,omitempty"`
}

// RespuestaEstado contiene la información de estado de un documento
type RespuestaEstado struct {
	XMLName    xml.Name `xml:"RESPUESTA"`
	Estado     string   `xml:"ESTADO"`
	Glosa      string   `xml:"GLOSA"`
	NumTrack   string   `xml:"NUMTRACK,omitempty"`
	RutEmisor  string   `xml:"RUTEMISOR,omitempty"`
	RutReceptor string `xml:"RUTRECEPTOR,omitempty"`
	TipoDTE    string   `xml:"TIPODTE,omitempty"`
	FolioDTE   string   `xml:"FOLIODTE,omitempty"`
	FechaEmision string `xml:"FECHAEMISION,omitempty"`
	MontoTotal  string   `xml:"MONTOTOTAL,omitempty"`
}

// SoapEnvelope es la estructura para mensajes SOAP
type SoapEnvelope struct {
	XMLName xml.Name    `xml:"Envelope"`
	XMLNS   string      `xml:"xmlns,attr"`
	Body    SoapBody    `xml:"Body"`
}

// SoapBody es el cuerpo del mensaje SOAP
type SoapBody struct {
	XMLName xml.Name    `xml:"Body"`
	Content interface{} `xml:",any"`
	Fault   *SoapFault  `xml:"Fault,omitempty"`
}

// SoapFault representa un error en el mensaje SOAP
type SoapFault struct {
	XMLName     xml.Name `xml:"Fault"`
	FaultCode   string   `xml:"faultcode"`
	FaultString string   `xml:"faultstring"`
} 