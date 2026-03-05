<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FlightHours - Procesando</title>
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
            background: linear-gradient(135deg, #10B981 0%, #059669 100%);
            border-radius: 50%;
            display: flex; align-items: center; justify-content: center;
            margin: 0 auto 1.5rem;
            font-size: 40px; color: white;
            box-shadow: 0 4px 16px rgba(16, 185, 129, 0.3);
        }
        .icon.password {
            background: linear-gradient(135deg, #0047AB 0%, #4A90E2 100%);
            box-shadow: 0 4px 16px rgba(0, 71, 171, 0.3);
        }
        .icon.success {
            background: linear-gradient(135deg, #10B981 0%, #059669 100%);
        }
        .spinner {
            width: 50px; height: 50px;
            border: 4px solid #E1E8ED;
            border-top: 4px solid #0047AB;
            border-radius: 50%;
            animation: spin 1s linear infinite;
            margin: 0 auto 1.5rem;
        }
        @keyframes spin { 100% { transform: rotate(360deg); } }
        h1 { color: #2C3E50; font-size: 24px; font-weight: 600; margin-bottom: 1rem; }
        .message { color: #64748B; font-size: 15px; line-height: 1.6; margin-bottom: 1.5rem; }
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

        /* Form Styles */
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
        .password-requirements p { font-size: 13px; color: #6B7280; margin-bottom: 4px; }
        .password-requirements ul { margin: 0; padding-left: 20px; }
        .password-requirements li { font-size: 12px; color: #6B7280; margin-bottom: 2px; }

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
            text-decoration: none;
            text-align: center;
        }
        .btn:hover:not(:disabled) {
            transform: translateY(-2px);
            box-shadow: 0 4px 16px rgba(0, 71, 171, 0.3);
        }
        .btn:disabled { opacity: 0.7; cursor: not-allowed; }
        .btn.loading::after {
            content: '';
            display: inline-block;
            width: 16px; height: 16px;
            border: 2px solid white;
            border-top-color: transparent;
            border-radius: 50%;
            margin-left: 8px;
            animation: spin 0.8s linear infinite;
            vertical-align: middle;
        }
        .error-message {
            background: #FEE2E2;
            color: #DC2626;
            padding: 12px;
            border-radius: 8px;
            margin-bottom: 1rem;
            font-size: 14px;
            display: none;
        }
        .footer { margin-top: 1.5rem; color: #64748B; font-size: 13px; }
        .action-list {
            text-align: left;
            background: #F3F4F6;
            border-radius: 8px;
            padding: 1rem;
            margin-bottom: 1.5rem;
        }
        .action-list ul { margin: 0.5rem 0 0 1.5rem; color: #374151; }
        .hide { display: none !important; }
    </style>
</head>
<body>
    <div class="container">
        <!-- Loading State -->
        <div class="card" id="loadingCard">
            <div class="logo">FlightHours</div>
            <div class="spinner"></div>
            <h1>Procesando...</h1>
            <p class="message">Por favor espera un momento.</p>
        </div>

        <!-- Password Form Card -->
        <div class="card hide" id="passwordFormCard">
            <div class="logo">FlightHours</div>
            <div class="icon password">🔐</div>

            <h1>Actualizar Contraseña</h1>

            <p class="message">
                Ingresa tu nueva contraseña para completar la recuperación de tu cuenta.
            </p>

            <div class="error-message" id="errorMessage"></div>

            <form id="passwordForm">
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

                <button type="submit" class="btn" id="submitBtn">Actualizar Contraseña</button>
            </form>

            <div class="footer">
                <p>© ${.now?string('yyyy')} FlightHours</p>
            </div>
        </div>

        <!-- Success Card (for password update) -->
        <div class="card hide" id="passwordSuccessCard">
            <div class="logo">FlightHours</div>
            <div class="icon success">✓</div>

            <h1 style="color: #10B981;">¡Contraseña Actualizada!</h1>

            <p class="message">
                Tu contraseña ha sido actualizada exitosamente.
            </p>

            <div class="app-box">
                <div class="app-icon">📱</div>
                <h2>Abre la aplicación FlightHours</h2>
                <p>Ya puedes iniciar sesión con tu nueva contraseña.</p>
            </div>

            <p class="message" style="color: #10B981; font-weight: 600; margin-bottom: 0;">
                Puedes cerrar esta ventana.
            </p>

            <div class="footer">
                <p>© ${.now?string('yyyy')} FlightHours</p>
            </div>
        </div>

        <!-- Generic Success Card (for other actions) -->
        <div class="card hide" id="genericSuccessCard">
            <div class="logo">FlightHours</div>
            <div class="icon success">✓</div>

            <h1>¡Acción completada!</h1>

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
                    La operación se ha completado correctamente.
                </p>
            </#if>

            <#if requiredActions?? && (requiredActions?size > 0)>
                <div class="action-list">
                    <strong>Acciones completadas:</strong>
                    <ul>
                        <#list requiredActions as action>
                            <li>${action}</li>
                        </#list>
                    </ul>
                </div>
            </#if>

            <div class="app-box">
                <div class="app-icon">📱</div>
                <h2>Abre la aplicación FlightHours</h2>
                <p>Ya puedes continuar usando la aplicación.</p>
            </div>

            <#if actionUri??>
                <a href="${actionUri}" class="btn">Continuar</a>
            <#elseif pageRedirectUri??>
                <a href="${pageRedirectUri}" class="btn">Volver a la aplicación</a>
            <#elseif client?? && client.baseUrl??>
                <a href="${client.baseUrl}" class="btn">Volver a la aplicación</a>
            <#else>
                <p class="message" style="color: #10B981; font-weight: 600; margin-bottom: 0;">
                    Puedes cerrar esta ventana.
                </p>
            </#if>

            <div class="footer">
                © ${.now?string('yyyy')} FlightHours
            </div>
        </div>
    </div>

    <script>
        (function() {
            console.log('[FlightHours] info.ftl cargado');

            var BACKEND_URL = 'http://localhost:8081/flighthours/api/v1/auth/update-password';

            var currentUrl = window.location.href;
            console.log('[FlightHours] URL:', currentUrl);

            // Detectar si es UPDATE_PASSWORD
            var isUpdatePassword = currentUrl.includes('UPDATE_PASSWORD') ||
                                   currentUrl.includes('update-password') ||
                                   currentUrl.includes('reset-credentials');

            // Buscar en requiredActions si existe (variable FreeMarker)
            var requiredActionsHtml = document.querySelector('.action-list');
            if (requiredActionsHtml && requiredActionsHtml.innerHTML.includes('UPDATE_PASSWORD')) {
                isUpdatePassword = true;
            }

            console.log('[FlightHours] isUpdatePassword:', isUpdatePassword);

            // Extraer el token 'key' de la URL
            function getTokenFromUrl() {
                var urlParams = new URLSearchParams(window.location.search);
                return urlParams.get('key') || '';
            }

            var token = getTokenFromUrl();
            console.log('[FlightHours] Token encontrado:', token ? 'Sí (' + token.substring(0, 20) + '...)' : 'No');

            // Elementos
            var loadingCard = document.getElementById('loadingCard');
            var passwordFormCard = document.getElementById('passwordFormCard');
            var passwordSuccessCard = document.getElementById('passwordSuccessCard');
            var genericSuccessCard = document.getElementById('genericSuccessCard');
            var errorMessage = document.getElementById('errorMessage');
            var passwordForm = document.getElementById('passwordForm');
            var submitBtn = document.getElementById('submitBtn');

            function showCard(cardId) {
                loadingCard.classList.add('hide');
                passwordFormCard.classList.add('hide');
                passwordSuccessCard.classList.add('hide');
                genericSuccessCard.classList.add('hide');
                document.getElementById(cardId).classList.remove('hide');
            }

            function showError(message) {
                errorMessage.textContent = message;
                errorMessage.style.display = 'block';
            }

            function hideError() {
                errorMessage.style.display = 'none';
            }

            function setLoading(isLoading) {
                if (isLoading) {
                    submitBtn.classList.add('loading');
                    submitBtn.disabled = true;
                    submitBtn.textContent = 'Actualizando...';
                } else {
                    submitBtn.classList.remove('loading');
                    submitBtn.disabled = false;
                    submitBtn.textContent = 'Actualizar Contraseña';
                }
            }

            // Si es UPDATE_PASSWORD y tenemos token, mostrar el formulario
            if (isUpdatePassword && token) {
                console.log('[FlightHours] Mostrando formulario de contraseña');
                showCard('passwordFormCard');

                // Manejar el submit del formulario
                passwordForm.addEventListener('submit', function(e) {
                    e.preventDefault();
                    hideError();

                    var newPassword = document.getElementById('password-new').value;
                    var confirmPassword = document.getElementById('password-confirm').value;

                    // Validaciones
                    if (newPassword.length < 8) {
                        showError('La contraseña debe tener al menos 8 caracteres.');
                        return;
                    }

                    if (newPassword !== confirmPassword) {
                        showError('Las contraseñas no coinciden.');
                        return;
                    }

                    setLoading(true);
                    console.log('[FlightHours] Enviando nueva contraseña al backend...');

                    // Enviar al backend
                    fetch(BACKEND_URL, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json'
                        },
                        body: JSON.stringify({
                            token: token,
                            new_password: newPassword,
                            confirm_password: confirmPassword
                        })
                    })
                    .then(function(response) {
                        return response.json().then(function(data) {
                            return { status: response.status, data: data };
                        });
                    })
                    .then(function(result) {
                        console.log('[FlightHours] Respuesta del backend:', result);

                        if (result.status >= 200 && result.status < 300 && result.data.success) {
                            console.log('[FlightHours] Contraseña actualizada exitosamente');
                            showCard('passwordSuccessCard');
                        } else {
                            var errorMsg = 'Error al actualizar la contraseña.';
                            if (result.data && result.data.message) {
                                errorMsg = result.data.message;
                            } else if (result.data && result.data.data && result.data.data.message) {
                                errorMsg = result.data.data.message;
                            }
                            console.error('[FlightHours] Error:', errorMsg);
                            setLoading(false);
                            showError(errorMsg);
                        }
                    })
                    .catch(function(error) {
                        console.error('[FlightHours] Error de red:', error);
                        setLoading(false);
                        showError('Error de conexión. Por favor intenta nuevamente.');
                    });
                });

            } else if (isUpdatePassword && !token) {
                // UPDATE_PASSWORD pero sin token - mostrar éxito genérico (ya se procesó antes)
                console.log('[FlightHours] UPDATE_PASSWORD sin token - mostrando éxito genérico');
                showCard('genericSuccessCard');

            } else {
                // Otra acción - mostrar éxito genérico
                console.log('[FlightHours] Acción genérica - mostrando éxito');
                showCard('genericSuccessCard');
            }
        })();
    </script>
</body>
</html>
