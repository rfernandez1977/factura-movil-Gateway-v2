package handlers

type NotificationHandlers struct {
    client *api.FacturaMovilClient
}

type Notification struct {
    Type        string    `json:"type"`     // EMAIL, SMS, PUSH
    Priority    string    `json:"priority"` // HIGH, MEDIUM, LOW
    Template    string    `json:"template"`
    Recipients  []string  `json:"recipients"`
    Data        interface{} `json:"data"`
}

func (h *NotificationHandlers) SendNotificationHandler(c *gin.Context) {
    // Envío de notificaciones
    // Gestión de plantillas
    // Seguimiento de entregas
}