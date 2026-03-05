<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FlightHours - Recuperar Contraseña</title>
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
            background: linear-gradient(135deg, #0047AB 0%, #4A90E2 100%);
            border-radius: 50%;
            display: flex; align-items: center; justify-content: center;
            margin: 0 auto 1.5rem;
            font-size: 40px; color: white;
            box-shadow: 0 4px 16px rgba(0, 71, 171, 0.3);
        }
        h1 { color: #2C3E50; font-size: 24px; font-weight: 600; margin-bottom: 1rem; }
        .message { color: #64748B; font-size: 15px; line-height: 1.6; margin-bottom: 1.5rem; }

        .form-group {
            margin-bottom: 1.25rem;
            text-align: left;
        }
        .form-group label {
            display: block;
            color: #374151;
            font-size: 14px;
            font-weight: 500;
            margin-bottom: 0.5rem;
        }
        .form-group input {
            width: 100%;
            padding: 12px 16px;
            border: 2px solid #E5E7EB;
            border-radius: 10px;
            font-size: 16px;
            transition: border-color 0.2s, box-shadow 0.2s;
        }
        .form-group input:focus {
            outline: none;
            border-color: #0047AB;
            box-shadow: 0 0 0 3px rgba(0, 71, 171, 0.1);
        }

        .btn {
            display: block;
            width: 100%;
            padding: 14px 24px;
            font-size: 16px;
            font-weight: 600;
            color: white;
            background: linear-gradient(135deg, #0047AB 0%, #4A90E2 100%);
            border: none;
            border-radius: 10px;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            margin-bottom: 1rem;
        }
        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 16px rgba(0, 71, 171, 0.3);
        }

        .btn-link {
            display: inline-block;
            color: #0047AB;
            text-decoration: none;
            font-size: 14px;
            font-weight: 500;
        }
        .btn-link:hover {
            text-decoration: underline;
        }

        .error-message {
            background: #FEE2E2;
            color: #DC2626;
            padding: 12px;
            border-radius: 8px;
            margin-bottom: 1rem;
            font-size: 14px;
        }

        .footer { margin-top: 1.5rem; color: #64748B; font-size: 13px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="card">
            <div class="logo">FlightHours</div>
            <div class="icon">🔐</div>

            <h1>Recuperar Contraseña</h1>

            <p class="message">
                Ingresa tu correo electrónico y te enviaremos un enlace para restablecer tu contraseña.
            </p>

            <#if message?has_content && (message.type = 'error' || message.type = 'warning')>
                <div class="error-message">
                    ${kcSanitize(message.summary)?no_esc}
                </div>
            </#if>

            <form action="${url.loginAction}" method="post">
                <div class="form-group">
                    <label for="username">Correo electrónico</label>
                    <input type="email" id="username" name="username"
                           placeholder="tu@email.com"
                           autofocus required
                           <#if auth?has_content && auth.showUsername()>
                               value="${auth.attemptedUsername}"
                           </#if>
                    >
                </div>

                <button type="submit" class="btn">Enviar enlace de recuperación</button>

                <a href="${url.loginUrl}" class="btn-link">← Volver al inicio de sesión</a>
            </form>

            <div class="footer">
                <p>© ${.now?string('yyyy')} FlightHours</p>
            </div>
        </div>
    </div>
</body>
</html>
