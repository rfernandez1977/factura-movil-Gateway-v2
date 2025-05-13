package services

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cursor/FMgo/models"

	"crypto/tls"

	"github.com/jlaffaye/ftp"
	"github.com/pkg/sftp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/ssh"
)

// LegacyService maneja la integración con sistemas legacy
type LegacyService struct {
	db *mongo.Database
}

// NewLegacyService crea una nueva instancia del servicio legacy
func NewLegacyService(db *mongo.Database) *LegacyService {
	return &LegacyService{
		db: db,
	}
}

// RegistrarConfiguracionArchivoPlano registra una nueva configuración de archivo plano
func (s *LegacyService) RegistrarConfiguracionArchivoPlano(ctx context.Context, config *models.ConfiguracionArchivoPlano) error {
	config.ID = generateID()
	config.FechaCreacion = time.Now()
	config.FechaActualizacion = time.Now()

	_, err := s.db.Collection("configuraciones_archivos_planos").InsertOne(ctx, config)
	if err != nil {
		return fmt.Errorf("error al registrar configuración de archivo plano: %v", err)
	}

	return nil
}

// RegistrarConfiguracionProtocolo registra una nueva configuración de protocolo
func (s *LegacyService) RegistrarConfiguracionProtocolo(ctx context.Context, config *models.ConfiguracionProtocolo) error {
	config.ID = generateID()
	config.FechaCreacion = time.Now()
	config.FechaActualizacion = time.Now()

	_, err := s.db.Collection("configuraciones_protocolos").InsertOne(ctx, config)
	if err != nil {
		return fmt.Errorf("error al registrar configuración de protocolo: %v", err)
	}

	return nil
}

// RegistrarTransformacionLegacy registra una nueva transformación legacy
func (s *LegacyService) RegistrarTransformacionLegacy(ctx context.Context, transformacion *models.TransformacionLegacy) error {
	transformacion.ID = generateID()
	transformacion.FechaCreacion = time.Now()
	transformacion.FechaActualizacion = time.Now()

	_, err := s.db.Collection("transformaciones_legacy").InsertOne(ctx, transformacion)
	if err != nil {
		return fmt.Errorf("error al registrar transformación legacy: %v", err)
	}

	return nil
}

// ProcesarArchivoPlano procesa un archivo plano según su configuración
func (s *LegacyService) ProcesarArchivoPlano(ctx context.Context, erpID string, archivo string) error {
	// Obtener configuración
	var config models.ConfiguracionArchivoPlano
	err := s.db.Collection("configuraciones_archivos_planos").FindOne(ctx, bson.M{"erp_id": erpID}).Decode(&config)
	if err != nil {
		return fmt.Errorf("error al obtener configuración de archivo plano: %v", err)
	}

	// Abrir archivo
	file, err := os.Open(archivo)
	if err != nil {
		return fmt.Errorf("error al abrir archivo: %v", err)
	}
	defer file.Close()

	// Procesar según formato
	switch config.Formato {
	case models.FormatoCSV:
		return s.procesarCSV(file, config)
	case models.FormatoTXT:
		return s.procesarTXT(file, config)
	case models.FormatoFIXED:
		return s.procesarFixed(file, config)
	case models.FormatoXML:
		return s.procesarXML(file, config)
	case models.FormatoJSON:
		return s.procesarJSON(file, config)
	default:
		return fmt.Errorf("formato no soportado: %s", config.Formato)
	}
}

// procesarCSV procesa un archivo CSV
func (s *LegacyService) procesarCSV(file io.Reader, config models.ConfiguracionArchivoPlano) error {
	reader := csv.NewReader(file)
	reader.Comma = rune(config.Delimitador[0])

	// Leer cabecera si existe
	if config.IncluirCabecera {
		_, err := reader.Read()
		if err != nil {
			return fmt.Errorf("error al leer cabecera: %v", err)
		}
	}

	// Procesar registros
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error al leer registro: %v", err)
		}

		// Transformar y validar datos
		datos, err := s.transformarDatos(record, config)
		if err != nil {
			log.Printf("Error al transformar datos: %v", err)
			continue
		}

		// Guardar en base de datos
		_, err = s.db.Collection("datos_legacy").InsertOne(context.Background(), datos)
		if err != nil {
			log.Printf("Error al guardar datos: %v", err)
		}
	}

	return nil
}

// procesarTXT procesa un archivo de texto plano
func (s *LegacyService) procesarTXT(file io.Reader, config models.ConfiguracionArchivoPlano) error {
	scanner := bufio.NewScanner(file)
	linea := 0

	for scanner.Scan() {
		linea++
		texto := scanner.Text()

		// Si es la primera línea y hay cabecera, la saltamos
		if linea == 1 && config.IncluirCabecera {
			continue
		}

		// Procesar la línea según el delimitador
		campos := strings.Split(texto, config.Delimitador)

		// Transformar y validar datos
		datos, err := s.transformarDatos(campos, config)
		if err != nil {
			log.Printf("Error al transformar datos en línea %d: %v", linea, err)
			continue
		}

		// Guardar en base de datos
		_, err = s.db.Collection("datos_legacy").InsertOne(context.Background(), datos)
		if err != nil {
			log.Printf("Error al guardar datos en línea %d: %v", linea, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error al leer archivo TXT: %v", err)
	}

	return nil
}

// procesarFixed procesa un archivo de ancho fijo
func (s *LegacyService) procesarFixed(file io.Reader, config models.ConfiguracionArchivoPlano) error {
	scanner := bufio.NewScanner(file)
	linea := 0

	// Obtener mapeo de campos y sus posiciones
	mapeoCampos, err := s.obtenerMapeoCamposFixed(config)
	if err != nil {
		return fmt.Errorf("error al obtener mapeo de campos: %v", err)
	}

	for scanner.Scan() {
		linea++
		texto := scanner.Text()

		// Si es la primera línea y hay cabecera, la saltamos
		if linea == 1 && config.IncluirCabecera {
			continue
		}

		// Extraer campos según el mapeo
		campos := make([]string, len(mapeoCampos))
		for i, mapeo := range mapeoCampos {
			if mapeo.Inicio+mapeo.Longitud <= len(texto) {
				campos[i] = strings.TrimSpace(texto[mapeo.Inicio : mapeo.Inicio+mapeo.Longitud])
			}
		}

		// Transformar y validar datos
		datos, err := s.transformarDatos(campos, config)
		if err != nil {
			log.Printf("Error al transformar datos en línea %d: %v", linea, err)
			continue
		}

		// Guardar en base de datos
		_, err = s.db.Collection("datos_legacy").InsertOne(context.Background(), datos)
		if err != nil {
			log.Printf("Error al guardar datos en línea %d: %v", linea, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error al leer archivo de ancho fijo: %v", err)
	}

	return nil
}

// procesarXML procesa un archivo XML
func (s *LegacyService) procesarXML(file io.Reader, config models.ConfiguracionArchivoPlano) error {
	decoder := xml.NewDecoder(file)

	for {
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error al decodificar XML: %v", err)
		}

		switch se := token.(type) {
		case xml.StartElement:
			if se.Name.Local == "registro" {
				var datos map[string]interface{}
				if err := decoder.DecodeElement(&datos, &se); err != nil {
					log.Printf("Error al decodificar elemento XML: %v", err)
					continue
				}

				// Transformar y validar datos
				datosTransformados, err := s.transformarDatosXML(datos, config)
				if err != nil {
					log.Printf("Error al transformar datos XML: %v", err)
					continue
				}

				// Guardar en base de datos
				_, err = s.db.Collection("datos_legacy").InsertOne(context.Background(), datosTransformados)
				if err != nil {
					log.Printf("Error al guardar datos XML: %v", err)
				}
			}
		}
	}

	return nil
}

// procesarJSON procesa un archivo JSON
func (s *LegacyService) procesarJSON(file io.Reader, config models.ConfiguracionArchivoPlano) error {
	decoder := json.NewDecoder(file)

	// Verificar si es un array o un objeto
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("error al leer token JSON: %v", err)
	}

	if delim, ok := token.(json.Delim); ok && delim == '[' {
		// Es un array
		for decoder.More() {
			var datos map[string]interface{}
			if err := decoder.Decode(&datos); err != nil {
				log.Printf("Error al decodificar JSON: %v", err)
				continue
			}

			// Transformar y validar datos
			datosTransformados, err := s.transformarDatosJSON(datos, config)
			if err != nil {
				log.Printf("Error al transformar datos JSON: %v", err)
				continue
			}

			// Guardar en base de datos
			_, err = s.db.Collection("datos_legacy").InsertOne(context.Background(), datosTransformados)
			if err != nil {
				log.Printf("Error al guardar datos JSON: %v", err)
			}
		}
	} else {
		// Es un objeto
		var datos map[string]interface{}
		if err := decoder.Decode(&datos); err != nil {
			return fmt.Errorf("error al decodificar JSON: %v", err)
		}

		// Transformar y validar datos
		datosTransformados, err := s.transformarDatosJSON(datos, config)
		if err != nil {
			return fmt.Errorf("error al transformar datos JSON: %v", err)
		}

		// Guardar en base de datos
		_, err = s.db.Collection("datos_legacy").InsertOne(context.Background(), datosTransformados)
		if err != nil {
			return fmt.Errorf("error al guardar datos JSON: %v", err)
		}
	}

	return nil
}

// obtenerMapeoCamposFixed obtiene el mapeo de campos para archivos de ancho fijo
func (s *LegacyService) obtenerMapeoCamposFixed(config models.ConfiguracionArchivoPlano) ([]struct {
	Nombre   string
	Inicio   int
	Longitud int
}, error) {
	// Obtener mapeo de la base de datos
	var mapeo []struct {
		Nombre   string
		Inicio   int
		Longitud int
	}

	cursor, err := s.db.Collection("mapeo_campos_fixed").Find(context.Background(), bson.M{"erp_id": config.ERPID})
	if err != nil {
		return nil, fmt.Errorf("error al obtener mapeo de campos: %v", err)
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &mapeo); err != nil {
		return nil, fmt.Errorf("error al decodificar mapeo de campos: %v", err)
	}

	return mapeo, nil
}

// transformarDatos transforma los datos según la configuración
func (s *LegacyService) transformarDatos(record []string, config models.ConfiguracionArchivoPlano) (map[string]interface{}, error) {
	// Obtener transformaciones
	var transformaciones []models.TransformacionLegacy
	cursor, err := s.db.Collection("transformaciones_legacy").Find(context.Background(), bson.M{"erp_id": config.ERPID})
	if err != nil {
		return nil, fmt.Errorf("error al obtener transformaciones: %v", err)
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &transformaciones); err != nil {
		return nil, fmt.Errorf("error al decodificar transformaciones: %v", err)
	}

	// Aplicar transformaciones
	resultado := make(map[string]interface{})
	for _, t := range transformaciones {
		valor, err := s.aplicarTransformacion(record, t)
		if err != nil {
			return nil, fmt.Errorf("error al transformar campo %s: %v", t.CampoOrigen, err)
		}
		resultado[t.CampoDestino] = valor
	}

	return resultado, nil
}

// aplicarTransformacion aplica una transformación específica
func (s *LegacyService) aplicarTransformacion(record []string, t models.TransformacionLegacy) (interface{}, error) {
	// Implementar lógica de transformación según tipo
	switch t.TipoTransformacion {
	case "FECHA":
		return s.transformarFecha(record[0], t.Parametros)
	case "NUMERO":
		return s.transformarNumero(record[0], t.Parametros)
	case "TEXTO":
		return s.transformarTexto(record[0], t.Parametros)
	default:
		return record[0], nil
	}
}

// transformarFecha transforma una fecha según el formato
func (s *LegacyService) transformarFecha(valor string, parametros map[string]string) (time.Time, error) {
	formato := parametros["formato"]
	if formato == "" {
		formato = "2006-01-02"
	}
	return time.Parse(formato, valor)
}

// transformarNumero transforma un número según el formato
func (s *LegacyService) transformarNumero(valor string, parametros map[string]string) (float64, error) {
	// Implementar transformación de números
	return 0, fmt.Errorf("transformación de números no implementada")
}

// transformarTexto transforma un texto según las reglas
func (s *LegacyService) transformarTexto(valor string, parametros map[string]string) (string, error) {
	// Implementar transformación de texto
	return valor, nil
}

// transformarDatosXML transforma los datos XML según la configuración
func (s *LegacyService) transformarDatosXML(datos map[string]interface{}, config models.ConfiguracionArchivoPlano) (map[string]interface{}, error) {
	resultado := make(map[string]interface{})

	// Obtener transformaciones
	var transformaciones []models.TransformacionLegacy
	cursor, err := s.db.Collection("transformaciones_legacy").Find(context.Background(), bson.M{"erp_id": config.ERPID})
	if err != nil {
		return nil, fmt.Errorf("error al obtener transformaciones: %v", err)
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &transformaciones); err != nil {
		return nil, fmt.Errorf("error al decodificar transformaciones: %v", err)
	}

	// Aplicar transformaciones
	for _, t := range transformaciones {
		if valor, ok := datos[t.CampoOrigen]; ok {
			valorTransformado, err := s.aplicarTransformacion([]string{fmt.Sprintf("%v", valor)}, t)
			if err != nil {
				return nil, fmt.Errorf("error al transformar campo %s: %v", t.CampoOrigen, err)
			}
			resultado[t.CampoDestino] = valorTransformado
		}
	}

	return resultado, nil
}

// transformarDatosJSON transforma los datos JSON según la configuración
func (s *LegacyService) transformarDatosJSON(datos map[string]interface{}, config models.ConfiguracionArchivoPlano) (map[string]interface{}, error) {
	resultado := make(map[string]interface{})

	// Obtener transformaciones
	var transformaciones []models.TransformacionLegacy
	cursor, err := s.db.Collection("transformaciones_legacy").Find(context.Background(), bson.M{"erp_id": config.ERPID})
	if err != nil {
		return nil, fmt.Errorf("error al obtener transformaciones: %v", err)
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &transformaciones); err != nil {
		return nil, fmt.Errorf("error al decodificar transformaciones: %v", err)
	}

	// Aplicar transformaciones
	for _, t := range transformaciones {
		if valor, ok := datos[t.CampoOrigen]; ok {
			valorTransformado, err := s.aplicarTransformacion([]string{fmt.Sprintf("%v", valor)}, t)
			if err != nil {
				return nil, fmt.Errorf("error al transformar campo %s: %v", t.CampoOrigen, err)
			}
			resultado[t.CampoDestino] = valorTransformado
		}
	}

	return resultado, nil
}

// TransferirArchivo transfiere un archivo usando el protocolo configurado
func (s *LegacyService) TransferirArchivo(ctx context.Context, erpID string, archivo string) error {
	// Obtener configuración de protocolo
	var config models.ConfiguracionProtocolo
	err := s.db.Collection("configuraciones_protocolos").FindOne(ctx, bson.M{"erp_id": erpID}).Decode(&config)
	if err != nil {
		return fmt.Errorf("error al obtener configuración de protocolo: %v", err)
	}

	// Transferir según protocolo
	switch config.Protocolo {
	case models.ProtocoloFTP:
		return s.transferirFTP(archivo, config)
	case models.ProtocoloSFTP:
		return s.transferirSFTP(archivo, config)
	case models.ProtocoloFTPS:
		return s.transferirFTPS(archivo, config)
	default:
		return fmt.Errorf("protocolo no soportado: %s", config.Protocolo)
	}
}

// transferirFTP transfiere un archivo por FTP
func (s *LegacyService) transferirFTP(archivo string, config models.ConfiguracionProtocolo) error {
	client, err := ftp.Dial(fmt.Sprintf("%s:%d", config.Host, config.Puerto))
	if err != nil {
		return fmt.Errorf("error al conectar a FTP: %v", err)
	}
	defer client.Quit()

	err = client.Login(config.Usuario, config.Contrasena)
	if err != nil {
		return fmt.Errorf("error al autenticar en FTP: %v", err)
	}

	file, err := os.Open(archivo)
	if err != nil {
		return fmt.Errorf("error al abrir archivo: %v", err)
	}
	defer file.Close()

	err = client.Stor(filepath.Join(config.Ruta, filepath.Base(archivo)), file)
	if err != nil {
		return fmt.Errorf("error al transferir archivo: %v", err)
	}

	return nil
}

// transferirSFTP transfiere un archivo por SFTP
func (s *LegacyService) transferirSFTP(archivo string, config models.ConfiguracionProtocolo) error {
	// Crear configuración del cliente SFTP
	clientConfig := &ssh.ClientConfig{
		User: config.Usuario,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Contrasena),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Duration(config.Timeout) * time.Second,
	}

	// Conectar al servidor SFTP
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Puerto), clientConfig)
	if err != nil {
		return fmt.Errorf("error al conectar a SFTP: %v", err)
	}
	defer conn.Close()

	// Crear cliente SFTP
	client, err := sftp.NewClient(conn)
	if err != nil {
		return fmt.Errorf("error al crear cliente SFTP: %v", err)
	}
	defer client.Close()

	// Abrir archivo local
	srcFile, err := os.Open(archivo)
	if err != nil {
		return fmt.Errorf("error al abrir archivo local: %v", err)
	}
	defer srcFile.Close()

	// Crear archivo remoto
	dstFile, err := client.Create(filepath.Join(config.Ruta, filepath.Base(archivo)))
	if err != nil {
		return fmt.Errorf("error al crear archivo remoto: %v", err)
	}
	defer dstFile.Close()

	// Copiar archivo
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("error al copiar archivo: %v", err)
	}

	return nil
}

// transferirFTPS transfiere un archivo por FTPS
func (s *LegacyService) transferirFTPS(archivo string, config models.ConfiguracionProtocolo) error {
	// Crear configuración TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Crear cliente FTPS
	client, err := ftp.Dial(fmt.Sprintf("%s:%d", config.Host, config.Puerto), ftp.DialWithTLS(tlsConfig))
	if err != nil {
		return fmt.Errorf("error al conectar a FTPS: %v", err)
	}
	defer client.Quit()

	// Autenticar
	err = client.Login(config.Usuario, config.Contrasena)
	if err != nil {
		return fmt.Errorf("error al autenticar en FTPS: %v", err)
	}

	// Abrir archivo local
	file, err := os.Open(archivo)
	if err != nil {
		return fmt.Errorf("error al abrir archivo local: %v", err)
	}
	defer file.Close()

	// Transferir archivo
	err = client.Stor(filepath.Join(config.Ruta, filepath.Base(archivo)), file)
	if err != nil {
		return fmt.Errorf("error al transferir archivo: %v", err)
	}

	return nil
}
