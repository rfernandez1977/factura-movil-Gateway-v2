package models

// Municipality representa una municipalidad
type Municipality struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// Activity representa una actividad económica
type Activity struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

// AdditionalAddress representa una dirección adicional
type AdditionalAddress struct {
	ID           int          `json:"id"`
	Address      string       `json:"address"`
	Municipality Municipality `json:"municipality"`
}

// Client representa un cliente en el sistema
// NOTA: Esta es la definición principal y única de Client para todo el sistema. No duplicar en otros archivos.
type Client struct {
	ID                int                 `json:"id"`
	Code              string              `json:"code"`
	Name              string              `json:"name"`
	Address           string              `json:"address"`
	AdditionalAddress []AdditionalAddress `json:"additionalAddress"`
	Email             string              `json:"email,omitempty"`
	Municipality      Municipality        `json:"municipality"`
	Activity          Activity            `json:"activity"`
	Line              string              `json:"line"`
}

// ClientResponse representa la respuesta de la API para clientes
type ClientResponse struct {
	Clients []Client `json:"clients"`
}

// SearchClientParams representa los parámetros de búsqueda de clientes
type SearchClientParams struct {
	SearchTerm string `json:"search_term"`
}
