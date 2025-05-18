package parser

import (
	"encoding/xml"
	"fmt"

	"FMgo/core/sii/models"
)

// ResponseParser parsea las respuestas XML del SII
type ResponseParser struct{}

// NewResponseParser crea una nueva instancia del parser
func NewResponseParser() *ResponseParser {
	return &ResponseParser{}
}

// ParseSemilla parsea la respuesta de semilla
func (p *ResponseParser) ParseSemilla(data []byte) (string, error) {
	var resp struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			XMLName     xml.Name `xml:"Body"`
			GetSeedResp struct {
				XMLName xml.Name `xml:"getSeedResponse"`
				Seed    string   `xml:"seed"`
			}
		}
	}

	if err := xml.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("error al parsear respuesta de semilla: %w", err)
	}

	return resp.Body.GetSeedResp.Seed, nil
}

// ParseToken parsea la respuesta de token
func (p *ResponseParser) ParseToken(data []byte) (string, error) {
	var resp struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			XMLName      xml.Name `xml:"Body"`
			GetTokenResp struct {
				XMLName xml.Name `xml:"getTokenResponse"`
				Token   string   `xml:"token"`
			}
		}
	}

	if err := xml.Unmarshal(data, &resp); err != nil {
		return "", fmt.Errorf("error al parsear respuesta de token: %w", err)
	}

	return resp.Body.GetTokenResp.Token, nil
}

// ParseEstadoDTE parsea la respuesta de estado de DTE
func (p *ResponseParser) ParseEstadoDTE(data []byte) (*models.EstadoSII, error) {
	var resp struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			XMLName       xml.Name `xml:"Body"`
			GetEstDteResp struct {
				XMLName xml.Name `xml:"getEstDteResponse"`
				Estado  string   `xml:"estado"`
				Glosa   string   `xml:"glosa"`
				TrackID string   `xml:"trackid"`
			}
		}
	}

	if err := xml.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error al parsear respuesta de estado DTE: %w", err)
	}

	return &models.EstadoSII{
		Estado:  resp.Body.GetEstDteResp.Estado,
		Glosa:   resp.Body.GetEstDteResp.Glosa,
		TrackID: resp.Body.GetEstDteResp.TrackID,
	}, nil
}

// ParseRespuestaEnvio parsea la respuesta de envío de DTE
func (p *ResponseParser) ParseRespuestaEnvio(data []byte) (*models.RespuestaSII, error) {
	var resp struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			XMLName     xml.Name `xml:"Body"`
			SendDTEResp struct {
				XMLName xml.Name `xml:"sendDTEResponse"`
				Estado  string   `xml:"estado"`
				Glosa   string   `xml:"glosa"`
				TrackID string   `xml:"trackid"`
			}
		}
	}

	if err := xml.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("error al parsear respuesta de envío: %w", err)
	}

	return &models.RespuestaSII{
		Estado:  resp.Body.SendDTEResp.Estado,
		Glosa:   resp.Body.SendDTEResp.Glosa,
		TrackID: resp.Body.SendDTEResp.TrackID,
	}, nil
}
