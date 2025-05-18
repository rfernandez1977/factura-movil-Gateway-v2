package client

import (
	"context"
	"fmt"
	"time"

	"FMgo/core/sii/models/siimodels"
)

// DTEClient implementa el cliente para operaciones con DTE
type DTEClient struct {
	soapClient *SOAPClient
	config     *siimodels.ConfigSII
	authClient *AuthClient
}

// NewDTEClient crea una nueva instancia del cliente DTE
func NewDTEClient(config *siimodels.ConfigSII) (*DTEClient, error) {
	soapClient, err := NewSOAPClient(config)
	if err != nil {
		return nil, fmt.Errorf("error al crear cliente SOAP: %w", err)
	}

	authClient := NewAuthClient(soapClient, config)

	return &DTEClient{
		soapClient: soapClient,
		config:     config,
		authClient: authClient,
	}, nil
}

// EnviarDTE envía un DTE al SII
func (c *DTEClient) EnviarDTE(ctx context.Context, dte *siimodels.DTE) (*siimodels.RespuestaEnvio, error) {
	// Obtener token de autenticación
	token, err := c.authClient.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener token: %w", err)
	}

	// Preparar el sobre de envío
	sobre := &siimodels.EnvioDTE{
		SetDTE: siimodels.SetDTE{
			ID:  fmt.Sprintf("SetDoc_%s", time.Now().Format("20060102150405")),
			DTE: *dte,
			Caratula: siimodels.Caratula{
				RutEmisor:  c.config.RutEmpresa,
				RutEnvia:   c.config.RutCertificado,
				FechaEnvio: time.Now(),
				Version:    "1.0",
			},
		},
		Version: "1.0",
	}

	// Preparar la respuesta
	respuesta := &siimodels.RespuestaEnvio{}

	// Crear contexto con token
	ctxWithToken := context.WithValue(ctx, "token", token)

	// Enviar la solicitud
	err = c.soapClient.Call(ctxWithToken, siimodels.EndpointEnvioCert, sobre, respuesta)
	if err != nil {
		return nil, fmt.Errorf("error al enviar DTE: %w", err)
	}

	return respuesta, nil
}

// ConsultarEstadoDTE consulta el estado de un DTE
func (c *DTEClient) ConsultarEstadoDTE(ctx context.Context, rutEmisor string, tipoDTE int, folio int64) (*siimodels.EstadoDTE, error) {
	// Obtener token de autenticación
	token, err := c.authClient.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener token: %w", err)
	}

	// Preparar la consulta
	consulta := &siimodels.ConsultaDTE{
		RutEmisor:     rutEmisor,
		TipoDTE:       tipoDTE,
		Folio:         folio,
		RutConsulta:   c.config.RutEmpresa,
		Token:         token,
		FechaConsulta: time.Now(),
	}

	// Preparar la respuesta
	respuesta := &siimodels.EstadoDTE{}

	// Crear contexto con token
	ctxWithToken := context.WithValue(ctx, "token", token)

	// Enviar la solicitud
	err = c.soapClient.Call(ctxWithToken, siimodels.EndpointConsultaCert, consulta, respuesta)
	if err != nil {
		return nil, fmt.Errorf("error al consultar estado DTE: %w", err)
	}

	return respuesta, nil
}

// ConsultarEstadoEnvio consulta el estado de un envío de DTE
func (c *DTEClient) ConsultarEstadoEnvio(ctx context.Context, trackID string) (*siimodels.EstadoEnvio, error) {
	// Obtener token de autenticación
	token, err := c.authClient.GetToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al obtener token: %w", err)
	}

	// Preparar la consulta
	consulta := &siimodels.ConsultaTrackID{
		RutEmpresa:    c.config.RutEmpresa,
		TrackID:       trackID,
		Token:         token,
		FechaConsulta: time.Now(),
	}

	// Preparar la respuesta
	respuesta := &siimodels.EstadoEnvio{}

	// Crear contexto con token
	ctxWithToken := context.WithValue(ctx, "token", token)

	// Enviar la solicitud
	err = c.soapClient.Call(ctxWithToken, siimodels.EndpointConsultaCert, consulta, respuesta)
	if err != nil {
		return nil, fmt.Errorf("error al consultar estado de envío: %w", err)
	}

	return respuesta, nil
}
