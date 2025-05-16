# Políticas de Permisos y Accesos - FMgo

## Estructura de Permisos

### Directorios Críticos

#### Certificados y Llaves (0600)
```bash
dev/config/certs/sii/
dev/config/certs/firma/
test/config/certs/sii/
test/config/certs/firma/
```
- Solo lectura para el proceso de la aplicación
- Sin acceso para otros usuarios
- Respaldo cifrado obligatorio

#### Archivos CAF (0640)
```bash
dev/config/caf/
test/config/caf/
```
- Lectura y escritura para el proceso de la aplicación
- Lectura para procesos de respaldo
- Sin acceso para otros usuarios

#### Configuración (0644)
```bash
dev/config/
test/config/
```
- Lectura para todos los usuarios del grupo
- Escritura solo para administradores
- Variables de entorno protegidas

### Directorios de Desarrollo

#### Código Fuente (0644)
```bash
dev/dte/
dev/firma/
dev/sii/
```
- Lectura para todos los desarrolladores
- Escritura controlada por git

#### Scripts (0755)
```bash
dev/scripts/
```
- Ejecutable para desarrolladores
- Modificación solo por administradores

### Directorios de Pruebas

#### Datos de Prueba (0644)
```bash
test/test_data/
```
- Acceso completo para desarrolladores
- Lectura para procesos de CI/CD

## Grupos y Roles

### Desarrolladores
- Acceso de lectura a todo el código
- Acceso de escritura a través de git
- Sin acceso directo a certificados y llaves

### Administradores
- Acceso completo a todos los directorios
- Gestión de certificados y llaves
- Gestión de configuración

### Procesos de Aplicación
- Acceso de lectura a certificados
- Acceso de lectura/escritura a CAF
- Acceso de lectura a configuración

## Políticas de Seguridad

### Certificados y Llaves
1. Almacenamiento cifrado en reposo
2. Rotación programada
3. Respaldos seguros
4. Registro de accesos

### Archivos CAF
1. Monitoreo de uso
2. Respaldo automático
3. Rotación programada
4. Verificación de integridad

### Configuración
1. Variables sensibles cifradas
2. Sin credenciales en código
3. Auditoría de cambios
4. Control de versiones

## Implementación

### Pasos de Configuración
1. Crear grupos de usuarios
2. Establecer permisos base
3. Configurar ACLs
4. Verificar permisos
5. Documentar cambios

### Monitoreo
1. Registro de accesos
2. Alertas de modificación
3. Auditoría periódica
4. Reportes de seguridad

## Mantenimiento

### Tareas Periódicas
1. Revisión de permisos
2. Actualización de grupos
3. Rotación de credenciales
4. Verificación de integridad

### Respaldos
1. Certificados y llaves
2. Archivos CAF
3. Configuración
4. Logs de acceso 