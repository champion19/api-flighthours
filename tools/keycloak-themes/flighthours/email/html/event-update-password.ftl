<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FlightHours - Contraseña Actualizada</title>
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
                <h1 style="color: #2C3E50; text-align: center; margin-bottom: 1.5rem; font-size: 24px; font-weight: 600;">Contraseña Actualizada</h1>

                <p style="margin: 1rem 0; color: #64748B; font-size: 15px;">Hola<#if user.firstName??> ${user.firstName}</#if>,</p>

                <p style="margin: 1rem 0; color: #64748B; font-size: 15px;">Te confirmamos que tu contraseña de FlightHours ha sido actualizada exitosamente.</p>

                <div style="background-color: #EFF6FF; border-left: 4px solid #0047AB; padding: 1.25rem; margin: 1.5rem 0; border-radius: 8px;">
                    <p style="margin: 0; color: #1E40AF; font-size: 14px;"><strong>Fecha:</strong> ${.now?string('dd/MM/yyyy, HH:mm')}</p>
                    <p style="margin: 0.5rem 0 0 0; color: #1E40AF; font-size: 14px;"><strong>Cambio:</strong> Contraseña actualizada</p>
                </div>

                <p style="margin: 1rem 0; color: #64748B; font-size: 15px;">A partir de ahora, podrás iniciar sesión con tu nueva contraseña.</p>

                <div style="background-color: #FEF2F2; border-left: 4px solid #EF4444; padding: 1.25rem; margin: 1.5rem 0; border-radius: 8px;">
                    <p style="margin: 0; color: #991B1B; font-size: 14px;"><strong>¿NO FUISTE TÚ?</strong><br>
                    Si NO realizaste este cambio, tu cuenta podría estar comprometida.<br>
                    Contacta a soporte INMEDIATAMENTE: soporte@flighthours.com</p>
                </div>

                <p style="color: #64748B; font-size: 13px; margin-top: 30px;">
                    <strong>Recomendaciones de seguridad:</strong><br>
                    - Usa una contraseña única y fuerte<br>
                    - Nunca compartas tus credenciales<br>
                    - Habilita la autenticación de dos factores si está disponible
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
