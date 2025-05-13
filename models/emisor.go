package models

// Emisor representa al emisor de un documento tributario
type Emisor struct {
	RUT         string   `json:"rut" bson:"rut"`
	RazonSocial string   `json:"razon_social" bson:"razon_social"`
	Giro        string   `json:"giro" bson:"giro"`
	Direccion   string   `json:"direccion" bson:"direccion"`
	Comuna      string   `json:"comuna" bson:"comuna"`
	Ciudad      string   `json:"ciudad" bson:"ciudad"`
	Correo      string   `json:"correo" bson:"correo"`
	Actecos     []string `json:"actecos" bson:"actecos"`
}
