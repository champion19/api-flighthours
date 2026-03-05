<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FlightHours - Actualizar Contraseña</title>
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
        .password-requirements {
            background: #F3F4F6;
            border-radius: 8px;
            padding: 12px;
            margin-bottom: 1.5rem;
            text-align: left;
        }
        .password-requirements p {
            font-size: 13px;
            color: #6B7280;
            margin-bottom: 4px;
        }
        .password-requirements ul {
            margin: 0;
            padding-left: 20px;
        }
        .password-requirements li {
            font-size: 12px;
            color: #6B7280;
            margin-bottom: 2px;
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
        }
        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 16px rgba(0, 71, 171, 0.3);
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

            <h1>Actualizar Contraseña</h1>

            <p class="message">
                Ingresa tu nueva contraseña para completar la actualización de tu cuenta.
            </p>

            <#if message?? && (message.type!'') = 'error'>
                <div class="error-message">
                    ${message.summary!'Error al procesar la solicitud'}
                </div>
            </#if>

            <#-- Formulario nativo de Keycloak -->
            <form action="${url.loginAction}" method="post">
                <div class="form-group">
                    <label for="password-new">Nueva contraseña</label>
                    <input type="password" id="password-new" name="password-new"
                           required minlength="8"
                           placeholder="Mínimo 8 caracteres"
                           autofocus autocomplete="new-password">
                </div>

                <div class="form-group">
                    <label for="password-confirm">Confirmar contraseña</label>
                    <input type="password" id="password-confirm" name="password-confirm"
                           required minlength="8"
                           placeholder="Repite la contraseña"
                           autocomplete="new-password">
                </div>

                <div class="password-requirements">
                    <p><strong>Requisitos de la contraseña:</strong></p>
                    <ul>
                        <li>Mínimo 8 caracteres</li>
                        <li>Se recomienda incluir mayúsculas, minúsculas y números</li>
                    </ul>
                </div>

                <button type="submit" class="btn">Actualizar Contraseña</button>
            </form>

            <div class="footer">
                <p>© ${.now?string('yyyy')} FlightHours</p>
            </div>
        </div>
    </div>
</body>
</html>
