<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FlightHours - Verificar Email</title>
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
                <h1 style="color: #2C3E50; text-align: center; margin-bottom: 1.5rem; font-size: 24px; font-weight: 600;">Bienvenido a FlightHours</h1>

                <p style="margin: 1rem 0; color: #64748B; font-size: 15px;">Gracias por registrarte en nuestra plataforma.</p>

                <p style="margin: 1rem 0; color: #64748B; font-size: 15px;">Para completar tu registro, por favor haz clic en el siguiente enlace:</p>

                <div style="text-align: center; margin: 1.5rem 0;">
                    <a href="${link}" style="display: inline-block; padding: 14px 28px; background-color: #0047AB; color: #FFFFFF; text-decoration: none; border-radius: 10px; font-weight: 600; font-size: 16px;">
                        Verificar mi email
                    </a>
                </div>

                <p style="color: #64748B; font-size: 14px; word-break: break-all;">Si no puedes hacer clic en el botón, copia y pega este enlace en tu navegador:</p>
                <p style="color: #64748B; font-size: 14px; word-break: break-all;">${link}</p>

                <div style="background-color: #FFFBEB; border-left: 4px solid #F59E0B; padding: 1.25rem; margin: 1.5rem 0; border-radius: 8px;">
                    <p style="margin: 0; color: #92400E; font-size: 14px;"><strong>IMPORTANTE:</strong> Este enlace expirará en 15 minutos por seguridad.</p>
                </div>

                <p style="color: #64748B; font-size: 13px; margin-top: 30px;">
                    Si no solicitaste este registro, puedes ignorar este correo de forma segura.
                </p>
            </div>

            <!-- Footer -->
            <div style="margin-top: 2rem; padding-top: 1.5rem; border-top: 1px solid #151ec3ff; text-align: center; font-size: 13px; color: #64748B;">
                <p><a href="mailto:soporte@flighthours.com" style="color: #0047AB; text-decoration: none; font-weight: 500;">soporte@flighthours.com</a></p>
                <p>&copy; ${.now?string('yyyy')} FlightHours. Todos los derechos reservados.</p>
            </div>

        </div>
    </div>
</body>
</html>
