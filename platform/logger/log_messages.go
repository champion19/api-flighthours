package logger

// ============================================
// CENTRALIZED LOG MESSAGES
// ============================================
// All log messages are defined here as constants to:
// 1. Avoid hardcoded strings across the codebase
// 2. Enable easy maintenance and updates
// 3. Support future internationalization
// 4. Ensure consistency

// ============================================
// APPLICATION LIFECYCLE
// ============================================
const (
	LogAppStarting          = "Iniciando aplicación FlighHours Backend"
	LogAppConfigLoaded      = "Configuración cargada exitosamente"
	LogAppConfigError       = "Error cargando configuración"
	LogAppDatabaseConnected = "Conexión a base de datos establecida"
	LogAppDatabaseError     = "Error conectando a base de datos"
	LogAppDatabasePingOK    = "Ping a base de datos exitoso"
	LogAppDatabasePingError = "Error haciendo ping a base de datos"
	LogAppServerStarting    = "Servidor iniciando"
	LogAppServerListening   = "Servidor escuchando"
	LogAppServerStartError  = "Error iniciando servidor"
	LogAppShuttingDown      = "Apagando aplicación gracefully"
)

// ============================================
// REGISTRATION / AUTHENTICATION
// ============================================
const (
	LogRegRequestReceived   = "Nueva solicitud de registro recibida"
	LogRegProcessing        = "Procesando registro de usuario"
	LogRegJSONParseError    = "Error parseando JSON del request"
	LogRegJSONBindError     = "Error vinculando JSON del request"
	LogRegProcessError      = "Error en proceso de registro"
	LogRegIDEncodeError     = "Error ofuscando ID de usuario"
	LogRegSuccess           = "Usuario registrado exitosamente"
	LogRegKeycloakSync      = "Sincronizando usuario con Keycloak"
	LogRegKeycloakSyncError = "Error sincronizando con Keycloak"
)

// ============================================
// DATABASE OPERATIONS
// ============================================
const (
	LogDBQueryExecuting       = "Ejecutando query de base de datos"
	LogDBQuerySuccess         = "Query ejecutado exitosamente"
	LogDBQueryError           = "Error ejecutando query"
	LogDBTransactionStart     = "Iniciando transacción de base de datos"
	LogDBTransactionBegin     = "Comenzando transacción"
	LogDBTransactionBeginErr  = "Error comenzando transacción"
	LogDBTransactionCommit    = "Commit de transacción exitoso"
	LogDBTransactionCommitErr = "Error haciendo commit de transacción"
	LogDBTransactionRollback  = "Rollback de transacción"
	LogDBConnectionPoolInfo   = "Información de connection pool"
	LogDBConnectionError      = "Error conectando a base de datos"
	LogDBPoolConfig           = ""
	LogDBConnecting           = "Conectando a base de datos"
	LogDBSSLEnabled           = "SSL habilitado"
	LogDBPinging              = ""
	LogDBPingError            = "Error haciendo ping a base de datos"
	LogDBConnected            = "Base de datos conectada exitosamente"
)

// ============================================
// MESSAGING / CACHE
// ============================================
const (
	LogMsgCacheInit            = "Inicializando cache de mensajes"
	LogMsgCacheLoaded          = "Mensajes del sistema cargados en cache"
	LogMsgCacheLoadError       = "Error cargando mensajes del sistema desde BD"
	LogMsgCacheRefreshStart    = "Iniciando auto-refresh de cache de mensajes"
	LogMsgCacheRefreshDisabled = "Auto-refresh de cache deshabilitado"
	LogMsgCacheRefreshing      = "Auto-refrescando cache de mensajes desde BD"
	LogMsgCacheRefreshOK       = "Cache de mensajes refrescado exitosamente"
	LogMsgCacheRefreshError    = "Error durante auto-refresh de cache"
	LogMsgCacheRefreshStop     = "Deteniendo auto-refresh de cache de mensajes"
	LogMsgNotInCache           = "Mensaje no encontrado en cache, cargando desde BD"
	LogMsgNotInDB              = "Mensaje no encontrado en base de datos"
	LogMsgCachedFromDB         = "Mensaje cargado desde BD y cacheado"
	LogMsgInactive             = "Mensaje existe pero está desactivado"
)

// ============================================
// ROUTING / MIDDLEWARE
// ============================================
const (
	LogRouteConfiguring      = "Configurando rutas de la aplicación"
	LogRouteConfigured       = "Rutas configuradas correctamente"
	LogRouteValidatorInit    = "Inicializando validador de schemas"
	LogRouteValidatorOK      = "Validador de schema inicializado"
	LogRouteValidatorError   = "Error creando validador de schema"
	LogMiddlewareErrorCaught = "Error de negocio capturado"
	LogMiddlewareInternalErr = "Error interno del servidor"
)

// ============================================
// VALIDATION
// ============================================
const (
	LogValidationStart   = "Iniciando validación de request"
	LogValidationOK      = "Validación de request exitosa"
	LogValidationFailed  = "Validación de request fallida"
	LogValidationDetails = "Detalles de validación"
)

// ============================================
// KEYCLOAK / EXTERNAL SERVICES
// ============================================
const (
	LogKeycloakClientInit                 = "Inicializando cliente Keycloak"
	LogKeycloakClientOK                   = "Cliente Keycloak inicializado correctamente"
	LogKeycloakClientError                = "Error inicializando cliente Keycloak"
	LogKeycloakClientCreated              = "Cliente Keycloak creado exitosamente"
	LogKeycloakConfigNil                  = "Configuración de Keycloak no puede ser nil"
	LogKeycloakAdminAuth                  = "Autenticando admin de Keycloak"
	LogKeycloakAdminAuthError             = "Error autenticando admin de Keycloak"
	LogKeycloakAdminTokenInit             = "Inicializando token de admin"
	LogKeycloakAdminTokenInitError        = "Error inicializando token de admin"
	LogKeycloakTokenRefresh               = "Refrescando token de admin de Keycloak"
	LogKeycloakTokenRefreshOK             = "Token de admin refrescado exitosamente"
	LogKeycloakTokenRefreshErr            = "Error refrescando token de admin de Keycloak"
	LogKeycloakTokenEnsure                = "Asegurando token válido"
	LogKeycloakTokenEnsureError           = "Error asegurando token válido"
	LogKeycloakUserLogin                  = "Intentando login de usuario"
	LogKeycloakUserLoginOK                = "Login de usuario exitoso"
	LogKeycloakUserLoginError             = "Error en login de usuario"
	LogKeycloakUserLoginFailed            = "Login de usuario falló"
	LogKeycloakUserCreate                 = "Creando usuario en Keycloak"
	LogKeycloakUserCreateOK               = "Usuario creado en Keycloak"
	LogKeycloakUserCreateError            = "Error creando usuario en Keycloak"
	LogKeycloakUserNil                    = "Usuario no puede ser nil"
	LogKeycloakUserGet                    = "Obteniendo usuario de Keycloak"
	LogKeycloakUserGetByEmail             = "Obteniendo usuario por email"
	LogKeycloakUserGetByEmailError        = "Error obteniendo usuario por email"
	LogKeycloakUserGetByID                = "Obteniendo usuario por ID"
	LogKeycloakUserGetByIDError           = "Error obteniendo usuario por ID"
	LogKeycloakUserGetError               = "Error obteniendo usuario de Keycloak"
	LogKeycloakUserNotFound               = "Usuario no encontrado en Keycloak"
	LogKeycloakUserDelete                 = "Eliminando usuario de Keycloak"
	LogKeycloakUserDeleteOK               = "Usuario eliminado de Keycloak"
	LogKeycloakUserDeleteError            = "Error eliminando usuario de Keycloak"
	LogKeycloakPasswordSet                = "Configurando password para usuario"
	LogKeycloakPasswordSetOK              = "Password configurado exitosamente"
	LogKeycloakPasswordSetError           = "Error configurando password"
	LogKeycloakRoleGet                    = "Obteniendo rol"
	LogKeycloakRoleGetError               = "Error obteniendo rol"
	LogKeycloakRoleAssign                 = "Asignando rol a usuario"
	LogKeycloakRoleAssignOK               = "Rol asignado exitosamente"
	LogKeycloakRoleAssignError            = "Error asignando rol a usuario"
	LogKeycloakUserTokenRefresh           = "Refrescando token de usuario"
	LogKeycloakUserTokenRefreshOK         = "Token de usuario refrescado exitosamente"
	LogKeycloakUserTokenRefreshErr        = "Error refrescando token de usuario"
	LogKeycloakEmailEmpty                 = "Email no puede estar vacío"
	LogKeycloakUserIDEmpty                = "ID de usuario no puede estar vacío"
	LogKeycloakPasswordEmpty              = "Password no puede estar vacío"
	LogKeycloakRoleNameEmpty              = "Nombre de rol no puede estar vacío"
	LogKeycloakRefreshTokenEmpty          = "Refresh token no puede estar vacío"
	LogKeycloakUsernameEmpty              = "Nombre de usuario no puede estar vacío"
	LogKeycloakSendVerificationEmail      = "Enviando email de verificación"
	LogKeycloakSendVerificationEmailOK    = "Email de verificación enviado exitosamente"
	LogKeycloakSendVerificationEmailError = "Error enviando email de verificación"
	LogKeycloakSendPasswordReset          = "Enviando email de recuperación de contraseña"
	LogKeycloakSendPasswordResetOK        = "Email de recuperación enviado exitosamente"
	LogKeycloakSendPasswordResetError     = "Error enviando email de recuperación"
	LogKeycloakSearchUserByEmail          = "Buscando usuario por email en Keycloak"
	LogKeycloakSearchUserByEmailOK        = "Usuario encontrado por email en Keycloak"
	LogKeycloakEmailVerify                = "Verificando email de usuario"
	LogKeycloakEmailVerifyOK              = "Email verificado exitosamente"
	LogKeycloakEmailVerifyError           = "Error verificando email"
	LogKeycloakEmailAlreadyVerified       = "Email ya ha sido verificado"
	LogKeycloakPasswordUpdate             = "Actualizando contraseña de usuario"
	LogKeycloakPasswordUpdateOK           = "Contraseña actualizada exitosamente"
	LogKeycloakPasswordUpdateError        = "Error actualizando contraseña"
	LogKeycloakPasswordMismatch           = "Las contraseñas no coinciden"
	LogKeycloakPasswordTokenValidation    = "Validando token de actualización de contraseña"
	LogKeycloakPasswordTokenValidOK       = "Token de actualización validado"
	LogKeycloakPasswordTokenInvalid       = "Token de actualización inválido"
	LogKeycloakTokenExpirationValidated   = "Token expiration validated"
	LogKeycloakTokenMissingExpiration     = "Token missing expiration claim, proceeding with caution"
	LogKeycloakTokenIssuerValidated       = "Token issuer validated"
	LogKeycloakTokenMissingIssuer         = "Token missing issuer claim"
	LogKeycloakTokenType                  = "Token type"
	// Login with verification
	LogKeycloakLoginCheckingVerification    = "Verificando estado de email antes de login"
	LogKeycloakLoginEmailNotVerified        = "Login rechazado: email no verificado"
	LogKeycloakLoginEmailVerified           = "Email verificado, procediendo con login"
	LogKeycloakLoginResendingVerification   = "Reenviando email de verificación automáticamente"
	LogKeycloakLoginResendVerificationOK    = "Email de verificación reenviado exitosamente"
	LogKeycloakLoginResendVerificationError = "Error reenviando email de verificación"
	// Change Password (authenticated user)
	LogKeycloakChangePassword           = "Cambiando contraseña de usuario autenticado"
	LogKeycloakChangePasswordValidating = "Validando contraseña actual del usuario"
	LogKeycloakChangePasswordInvalid    = "Contraseña actual incorrecta"
	LogKeycloakChangePasswordOK         = "Contraseña cambiada exitosamente"
	LogKeycloakChangePasswordError      = "Error cambiando contraseña"
	LogKeycloakChangePasswordMismatch   = "Las contraseñas nuevas no coinciden"
)

// ============================================
// KEYCLOAK AVAILABILITY
// ============================================
const (
	LogKeycloakAvailabilityCheck = "Verificando disponibilidad de Keycloak"
	LogKeycloakAvailable         = "Keycloak disponible y respondiendo"
	LogKeycloakUnavailable       = "Keycloak no disponible"
	LogKeycloakConnectionError   = "Error de conexión con Keycloak"
	LogKeycloakTimeoutError      = "Timeout en conexión con Keycloak"
)

// ============================================
// DATABASE AVAILABILITY
// ============================================
const (
	LogDatabaseAvailabilityCheck = "Verificando disponibilidad de base de datos"
	LogDatabaseAvailable         = "Base de datos disponible y respondiendo"
	LogDatabaseUnavailable       = "Base de datos no disponible"
	LogDatabaseConnectionError   = "Error de conexión con base de datos"
)

// ============================================
// DUAL SYSTEM VALIDATION
// ============================================
const (
	LogDualSystemCheck          = "Validando existencia en ambos sistemas"
	LogUserExistsInBoth         = "Usuario existe en ambos sistemas"
	LogUserExistsOnlyInDB       = "Usuario existe solo en base de datos"
	LogUserExistsOnlyInKeycloak = "Usuario existe solo en Keycloak"
	LogUserNotFoundInEither     = "Usuario no encontrado en ningún sistema"
	LogInconsistentStateDetect  = "Estado inconsistente detectado entre sistemas"
)

// ============================================
// REPOSITORY / MESSAGE REPOSITORY
// ============================================
const (
	LogRepoMsgInit           = "Inicializando repositorio de mensajes"
	LogRepoMsgInitOK         = "Repositorio de mensajes inicializado"
	LogRepoMsgInitError      = "Error inicializando repositorio de mensajes"
	LogEmployeeRepoInit      = "Inicializando repositorio de empleados"
	LogEmployeeRepoInitOK    = "Repositorio de empleados inicializado"
	LogEmployeeRepoInitError = "Error creando repositorio de empleados"
	LogIDEncoderInit         = "Inicializando ID encoder"
	LogIDEncoderInitError    = "Error inicializando ID encoder"
)

// ============================================
// HTTP REQUESTS
// ============================================
const (
	LogHTTPRequestIncoming = "Request HTTP entrante"
	LogHTTPResponseSent    = "Respuesta HTTP enviada"
	LogHTTPClientIP        = "IP del cliente"
	LogHTTPMethod          = "Método HTTP"
	LogHTTPPath            = "Path HTTP"
	LogHTTPStatus          = "Status HTTP"
)

// ============================================
// EMPLOYEE / USER SERVICES
// ============================================
const (
	LogEmployeeCreating              = "Creando empleado en base de datos"
	LogEmployeeCreated               = "Empleado creado exitosamente"
	LogEmployeeCreateError           = "Error creando empleado"
	LogEmployeeSaving                = "Guardando empleado en base de datos"
	LogEmployeeSaveError             = "Error guardando empleado en base de datos"
	LogEmployeeUpdating              = "Actualizando empleado"
	LogEmployeeUpdated               = "Empleado actualizado"
	LogEmployeeUpdateError           = "Error actualizando empleado"
	LogEmployeeUpdateKeycloakID      = "Actualizando ID de Keycloak del empleado"
	LogEmployeeUpdateKeycloakIDError = "Error actualizando ID de Keycloak en base de datos"
	LogEmployeeSearching             = "Buscando empleado"
	LogEmployeeGetByEmail            = "Obteniendo empleado por email"
	LogEmployeeGetByEmailError       = "Error obteniendo empleado por email"
	LogEmployeeGetByID               = "Obteniendo empleado por ID"
	LogEmployeeGetByIDError          = "Error obteniendo empleado por ID"
	LogEmployeeGetByIDOK             = "Empleado obtenido exitosamente por ID"
	LogEmployeeFound                 = "Empleado encontrado"
	LogEmployeeNotFound              = "Empleado no encontrado"
	LogEmployeeSearchError           = "Error buscando empleado"
	LogEmployeeRegistering           = "Registrando empleado"
	LogEmployeeRegisterError         = "Error registrando empleado"
	LogEmployeeRegisterSuccess       = "Empleado registrado exitosamente"
	LogEmployeeExists                = "El empleado ya existe"
	LogEmployeeDeleting              = "Eliminando empleado de base de datos"
	LogEmployeeDeleteError           = "Error eliminando empleado de base de datos"
	LogEmployeeLocating              = "Localizando empleado"
	LogEmployeeLocateError           = "Error localizando empleado"
	// Update specific logs
	LogEmployeeUpdateRequest         = "Solicitud de actualización de empleado recibida"
	LogEmployeeUpdateComplete        = "Empleado actualizado exitosamente"
	LogEmployeeKeycloakStatusUpdate  = "Actualizando estado del usuario en Keycloak"
	LogEmployeeKeycloakStatusUpdated = "Estado del usuario actualizado en Keycloak"
	LogEmployeeKeycloakRoleUpdate    = "Actualizando rol del usuario en Keycloak"
	LogEmployeeKeycloakRoleUpdated   = "Rol del usuario actualizado en Keycloak"
	// Delete specific logs
	LogEmployeeDeletingKeycloak    = "Eliminando usuario de Keycloak"
	LogEmployeeDeleteKeycloakError = "Error eliminando usuario de Keycloak"
	LogEmployeeDeletedKeycloak     = "Usuario eliminado de Keycloak exitosamente"
	LogEmployeeDeletingDB          = "Eliminando empleado de base de datos"
	LogEmployeeDeleteDBError       = "Error eliminando empleado de base de datos"
	LogEmployeeDeleteComplete      = "Empleado eliminado exitosamente"
	// Repository layer - Employee Update
	LogRepoEmployeeUpdateStart   = "Iniciando actualización de empleado en repositorio"
	LogRepoEmployeeUpdateTxOK    = "Transacción de repositorio validada correctamente"
	LogRepoEmployeeUpdateQueryOK = "Query de actualización ejecutado exitosamente"
	LogRepoEmployeeUpdateTxErr   = "Error: tipo de transacción inválido"
	// Registration success
	LogEmployeeRegisterSuccessLog = "register employee success"
)

// ============================================
// DEPENDENCY INJECTION
// ============================================
const (
	LogDepInit          = "Inicializando dependencias"
	LogDepInitComplete  = "Dependencias inicializadas exitosamente"
	LogDepInitError     = "Error inicializando dependencias"
	LogDepWiringService = "Inyectando servicio"
	LogDepConfigError   = "Error cargando configuración de dependencias"
	// JWKS Validator
	LogJWKSValidatorInitFailed = "JWKS validator initialization failed, using fallback validation"
	LogJWKSValidatorInitOK     = "JWKS validator initialized"
)

// ============================================
// DTO RESPONSE MESSAGES
// ============================================
const (
	DtoMsgUserRegistered  = "user registered successfully"
	DtoMsgEmployeeUpdated = "employee updated successfully"
	DtoMsgEmployeeLocated = "Employee located successfully"
)

// ============================================
// SECURITY / ENCODING
// ============================================
const (
	LogIDEncode               = "Codificando ID"
	LogIDEncodeOK             = "ID codificado exitosamente"
	LogIDEncodeError          = "Error codificando ID"
	LogIDDecode               = "Decodificando ID"
	LogIDDecodeOK             = "ID decodificado exitosamente"
	LogIDDecodeError          = "Error decodificando ID"
	LogIDEncoderInvalidUUID   = "UUID inválido"
	LogIDEncoderHashidsCreate = "Error creando hashids"
	LogIDEncoderEncodingError = "Error encodeando UUID"
	LogIDEncoderEmptyID       = "ID ofuscado no puede estar vacío"
	LogIDEncoderDecodingError = "Error decodeando ID ofuscado"
	LogIDEncoderInvalidFormat = "ID ofuscado tiene formato incorrecto"
	LogIDEncoderUUIDError     = "Error reconstruyendo UUID"
	LogIDEncoderMinLengthWarn = "MinLength es igual a 36, lo cual es el valor por defecto"
)

// ============================================
// UTILS / HELPERS
// ============================================
const (
	LogUtilsModuleRootSearch = "Buscando raíz del módulo"
	LogUtilsModuleRootFound  = "Raíz del módulo encontrada"
	LogUtilsModuleRootError  = "No se pudo encontrar la raíz del módulo"
	LogUtilsCurrentDirError  = "No se pudo determinar el directorio actual"
	LogUtilsPathResolved     = "Ruta resuelta exitosamente"
	LogUtilsPathError        = "Error resolviendo ruta"
)

// ============================================
// MIDDLEWARE / VALIDATORS
// ============================================
const (
	LogMiddlewareValidationStart    = "Iniciando validación de request body"
	LogMiddlewareValidationOK       = "Validación exitosa"
	LogMiddlewareValidationFailed   = "Validación de request fallida"
	LogMiddlewareBodyReadError      = "Error leyendo body del request"
	LogMiddlewareJSONParseError     = "Error parseando JSON del body"
	LogMiddlewareSchemaError        = "Error de validación de schema"
	LogMiddlewareResponseCacheError = "Error obteniendo mensaje de cache"
	LogMiddlewareResponseSuccess    = "Respuesta enviada exitosamente"
	LogMiddlewareErrorHandling      = "Manejando error de middleware"
	LogMiddlewareNotFound           = "Endpoint no encontrado"
)

// ============================================
// SCHEMA VALIDATION
// ============================================
const (
	LogSchemaValidatorInit      = "Inicializando validador de esquemas"
	LogSchemaValidatorInitOK    = "Validador de esquemas inicializado"
	LogSchemaValidatorInitError = "Error inicializando validador de esquemas"
	LogSchemaReading            = "Leyendo esquema JSON"
	LogSchemaReadError          = "Error leyendo esquema JSON"
	LogSchemaEmpty              = "Esquema JSON vacío o nulo"
	LogSchemaCompiling          = "Compilando esquema JSON"
	LogSchemaCompileError       = "Error compilando esquema JSON"
	LogSchemaValidating         = "Validando datos contra esquema"
	LogSchemaValidationOK       = "Validación de esquema exitosa"
	LogSchemaValidationError    = "Error en validación de esquema"
	LogSchemaFileNotFound       = "Archivo de esquema no encontrado"
	LogSchemaNotInitialized     = "Validador no inicializado"
)

// ============================================
// CONFIG LOADING
// ============================================
const (
	LogConfigFileNotFound    = "Archivo de configuración no encontrado"
	LogConfigFallback        = "Usando configuración por defecto"
	LogConfigReading         = "Leyendo archivo de configuración"
	LogConfigReadError       = "Error leyendo archivo de configuración"
	LogConfigParsing         = "Parseando configuración JSON"
	LogConfigParseError      = "Error parseando configuración JSON"
	LogConfigEnvOverride     = "Sobreescribiendo configuración con variables de entorno"
	LogConfigValidating      = "Validando configuración"
	LogConfigValidationError = "Error validando configuración"
)

// ============================================
// MESSAGES MODULE / SYSTEM MESSAGES
// ============================================
const (
	LogMessageServiceInit       = "Inicializando servicio de mensajes"
	LogMessageServiceInitOK     = "Servicio de mensajes inicializado"
	LogMessageCreate            = "Creando mensaje del sistema"
	LogMessageCreateOK          = "Mensaje creado exitosamente"
	LogMessageCreateError       = "Error creando mensaje"
	LogMessageCreateProcessing  = "Procesando creación de mensaje"
	LogMessageUpdate            = "Actualizando mensaje del sistema"
	LogMessageUpdateOK          = "Mensaje actualizado exitosamente"
	LogMessageUpdateError       = "Error actualizando mensaje"
	LogMessageUpdateProcessing  = "Procesando actualización de mensaje"
	LogMessageDelete            = "Eliminando mensaje del sistema"
	LogMessageDeleteOK          = "Mensaje eliminado exitosamente"
	LogMessageDeleteError       = "Error eliminando mensaje"
	LogMessageDeleteProcessing  = "Procesando eliminación de mensaje"
	LogMessageGet               = "Obteniendo mensaje del sistema"
	LogMessageGetOK             = "Mensaje obtenido exitosamente"
	LogMessageGetError          = "Error obteniendo mensaje"
	LogMessageList              = "Listando mensajes del sistema"
	LogMessageListOK            = "Mensajes listados exitosamente"
	LogMessageListError         = "Error listando mensajes"
	LogMessageValidation        = "Validando mensaje"
	LogMessageValidationOK      = "Mensaje validado exitosamente"
	LogMessageValidationError   = "Error validando mensaje"
	LogMessageCodeDuplicate     = "Código de mensaje duplicado"
	LogMessageTxBegin           = "Iniciando transacción para mensaje"
	LogMessageTxBeginOK         = "Transacción iniciada para mensaje"
	LogMessageTxCommit          = "Confirmando transacción de mensaje"
	LogMessageTxCommitOK        = "Transacción confirmada exitosamente"
	LogMessageTxCommitError     = "Error confirmando transacción"
	LogMessageTxRollback        = "Ejecutando rollback de transacción"
	LogMessageTxRollbackOK      = "Rollback ejecutado exitosamente"
	LogMessageTxRollbackError   = "Error ejecutando rollback"
	LogMessageCacheRefresh      = "Refrescando cache de mensajes"
	LogMessageCacheRefreshOK    = "Cache de mensajes refrescado"
	LogMessageCacheRefreshError = "Error refrescando cache de mensajes"
	LogMessageInvalidID         = "ID inválido"
	LogMessageIDEncodeError     = "Error ofuscando ID"
	LogMessageIDDecodeError     = "Error decodificando ID"
)

// ============================================
// PROMETHEUS / OBSERVABILITY
// ============================================
const (
	LogPrometheusInit          = "Inicializando métricas de Prometheus"
	LogPrometheusInitOK        = "Métricas de Prometheus inicializadas correctamente"
	LogPrometheusInitError     = "Error inicializando métricas de Prometheus"
	LogPrometheusMetricRecord  = "Registrando métrica"
	LogPrometheusMetricError   = "Error registrando métrica"
	LogPrometheusScrapeSuccess = "Scraping de métricas exitoso"
	LogPrometheusScrapeError   = "Error durante scraping de métricas"
)

// ============================================
// EMPLOYEE SERVICES
// ============================================
const (
	LogEmployeeServiceSearchByEmail             = "Buscando persona por email"
	LogEmployeeServiceSearchByID                = "Buscando persona por ID"
	LogEmployeeServiceFoundByEmail              = "Persona encontrada por email"
	LogEmployeeServiceFoundByID                 = "Persona encontrada por ID"
	LogEmployeeServiceErrorByEmail              = "Error buscando persona por email"
	LogEmployeeServiceErrorByID                 = "Error buscando persona por ID"
	LogEmployeeServiceValidationStart           = "Iniciando validaciones de registro"
	LogEmployeeServiceValidationComplete        = "Validaciones de registro completadas"
	LogEmployeeServiceDuplicateEmail            = "Intento de registro con email duplicado"
	LogEmployeeServiceSavingToDB                = "Guardando persona en base de datos"
	LogEmployeeServiceSavedToDB                 = "Persona guardada en base de datos"
	LogEmployeeServiceSaveError                 = "Error guardando persona en BD"
	LogEmployeeServiceCreatingKeycloak          = "Creando usuario en Keycloak"
	LogEmployeeServiceCreatedKeycloak           = "Usuario creado en Keycloak"
	LogEmployeeServiceKeycloakError             = "Error creando usuario en Keycloak"
	LogEmployeeServicePasswordSet               = "Configurando password de usuario"
	LogEmployeeServicePasswordSetOK             = "Password configurado"
	LogEmployeeServicePasswordError             = "Error configurando password"
	LogEmployeeServiceRoleAssigning             = "Asignando rol a usuario"
	LogEmployeeServiceRoleAssigned              = "Rol asignado"
	LogEmployeeServiceRoleError                 = "Error asignando rol"
	LogEmployeeServiceKeycloakIDUpdate          = "Actualizando keycloak_user_id en BD"
	LogEmployeeServiceKeycloakIDUpdated         = "Keycloak_user_id actualizado"
	LogEmployeeServiceKeycloakIDUpdateError     = "Error actualizando keycloak_user_id"
	LogEmployeeServiceRollbackPerson            = "Ejecutando rollback: eliminando persona de BD"
	LogEmployeeServiceRollbackPersonError       = "Error en rollback de persona"
	LogEmployeeServiceRollbackPersonComplete    = "Rollback de persona completado"
	LogEmployeeServiceRollbackKeycloak          = "Ejecutando rollback: eliminando usuario de Keycloak"
	LogEmployeeServiceRollbackKeycloakError     = "Error en rollback de usuario Keycloak"
	LogEmployeeServiceRollbackKeycloakComplete  = "Rollback de usuario Keycloak completado"
	LogEmployeeServiceInconsistentStateDetected = "Estado inconsistente detectado entre Keycloak y BD"
	LogEmployeeServiceCleaningOrphan            = "Limpiando usuario huérfano"
	LogEmployeeServiceOrphanCleaned             = "Usuario huérfano eliminado exitosamente"
	LogEmployeeServiceOrphanCleanError          = "Error limpiando usuario huérfano"
	LogEmployeeServiceRollbackEmployee          = "Ejecutando rollback: eliminando empleado de BD"
	LogEmployeeServiceRollbackEmployeeComplete  = "Rollback de empleado completado"
	LogEmployeeServiceRollbackEmployeeError     = "Error en rollback de empleado"
)

// ============================================
// EMPLOYEE INTERACTOR
// ============================================
const (
	LogEmployeeInteractorRegStart             = "Iniciando proceso de registro"
	LogEmployeeInteractorStep1_Error          = "[PASO 1/8] Validaciones fallidas"
	LogEmployeeInteractorStep1_OK             = "[PASO 1/8] Validaciones completadas"
	LogEmployeeInteractorIDGenerated          = "ID generado para persona"
	LogEmployeeInteractorStep15_Error         = "[PASO 1.5/8] Estado inconsistente detectado y limpiado"
	LogEmployeeInteractorStep15_OK            = "[PASO 1.5/8] Estado consistente verificado"
	LogEmployeeInteractorStep2_Error          = "[PASO 2/8] Error iniciando transacción"
	LogEmployeeInteractorStep2_OK             = "[PASO 2/8] Transacción iniciada"
	LogEmployeeInteractorStep3_Error          = "[PASO 3/8] Error guardando persona"
	LogEmployeeInteractorStep3_OK             = "[PASO 3/8] Persona guardada en BD"
	LogEmployeeInteractorStep4_Error          = "[PASO 4/8] Error creando usuario en Keycloak"
	LogEmployeeInteractorStep4_OK             = "[PASO 4/8] Usuario creado en Keycloak"
	LogEmployeeInteractorStep5_Error          = "[PASO 5/8] Error configurando password"
	LogEmployeeInteractorStep5_OK             = "[PASO 5/8] Password configurado"
	LogEmployeeInteractorStep6_Error          = "[PASO 6/8] Error asignando rol"
	LogEmployeeInteractorStep6_OK             = "[PASO 6/8] Rol asignado"
	LogEmployeeInteractorStep7_Error          = "[PASO 7/8] Error actualizando Keycloak ID en BD"
	LogEmployeeInteractorStep7_OK             = "[PASO 7/8] Keycloak_user_id actualizado en BD"
	LogEmployeeInteractorCommit_Error         = "COMMIT FALLÓ - ALERTA CRÍTICA"
	LogEmployeeInteractorCommit_OK            = "Transacción confirmada exitosamente"
	LogEmployeeInteractorRegComplete          = "Registro completado exitosamente"
	LogEmployeeInteractorRollbackDB_Error     = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA"
	LogEmployeeInteractorRollbackDB_OK        = "Rollback BD ejecutado correctamente"
	LogEmployeeInteractorRollbackKeycloak_Err = "ROLLBACK KEYCLOAK FALLÓ - ALERTA CRÍTICA"
	LogEmployeeInteractorRollbackKeycloak_OK  = "Rollback Keycloak ejecutado correctamente"
	LogEmployeeInteractorIncompleteDetected   = "Registro incompleto detectado"
	LogEmployeeInteractorCleanup_Error        = "Error limpiando estado inconsistente"
	LogEmployeeInteractorCleanup_OK           = "Estado inconsistente limpiado exitosamente"
)

// DEPENDENCY INITIALIZATION
const (
	LogDependencyMessageRepoInit = "Dependencia de repositorio de mensajes inicializada exitosamente"
	LogDependencyMessageIntInit  = "Error inicializando dependencia de repositorio de mensajes"
)

// ============================================
// MESSAGE INTERACTOR
// ============================================
const (
	// CREATE flow
	LogMessageInteractorCreateStep1Error = "[PASO 1/3] Validación de mensaje fallida"
	LogMessageInteractorCreateStep1OK    = "[PASO 1/3] Validación de mensaje completada"
	LogMessageInteractorCreateStep2Error = "[PASO 2/3] Error iniciando transacción"
	LogMessageInteractorCreateStep2OK    = "[PASO 2/3] Transacción iniciada"
	LogMessageInteractorCreateStep3Error = "[PASO 3/3] Error guardando mensaje"
	LogMessageInteractorCreateStep3OK    = "[PASO 3/3] Mensaje guardado en BD"
	LogMessageInteractorCreateCommitErr  = "COMMIT FALLÓ - ALERTA CRÍTICA"
	LogMessageInteractorCreateCommitOK   = "Transacción confirmada exitosamente"
	LogMessageInteractorCreateComplete   = "Mensaje creado exitosamente"

	// UPDATE flow
	LogMessageInteractorUpdateStep1Error = "[PASO 1/4] Mensaje no encontrado"
	LogMessageInteractorUpdateStep1OK    = "[PASO 1/4] Mensaje encontrado"
	LogMessageInteractorUpdateStep2Error = "[PASO 2/4] Validación de mensaje fallida"
	LogMessageInteractorUpdateStep2OK    = "[PASO 2/4] Validación de mensaje completada"
	LogMessageInteractorUpdateStep3Error = "[PASO 3/4] Error iniciando transacción"
	LogMessageInteractorUpdateStep3OK    = "[PASO 3/4] Transacción iniciada"
	LogMessageInteractorUpdateStep4Error = "[PASO 4/4] Error actualizando mensaje"
	LogMessageInteractorUpdateStep4OK    = "[PASO 4/4] Mensaje actualizado en BD"
	LogMessageInteractorUpdateCommitErr  = "COMMIT FALLÓ - ALERTA CRÍTICA"
	LogMessageInteractorUpdateCommitOK   = "Transacción confirmada exitosamente"
	LogMessageInteractorUpdateComplete   = "Mensaje actualizado exitosamente"

	// DELETE flow
	LogMessageInteractorDeleteStep1Error = "[PASO 1/3] Mensaje no encontrado"
	LogMessageInteractorDeleteStep1OK    = "[PASO 1/3] Mensaje encontrado"
	LogMessageInteractorDeleteStep2Error = "[PASO 2/3] Error iniciando transacción"
	LogMessageInteractorDeleteStep2OK    = "[PASO 2/3] Transacción iniciada"
	LogMessageInteractorDeleteStep3Error = "[PASO 3/3] Error eliminando mensaje"
	LogMessageInteractorDeleteStep3OK    = "[PASO 3/3] Mensaje eliminado de BD"
	LogMessageInteractorDeleteCommitErr  = "COMMIT FALLÓ - ALERTA CRÍTICA"
	LogMessageInteractorDeleteCommitOK   = "Transacción confirmada exitosamente"
	LogMessageInteractorDeleteComplete   = "Mensaje eliminado exitosamente"

	// Common rollback
	LogMessageInteractorRollbackError = "ROLLBACK BD FALLÓ - ALERTA CRÍTICA"
	LogMessageInteractorRollbackOK    = "Rollback BD ejecutado correctamente"

	// Controller success logs
	LogMessageCreatedSuccess = "Mensaje creado exitosamente"
	LogMessageUpdatedSuccess = "Mensaje actualizado exitosamente"
	LogMessageDeletedSuccess = "Mensaje eliminado exitosamente"

	// Cache reload logs
	LogMessageCacheReloading     = "Recargando caché de mensajes"
	LogMessageCacheReloadError   = "Error al recargar caché de mensajes"
	LogMessageCacheReloadSuccess = "Caché de mensajes recargado exitosamente"
	LogMessageCacheReloadedMsg   = "Caché de mensajes recargado exitosamente desde la base de datos"
)

// ============================================
// AIRLINE INTERACTOR
// ============================================
const (
	LogAirlineGet             = "Obteniendo información de aerolínea"
	LogAirlineGetOK           = "Aerolínea obtenida exitosamente"
	LogAirlineGetError        = "Error obteniendo aerolínea"
	LogAirlineNotFound        = "Aerolínea no encontrada"
	LogAirlineActivate        = "Activando aerolínea"
	LogAirlineActivateOK      = "Aerolínea activada exitosamente"
	LogAirlineActivateError   = "Error activando aerolínea"
	LogAirlineDeactivate      = "Desactivando aerolínea"
	LogAirlineDeactivateOK    = "Aerolínea desactivada exitosamente"
	LogAirlineDeactivateError = "Error desactivando aerolínea"
	LogAirlineList            = "Listando aerolíneas"
	LogAirlineListOK          = "Aerolíneas listadas exitosamente"
	LogAirlineListError       = "Error listando aerolíneas"
	LogAirlineRepoInit        = "Inicializando repositorio de aerolíneas"
	LogAirlineRepoInitOK      = "Repositorio de aerolíneas inicializado"
	LogAirlineRepoInitError   = "Error inicializando repositorio de aerolíneas"
)

// ============================================
// AIRPORT INTERACTOR
// ============================================
const (
	LogAirportGet             = "Obteniendo información de aeropuerto"
	LogAirportGetOK           = "Aeropuerto obtenido exitosamente"
	LogAirportGetError        = "Error obteniendo aeropuerto"
	LogAirportNotFound        = "Aeropuerto no encontrado"
	LogAirportActivate        = "Activando aeropuerto"
	LogAirportActivateOK      = "Aeropuerto activado exitosamente"
	LogAirportActivateError   = "Error activando aeropuerto"
	LogAirportDeactivate      = "Desactivando aeropuerto"
	LogAirportDeactivateOK    = "Aeropuerto desactivado exitosamente"
	LogAirportDeactivateError = "Error desactivando aeropuerto"
	LogAirportRepoInit        = "Inicializando repositorio de aeropuertos"
	LogAirportRepoInitOK      = "Repositorio de aeropuertos inicializado"
	LogAirportRepoInitError   = "Error inicializando repositorio de aeropuertos"
	LogAirportList            = "Listando aeropuertos"
	LogAirportListOK          = "Aeropuertos listados exitosamente"
	LogAirportListError       = "Error listando aeropuertos"
)

// ============================================
// DAILY LOGBOOK INTERACTOR
// ============================================
const (
	LogDailyLogbookGet             = "Obteniendo información de bitácora diaria"
	LogDailyLogbookGetOK           = "Bitácora diaria obtenida exitosamente"
	LogDailyLogbookGetError        = "Error obteniendo bitácora diaria"
	LogDailyLogbookNotFound        = "Bitácora diaria no encontrada"
	LogDailyLogbookCreate          = "Creando bitácora diaria"
	LogDailyLogbookCreateOK        = "Bitácora diaria creada exitosamente"
	LogDailyLogbookCreateError     = "Error creando bitácora diaria"
	LogDailyLogbookDelete          = "Eliminando bitácora diaria"
	LogDailyLogbookDeleteOK        = "Bitácora diaria eliminada exitosamente"
	LogDailyLogbookDeleteError     = "Error eliminando bitácora diaria"
	LogDailyLogbookActivate        = "Activando bitácora diaria"
	LogDailyLogbookActivateOK      = "Bitácora diaria activada exitosamente"
	LogDailyLogbookActivateError   = "Error activando bitácora diaria"
	LogDailyLogbookDeactivate      = "Desactivando bitácora diaria"
	LogDailyLogbookDeactivateOK    = "Bitácora diaria desactivada exitosamente"
	LogDailyLogbookDeactivateError = "Error desactivando bitácora diaria"
	LogDailyLogbookList            = "Listando bitácoras diarias"
	LogDailyLogbookListOK          = "Bitácoras diarias listadas exitosamente"
	LogDailyLogbookListError       = "Error listando bitácoras diarias"
	LogDailyLogbookRepoInit        = "Inicializando repositorio de bitácoras diarias"
	LogDailyLogbookRepoInitOK      = "Repositorio de bitácoras diarias inicializado"
	LogDailyLogbookRepoInitError   = "Error inicializando repositorio de bitácoras diarias"
)

// ============================================
// AIRCRAFT REGISTRATION INTERACTOR
// ============================================
const (
	LogLicensePlateGet                  = "Obteniendo información de matrícula"
	LogLicensePlateGetOK                = "Matrícula obtenida exitosamente"
	LogLicensePlateGetError             = "Error obteniendo matrícula"
	LogLicensePlateNotFound             = "Matrícula no encontrada"
	LogLicensePlateCreate               = "Creando matrícula"
	LogLicensePlateCreateOK             = "Matrícula creada exitosamente"
	LogLicensePlateCreateError          = "Error creando matrícula"
	LogLicensePlateUpdate               = "Actualizando matrícula"
	LogLicensePlateUpdateOK             = "Matrícula actualizada exitosamente"
	LogLicensePlateUpdateError          = "Error actualizando matrícula"
	LogLicensePlateList                 = "Listando matrículas"
	LogLicensePlateListOK               = "Matrículas listadas exitosamente"
	LogLicensePlateListError            = "Error listando matrículas"
	LogLicensePlateRepoInit             = "Inicializando repositorio de matrículas"
	LogLicensePlateRepoInitOK           = "Repositorio de matrículas inicializado"
	LogLicensePlateRepoInitError        = "Error inicializando repositorio de matrículas"
	LogLicensePlateInvalidModelID       = "Formato inválido de aircraft_model_id"
	LogLicensePlateInvalidAirlineID     = "Formato inválido de airline_id"
	LogLicensePlateFetchAfterWriteError = "No se pudo obtener matrícula después de guardar"
)

// ============================================
// AIRCRAFT MODEL INTERACTOR
// ============================================
const (
	LogAircraftModelGet             = "Obteniendo información de modelo de aeronave"
	LogAircraftModelGetOK           = "Modelo de aeronave obtenido exitosamente"
	LogAircraftModelGetError        = "Error obteniendo modelo de aeronave"
	LogAircraftModelNotFound        = "Modelo de aeronave no encontrado"
	LogAircraftModelList            = "Listando modelos de aeronave"
	LogAircraftModelListOK          = "Modelos de aeronave listados exitosamente"
	LogAircraftModelListError       = "Error listando modelos de aeronave"
	LogAircraftModelRepoInit        = "Inicializando repositorio de modelos de aeronave"
	LogAircraftModelRepoInitOK      = "Repositorio de modelos de aeronave inicializado"
	LogAircraftModelRepoInitError   = "Error inicializando repositorio de modelos de aeronave"
	LogAircraftModelActivate        = "Activando modelo de aeronave"
	LogAircraftModelActivateOK      = "Modelo de aeronave activado exitosamente"
	LogAircraftModelActivateError   = "Error activando modelo de aeronave"
	LogAircraftModelDeactivate      = "Desactivando modelo de aeronave"
	LogAircraftModelDeactivateOK    = "Modelo de aeronave desactivado exitosamente"
	LogAircraftModelDeactivateError = "Error desactivando modelo de aeronave"
)

// ============================================
// AIRCRAFT FAMILY
// ============================================
const (
	LogAircraftFamilyGet      = "Obteniendo información de familia de aeronave"
	LogAircraftFamilyGetOK    = "Familia de aeronave obtenida exitosamente"
	LogAircraftFamilyGetError = "Error obteniendo familia de aeronave"
	LogAircraftFamilyNotFound = "Familia de aeronave no encontrada"
)

// ============================================
// ROUTE INTERACTOR
// ============================================
const (
	LogRouteGet           = "Obteniendo información de ruta"
	LogRouteGetOK         = "Ruta obtenida exitosamente"
	LogRouteGetError      = "Error obteniendo ruta"
	LogRouteNotFound      = "Ruta no encontrada"
	LogRouteList          = "Listando rutas"
	LogRouteListOK        = "Rutas listadas exitosamente"
	LogRouteListError     = "Error listando rutas"
	LogRouteRepoInit      = "Inicializando repositorio de rutas"
	LogRouteRepoInitOK    = "Repositorio de rutas inicializado"
	LogRouteRepoInitError = "Error inicializando repositorio de rutas"
)

// ============================================
// AIRLINE ROUTE INTERACTOR
// ============================================
const (
	LogAirlineRouteGet             = "Obteniendo información de ruta aerolínea"
	LogAirlineRouteGetOK           = "Ruta aerolínea obtenida exitosamente"
	LogAirlineRouteGetError        = "Error obteniendo ruta aerolínea"
	LogAirlineRouteNotFound        = "Ruta aerolínea no encontrada"
	LogAirlineRouteList            = "Listando rutas aerolínea"
	LogAirlineRouteListOK          = "Rutas aerolínea listadas exitosamente"
	LogAirlineRouteListError       = "Error listando rutas aerolínea"
	LogAirlineRouteActivate        = "Activando ruta aerolínea"
	LogAirlineRouteActivateOK      = "Ruta aerolínea activada exitosamente"
	LogAirlineRouteActivateError   = "Error activando ruta aerolínea"
	LogAirlineRouteDeactivate      = "Desactivando ruta aerolínea"
	LogAirlineRouteDeactivateOK    = "Ruta aerolínea desactivada exitosamente"
	LogAirlineRouteDeactivateError = "Error desactivando ruta aerolínea"
	LogAirlineRouteRepoInit        = "Inicializando repositorio de rutas aerolínea"
	LogAirlineRouteRepoInitOK      = "Repositorio de rutas aerolínea inicializado"
	LogAirlineRouteRepoInitError   = "Error inicializando repositorio de rutas aerolínea"
)

// ============================================
// FLIGHT INTERACTOR
// ============================================
const (
	LogFlightGet           = "Obteniendo información de vuelo"
	LogFlightGetOK         = "Vuelo obtenido exitosamente"
	LogFlightGetError      = "Error obteniendo vuelo"
	LogFlightNotFound      = "Vuelo no encontrado"
	LogFlightList          = "Listando vuelos"
	LogFlightListOK        = "Vuelos listados exitosamente"
	LogFlightListError     = "Error listando vuelos"
	LogFlightCreate        = "Creando vuelo"
	LogFlightCreateOK      = "Vuelo creado exitosamente"
	LogFlightCreateError   = "Error creando vuelo"
	LogFlightUpdate        = "Actualizando vuelo"
	LogFlightUpdateOK      = "Vuelo actualizado exitosamente"
	LogFlightUpdateError   = "Error actualizando vuelo"
	LogFlightDelete        = "Eliminando vuelo"
	LogFlightDeleteOK      = "Vuelo eliminado exitosamente"
	LogFlightDeleteError   = "Error eliminando vuelo"
	LogFlightRepoInit      = "Inicializando repositorio de vuelos"
	LogFlightRepoInitOK    = "Repositorio de vuelos inicializado"
	LogFlightRepoInitError = "Error inicializando repositorio de vuelos"
)

// ============================================
// DAILY LOGBOOK DETAIL INTERACTOR
// ============================================
const (
	LogDailyLogbookDetailGet           = "Obteniendo información de detalle de bitácora"
	LogDailyLogbookDetailGetOK         = "Detalle de bitácora obtenido exitosamente"
	LogDailyLogbookDetailGetError      = "Error obteniendo detalle de bitácora"
	LogDailyLogbookDetailNotFound      = "Detalle de bitácora no encontrado"
	LogDailyLogbookDetailList          = "Listando detalles de bitácora"
	LogDailyLogbookDetailListOK        = "Detalles de bitácora listados exitosamente"
	LogDailyLogbookDetailListError     = "Error listando detalles de bitácora"
	LogDailyLogbookDetailCreate        = "Creando detalle de bitácora"
	LogDailyLogbookDetailCreateOK      = "Detalle de bitácora creado exitosamente"
	LogDailyLogbookDetailCreateError   = "Error creando detalle de bitácora"
	LogDailyLogbookDetailUpdate        = "Actualizando detalle de bitácora"
	LogDailyLogbookDetailUpdateOK      = "Detalle de bitácora actualizado exitosamente"
	LogDailyLogbookDetailUpdateError   = "Error actualizando detalle de bitácora"
	LogDailyLogbookDetailDelete        = "Eliminando detalle de bitácora"
	LogDailyLogbookDetailDeleteOK      = "Detalle de bitácora eliminado exitosamente"
	LogDailyLogbookDetailDeleteError   = "Error eliminando detalle de bitácora"
	LogDailyLogbookDetailRepoInit      = "Inicializando repositorio de detalles de bitácora"
	LogDailyLogbookDetailRepoInitOK    = "Repositorio de detalles de bitácora inicializado"
	LogDailyLogbookDetailRepoInitError = "Error inicializando repositorio de detalles de bitácora"
	LogDailyLogbookDetailDuplicate     = "Detalle de bitácora duplicado"
)

// ============================================
// ENGINE INTERACTOR
// ============================================
const (
	LogEngineGet           = "Obteniendo información de motor"
	LogEngineGetOK         = "Motor obtenido exitosamente"
	LogEngineGetError      = "Error obteniendo motor"
	LogEngineNotFound      = "Motor no encontrado"
	LogEngineList          = "Listando motores"
	LogEngineListOK        = "Motores listados exitosamente"
	LogEngineListError     = "Error listando motores"
	LogEngineRepoInit      = "Inicializando repositorio de motores"
	LogEngineRepoInitOK    = "Repositorio de motores inicializado"
	LogEngineRepoInitError = "Error inicializando repositorio de motores"
)

// ============================================
// MANUFACTURER INTERACTOR
// ============================================
const (
	LogManufacturerGet           = "Obteniendo información de fabricante"
	LogManufacturerGetOK         = "Fabricante obtenido exitosamente"
	LogManufacturerGetError      = "Error obteniendo fabricante"
	LogManufacturerNotFound      = "Fabricante no encontrado"
	LogManufacturerList          = "Listando fabricantes"
	LogManufacturerListOK        = "Fabricantes listados exitosamente"
	LogManufacturerListError     = "Error listando fabricantes"
	LogManufacturerRepoInit      = "Inicializando repositorio de fabricantes"
	LogManufacturerRepoInitOK    = "Repositorio de fabricantes inicializado"
	LogManufacturerRepoInitError = "Error inicializando repositorio de fabricantes"
)

// ============================================
// CITY INTERACTOR (- Virtual Entity pattern)
// ============================================
const (
	LogCityGet      = "Obteniendo aeropuertos por ciudad"
	LogCityGetOK    = "Aeropuertos por ciudad obtenidos exitosamente"
	LogCityGetError = "Error obteniendo aeropuertos por ciudad"
	LogCityNotFound = "Ciudad no encontrada (sin aeropuertos)"
)

// ============================================
// COUNTRY INTERACTOR (- Virtual Entity pattern)
// ============================================
const (
	LogCountryGet      = "Obteniendo aeropuertos por país"
	LogCountryGetOK    = "Aeropuertos por país obtenidos exitosamente"
	LogCountryGetError = "Error obteniendo aeropuertos por país"
	LogCountryNotFound = "País no encontrado (sin aeropuertos)"
)

// ============================================
// AIRPORT TYPE INTERACTOR (HU46 - Virtual Entity pattern)
// ============================================
const (
	LogAirportTypeGet      = "Obteniendo aeropuertos por tipo"
	LogAirportTypeGetOK    = "Aeropuertos por tipo obtenidos exitosamente"
	LogAirportTypeGetError = "Error obteniendo aeropuertos por tipo"
	LogAirportTypeNotFound = "Tipo de aeropuerto no encontrado (sin aeropuertos)"
)

// ============================================
// CREW MEMBER TYPE INTERACTOR (HU47 - Virtual Entity pattern)
// ============================================
const (
	LogCrewMemberTypeGet      = "Obteniendo empleados por tipo de rol"
	LogCrewMemberTypeGetOK    = "Empleados por tipo de rol obtenidos exitosamente"
	LogCrewMemberTypeGetError = "Error obteniendo empleados por tipo de rol"
	LogCrewMemberTypeNotFound = "Tipo de integrante no encontrado (sin empleados)"
)

// ============================================
// AIRLINE EMPLOYEE INTERACTOR (Release 15)
// ============================================
const (
	LogAirlineEmployeeGet             = "Obteniendo información de empleado aerolínea"
	LogAirlineEmployeeGetOK           = "Empleado aerolínea obtenido exitosamente"
	LogAirlineEmployeeGetError        = "Error obteniendo empleado aerolínea"
	LogAirlineEmployeeNotFound        = "Empleado aerolínea no encontrado"
	LogAirlineEmployeeList            = "Listando empleados aerolínea"
	LogAirlineEmployeeListOK          = "Empleados aerolínea listados exitosamente"
	LogAirlineEmployeeListError       = "Error listando empleados aerolínea"
	LogAirlineEmployeeCreate          = "Creando empleado aerolínea"
	LogAirlineEmployeeCreateOK        = "Empleado aerolínea creado exitosamente"
	LogAirlineEmployeeCreateError     = "Error creando empleado aerolínea"
	LogAirlineEmployeeUpdate          = "Actualizando empleado aerolínea"
	LogAirlineEmployeeUpdateOK        = "Empleado aerolínea actualizado exitosamente"
	LogAirlineEmployeeUpdateError     = "Error actualizando empleado aerolínea"
	LogAirlineEmployeeAdd             = "Agregando información de aerolínea a empleado"
	LogAirlineEmployeeAddOK           = "Información de aerolínea agregada exitosamente"
	LogAirlineEmployeeAddError        = "Error agregando información de aerolínea"
	LogAirlineEmployeeActivate        = "Activando empleado aerolínea"
	LogAirlineEmployeeActivateOK      = "Empleado aerolínea activado exitosamente"
	LogAirlineEmployeeActivateError   = "Error activando empleado aerolínea"
	LogAirlineEmployeeDeactivate      = "Desactivando empleado aerolínea"
	LogAirlineEmployeeDeactivateOK    = "Empleado aerolínea desactivado exitosamente"
	LogAirlineEmployeeDeactivateError = "Error desactivando empleado aerolínea"
	LogAirlineEmployeeRepoInit        = "Inicializando repositorio de empleados aerolínea"
	LogAirlineEmployeeRepoInitOK      = "Repositorio de empleados aerolínea inicializado"
	LogAirlineEmployeeRepoInitError   = "Error inicializando repositorio de empleados aerolínea"
)

// ============================================
// ERROR DETAIL MESSAGES (for log attributes)
// ============================================
const (
	LogErrAuthUserNotInContext  = "authenticated user not in context"
	LogErrInvalidAirlineID      = "invalid airline id"
	LogErrAirlineNotFound       = "airline not found"
	LogErrInvalidStartDate      = "invalid start date"
	LogErrInvalidEndDate        = "invalid end date"
	LogErrInvalidTransaction    = "invalid transaction type"
	LogErrEmptyIDParam          = "empty id parameter"
	LogErrEmployeeNoAirlineInfo = "employee has no airline info assigned"
)
