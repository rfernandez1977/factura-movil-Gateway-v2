package supabase

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/cursor/FMgo/config"
	supa "github.com/supabase-community/supabase-go"
)

// Client represents a Supabase client
type Client struct {
	client *supa.Client
	config *config.Config
}

// NewClient creates a new Supabase client
func NewClient(cfg *config.Config) (*Client, error) {
	// Mostrar información de depuración
	log.Printf("Creando cliente Supabase con URL: %s", cfg.Supabase.URL)
	log.Printf("Longitud de AnonKey: %d caracteres", len(cfg.Supabase.AnonKey))
	log.Printf("Primeros 10 caracteres de AnonKey: %s...", cfg.Supabase.AnonKey[:10])

	client, err := supa.NewClient(cfg.Supabase.URL, cfg.Supabase.AnonKey, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating Supabase client: %w", err)
	}

	return &Client{
		client: client,
		config: cfg,
	}, nil
}

// Ping checks if the connection to Supabase is working
func (c *Client) Ping(ctx context.Context) error {
	// Make a request to the health check endpoint
	healthURL := fmt.Sprintf("%s/rest/v1/", c.config.Supabase.URL)
	log.Printf("Enviando solicitud a: %s", healthURL)

	req, err := http.NewRequestWithContext(ctx, "GET", healthURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	log.Printf("Estableciendo encabezados con AnonKey: %s...", c.config.Supabase.AnonKey[:10])
	req.Header.Set("apikey", c.config.Supabase.AnonKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.Supabase.AnonKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error pinging Supabase: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("Respuesta recibida, código: %d", resp.StatusCode)

	// Leer y mostrar el cuerpo de la respuesta para depuración
	body, _ := io.ReadAll(resp.Body)
	log.Printf("Respuesta: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from Supabase: %d", resp.StatusCode)
	}

	return nil
}

// GetClient returns the underlying Supabase client
func (c *Client) GetClient() *supa.Client {
	return c.client
}

// GetConfig returns the configuration used by this client
func (c *Client) GetConfig() *config.Config {
	return c.config
}
