<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Flighthours - Verificando correo</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            background: linear-gradient(135deg, #0047AB 0%, #4A90E2 100%);
            background-attachment: fixed;
            min-height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            padding: 20px;
        }
        .container {
            background: white;
            max-width: 480px;
            width: 100%;
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
            width: 80px;
            height: 80px;
            background: linear-gradient(135deg, #0047AB 0%, #4A90E2 100%);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            margin: 0 auto 1.5rem;
            font-size: 40px;
            color: white;
            box-shadow: 0 4px 16px rgba(0, 71, 171, 0.3);
        }
        h1 {
            color: #2C3E50;
            font-size: 24px;
            margin-bottom: 1rem;
            font-weight: 600;
        }
        .message {
            color: #64748B;
            font-size: 15px;
            line-height: 1.6;
            margin-bottom: 1.5rem;
        }
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
        .btn {
            display: inline-block;
            padding: 14px 32px;
            background: linear-gradient(135deg, #0047AB 0%, #1E88E5 100%);
            color: white;
            text-decoration: none;
            border-radius: 10px;
            font-weight: 600;
            font-size: 16px;
            border: none;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            box-shadow: 0 4px 12px rgba(0, 71, 171, 0.3);
        }
        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(0, 71, 171, 0.4);
        }
        .footer {
            margin-top: 2rem;
            color: #64748B;
            font-size: 13px;
        }
        .info-box {
            background: #FFFBEB;
            border-left: 4px solid #F59E0B;
            border-radius: 8px;
            padding: 1rem;
            margin-bottom: 1.5rem;
            text-align: left;
        }
        .info-box p {
            color: #92400E;
            font-size: 14px;
            margin: 0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">Flighthours</div>

        <div class="icon">📧</div>

        <h1>Verificar correo electrónico</h1>

        <p class="message">
            Te hemos enviado un correo electrónico con instrucciones para verificar tu cuenta.
        </p>

        <div class="info-box">
            <p><strong>Revisa tu bandeja de entrada</strong></p>
            <p>Si no ves el correo, revisa también la carpeta de spam.</p>
        </div>

        <div class="app-box">
            <div class="app-icon">📱</div>
            <h2>¿Ya verificaste tu correo?</h2>
            <p>Vuelve a la aplicación FlightHours e inicia sesión.</p>
        </div>

        <#if url?? && url.loginAction??>
            <form action="${url.loginAction}" method="post">
                <button type="submit" class="btn">Ya verifiqué mi correo</button>
            </form>
        </#if>

        <div class="footer">
            © ${.now?string('yyyy')} Flighthours
        </div>
    </div>
</body>
</html>
