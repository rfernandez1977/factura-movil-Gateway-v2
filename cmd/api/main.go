package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fmgo/config"
	"github.com/fmgo/controllers"
	"github.com/fmgo/repository"
	"github.com/fmgo/services"
	"github.com/fmgo/sii"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Configurar MongoDB
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("fmgodb")

	// Configurar Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Configurar Supabase
	supabaseConfig := config.NewSupabaseConfig(
		os.Getenv("SUPABASE_URL"),
		os.Getenv("SUPABASE_KEY"),
		os.Getenv("SUPABASE_TOKEN"),
		os.Getenv("SII_AMBIENTE"),
		os.Getenv("SII_BASE_URL"),
	)

	// Configurar SII
	siiService := sii.NewSIIService(
		os.Getenv("SII_BASE_URL"),
		os.Getenv("SII_CERT_FILE"),
		os.Getenv("SII_KEY_FILE"),
		os.Getenv("SII_AMBIENTE"),
	)

	// Inicializar repositorios
	docRepo := repository.NewDocumentRepository(db)

	// Inicializar servicios
	validationSvc := services.NewValidationService()
	cafSvc := services.NewCAFService(
		db,
		redisClient,
		siiService,
		os.Getenv("SII_CERT_FILE"),
		os.Getenv("SII_KEY_FILE"),
		supabaseConfig,
	)
	auditSvc := services.NewAuditService(db)
	docService := services.NewDocumentService(docRepo, validationSvc, cafSvc, auditSvc)

	// Inicializar controladores
	docController := controllers.NewDocumentController(docService)

	// Configurar router
	router := gin.Default()

	// Rutas de documentos
	router.POST("/api/documentos", docController.CrearDocumento)
	router.GET("/api/documentos/:tipo/:folio", docController.ObtenerDocumento)
	router.PUT("/api/documentos/:tipo/:folio", docController.ActualizarDocumento)
	router.PATCH("/api/documentos/:id/estado/:estado", docController.CambiarEstadoDocumento)
	router.POST("/api/documentos/referencias", docController.AgregarReferencia)
	router.GET("/api/documentos/:tipo/:folio/referencias", docController.ObtenerReferencias)

	// Configurar servidor
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Iniciar servidor en una goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error al iniciar el servidor: %v\n", err)
		}
	}()

	// Esperar señal de interrupción
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Apagando servidor...")

	// Dar tiempo para que las conexiones se cierren
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Error al apagar el servidor:", err)
	}

	log.Println("Servidor apagado correctamente")
}

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
