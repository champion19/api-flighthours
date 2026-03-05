<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Flighthours - Acción Requerida</title>
</head>
<body style="font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #F5F7FA; margin: 0; padding: 2rem;">
    <div style="max-width: 600px; margin: 0 auto;">
        <div style="background: white; padding: 2.5rem; border-radius: 16px; box-shadow: 0 4px 16px rgba(0, 71, 171, 0.12);">

            <!-- Header -->
            <div style="text-align: center; margin-bottom: 2rem; padding-bottom: 1.5rem; border-bottom: 2px solid #0047AB;">
                <div style="color: #0047AB; font-size: 28px; font-weight: 700;">✈️ Flighthours</div>
            </div>

            <!-- Content -->
            <div style="color: #2C3E50; line-height: 1.6;">
                <h1 style="color: #2C3E50; text-align: center; margin-bottom: 1.5rem; font-size: 24px; font-weight: 600;">Acción requerida en tu cuenta</h1>

                <p style="margin: 1rem 0; color: #64748B; font-size: 15px;">Hola<#if user.firstName??> ${user.firstName}</#if>,</p>

                <p style="margin: 1rem 0; color: #64748B; font-size: 15px;">Necesitamos que completes algunas acciones en tu cuenta de Flighthours para continuar.</p>

                <#if requiredActions??>
                <div style="background-color: #EFF6FF; border-left: 4px solid #0047AB; padding: 1.25rem; margin: 1.5rem 0; border-radius: 8px;">
                    <p style="margin: 0; color: #1E40AF; font-size: 14px;"><strong>Acciones pendientes:</strong></p>
                    <ul style="margin: 8px 0; padding-left: 20px; color: #1E40AF;">
                    <#list requiredActions as action>
                        <li>
                            <#if action == "VERIFY_EMAIL">
                                Verificar correo electrónico
                            <#elseif action == "UPDATE_PASSWORD">
                                Actualizar contraseña
                            <#elseif action == "CONFIGURE_TOTP">
                                Configurar autenticación de dos factores
                            <#elseif action == "UPDATE_PROFILE">
                                Actualizar perfil
                            <#elseif action == "TERMS_AND_CONDITIONS">
                                Aceptar términos y condiciones
                            <#else>
                                ${action}
                            </#if>
                        </li>
                    </#list>
                    </ul>
                </div>
                </#if>

                <#--
                    El enlace apunta a Keycloak, pero la página de Keycloak
                    (login-verify-email.ftl) enviará el token a nuestro backend
                -->
                <div style="text-align: center; margin: 1.5rem 0;">
                    <a href="${link}" style="display: inline-block; padding: 14px 28px; background-color: #0047AB; color: #FFFFFF; text-decoration: none; border-radius: 10px; font-weight: 600; font-size: 16px;">
                        Completar acciones
                    </a>
                </div>

                <div style="background-color: #FFFBEB; border-left: 4px solid #F59E0B; padding: 1.25rem; margin: 1.5rem 0; border-radius: 8px;">
                    <p style="margin: 0; color: #92400E; font-size: 14px;"><strong>IMPORTANTE:</strong> Este enlace expira en 15 minutos.</p>
                </div>

                <p style="color: #64748B; font-size: 14px; word-break: break-all;">Si el botón no funciona, copia y pega este enlace en tu navegador:</p>
                <p style="color: #64748B; font-size: 14px; word-break: break-all;">${link}</p>

                <div style="background-color: #EFF6FF; border-left: 4px solid #0047AB; padding: 1.25rem; margin: 1.5rem 0; border-radius: 8px;">
                    <p style="margin: 0; color: #1E40AF; font-size: 14px;"><strong>Nota:</strong> Si no reconoces esta solicitud, contacta a nuestro equipo de soporte de inmediato.</p>
                </div>
            </div>

            <!-- Footer -->
            <div style="margin-top: 2rem; padding-top: 1.5rem; border-top: 1px solid #E1E8ED; text-align: center; font-size: 13px; color: #64748B;">
                <p><a href="mailto:soporte@flighthours.com" style="color: #0047AB; text-decoration: none; font-weight: 500;">soporte@flighthours.com</a></p>
                <p>&copy; ${.now?string('yyyy')} Flighthours. Todos los derechos reservados.</p>
            </div>

        </div>
    </div>
</body>
</html>
