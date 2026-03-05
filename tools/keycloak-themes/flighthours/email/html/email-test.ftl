<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FlightHours - Test SMTP</title>
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
                <h1 style="color: #2C3E50; text-align: center; margin-bottom: 1.5rem; font-size: 24px; font-weight: 600;">Test de Configuración SMTP</h1>

                <p style="text-align: center; font-size: 16px; color: #64748B;">
                    Si estás viendo este mensaje, significa que FlightHours está enviando correos electrónicos correctamente.
                </p>

                <div style="background-color: #EFF6FF; border-left: 4px solid #0047AB; padding: 1.25rem; margin: 1.5rem 0; border-radius: 8px;">
                    <p style="margin: 0; color: #1E40AF; font-size: 14px;"><strong>Configuración SMTP exitosa</strong></p>
                    <p style="margin: 0.5rem 0 0 0; color: #1E40AF; font-size: 14px;">Tu servidor de correo está funcionando perfectamente con Keycloak.</p>
                </div>

                <div style="background-color: #F0FDF4; padding: 24px; border-radius: 12px; text-align: center; margin: 24px 0; border: 2px solid #10B981;">
                    <p style="color: #065F46; font-size: 18px; font-weight: 600; margin: 0;">
                        ✓ Todo está listo para enviar notificaciones
                    </p>
                </div>

                <p style="color: #64748B; font-size: 14px; text-align: center;">
                    Este es un correo de prueba generado automáticamente por Keycloak.
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
