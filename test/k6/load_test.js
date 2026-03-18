import http from 'k6/http';
import { check, sleep, group } from 'k6';
import { Rate, Trend, Counter } from 'k6/metrics';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import exec from 'k6/execution';

// ============================================
// FlightHours API — Pruebas de Carga v1.1
// ============================================
// ESCENARIOS:
//   smoke:      Sanity check (1 usuario, 30s)
//   load:       Carga normal progresiva (hasta 100 VUs)
//   stress:     Estrés controlado (hasta 800 VUs)
//   spike:      Pico repentino moderado (hasta 150 VUs)
//   endurance:  Prueba de resistencia (2 horas, carga media)
//   breakpoint: Buscar límite real (incrementos de 50 VUs)
//   thousand:   Carga masiva (1000 VUs concurrentes)
//
// Uso:
//   k6 run test/k6/load_test.js
//   k6 run --env SCENARIO=smoke test/k6/load_test.js
//   k6 run --env SCENARIO=thousand test/k6/load_test.js
//
// Con Grafana (Prometheus remote write):
//   k6 run -o experimental-prometheus-rw \
//     --env K6_PROMETHEUS_RW_SERVER_URL=http://localhost:9090/api/v1/write \
//     --env SCENARIO=thousand test/k6/load_test.js
//
//   O simplemente: ./test/k6/run_k6_grafana.sh [scenario]

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8081';
const API = `${BASE_URL}/flighthours/api/v1`;
const USERNAME = __ENV.K6_USERNAME || 'josedalogo@yopmail.com';
const PASSWORD = __ENV.K6_PASSWORD || 'Jose2026!';
const SCENARIO = __ENV.SCENARIO || 'load';

// Health check abort: máximo de fallos consecutivos antes de abortar
const MAX_HEALTH_FAILURES = Number.parseInt(__ENV.MAX_HEALTH_FAILURES || '3', 10);

// Token refresh: refrescar 2 minutos antes de expirar (margen de seguridad)
const TOKEN_REFRESH_MARGIN_S = 120;

// ── Métricas de negocio personalizadas ─────────────────
const loginErrors = new Rate('login_errors');
const catalogTime = new Trend('catalog_response_time');
const authTime = new Trend('auth_response_time');
const flightQueryTime = new Trend('flight_query_response_time');
const businessErrors = new Rate('business_errors');
const throughput = new Counter('requests_total');
const tokenRefreshes = new Counter('token_refreshes');

// ── Estado per-VU ──────────────────────────────────────
let consecutiveHealthFailures = 0;
let vuToken = null;
let vuTokenExpiresAt = 0; // timestamp en segundos

// ── ESCENARIOS REALISTAS ───────────────────────────────

const scenarios = {
  smoke: {
    smoke: {
      executor: 'constant-vus',
      vus: 1,
      duration: '30s',
      tags: { test_type: 'smoke' },
    },
  },

  load: {
    morning_rush: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 50 },
        { duration: '5m', target: 50 },
        { duration: '3m', target: 100 },
        { duration: '5m', target: 100 },
        { duration: '2m', target: 80 },
      ],
      tags: { test_type: 'load', period: 'morning' },
    },
    afternoon_normal: {
      executor: 'ramping-vus',
      startVUs: 0,
      startTime: '17m',
      stages: [
        { duration: '2m', target: 60 },
        { duration: '10m', target: 60 },
        { duration: '2m', target: 40 },
      ],
      tags: { test_type: 'load', period: 'afternoon' },
    },
    evening_low: {
      executor: 'ramping-vus',
      startVUs: 0,
      startTime: '31m',
      stages: [
        { duration: '2m', target: 20 },
        { duration: '5m', target: 20 },
        { duration: '2m', target: 0 },
      ],
      tags: { test_type: 'load', period: 'evening' },
    },
  },

  stress: {
    stress_gradual: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '3m', target: 100 },
        { duration: '2m', target: 100 },
        { duration: '3m', target: 200 },
        { duration: '2m', target: 200 },
        { duration: '3m', target: 300 },
        { duration: '2m', target: 300 },
        { duration: '3m', target: 400 },
        { duration: '2m', target: 400 },
        { duration: '3m', target: 500 },
        { duration: '2m', target: 500 },
        { duration: '3m', target: 600 },
        { duration: '2m', target: 600 },
        { duration: '3m', target: 700 },
        { duration: '2m', target: 700 },
        { duration: '3m', target: 800 },
        { duration: '2m', target: 800 },
        { duration: '3m', target: 0 },
      ],
      tags: { test_type: 'stress' },
    },
  },

  spike: {
    spike_moderate: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '5m', target: 50 },
        { duration: '30s', target: 150 },
        { duration: '3m', target: 150 },
        { duration: '1m', target: 80 },
        { duration: '5m', target: 80 },
        { duration: '2m', target: 0 },
      ],
      tags: { test_type: 'spike' },
    },
  },

  endurance: {
    endurance_test: {
      executor: 'constant-vus',
      vus: 75,
      duration: '2h',
      tags: { test_type: 'endurance' },
    },
  },

  breakpoint: {
    breakpoint_search: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 50 },
        { duration: '3m', target: 50 },
        { duration: '2m', target: 100 },
        { duration: '3m', target: 100 },
        { duration: '2m', target: 150 },
        { duration: '3m', target: 150 },
        { duration: '2m', target: 200 },
        { duration: '3m', target: 200 },
        { duration: '2m', target: 250 },
        { duration: '3m', target: 250 },
        { duration: '2m', target: 300 },
        { duration: '3m', target: 300 },
        { duration: '2m', target: 350 },
        { duration: '3m', target: 350 },
        { duration: '2m', target: 400 },
        { duration: '5m', target: 400 },
      ],
      tags: { test_type: 'breakpoint' },
    },
  },

  // ⭐ Escenario de 1000 usuarios concurrentes
  thousand: {
    thousand_ramp: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '2m', target: 200 },   // warm-up
        { duration: '3m', target: 200 },   // hold 200
        { duration: '2m', target: 500 },   // escalada media
        { duration: '3m', target: 500 },   // hold 500
        { duration: '2m', target: 800 },   // escalada alta
        { duration: '3m', target: 800 },   // hold 800
        { duration: '2m', target: 1000 },  // pico máximo 🔥
        { duration: '5m', target: 1000 },  // hold sostenido
        { duration: '3m', target: 0 },     // cool-down
      ],
      tags: { test_type: 'thousand' },
    },
  },
};

export const options = {
  scenarios: scenarios[SCENARIO] || scenarios.load,

  thresholds: {
    http_req_duration: ['p(95)<1000', 'p(99)<2000'],
    http_req_failed: ['rate<0.10'],
    login_errors: ['rate<0.05'],
    business_errors: ['rate<0.05'],
  },
};

// ── SETUP: Login inicial ───────────────────────────────
export function setup() {
  let token = null;
  let refreshToken = null;
  let expiresIn = 0;
  let attempts = 0;
  const maxAttempts = 3;

  while (!token && attempts < maxAttempts) {
    attempts++;
    console.log(`🔑 Intentando login (intento ${attempts}/${maxAttempts})...`);

    const loginRes = http.post(
      `${API}/login`,
      JSON.stringify({ email: USERNAME, password: PASSWORD }),
      {
        headers: { 'Content-Type': 'application/json' },
        timeout: '10s',
        tags: { name: 'login' },
      }
    );

    const success = check(loginRes, {
      'login status 200': (r) => r.status === 200,
      'login response time < 2s': (r) => r.timings.duration < 2000,
    });

    if (success) {
      try {
        const body = JSON.parse(loginRes.body);
        // FlightHours envuelve la respuesta en data
        const data = body.data || body;
        token = data.access_token;
        refreshToken = data.refresh_token;
        expiresIn = data.expires_in || 900; // default 15 min
        console.log(`✅ Login OK en intento ${attempts} — token expira en ${expiresIn}s`);
        break;
      } catch (e) {
        console.error('❌ Error parseando respuesta de login:', e);
      }
    } else {
      console.warn(`⚠️ Login falló (status ${loginRes.status}): ${loginRes.body?.substring(0, 100)}`);
      sleep(2);
    }
  }

  if (!token) {
    console.error('❌❌ Login falló después de todos los intentos. Abortando tests autenticados.');
  }

  const now = Math.floor(Date.now() / 1000);
  return {
    token,
    refreshToken,
    loginSuccess: !!token,
    expiresIn,
    tokenObtainedAt: now,
  };
}

// ── FUNCIONES AUXILIARES ───────────────────────────────

function categorizeError(status) {
  if (status >= 500) {
    businessErrors.add(1);
    return 'server_error';
  }
  // Registrar éxito para que el Rate refleje la proporción real
  businessErrors.add(0);
  if (status >= 400) {
    return 'client_error';
  }
  return 'success';
}

// ⭐ Refresh PROACTIVO de token (por tiempo, no por 401)
// Cada VU refresca su propio token antes de que expire
function ensureFreshToken(data) {
  const now = Math.floor(Date.now() / 1000);

  // Si el VU ya tiene un token vigente, usarlo
  if (vuToken && now < vuTokenExpiresAt) {
    return vuToken;
  }

  // Si el token del setup aún es vigente, usarlo
  const setupTokenExpiresAt = data.tokenObtainedAt + data.expiresIn - TOKEN_REFRESH_MARGIN_S;
  if (!vuToken && now < setupTokenExpiresAt) {
    return data.token;
  }

  // Intentar refresh con refresh_token primero
  if (data.refreshToken) {
    const refreshRes = http.post(
      `${API}/auth/refresh`,
      JSON.stringify({ refresh_token: data.refreshToken }),
      {
        headers: { 'Content-Type': 'application/json' },
        timeout: '10s',
        tags: { name: 'token_refresh' },
      }
    );

    if (refreshRes.status === 200) {
      try {
        const body = JSON.parse(refreshRes.body);
        const respData = body.data || body;
        vuToken = respData.access_token;
        const expiresIn = respData.expires_in || 900;
        vuTokenExpiresAt = now + expiresIn - TOKEN_REFRESH_MARGIN_S;
        tokenRefreshes.add(1);
        return vuToken;
      } catch (e) {
        console.warn('⚠️ Error parseando refresh response:', e);
      }
    }
  }

  // Fallback: re-login
  const loginRes = http.post(
    `${API}/login`,
    JSON.stringify({ email: USERNAME, password: PASSWORD }),
    {
      headers: { 'Content-Type': 'application/json' },
      timeout: '10s',
      tags: { name: 'token_refresh' },
    }
  );

  if (loginRes.status === 200) {
    try {
      const body = JSON.parse(loginRes.body);
      const respData = body.data || body;
      vuToken = respData.access_token;
      const expiresIn = respData.expires_in || 900;
      vuTokenExpiresAt = now + expiresIn - TOKEN_REFRESH_MARGIN_S;
      tokenRefreshes.add(1);
      return vuToken;
    } catch (e) {
      console.warn('⚠️ Error parseando re-login response:', e);
    }
  }

  // Fallback: devolver el último token disponible en caso de fallo de refresh y re-login
  return vuToken || data.token;
}

// ── HELPER FUNCTIONS (extracted to reduce cognitive complexity) ──

function runHealthCheck() {
  group('01_Health', () => {
    const res = http.get(`${BASE_URL}/health`, {
      timeout: '5s',
      tags: { name: 'health_check' },
    });
    const ok = check(res, {
      'health 200': (r) => r.status === 200,
      'health fast < 500ms': (r) => r.timings.duration < 500,
    });

    if (ok) {
      consecutiveHealthFailures = 0;
      return;
    }

    consecutiveHealthFailures++;
    console.warn(
      `💔 Health check falló (${consecutiveHealthFailures}/${MAX_HEALTH_FAILURES}): ` +
      `status=${res.status} (${res.timings.duration}ms)`
    );
    categorizeError(res.status);

    if (consecutiveHealthFailures >= MAX_HEALTH_FAILURES) {
      console.error(
        `🚨🚨🚨 ABORT: El servidor no responde después de ${MAX_HEALTH_FAILURES} intentos consecutivos. ` +
        `Abortando test para evitar miles de errores innecesarios.`
      );
      exec.test.abort(
        `Server health check failed ${MAX_HEALTH_FAILURES} consecutive times — server likely crashed`
      );
    }
  });
}

function runCatalogQueries() {
  group('02_Catalogs', () => {
    const endpoints = [
      { path: '/airlines', tag: 'catalog_airlines' },
      { path: '/airports', tag: 'catalog_airports' },
      { path: '/engines', tag: 'catalog_engines' },
      { path: '/routes', tag: 'catalog_routes' },
      { path: '/manufacturers', tag: 'catalog_manufacturers' },
      { path: '/aircraft-models', tag: 'catalog_aircraft_models' },
      { path: '/crew-member-types', tag: 'catalog_crew_types' },
      { path: '/airline-routes', tag: 'catalog_airline_routes' },
    ];

    const count = randomIntBetween(3, 4);
    endpoints.sort(() => Math.random() - 0.5);
    const shuffled = endpoints.slice(0, count);

    for (const ep of shuffled) {
      const res = http.get(`${API}${ep.path}`, {
        timeout: '10s',
        tags: { name: ep.tag },
      });
      catalogTime.add(res.timings.duration);

      const ok = check(res, {
        [`${ep.path} status 200`]: (r) => r.status === 200,
        [`${ep.path} response < 1s`]: (r) => r.timings.duration < 1000,
      });

      if (!ok) {
        console.warn(`📚 Catalog ${ep.path} falló: ${res.status} (${res.timings.duration}ms)`);
        categorizeError(res.status);
      }

      sleep(Math.random() * 0.3 + 0.1);
    }
  });
}

function checkEndpoint(groupName, url, authHeaders, metricFn, label, checks, timeoutVal) {
  group(groupName, () => {
    const res = http.get(url, {
      ...authHeaders,
      timeout: timeoutVal || '5s',
      tags: { name: label },
    });
    metricFn.add(res.timings.duration);

    const ok = check(res, checks);
    if (!ok) {
      console.warn(`${label} falló: ${res.status} (${res.timings.duration}ms)`);
      categorizeError(res.status);
    }
    return res;
  });
}

function runAuthenticatedEndpoints(authHeaders) {
  if (Math.random() < 0.3) {
    checkEndpoint('03_Profile', `${API}/employees`, authHeaders, authTime, 'profile_get', {
      'employees 200': (r) => r.status === 200,
      'profile < 800ms': (r) => r.timings.duration < 800,
    });
    sleep(Math.random() * 0.5 + 0.2);
  }

  if (Math.random() < 0.4) {
    checkEndpoint('04_AirlineInfo', `${API}/employees/airline`, authHeaders, authTime, 'airline_info_get', {
      'airline info 200/204': (r) => r.status === 200 || r.status === 204,
      'airline info < 800ms': (r) => r.timings.duration < 800,
    });
    sleep(Math.random() * 0.5 + 0.3);
  }

  if (Math.random() < 0.5) {
    checkEndpoint('05_DailyLogbooks', `${API}/daily-logbooks`, authHeaders, flightQueryTime, 'daily_logbooks_list', {
      'daily-logbooks 200/204': (r) => r.status === 200 || r.status === 204,
      'daily-logbooks < 1.5s': (r) => r.timings.duration < 1500,
    }, '10s');
    sleep(Math.random() * 0.5 + 0.2);
  }

  if (Math.random() < 0.4) {
    checkEndpoint('06_MyAirlineRoutes', `${API}/employees/airline-routes`, authHeaders, flightQueryTime, 'my_airline_routes', {
      'my airline-routes 200/204': (r) => r.status === 200 || r.status === 204,
      'my airline-routes < 1s': (r) => r.timings.duration < 1000,
    }, '10s');
    sleep(Math.random() * 0.3 + 0.2);
  }

  if (Math.random() < 0.3) {
    checkEndpoint('07_TailNumbers', `${API}/tail-numbers`, authHeaders, authTime, 'tail_numbers_list', {
      'tail-numbers 200/204': (r) => r.status === 200 || r.status === 204,
      'tail-numbers < 800ms': (r) => r.timings.duration < 800,
    });
    sleep(Math.random() * 0.3 + 0.2);
  }

  // ⭐ Endpoint principal del negocio de FlightHours (80% usuarios)
  if (Math.random() < 0.8) {
    group('08_FlightSummary', () => {
      const res = http.get(`${API}/employees/flight-hours-summary`, {
        ...authHeaders,
        timeout: '15s',
        tags: { name: 'flight_hours_summary' },
      });
      flightQueryTime.add(res.timings.duration);

      const ok = check(res, {
        'flight-summary 200/204': (r) => r.status === 200 || r.status === 204,
        'flight-summary < 2s': (r) => r.timings.duration < 2000,
      });

      if (!ok) {
        console.warn(`⏱️ Flight summary falló: ${res.status} (${res.timings.duration}ms)`);
        categorizeError(res.status);
      } else if (res.timings.duration > 1000) {
        console.log(`🐌 Query de horas de vuelo lenta: ${res.timings.duration}ms`);
      }
    });
  }

  if (Math.random() < 0.6) {
    checkEndpoint('09_FlightAlerts', `${API}/employees/flight-alerts`, authHeaders, flightQueryTime, 'flight_alerts', {
      'flight-alerts 200/204': (r) => r.status === 200 || r.status === 204,
      'flight-alerts < 1s': (r) => r.timings.duration < 1000,
    }, '10s');
    sleep(Math.random() * 0.3 + 0.2);
  }

  if (Math.random() < 0.5) {
    checkEndpoint('10_RecentFlights', `${API}/employees/recent-flights`, authHeaders, flightQueryTime, 'recent_flights', {
      'recent-flights 200/204': (r) => r.status === 200 || r.status === 204,
      'recent-flights < 1.5s': (r) => r.timings.duration < 1500,
    }, '10s');
  }
}

function getThinkTime() {
  const rand = Math.random();
  if (rand < 0.6) return Math.random() * 2 + 1;
  if (rand < 0.9) return Math.random() * 2 + 3;
  return Math.random() * 5 + 5;
}

// ── MAIN TEST FUNCTION ─────────────────────────────────
export default function mainTest(data) {
  const activeToken = data.loginSuccess ? ensureFreshToken(data) : null;
  const authHeaders = activeToken
    ? { headers: { Authorization: `Bearer ${activeToken}` } }
    : {};

  throughput.add(1);

  runHealthCheck();
  sleep(Math.random() * 0.5 + 0.2);

  if (Math.random() < 0.7) {
    runCatalogQueries();
  }

  sleep(Math.random() * 0.5 + 0.3);

  if (activeToken && data.loginSuccess) {
    runAuthenticatedEndpoints(authHeaders);
  } else {
    loginErrors.add(1);
    sleep(1);
  }

  sleep(getThinkTime());
}

// ── TEARDOWN ───────────────────────────────────────────
export function teardown(data) {
  console.log('\n🏁 RESUMEN DE PRUEBA');
  console.log('===================');

  if (data.loginSuccess) {
    console.log('✅ Autenticación: OK');
  } else {
    console.log('❌ Autenticación: FALLÓ - Solo se probaron endpoints públicos');
  }

  console.log(`📊 Escenario ejecutado: ${SCENARIO}`);
  console.log('📈 Revisa Grafana/dashboard para métricas detalladas');
  console.log('💡 Tip: Si ves muchos errores 5xx, revisa logs del servidor Go');
}
