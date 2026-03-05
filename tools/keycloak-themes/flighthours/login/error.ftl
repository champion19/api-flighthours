<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Flighthours - Error</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #0047AB 0%, #4A90E2 100%);
            background-attachment: fixed;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container { max-width: 480px; width: 100%; }
        .card {
            background: white;
            padding: 3rem 2.5rem;
            border-radius: 20px;
            box-shadow: 0 8px 32px rgba(0, 71, 171, 0.16);
            text-align: center;
        }
        .logo {
            color: #0047AB;
            font-size: 32px;
            font-weight: 700;
            margin-bottom: 2rem;
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 0.5rem;
        }
        .logo::before { content: '✈️'; font-size: 28px; }
        .icon {
            width: 80px; height: 80px;
            background: linear-gradient(135deg, #EF4444 0%, #DC2626 100%);
            border-radius: 50%;
            display: flex; align-items: center; justify-content: center;
            margin: 0 auto 1.5rem;
            font-size: 40px; color: white;
            box-shadow: 0 4px 16px rgba(239, 68, 68, 0.3);
        }
        h1 { color: #2C3E50; font-size: 24px; font-weight: 600; margin-bottom: 1rem; }
        .message { color: #64748B; font-size: 15px; line-height: 1.6; margin-bottom: 1.5rem; }
        .error-box {
            background: #FEF2F2;
            border-left: 4px solid #EF4444;
            border-radius: 10px;
            padding: 1rem;
            margin-bottom: 1.5rem;
            text-align: left;
        }
        .error-box p { color: #991B1B; font-size: 14px; margin: 0; }
        .app-box {
            background: linear-gradient(135deg, #EFF6FF 0%, #DBEAFE 100%);
            border: 2px solid #0047AB;
            border-radius: 12px;
            padding: 1.5rem;
            margin-bottom: 1.5rem;
        }
        .app-box h2 { color: #0047AB; font-size: 18px; margin-bottom: 0.5rem; font-weight: 600; }
        .app-box p { color: #2C3E50; font-size: 14px; margin: 0; }
        .app-icon { font-size: 48px; margin-bottom: 0.5rem; }
        .footer { margin-top: 1.5rem; color: #64748B; font-size: 13px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <div class="logo">Flighthours</div>
            <div class="icon">✕</div>

            <h1>Enlace inválido o expirado</h1>

            <#if message??>
                <p class="message">
                    <#if message.summary??>
                        ${message.summary}
                    <#else>
                        ${message}
                    </#if>
                </p>
            <#else>
                <p class="message">
                    El enlace que utilizaste ya no es válido o ha expirado.
                </p>
            </#if>

            <div class="app-box">
                <div class="app-icon">📱</div>
                <h2>Solicita un nuevo enlace</h2>
                <p>Abre la aplicación FlightHours y solicita un nuevo correo.</p>
            </div>

            <div class="error-box">
                <p><strong>¿Por qué pasó esto?</strong></p>
                <p>• El enlace puede haber expirado<br>
                   • Ya fue utilizado anteriormente<br>
                   • El enlace está incompleto</p>
            </div>

            <p class="message" style="margin-bottom: 0;">
                Puedes cerrar esta ventana.
            </p>

            <div class="footer">
                <p>© ${.now?string('yyyy')} Flighthours</p>
            </div>
        </div>
    </div>
</body>
</html>
