# Tareas de GitHub - Fase de Certificación

## Issues Actuales

### Milestone: Certificación SII

#### Set de Pruebas (#1)
- [ ] Implementar casos de prueba según guía SII
- [ ] Crear datos de prueba
- [ ] Documentar procedimientos
- [ ] Validar resultados

#### Ambiente de Certificación (#2)
- [ ] Configurar ambiente
- [ ] Instalar certificados
- [ ] Validar conectividad
- [ ] Configurar monitoreo

#### Validaciones (#3)
- [ ] Implementar validaciones de negocio
- [ ] Validar folios y CAF
- [ ] Verificar timbraje
- [ ] Documentar reglas

#### Documentación (#4)
- [ ] Crear guía de certificación
- [ ] Documentar procedimientos
- [ ] Crear manual de operación
- [ ] Actualizar README

## Pull Requests

### En Revisión
1. PR #101: Implementación de pruebas unitarias
   - ✅ Tests del cliente SOAP
   - ✅ Tests de autenticación
   - ✅ Tests del cliente DTE

2. PR #102: Validación XSD
   - ✅ ValidatorService
   - ✅ Middleware
   - ✅ Documentación

### Pendientes
1. PR #103: Set de pruebas de certificación
   - 🔄 Casos de prueba SII
   - ⏳ Pruebas de integración
   - ⏳ Documentación

2. PR #104: Ambiente de certificación
   - 🔄 Configuración
   - ⏳ Certificados
   - ⏳ Monitoreo

## Labels

- `certificacion`: Issues relacionados con el proceso de certificación
- `pruebas`: Issues de implementación de pruebas
- `documentacion`: Issues de documentación
- `validacion`: Issues de validación
- `ambiente`: Issues de configuración de ambiente

## Projects

### Tablero: Certificación SII
- **To Do**
  - Configurar ambiente de certificación
  - Implementar validaciones pendientes
  - Crear documentación

- **In Progress**
  - Set de pruebas de certificación
  - Configuración de certificados
  - Validaciones de negocio

- **Review**
  - Pruebas unitarias
  - Validación XSD
  - Documentación base

- **Done**
  - Cliente SOAP
  - Cliente de autenticación
  - Cliente DTE
  - Firma digital

## Workflow

1. **Creación de Issues**
   - Usar template correspondiente
   - Asignar milestone "Certificación"
   - Aplicar labels relevantes
   - Asignar al proyecto

2. **Pull Requests**
   - Referenciar issue relacionado
   - Incluir tests
   - Actualizar documentación
   - Solicitar review

3. **Code Review**
   - Verificar estándares
   - Validar tests
   - Revisar documentación
   - Aprobar cambios

4. **Merge**
   - Squash and merge
   - Eliminar rama
   - Cerrar issues relacionados
   - Actualizar proyecto

## Timeline

1. **Semana 1-2**
   - [ ] Issues #1, #2
   - [ ] PRs #101, #102
   - [ ] Configuración inicial

2. **Semana 3-4**
   - [ ] Issues #3
   - [ ] PR #103
   - [ ] Implementación

3. **Semana 5-6**
   - [ ] Issue #4
   - [ ] PR #104
   - [ ] Validación

4. **Semana 7-8**
   - [ ] Certificación
   - [ ] Documentación final
   - [ ] Release 