<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FlightHours - Recuperar Contraseña</title>
</head>
<body style="font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #F5F7FA; margin: 0; padding: 2rem;">
    <div style="max-width: 600px; margin: 0 auto;">
        <div style="background: white; padding: 2.5rem; border-radius: 16px; box-shadow: 0 4px 16px rgba(0, 71, 171, 0.12);">

            <!-- Header -->
            <div style="text-align: center; margin-bottom: 2rem; padding-bottom: 1.5rem; border-bottom: 2px solid #0047AB;">
                <div style="color: #0047AB; font-size: 28px; font-weight: 700;">✈️ FlightHours</div>
            </div>

            <!-- Content -->
            <div style="color: #2C3E50; line-height: 1.6;">
                <h1 style="color: #2C3E50; text-align: center; margin-bottom: 1.5rem; font-size: 24px; font-weight: 600;">Recuperación de Contraseña</h1>

                <p style="margin: 1rem 0; color: #64748B; font-size: 15px;">Hola<#if user.firstName??> ${user.firstName}</#if>,</p>

                <p style="margin: 1rem 0; color: #64748B; font-size: 15px;">Recibimos una solicitud para restablecer la contraseña de tu cuenta en FlightHours.</p>

                <p style="margin: 1rem 0; color: #64748B; font-size: 15px;">Si fuiste tú quien solicitó esto, haz clic en el siguiente botón para crear una nueva contraseña:</p>

                <#--
                    Usamos el link de Keycloak pero con un redirect a nuestra página personalizada.
                    El link original ${link} va a Keycloak action-token.
                    Lo modificamos para redirigir a nuestra página de update-password que tiene JavaScript.
                -->
                <#assign customLink = "http://localhost:8080/realms/flighthours/protocol/openid-connect/auth?client_id=account&redirect_uri=" + link?url + "&response_type=code&scope=openid&kc_action=UPDATE_PASSWORD">

                <#-- Para simplicidad, seguimos usando el link original pero nuestra página lo interceptará -->
                <div style="text-align: center; margin: 1.5rem 0;">
                    <a href="${link}" style="display: inline-block; padding: 14px 28px; background-color: #0047AB; color: #FFFFFF; text-decoration: none; border-radius: 10px; font-weight: 600; font-size: 16px;">
                        Restablecer contraseña
                    </a>
                </div>

                <div style="background-color: #FFFBEB; border-left: 4px solid #F59E0B; padding: 1.25rem; margin: 1.5rem 0; border-radius: 8px;">
                    <p style="margin: 0; color: #92400E; font-size: 14px;"><strong>IMPORTANTE:</strong> Este enlace expira en 12 horas.</p>
                </div>

                <p style="color: #64748B; font-size: 14px; word-break: break-all;">Si no puedes hacer clic en el botón, copia y pega este enlace en tu navegador:</p>
                <p style="color: #64748B; font-size: 14px; word-break: break-all;">${link}</p>

                <div style="background-color: #EFF6FF; border-left: 4px solid #0047AB; padding: 1.25rem; margin: 1.5rem 0; border-radius: 8px;">
                    <p style="margin: 0; color: #1E40AF; font-size: 14px;"><strong>¿No solicitaste esto?</strong><br>
                    No te preocupes, tu contraseña actual permanece totalmente segura. Puedes ignorar este mensaje.</p>
                </div>

                <p style="color: #64748B; font-size: 13px; margin-top: 30px;">
                    <strong>Consejo de seguridad:</strong><br>
                    Nunca compartas tu contraseña con nadie. FlightHours nunca te pedirá tu contraseña por correo electrónico.
                </p>
            </div>

            <!-- Footer -->
            <div style="margin-top: 2rem; padding-top: 1.5rem; border-top: 1px solid #E1E8ED; text-align: center; font-size: 13px; color: #64748B;">
                <p><a href="mailto:soporte@flighthours.com" style="color: #0047AB; text-decoration: none; font-weight: 500;">soporte@flighthours.com</a></p>
                <p>&copy; ${.now?string('yyyy')} FlightHours. Todos los derechos reservados.</p>
            </div>

        </div>
    </div>
</body>
</html>
