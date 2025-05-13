package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cursor/FMgo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EcommercePlatform define los tipos de plataformas soportadas
type EcommercePlatform string

const (
	PlatformPrestashop  EcommercePlatform = "PRESTASHOP"
	PlatformShopify     EcommercePlatform = "SHOPIFY"
	PlatformWooCommerce EcommercePlatform = "WOOCOMMERCE"
	PlatformJumpseller  EcommercePlatform = "JUMPSELLER"
	PlatformWix         EcommercePlatform = "WIX"
)

// EcommerceConfig contiene la configuración de conexión a una plataforma
type EcommerceConfig struct {
	ID            string            `json:"id" bson:"_id"`
	Platform      EcommercePlatform `json:"platform" bson:"platform"`
	StoreName     string            `json:"store_name" bson:"store_name"`
	APIKey        string            `json:"api_key" bson:"api_key"`
	APISecret     string            `json:"api_secret" bson:"api_secret"`
	StoreURL      string            `json:"store_url" bson:"store_url"`
	WebhookURL    string            `json:"webhook_url" bson:"webhook_url"`
	SyncProducts  bool              `json:"sync_products" bson:"sync_products"`
	SyncOrders    bool              `json:"sync_orders" bson:"sync_orders"`
	SyncCustomers bool              `json:"sync_customers" bson:"sync_customers"`
	LastSync      time.Time         `json:"last_sync" bson:"last_sync"`
	CreatedAt     time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at" bson:"updated_at"`
}

// EcommerceService maneja la integración con plataformas de e-commerce
type EcommerceService struct {
	db    *mongo.Database
	cache *CacheManager
}

// CacheManager maneja el caché de documentos
type CacheManager struct {
	documentCache map[string]*models.DocumentoAlmacenado
	ttl           time.Duration
}

// NewEcommerceService crea una nueva instancia del servicio de e-commerce
func NewEcommerceService(db *mongo.Database) *EcommerceService {
	return &EcommerceService{
		db: db,
		cache: &CacheManager{
			documentCache: make(map[string]*models.DocumentoAlmacenado),
			ttl:           24 * time.Hour,
		},
	}
}

// GetDocument obtiene un documento del caché
func (cm *CacheManager) GetDocument(key string) (*models.DocumentoAlmacenado, error) {
	if doc, exists := cm.documentCache[key]; exists {
		if time.Now().Before(doc.CacheInfo.ExpiresAt) {
			return doc, nil
		}
		delete(cm.documentCache, key)
	}
	return nil, fmt.Errorf("documento no encontrado o expirado")
}

// SetDocument almacena un documento en el caché
func (cm *CacheManager) SetDocument(key string, doc *models.DocumentoAlmacenado) {
	doc.CacheInfo.CreatedAt = time.Now()
	doc.CacheInfo.ExpiresAt = time.Now().Add(cm.ttl)
	cm.documentCache[key] = doc
}

// RegisterStore registra una nueva tienda para sincronización
func (s *EcommerceService) RegisterStore(ctx context.Context, config *EcommerceConfig) error {
	// Comprobar si ya existe
	var existingConfig EcommerceConfig
	err := s.db.Collection("ecommerce_stores").FindOne(ctx, bson.M{
		"platform":   config.Platform,
		"store_name": config.StoreName,
	}).Decode(&existingConfig)

	// Si no existe error o no es un error "no encontrado", retornar
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	// Si ya existe, retornar error
	if err == nil {
		return fmt.Errorf("la tienda ya está registrada")
	}

	// Establecer valores por defecto
	config.ID = GenerateID()
	config.CreatedAt = time.Now()
	config.UpdatedAt = time.Now()

	// Insertar nueva configuración
	_, err = s.db.Collection("ecommerce_stores").InsertOne(ctx, config)
	return err
}

// SyncProducts sincroniza productos con la plataforma
func (s *EcommerceService) SyncProducts(ctx context.Context, storeID string) error {
	config, err := s.getStoreConfig(ctx, storeID)
	if err != nil {
		return err
	}

	switch config.Platform {
	case PlatformPrestashop:
		return s.syncPrestashopProducts(ctx, config)
	case PlatformShopify:
		return s.syncShopifyProducts(ctx, config)
	case PlatformWooCommerce:
		return s.syncWooCommerceProducts(ctx, config)
	case PlatformJumpseller:
		return s.syncJumpsellerProducts(ctx, config)
	case PlatformWix:
		return s.syncWixProducts(ctx, config)
	default:
		return fmt.Errorf("plataforma no soportada: %s", config.Platform)
	}
}

// SyncOrders sincroniza órdenes con la plataforma
func (s *EcommerceService) SyncOrders(ctx context.Context, storeID string) error {
	config, err := s.getStoreConfig(ctx, storeID)
	if err != nil {
		return err
	}

	switch config.Platform {
	case PlatformPrestashop:
		return s.syncPrestashopOrders(ctx, config)
	case PlatformShopify:
		return s.SyncShopifyOrders(ctx, config.ID)
	case PlatformWooCommerce:
		return s.SyncWooCommerceOrders(ctx, config.ID)
	case PlatformJumpseller:
		return s.syncJumpsellerOrders(ctx, config)
	case PlatformWix:
		return s.syncWixOrders(ctx, config)
	default:
		return fmt.Errorf("plataforma no soportada: %s", config.Platform)
	}
}

// getStoreConfig obtiene la configuración de una tienda
func (s *EcommerceService) getStoreConfig(ctx context.Context, storeID string) (*EcommerceConfig, error) {
	var config EcommerceConfig
	err := s.db.Collection("ecommerce_stores").FindOne(ctx, bson.M{"_id": storeID}).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("error al obtener configuración de tienda: %v", err)
	}
	return &config, nil
}

// PrestaShopProduct representa un producto de PrestaShop
type PrestaShopProduct struct {
	ID                 int64   `json:"id" bson:"id"`
	Name               string  `json:"name" bson:"name"`
	Description        string  `json:"description" bson:"description"`
	Price              float64 `json:"price" bson:"price"`
	Reference          string  `json:"reference" bson:"reference"`
	Active             bool    `json:"active" bson:"active"`
	AvailableForOrder  bool    `json:"available_for_order" bson:"available_for_order"`
	ShowPrice          bool    `json:"show_price" bson:"show_price"`
	Quantity           int     `json:"quantity" bson:"quantity"`
	MinimalQuantity    int     `json:"minimal_quantity" bson:"minimal_quantity"`
	IDCategoryDefault  int     `json:"id_category_default" bson:"id_category_default"`
	ManufacturerName   string  `json:"manufacturer_name" bson:"manufacturer_name"`
	SupplierReference  string  `json:"supplier_reference" bson:"supplier_reference"`
	EAN13              string  `json:"ean13" bson:"ean13"`
	UPC                string  `json:"upc" bson:"upc"`
	Width              float64 `json:"width" bson:"width"`
	Height             float64 `json:"height" bson:"height"`
	Depth              float64 `json:"depth" bson:"depth"`
	Weight             float64 `json:"weight" bson:"weight"`
	TaxRulesGroupID    int     `json:"id_tax_rules_group" bson:"id_tax_rules_group"`
	WholesalePrice     float64 `json:"wholesale_price" bson:"wholesale_price"`
	OnSale             bool    `json:"on_sale" bson:"on_sale"`
	OnlineOnly         bool    `json:"online_only" bson:"online_only"`
	Condition          string  `json:"condition" bson:"condition"`
	Customizable       bool    `json:"customizable" bson:"customizable"`
	UploadableFiles    bool    `json:"uploadable_files" bson:"uploadable_files"`
	TextFields         bool    `json:"text_fields" bson:"text_fields"`
	OutOfStock         int     `json:"out_of_stock" bson:"out_of_stock"`
	AdditionalShipping float64 `json:"additional_shipping_cost" bson:"additional_shipping_cost"`
	Unity              string  `json:"unity" bson:"unity"`
	UnitPriceRatio     float64 `json:"unit_price_ratio" bson:"unit_price_ratio"`
}

// PrestaShopProductResponse representa la respuesta de la API de PrestaShop
type PrestaShopProductResponse struct {
	Product PrestaShopProduct `json:"product"`
}

// syncPrestashopProducts sincroniza productos con PrestaShop
func (s *EcommerceService) syncPrestashopProducts(ctx context.Context, config *EcommerceConfig) error {
	// Crear cliente HTTP con autenticación básica
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/products", config.StoreURL), nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %v", err)
	}

	// Configurar autenticación básica
	req.SetBasicAuth(config.APIKey, "")

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta: %s", resp.Status)
	}

	// Decodificar la respuesta
	var products []PrestaShopProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	// Procesar cada producto
	for _, productResponse := range products {
		product := productResponse.Product

		// Verificar si el producto ya existe en nuestra base de datos
		var existingProduct PrestaShopProduct
		err := s.db.Collection("prestashop_products").FindOne(ctx, bson.M{"id": product.ID}).Decode(&existingProduct)
		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("error al buscar producto existente: %v", err)
		}

		// Si el producto no existe o ha sido modificado, actualizarlo
		if err == mongo.ErrNoDocuments {
			_, err = s.db.Collection("prestashop_products").UpdateOne(
				ctx,
				bson.M{"id": product.ID},
				bson.M{
					"$set": bson.M{
						"name":                product.Name,
						"description":         product.Description,
						"price":               product.Price,
						"reference":           product.Reference,
						"active":              product.Active,
						"available_for_order": product.AvailableForOrder,
						"show_price":          product.ShowPrice,
						"quantity":            product.Quantity,
						"minimal_quantity":    product.MinimalQuantity,
						"id_category_default": product.IDCategoryDefault,
						"manufacturer_name":   product.ManufacturerName,
						"supplier_reference":  product.SupplierReference,
						"ean13":               product.EAN13,
						"upc":                 product.UPC,
						"width":               product.Width,
						"height":              product.Height,
						"depth":               product.Depth,
						"weight":              product.Weight,
						"id_tax_rules_group":  product.TaxRulesGroupID,
						"wholesale_price":     product.WholesalePrice,
						"on_sale":             product.OnSale,
						"online_only":         product.OnlineOnly,
						"condition":           product.Condition,
						"customizable":        product.Customizable,
						"uploadable_files":    product.UploadableFiles,
						"text_fields":         product.TextFields,
						"out_of_stock":        product.OutOfStock,
						"additional_shipping": product.AdditionalShipping,
						"unity":               product.Unity,
						"unit_price_ratio":    product.UnitPriceRatio,
						"updated_at":          time.Now(),
					},
				},
			)
			if err != nil {
				return fmt.Errorf("error al actualizar producto: %v", err)
			}
		}
	}

	// Actualizar la última sincronización
	_, err = s.db.Collection("ecommerce_stores").UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{"$set": bson.M{"last_sync": time.Now()}},
	)
	if err != nil {
		return fmt.Errorf("error al actualizar última sincronización: %v", err)
	}

	return nil
}

// ShopifyProduct representa un producto de Shopify
type ShopifyProduct struct {
	ID                int64     `json:"id" bson:"id"`
	Title             string    `json:"title" bson:"title"`
	BodyHTML          string    `json:"body_html" bson:"body_html"`
	Vendor            string    `json:"vendor" bson:"vendor"`
	ProductType       string    `json:"product_type" bson:"product_type"`
	Handle            string    `json:"handle" bson:"handle"`
	Status            string    `json:"status" bson:"status"`
	PublishedScope    string    `json:"published_scope" bson:"published_scope"`
	Tags              string    `json:"tags" bson:"tags"`
	CreatedAt         time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" bson:"updated_at"`
	PublishedAt       time.Time `json:"published_at" bson:"published_at"`
	TemplateSuffix    string    `json:"template_suffix" bson:"template_suffix"`
	AdminGraphqlAPIID string    `json:"admin_graphql_api_id" bson:"admin_graphql_api_id"`
	Variants          []struct {
		ID                   int64     `json:"id" bson:"id"`
		ProductID            int64     `json:"product_id" bson:"product_id"`
		Title                string    `json:"title" bson:"title"`
		Price                string    `json:"price" bson:"price"`
		SKU                  string    `json:"sku" bson:"sku"`
		Position             int       `json:"position" bson:"position"`
		InventoryPolicy      string    `json:"inventory_policy" bson:"inventory_policy"`
		CompareAtPrice       string    `json:"compare_at_price" bson:"compare_at_price"`
		FulfillmentService   string    `json:"fulfillment_service" bson:"fulfillment_service"`
		InventoryManagement  string    `json:"inventory_management" bson:"inventory_management"`
		Option1              string    `json:"option1" bson:"option1"`
		Option2              string    `json:"option2" bson:"option2"`
		Option3              string    `json:"option3" bson:"option3"`
		CreatedAt            time.Time `json:"created_at" bson:"created_at"`
		UpdatedAt            time.Time `json:"updated_at" bson:"updated_at"`
		Taxable              bool      `json:"taxable" bson:"taxable"`
		Barcode              string    `json:"barcode" bson:"barcode"`
		Grams                int       `json:"grams" bson:"grams"`
		ImageID              int64     `json:"image_id" bson:"image_id"`
		Weight               float64   `json:"weight" bson:"weight"`
		WeightUnit           string    `json:"weight_unit" bson:"weight_unit"`
		InventoryItemID      int64     `json:"inventory_item_id" bson:"inventory_item_id"`
		InventoryQuantity    int       `json:"inventory_quantity" bson:"inventory_quantity"`
		OldInventoryQuantity int       `json:"old_inventory_quantity" bson:"old_inventory_quantity"`
		RequiresShipping     bool      `json:"requires_shipping" bson:"requires_shipping"`
		AdminGraphqlAPIID    string    `json:"admin_graphql_api_id" bson:"admin_graphql_api_id"`
	} `json:"variants" bson:"variants"`
	Options []struct {
		ID        int64    `json:"id" bson:"id"`
		ProductID int64    `json:"product_id" bson:"product_id"`
		Name      string   `json:"name" bson:"name"`
		Position  int      `json:"position" bson:"position"`
		Values    []string `json:"values" bson:"values"`
	} `json:"options" bson:"options"`
	Images []struct {
		ID                int64     `json:"id" bson:"id"`
		ProductID         int64     `json:"product_id" bson:"product_id"`
		Position          int       `json:"position" bson:"position"`
		CreatedAt         time.Time `json:"created_at" bson:"created_at"`
		UpdatedAt         time.Time `json:"updated_at" bson:"updated_at"`
		Alt               string    `json:"alt" bson:"alt"`
		Width             int       `json:"width" bson:"width"`
		Height            int       `json:"height" bson:"height"`
		Src               string    `json:"src" bson:"src"`
		VariantIDs        []int64   `json:"variant_ids" bson:"variant_ids"`
		AdminGraphqlAPIID string    `json:"admin_graphql_api_id" bson:"admin_graphql_api_id"`
	} `json:"images" bson:"images"`
	Image struct {
		ID                int64     `json:"id" bson:"id"`
		ProductID         int64     `json:"product_id" bson:"product_id"`
		Position          int       `json:"position" bson:"position"`
		CreatedAt         time.Time `json:"created_at" bson:"created_at"`
		UpdatedAt         time.Time `json:"updated_at" bson:"updated_at"`
		Alt               string    `json:"alt" bson:"alt"`
		Width             int       `json:"width" bson:"width"`
		Height            int       `json:"height" bson:"height"`
		Src               string    `json:"src" bson:"src"`
		VariantIDs        []int64   `json:"variant_ids" bson:"variant_ids"`
		AdminGraphqlAPIID string    `json:"admin_graphql_api_id" bson:"admin_graphql_api_id"`
	} `json:"image" bson:"image"`
}

// ShopifyProductResponse representa la respuesta de la API de Shopify para productos
type ShopifyProductResponse struct {
	Product ShopifyProduct `json:"product"`
}

// syncShopifyProducts sincroniza productos con Shopify
func (s *EcommerceService) syncShopifyProducts(ctx context.Context, config *EcommerceConfig) error {
	// Crear cliente HTTP con autenticación OAuth
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/admin/api/2023-01/products.json", config.StoreURL), nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %v", err)
	}

	// Configurar autenticación OAuth
	req.Header.Set("X-Shopify-Access-Token", config.APIKey)

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta: %s", resp.Status)
	}

	// Decodificar la respuesta
	var productsResponse struct {
		Products []ShopifyProduct `json:"products"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&productsResponse); err != nil {
		return fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	// Procesar cada producto
	for _, product := range productsResponse.Products {
		// Verificar si el producto ya existe en nuestra base de datos
		var existingProduct ShopifyProduct
		err := s.db.Collection("shopify_products").FindOne(ctx, bson.M{"id": product.ID}).Decode(&existingProduct)
		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("error al buscar producto existente: %v", err)
		}

		// Si el producto no existe o ha sido modificado, actualizarlo
		if err == mongo.ErrNoDocuments || existingProduct.UpdatedAt.Before(product.UpdatedAt) {
			_, err = s.db.Collection("shopify_products").UpdateOne(
				ctx,
				bson.M{"id": product.ID},
				bson.M{
					"$set": product,
				},
			)
			if err != nil {
				return fmt.Errorf("error al actualizar producto: %v", err)
			}
		}
	}

	// Actualizar la última sincronización
	_, err = s.db.Collection("ecommerce_stores").UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{"$set": bson.M{"last_sync": time.Now()}},
	)
	if err != nil {
		return fmt.Errorf("error al actualizar última sincronización: %v", err)
	}

	return nil
}

// WooCommerceProduct representa un producto de WooCommerce
type WooCommerceProduct struct {
	ID                int64     `json:"id" bson:"id"`
	Name              string    `json:"name" bson:"name"`
	Slug              string    `json:"slug" bson:"slug"`
	Permalink         string    `json:"permalink" bson:"permalink"`
	DateCreated       time.Time `json:"date_created" bson:"date_created"`
	DateModified      time.Time `json:"date_modified" bson:"date_modified"`
	Type              string    `json:"type" bson:"type"`
	Status            string    `json:"status" bson:"status"`
	Featured          bool      `json:"featured" bson:"featured"`
	CatalogVisibility string    `json:"catalog_visibility" bson:"catalog_visibility"`
	Description       string    `json:"description" bson:"description"`
	ShortDescription  string    `json:"short_description" bson:"short_description"`
	SKU               string    `json:"sku" bson:"sku"`
	Price             string    `json:"price" bson:"price"`
	RegularPrice      string    `json:"regular_price" bson:"regular_price"`
	SalePrice         string    `json:"sale_price" bson:"sale_price"`
	DateOnSaleFrom    time.Time `json:"date_on_sale_from" bson:"date_on_sale_from"`
	DateOnSaleTo      time.Time `json:"date_on_sale_to" bson:"date_on_sale_to"`
	PriceHTML         string    `json:"price_html" bson:"price_html"`
	OnSale            bool      `json:"on_sale" bson:"on_sale"`
	Purchasable       bool      `json:"purchasable" bson:"purchasable"`
	TotalSales        int       `json:"total_sales" bson:"total_sales"`
	Virtual           bool      `json:"virtual" bson:"virtual"`
	Downloadable      bool      `json:"downloadable" bson:"downloadable"`
	Downloads         []struct {
		ID   string `json:"id" bson:"id"`
		Name string `json:"name" bson:"name"`
		File string `json:"file" bson:"file"`
	} `json:"downloads" bson:"downloads"`
	DownloadLimit     int    `json:"download_limit" bson:"download_limit"`
	DownloadExpiry    int    `json:"download_expiry" bson:"download_expiry"`
	ExternalURL       string `json:"external_url" bson:"external_url"`
	ButtonText        string `json:"button_text" bson:"button_text"`
	TaxStatus         string `json:"tax_status" bson:"tax_status"`
	TaxClass          string `json:"tax_class" bson:"tax_class"`
	ManageStock       bool   `json:"manage_stock" bson:"manage_stock"`
	StockQuantity     int    `json:"stock_quantity" bson:"stock_quantity"`
	StockStatus       string `json:"stock_status" bson:"stock_status"`
	Backorders        string `json:"backorders" bson:"backorders"`
	BackordersAllowed bool   `json:"backorders_allowed" bson:"backorders_allowed"`
	Backordered       bool   `json:"backordered" bson:"backordered"`
	SoldIndividually  bool   `json:"sold_individually" bson:"sold_individually"`
	Weight            string `json:"weight" bson:"weight"`
	Dimensions        struct {
		Length string `json:"length" bson:"length"`
		Width  string `json:"width" bson:"width"`
		Height string `json:"height" bson:"height"`
	} `json:"dimensions" bson:"dimensions"`
	ShippingRequired bool   `json:"shipping_required" bson:"shipping_required"`
	ShippingTaxable  bool   `json:"shipping_taxable" bson:"shipping_taxable"`
	ShippingClass    string `json:"shipping_class" bson:"shipping_class"`
	ShippingClassID  int    `json:"shipping_class_id" bson:"shipping_class_id"`
	ReviewsAllowed   bool   `json:"reviews_allowed" bson:"reviews_allowed"`
	AverageRating    string `json:"average_rating" bson:"average_rating"`
	RatingCount      int    `json:"rating_count" bson:"rating_count"`
	RelatedIDs       []int  `json:"related_ids" bson:"related_ids"`
	UpsellIDs        []int  `json:"upsell_ids" bson:"upsell_ids"`
	CrossSellIDs     []int  `json:"cross_sell_ids" bson:"cross_sell_ids"`
	ParentID         int    `json:"parent_id" bson:"parent_id"`
	PurchaseNote     string `json:"purchase_note" bson:"purchase_note"`
	Categories       []struct {
		ID   int    `json:"id" bson:"id"`
		Name string `json:"name" bson:"name"`
		Slug string `json:"slug" bson:"slug"`
	} `json:"categories" bson:"categories"`
	Tags []struct {
		ID   int    `json:"id" bson:"id"`
		Name string `json:"name" bson:"name"`
		Slug string `json:"slug" bson:"slug"`
	} `json:"tags" bson:"tags"`
	Images []struct {
		ID           int64  `json:"id" bson:"id"`
		DateCreated  string `json:"date_created" bson:"date_created"`
		DateModified string `json:"date_modified" bson:"date_modified"`
		Src          string `json:"src" bson:"src"`
		Name         string `json:"name" bson:"name"`
		Alt          string `json:"alt" bson:"alt"`
		Position     int    `json:"position" bson:"position"`
	} `json:"images" bson:"images"`
	Attributes []struct {
		ID        int      `json:"id" bson:"id"`
		Name      string   `json:"name" bson:"name"`
		Position  int      `json:"position" bson:"position"`
		Visible   bool     `json:"visible" bson:"visible"`
		Variation bool     `json:"variation" bson:"variation"`
		Options   []string `json:"options" bson:"options"`
	} `json:"attributes" bson:"attributes"`
	DefaultAttributes []struct {
		ID     int    `json:"id" bson:"id"`
		Name   string `json:"name" bson:"name"`
		Option string `json:"option" bson:"option"`
	} `json:"default_attributes" bson:"default_attributes"`
	Variations      []int `json:"variations" bson:"variations"`
	GroupedProducts []int `json:"grouped_products" bson:"grouped_products"`
	MenuOrder       int   `json:"menu_order" bson:"menu_order"`
	MetaData        []struct {
		ID    int    `json:"id" bson:"id"`
		Key   string `json:"key" bson:"key"`
		Value string `json:"value" bson:"value"`
	} `json:"meta_data" bson:"meta_data"`
}

// WooCommerceOrder representa una orden de WooCommerce
type WooCommerceOrder struct {
	ID                int64     `json:"id" bson:"id"`
	ParentID          int       `json:"parent_id" bson:"parent_id"`
	Number            string    `json:"number" bson:"number"`
	OrderKey          string    `json:"order_key" bson:"order_key"`
	CreatedVia        string    `json:"created_via" bson:"created_via"`
	Version           string    `json:"version" bson:"version"`
	Status            string    `json:"status" bson:"status"`
	Currency          string    `json:"currency" bson:"currency"`
	DateCreated       time.Time `json:"date_created" bson:"date_created"`
	DateModified      time.Time `json:"date_modified" bson:"date_modified"`
	DiscountTotal     string    `json:"discount_total" bson:"discount_total"`
	DiscountTax       string    `json:"discount_tax" bson:"discount_tax"`
	ShippingTotal     string    `json:"shipping_total" bson:"shipping_total"`
	ShippingTax       string    `json:"shipping_tax" bson:"shipping_tax"`
	CartTax           string    `json:"cart_tax" bson:"cart_tax"`
	Total             string    `json:"total" bson:"total"`
	TotalTax          string    `json:"total_tax" bson:"total_tax"`
	PricesIncludeTax  bool      `json:"prices_include_tax" bson:"prices_include_tax"`
	CustomerID        int       `json:"customer_id" bson:"customer_id"`
	CustomerIPAddress string    `json:"customer_ip_address" bson:"customer_ip_address"`
	CustomerUserAgent string    `json:"customer_user_agent" bson:"customer_user_agent"`
	CustomerNote      string    `json:"customer_note" bson:"customer_note"`
	Billing           struct {
		FirstName string `json:"first_name" bson:"first_name"`
		LastName  string `json:"last_name" bson:"last_name"`
		Company   string `json:"company" bson:"company"`
		Address1  string `json:"address_1" bson:"address_1"`
		Address2  string `json:"address_2" bson:"address_2"`
		City      string `json:"city" bson:"city"`
		State     string `json:"state" bson:"state"`
		Postcode  string `json:"postcode" bson:"postcode"`
		Country   string `json:"country" bson:"country"`
		Email     string `json:"email" bson:"email"`
		Phone     string `json:"phone" bson:"phone"`
	} `json:"billing" bson:"billing"`
	Shipping struct {
		FirstName string `json:"first_name" bson:"first_name"`
		LastName  string `json:"last_name" bson:"last_name"`
		Company   string `json:"company" bson:"company"`
		Address1  string `json:"address_1" bson:"address_1"`
		Address2  string `json:"address_2" bson:"address_2"`
		City      string `json:"city" bson:"city"`
		State     string `json:"state" bson:"state"`
		Postcode  string `json:"postcode" bson:"postcode"`
		Country   string `json:"country" bson:"country"`
	} `json:"shipping" bson:"shipping"`
	PaymentMethod      string    `json:"payment_method" bson:"payment_method"`
	PaymentMethodTitle string    `json:"payment_method_title" bson:"payment_method_title"`
	TransactionID      string    `json:"transaction_id" bson:"transaction_id"`
	DatePaid           time.Time `json:"date_paid" bson:"date_paid"`
	DateCompleted      time.Time `json:"date_completed" bson:"date_completed"`
	CartHash           string    `json:"cart_hash" bson:"cart_hash"`
	LineItems          []struct {
		ID          int64  `json:"id" bson:"id"`
		Name        string `json:"name" bson:"name"`
		ProductID   int    `json:"product_id" bson:"product_id"`
		VariationID int    `json:"variation_id" bson:"variation_id"`
		Quantity    int    `json:"quantity" bson:"quantity"`
		TaxClass    string `json:"tax_class" bson:"tax_class"`
		Subtotal    string `json:"subtotal" bson:"subtotal"`
		SubtotalTax string `json:"subtotal_tax" bson:"subtotal_tax"`
		Total       string `json:"total" bson:"total"`
		TotalTax    string `json:"total_tax" bson:"total_tax"`
		Taxes       []struct {
			ID       int    `json:"id" bson:"id"`
			Total    string `json:"total" bson:"total"`
			Subtotal string `json:"subtotal" bson:"subtotal"`
		} `json:"taxes" bson:"taxes"`
		MetaData []struct {
			ID    int    `json:"id" bson:"id"`
			Key   string `json:"key" bson:"key"`
			Value string `json:"value" bson:"value"`
		} `json:"meta_data" bson:"meta_data"`
		SKU   string `json:"sku" bson:"sku"`
		Price int    `json:"price" bson:"price"`
	} `json:"line_items" bson:"line_items"`
	TaxLines []struct {
		ID               int     `json:"id" bson:"id"`
		RateCode         string  `json:"rate_code" bson:"rate_code"`
		RateID           int     `json:"rate_id" bson:"rate_id"`
		Label            string  `json:"label" bson:"label"`
		Compound         bool    `json:"compound" bson:"compound"`
		TaxTotal         string  `json:"tax_total" bson:"tax_total"`
		ShippingTaxTotal string  `json:"shipping_tax_total" bson:"shipping_tax_total"`
		RatePercent      float64 `json:"rate_percent" bson:"rate_percent"`
		MetaData         []struct {
			ID    int    `json:"id" bson:"id"`
			Key   string `json:"key" bson:"key"`
			Value string `json:"value" bson:"value"`
		} `json:"meta_data" bson:"meta_data"`
	} `json:"tax_lines" bson:"tax_lines"`
	ShippingLines []struct {
		ID          int64  `json:"id" bson:"id"`
		MethodTitle string `json:"method_title" bson:"method_title"`
		MethodID    string `json:"method_id" bson:"method_id"`
		InstanceID  string `json:"instance_id" bson:"instance_id"`
		Total       string `json:"total" bson:"total"`
		TotalTax    string `json:"total_tax" bson:"total_tax"`
		Taxes       []struct {
			ID       int    `json:"id" bson:"id"`
			Total    string `json:"total" bson:"total"`
			Subtotal string `json:"subtotal" bson:"subtotal"`
		} `json:"taxes" bson:"taxes"`
		MetaData []struct {
			ID    int    `json:"id" bson:"id"`
			Key   string `json:"key" bson:"key"`
			Value string `json:"value" bson:"value"`
		} `json:"meta_data" bson:"meta_data"`
	} `json:"shipping_lines" bson:"shipping_lines"`
	FeeLines []struct {
		ID        int64  `json:"id" bson:"id"`
		Name      string `json:"name" bson:"name"`
		TaxClass  string `json:"tax_class" bson:"tax_class"`
		TaxStatus string `json:"tax_status" bson:"tax_status"`
		Total     string `json:"total" bson:"total"`
		TotalTax  string `json:"total_tax" bson:"total_tax"`
		Taxes     []struct {
			ID       int    `json:"id" bson:"id"`
			Total    string `json:"total" bson:"total"`
			Subtotal string `json:"subtotal" bson:"subtotal"`
		} `json:"taxes" bson:"taxes"`
		MetaData []struct {
			ID    int    `json:"id" bson:"id"`
			Key   string `json:"key" bson:"key"`
			Value string `json:"value" bson:"value"`
		} `json:"meta_data" bson:"meta_data"`
	} `json:"fee_lines" bson:"fee_lines"`
	CouponLines []struct {
		ID          int64  `json:"id" bson:"id"`
		Code        string `json:"code" bson:"code"`
		Discount    string `json:"discount" bson:"discount"`
		DiscountTax string `json:"discount_tax" bson:"discount_tax"`
		MetaData    []struct {
			ID    int    `json:"id" bson:"id"`
			Key   string `json:"key" bson:"key"`
			Value string `json:"value" bson:"value"`
		} `json:"meta_data" bson:"meta_data"`
	} `json:"coupon_lines" bson:"coupon_lines"`
	Refunds []struct {
		ID     int64  `json:"id" bson:"id"`
		Reason string `json:"reason" bson:"reason"`
		Total  string `json:"total" bson:"total"`
	} `json:"refunds" bson:"refunds"`
	MetaData []struct {
		ID    int    `json:"id" bson:"id"`
		Key   string `json:"key" bson:"key"`
		Value string `json:"value" bson:"value"`
	} `json:"meta_data" bson:"meta_data"`
}

// syncWooCommerceProducts sincroniza productos con WooCommerce
func (s *EcommerceService) syncWooCommerceProducts(ctx context.Context, config *EcommerceConfig) error {
	// Crear cliente HTTP con autenticación básica
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/wp-json/wc/v3/products", config.StoreURL), nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %v", err)
	}

	// Configurar autenticación básica
	req.SetBasicAuth(config.APIKey, config.APISecret)

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta: %s", resp.Status)
	}

	// Decodificar la respuesta
	var products []WooCommerceProduct
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	// Procesar cada producto
	for _, product := range products {
		// Verificar si el producto ya existe en nuestra base de datos
		var existingProduct WooCommerceProduct
		err := s.db.Collection("woocommerce_products").FindOne(ctx, bson.M{"id": product.ID}).Decode(&existingProduct)
		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("error al buscar producto existente: %v", err)
		}

		// Si el producto no existe o ha sido modificado, actualizarlo
		if err == mongo.ErrNoDocuments || existingProduct.DateModified.Before(product.DateModified) {
			_, err = s.db.Collection("woocommerce_products").UpdateOne(
				ctx,
				bson.M{"id": product.ID},
				bson.M{
					"$set": product,
				},
			)
			if err != nil {
				return fmt.Errorf("error al actualizar producto: %v", err)
			}
		}
	}

	// Actualizar la última sincronización
	_, err = s.db.Collection("ecommerce_stores").UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{"$set": bson.M{"last_sync": time.Now()}},
	)
	if err != nil {
		return fmt.Errorf("error al actualizar última sincronización: %v", err)
	}

	return nil
}

// syncWooCommerceOrders sincroniza órdenes con WooCommerce
func (s *EcommerceService) syncWooCommerceOrders(ctx context.Context, config *EcommerceConfig) error {
	// Crear cliente HTTP con autenticación básica
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/wp-json/wc/v3/orders", config.StoreURL), nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %v", err)
	}

	// Configurar autenticación básica
	req.SetBasicAuth(config.APIKey, config.APISecret)

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta: %s", resp.Status)
	}

	// Decodificar la respuesta
	var orders []WooCommerceOrder
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	// Procesar cada orden
	for _, order := range orders {
		// Verificar si la orden ya existe en nuestra base de datos
		var existingOrder WooCommerceOrder
		err := s.db.Collection("woocommerce_orders").FindOne(ctx, bson.M{"id": order.ID}).Decode(&existingOrder)
		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("error al buscar orden existente: %v", err)
		}

		// Si la orden no existe o ha sido modificada, actualizarla
		if err == mongo.ErrNoDocuments || existingOrder.DateModified.Before(order.DateModified) {
			_, err = s.db.Collection("woocommerce_orders").UpdateOne(
				ctx,
				bson.M{"id": order.ID},
				bson.M{
					"$set": order,
				},
			)
			if err != nil {
				return fmt.Errorf("error al actualizar orden: %v", err)
			}
		}
	}

	// Actualizar la última sincronización
	_, err = s.db.Collection("ecommerce_stores").UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{"$set": bson.M{"last_sync": time.Now()}},
	)
	if err != nil {
		return fmt.Errorf("error al actualizar última sincronización: %v", err)
	}

	return nil
}

// JumpsellerProduct representa un producto de Jumpseller
type JumpsellerProduct struct {
	ID          int64     `json:"id" bson:"id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	Stock       int       `json:"stock" bson:"stock"`
	SKU         string    `json:"sku" bson:"sku"`
	Status      string    `json:"status" bson:"status"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// JumpsellerOrder representa una orden de Jumpseller
type JumpsellerOrder struct {
	ID            int64     `json:"id" bson:"id"`
	Number        string    `json:"number" bson:"number"`
	Status        string    `json:"status" bson:"status"`
	Total         float64   `json:"total" bson:"total"`
	Subtotal      float64   `json:"subtotal" bson:"subtotal"`
	Tax           float64   `json:"tax" bson:"tax"`
	Shipping      float64   `json:"shipping" bson:"shipping"`
	CustomerEmail string    `json:"customer_email" bson:"customer_email"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

// syncJumpsellerProducts sincroniza productos con Jumpseller
func (s *EcommerceService) syncJumpsellerProducts(ctx context.Context, config *EcommerceConfig) error {
	// Crear cliente HTTP
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/products", config.StoreURL), nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %v", err)
	}

	// Configurar autenticación con API key
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta: %s", resp.Status)
	}

	// Decodificar la respuesta
	var products []JumpsellerProduct
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	// Procesar cada producto
	for _, product := range products {
		// Verificar si el producto ya existe en nuestra base de datos
		var existingProduct JumpsellerProduct
		err = s.db.Collection("jumpseller_products").FindOne(ctx, bson.M{"id": product.ID}).Decode(&existingProduct)
		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("error al buscar producto existente: %v", err)
		}

		// Si el producto no existe o ha sido modificado, actualizarlo
		if err == mongo.ErrNoDocuments || existingProduct.UpdatedAt.Before(product.UpdatedAt) {
			_, err = s.db.Collection("jumpseller_products").UpdateOne(
				ctx,
				bson.M{"id": product.ID},
				bson.M{"$set": product},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return fmt.Errorf("error al actualizar producto: %v", err)
			}
		}
	}

	// Actualizar la última sincronización
	_, err = s.db.Collection("ecommerce_stores").UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{"$set": bson.M{"last_sync": time.Now()}},
	)
	if err != nil {
		return fmt.Errorf("error al actualizar última sincronización: %v", err)
	}

	return nil
}

// syncJumpsellerOrders sincroniza órdenes con Jumpseller
func (s *EcommerceService) syncJumpsellerOrders(ctx context.Context, config *EcommerceConfig) error {
	// Crear cliente HTTP
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/orders", config.StoreURL), nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %v", err)
	}

	// Configurar autenticación con API key
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta: %s", resp.Status)
	}

	// Decodificar la respuesta
	var orders []JumpsellerOrder
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	// Procesar cada orden
	for _, order := range orders {
		// Verificar si la orden ya existe en nuestra base de datos
		var existingOrder JumpsellerOrder
		err = s.db.Collection("jumpseller_orders").FindOne(ctx, bson.M{"id": order.ID}).Decode(&existingOrder)
		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("error al buscar orden existente: %v", err)
		}

		// Si la orden no existe o ha sido modificada, actualizarla
		if err == mongo.ErrNoDocuments || existingOrder.UpdatedAt.Before(order.UpdatedAt) {
			_, err = s.db.Collection("jumpseller_orders").UpdateOne(
				ctx,
				bson.M{"id": order.ID},
				bson.M{"$set": order},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return fmt.Errorf("error al actualizar orden: %v", err)
			}
		}
	}

	// Actualizar la última sincronización
	_, err = s.db.Collection("ecommerce_stores").UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{"$set": bson.M{"last_sync": time.Now()}},
	)
	if err != nil {
		return fmt.Errorf("error al actualizar última sincronización: %v", err)
	}

	return nil
}

// WixProduct representa un producto de Wix
type WixProduct struct {
	ID          string    `json:"id" bson:"id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	Price       float64   `json:"price" bson:"price"`
	Stock       int       `json:"stock" bson:"stock"`
	SKU         string    `json:"sku" bson:"sku"`
	Status      string    `json:"status" bson:"status"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// WixOrder representa una orden de Wix
type WixOrder struct {
	ID            string    `json:"id" bson:"id"`
	Number        string    `json:"number" bson:"number"`
	Status        string    `json:"status" bson:"status"`
	Total         float64   `json:"total" bson:"total"`
	Subtotal      float64   `json:"subtotal" bson:"subtotal"`
	Tax           float64   `json:"tax" bson:"tax"`
	Shipping      float64   `json:"shipping" bson:"shipping"`
	CustomerEmail string    `json:"customer_email" bson:"customer_email"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

// syncWixProducts sincroniza productos con Wix
func (s *EcommerceService) syncWixProducts(ctx context.Context, config *EcommerceConfig) error {
	// Crear cliente HTTP
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/products", config.StoreURL), nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %v", err)
	}

	// Configurar autenticación OAuth
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta: %s", resp.Status)
	}

	// Decodificar la respuesta
	var products []WixProduct
	if err := json.NewDecoder(resp.Body).Decode(&products); err != nil {
		return fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	// Procesar cada producto
	for _, product := range products {
		// Verificar si el producto ya existe en nuestra base de datos
		var existingProduct WixProduct
		err = s.db.Collection("wix_products").FindOne(ctx, bson.M{"id": product.ID}).Decode(&existingProduct)
		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("error al buscar producto existente: %v", err)
		}

		// Si el producto no existe o ha sido modificado, actualizarlo
		if err == mongo.ErrNoDocuments || existingProduct.UpdatedAt.Before(product.UpdatedAt) {
			_, err = s.db.Collection("wix_products").UpdateOne(
				ctx,
				bson.M{"id": product.ID},
				bson.M{"$set": product},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return fmt.Errorf("error al actualizar producto: %v", err)
			}
		}
	}

	// Actualizar la última sincronización
	_, err = s.db.Collection("ecommerce_stores").UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{"$set": bson.M{"last_sync": time.Now()}},
	)
	if err != nil {
		return fmt.Errorf("error al actualizar última sincronización: %v", err)
	}

	return nil
}

// syncWixOrders sincroniza órdenes con Wix
func (s *EcommerceService) syncWixOrders(ctx context.Context, config *EcommerceConfig) error {
	// Crear cliente HTTP
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v1/orders", config.StoreURL), nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %v", err)
	}

	// Configurar autenticación OAuth
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta: %s", resp.Status)
	}

	// Decodificar la respuesta
	var orders []WixOrder
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	// Procesar cada orden
	for _, order := range orders {
		// Verificar si la orden ya existe en nuestra base de datos
		var existingOrder WixOrder
		err = s.db.Collection("wix_orders").FindOne(ctx, bson.M{"id": order.ID}).Decode(&existingOrder)
		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("error al buscar orden existente: %v", err)
		}

		// Si la orden no existe o ha sido modificada, actualizarla
		if err == mongo.ErrNoDocuments || existingOrder.UpdatedAt.Before(order.UpdatedAt) {
			_, err = s.db.Collection("wix_orders").UpdateOne(
				ctx,
				bson.M{"id": order.ID},
				bson.M{"$set": order},
				options.Update().SetUpsert(true),
			)
			if err != nil {
				return fmt.Errorf("error al actualizar orden: %v", err)
			}
		}
	}

	// Actualizar la última sincronización
	_, err = s.db.Collection("ecommerce_stores").UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{"$set": bson.M{"last_sync": time.Now()}},
	)
	if err != nil {
		return fmt.Errorf("error al actualizar última sincronización: %v", err)
	}

	return nil
}

// syncPrestashopOrders sincroniza órdenes con PrestaShop
func (s *EcommerceService) syncPrestashopOrders(ctx context.Context, config *EcommerceConfig) error {
	// Crear cliente HTTP con autenticación básica
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/orders", config.StoreURL), nil)
	if err != nil {
		return fmt.Errorf("error al crear request: %v", err)
	}

	// Configurar autenticación básica
	req.SetBasicAuth(config.APIKey, "")

	// Realizar la petición
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error al realizar petición: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error en la respuesta: %s", resp.Status)
	}

	// Decodificar la respuesta
	var orders []struct {
		ID int64 `json:"id" bson:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return fmt.Errorf("error al decodificar respuesta: %v", err)
	}

	// Procesar cada orden
	for _, order := range orders {
		// Obtener los detalles de la orden
		detailsReq, err := http.NewRequest("GET", fmt.Sprintf("%s/api/orders/%d/order_details", config.StoreURL, order.ID), nil)
		if err != nil {
			return fmt.Errorf("error al crear request de detalles: %v", err)
		}
		detailsReq.SetBasicAuth(config.APIKey, "")

		detailsResp, err := client.Do(detailsReq)
		if err != nil {
			return fmt.Errorf("error al obtener detalles de la orden: %v", err)
		}
		defer detailsResp.Body.Close()

		// Verificar si la orden ya existe en nuestra base de datos
		var existingOrder struct {
			ID int64 `json:"id" bson:"id"`
		}
		err = s.db.Collection("prestashop_orders").FindOne(ctx, bson.M{"id": order.ID}).Decode(&existingOrder)
		if err != nil && err != mongo.ErrNoDocuments {
			return fmt.Errorf("error al buscar orden existente: %v", err)
		}

		// Si la orden no existe o ha sido modificada, actualizarla
		if err == mongo.ErrNoDocuments {
			_, err = s.db.Collection("prestashop_orders").InsertOne(ctx, order)
			if err != nil {
				return fmt.Errorf("error al insertar orden: %v", err)
			}
		}
	}

	// Actualizar la última sincronización
	_, err = s.db.Collection("ecommerce_stores").UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{"$set": bson.M{"last_sync": time.Now()}},
	)
	if err != nil {
		return fmt.Errorf("error al actualizar última sincronización: %v", err)
	}

	return nil
}

// SyncShopifyOrders sincroniza órdenes desde Shopify
func (s *EcommerceService) SyncShopifyOrders(ctx context.Context, shopID string) error {
	// Implementar lógica de sincronización con Shopify
	orders, err := s.fetchShopifyOrders(ctx, shopID)
	if err != nil {
		return fmt.Errorf("error obteniendo órdenes de Shopify: %v", err)
	}

	for _, order := range orders {
		if err := s.processOrder(ctx, order); err != nil {
			return fmt.Errorf("error procesando orden: %v", err)
		}
	}

	return nil
}

// SyncWooCommerceOrders sincroniza órdenes desde WooCommerce
func (s *EcommerceService) SyncWooCommerceOrders(ctx context.Context, storeID string) error {
	// Implementar lógica de sincronización con WooCommerce
	orders, err := s.fetchWooCommerceOrders(ctx, storeID)
	if err != nil {
		return fmt.Errorf("error obteniendo órdenes de WooCommerce: %v", err)
	}

	for _, order := range orders {
		if err := s.processOrder(ctx, order); err != nil {
			return fmt.Errorf("error procesando orden: %v", err)
		}
	}

	return nil
}

// SyncMercadoLibreOrders sincroniza órdenes desde MercadoLibre
func (s *EcommerceService) SyncMercadoLibreOrders(ctx context.Context, sellerID string) error {
	// Implementar lógica de sincronización con MercadoLibre
	orders, err := s.fetchMercadoLibreOrders(ctx, sellerID)
	if err != nil {
		return fmt.Errorf("error obteniendo órdenes de MercadoLibre: %v", err)
	}

	for _, order := range orders {
		if err := s.processOrder(ctx, order); err != nil {
			return fmt.Errorf("error procesando orden: %v", err)
		}
	}

	return nil
}

// fetchShopifyOrders obtiene órdenes desde Shopify
func (s *EcommerceService) fetchShopifyOrders(ctx context.Context, shopID string) ([]map[string]interface{}, error) {
	// Implementar lógica para obtener órdenes de Shopify
	return nil, nil
}

// fetchWooCommerceOrders obtiene órdenes desde WooCommerce
func (s *EcommerceService) fetchWooCommerceOrders(ctx context.Context, storeID string) ([]map[string]interface{}, error) {
	// Implementar lógica para obtener órdenes de WooCommerce
	return nil, nil
}

// fetchMercadoLibreOrders obtiene órdenes desde MercadoLibre
func (s *EcommerceService) fetchMercadoLibreOrders(ctx context.Context, sellerID string) ([]map[string]interface{}, error) {
	// Implementar lógica para obtener órdenes de MercadoLibre
	return nil, nil
}

// processOrder procesa una orden individual
func (s *EcommerceService) processOrder(ctx context.Context, order map[string]interface{}) error {
	// Validar la orden
	if err := s.validateOrder(order); err != nil {
		return fmt.Errorf("error validando orden: %v", err)
	}

	// Generar documento tributario
	doc, err := s.generateTaxDocument(ctx, order)
	if err != nil {
		return fmt.Errorf("error generando documento tributario: %v", err)
	}

	// Almacenar en caché
	s.cache.SetDocument(doc.ID.Hex(), &models.DocumentoAlmacenado{
		ID: doc.ID,
		Metadata: map[string]interface{}{
			"order_id": order["id"],
			"platform": order["platform"],
		},
		Contenido: doc.Contenido,
	})

	return nil
}

// validateOrder valida una orden
func (s *EcommerceService) validateOrder(order map[string]interface{}) error {
	// Implementar validaciones necesarias
	return nil
}

// generateTaxDocument genera un documento tributario a partir de una orden
func (s *EcommerceService) generateTaxDocument(ctx context.Context, order map[string]interface{}) (*models.DocumentoAlmacenado, error) {
	// Implementar generación de documento tributario
	return nil, nil
}

// FindDocument busca un documento por su ID
func (s *EcommerceService) FindDocument(ctx context.Context, id primitive.ObjectID) (*models.DocumentoAlmacenado, error) {
	// Intentar obtener del caché primero
	if doc, err := s.cache.GetDocument(id.Hex()); err == nil {
		return doc, nil
	}

	// Si no está en caché, buscar en la base de datos
	var doc models.DocumentoAlmacenado
	err := s.db.Collection("documentos").FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		return nil, fmt.Errorf("error buscando documento: %v", err)
	}

	// Almacenar en caché
	s.cache.SetDocument(id.Hex(), &doc)

	return &doc, nil
}

// UpdateDocument actualiza un documento existente
func (s *EcommerceService) UpdateDocument(ctx context.Context, id primitive.ObjectID, updates map[string]interface{}) error {
	// Actualizar en la base de datos
	update := bson.M{"$set": updates}
	_, err := s.db.Collection("documentos").UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("error actualizando documento: %v", err)
	}

	// Actualizar en caché si existe
	if doc, err := s.cache.GetDocument(id.Hex()); err == nil {
		for k, v := range updates {
			doc.Metadata[k] = v
		}
		s.cache.SetDocument(id.Hex(), doc)
	}

	return nil
}

// DeleteDocument elimina un documento
func (s *EcommerceService) DeleteDocument(ctx context.Context, id primitive.ObjectID) error {
	// Eliminar de la base de datos
	_, err := s.db.Collection("documentos").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("error eliminando documento: %v", err)
	}

	// Eliminar de caché
	delete(s.cache.documentCache, id.Hex())

	return nil
}

// ListDocuments lista documentos con filtros
func (s *EcommerceService) ListDocuments(ctx context.Context, filter map[string]interface{}, limit int64) ([]models.DocumentoAlmacenado, error) {
	cursor, err := s.db.Collection("documentos").
		Find(ctx, filter, options.Find().SetLimit(limit))
	if err != nil {
		return nil, fmt.Errorf("error listando documentos: %v", err)
	}
	defer cursor.Close(ctx)

	var docs []models.DocumentoAlmacenado
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, fmt.Errorf("error decodificando documentos: %v", err)
	}

	return docs, nil
}

// ExportToJSON exporta un documento a formato JSON
func (s *EcommerceService) ExportToJSON(doc *models.DocumentoAlmacenado) ([]byte, error) {
	return json.Marshal(doc)
}

// ImportFromJSON importa un documento desde JSON
func (s *EcommerceService) ImportFromJSON(data []byte) (*models.DocumentoAlmacenado, error) {
	var doc models.DocumentoAlmacenado
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	return &doc, nil
}

// CleanExpiredCache limpia el caché expirado
func (s *EcommerceService) CleanExpiredCache() {
	now := time.Now()
	for key, doc := range s.cache.documentCache {
		if now.After(doc.CacheInfo.ExpiresAt) {
			delete(s.cache.documentCache, key)
		}
	}
}

// GetCacheStats retorna estadísticas del caché
func (s *EcommerceService) GetCacheStats() map[string]interface{} {
	return map[string]interface{}{
		"total_documents": len(s.cache.documentCache),
		"ttl":             s.cache.ttl.String(),
	}
}

// Close limpia recursos del servicio
func (s *EcommerceService) Close() {
	s.CleanExpiredCache()
}
