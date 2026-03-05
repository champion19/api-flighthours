Acción requerida - Flighthours

Hola<#if user.firstName??> ${user.firstName}</#if>,

Necesitamos que completes algunas acciones en tu cuenta de Flighthours para continuar.

<#if requiredActions??>
Acciones pendientes:
<#list requiredActions as action>
- <#if action == "VERIFY_EMAIL">Verificar correo electrónico<#elseif action == "UPDATE_PASSWORD">Actualizar contraseña<#elseif action == "CONFIGURE_TOTP">Configurar autenticación de dos factores<#elseif action == "UPDATE_PROFILE">Actualizar perfil<#elseif action == "TERMS_AND_CONDITIONS">Aceptar términos y condiciones<#else>${action}</#if>
</#list>
</#if>

Enlace para completar acciones:
${link}

IMPORTANTE: Este enlace expira en 24 horas.

--
soporte@flighthours.com
© ${.now?string('yyyy')} Flighthours
