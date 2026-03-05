<#macro emailLayout>
<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>FlightHours</title>
    <style>
        body {
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background-color: #F5F7FA;
            margin: 0;
            padding: 2rem;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
        }
        .card {
            background: white;
            padding: 2.5rem;
            border-radius: 16px;
            box-shadow: 0 4px 16px rgba(0, 71, 171, 0.12);
        }
        .header {
            text-align: center;
            margin-bottom: 2rem;
            padding-bottom: 1.5rem;
            border-bottom: 2px solid #0047AB;
        }
        .header h1 {
            color: #2C3E50;
            font-size: 24px;
            font-weight: 600;
            margin: 0.5rem 0;
        }
        .logo {
            color: #0047AB;
            font-size: 28px;
            font-weight: 700;
            margin-bottom: 0.5rem;
            letter-spacing: -0.5px;
        }
        .logo::before {
            content: '✈️ ';
            font-size: 24px;
        }
        .content {
            color: #2C3E50;
            line-height: 1.6;
        }
        .content p {
            margin: 1rem 0;
            color: #64748B;
            font-size: 15px;
        }
        .button {
            display: inline-block;
            margin: 1.5rem 0;
            padding: 14px 28px;
            background-color: #0047AB;
            color: #FFFFFF !important;
            text-decoration: none;
            border-radius: 10px;
            font-weight: 600;
            font-size: 16px;
        }
        .button:hover {
            background-color: #003d96;
        }
        .link-text {
            color: #64748B;
            font-size: 14px;
            word-break: break-all;
        }
        .footer {
            margin-top: 2rem;
            padding-top: 1.5rem;
            border-top: 1px solid #E1E8ED;
            text-align: center;
            font-size: 13px;
            color: #64748B;
        }
        .footer a {
            color: #0047AB;
            text-decoration: none;
            font-weight: 500;
        }
        .info-box {
            background-color: #EFF6FF;
            border-left: 4px solid #0047AB;
            padding: 1.25rem;
            margin: 1.5rem 0;
            border-radius: 8px;
        }
        .info-box p {
            margin: 0;
            color: #1E40AF;
            font-size: 14px;
        }
        .warning-box {
            background-color: #FFFBEB;
            border-left: 4px solid #F59E0B;
            padding: 1.25rem;
            margin: 1.5rem 0;
            border-radius: 8px;
        }
        .warning-box p {
            margin: 0;
            color: #92400E;
            font-size: 14px;
        }
        @media only screen and (max-width: 600px) {
            body {
                padding: 1rem;
            }
            .card {
                padding: 2rem 1.5rem;
            }
            .button {
                display: block;
                text-align: center;
            }
        }
    </style>

</head>
<body>
    <div class="container">
        <div class="card">
            <div class="header">
                <div class="logo">FlightHours</div>
            </div>
            <div class="content">
                <#nested>
            </div>
            <div class="footer">
                <p><a href="mailto:soporte@flighthours.com">soporte@flighthours.com</a></p>
                <p>&copy; ${.now?string('yyyy')} FlightHours. Todos los derechos reservados.</p>
            </div>
        </div>
    </div>
</body>
</html>
</#macro>
