# TwitBet Backend

Backend transaccional para **ligas de apuestas simuladas**. Permite crear ligas entre amigos, administrar partidos y mercados, realizar apuestas con saldo virtual, cashout, bonos y recargas, todo con eventos en tiempo real vía WebSocket.

Desarrollado con **Go** (Gin Gonic), **PostgreSQL**, **Redis** y arquitectura **Hexagonal (Ports & Adapters)** con **Domain-Driven Design (DDD)**.

> **Repositorio:** `github.com/jmlc643/twitbet-backend`

---

## Funcionalidades

- **Autenticación y usuarios**
  - Registro con verificación de cuenta por OTP (email).
  - Login con JWT.
  - Recuperación de contraseña (OTP) y cambio de contraseña.
  - Perfil de usuario con actualización y subida de avatar (Cloudinary).
- **Ligas**
  - Creación de ligas privadas con código de invitación, saldo inicial y configuraciones (`maxRecharges`, `minBetsToQualify`, `hideStandings`).
  - Unirse a una liga, gestión de administradores, líderes (leaderboard) y estados.
- **Partidos y mercados**
  - Creación de partidos y mercados de apuesta (mínimo 2 opciones).
  - Actualización en vivo de estados y cuotas de mercados.
  - Resolución y cancelación de mercados (reembolso/anulación de apuestas).
- **Apuestas**
  - Colocar apuestas con saldo virtual (transacciones ACID), cashout en vivo.
  - Estados: `PENDING`, `ACCEPTED`, `WON`, `LOST`, `VOIDED`, `CASHOUT`, `REJECTED`.
- **Economía simulada**
  - Recargas de saldo, bonos de liga y apuestas con bono.
- **Tiempo real**
  - WebSocket (`/ws`) que difunde eventos de mercados a través de Redis Pub/Sub.

---

## Stack Tecnológico

| Capa        | Tecnología                                  |
| ----------- | ------------------------------------------- |
| Lenguaje    | Go 1.25                                     |
| Framework   | Gin Gonic                                   |
| Base de datos | PostgreSQL (GORM)                         |
| Cache / Pub-Sub | Redis (go-redis)                        |
| Tiempo real | Gorilla WebSocket                          |
| Email       | Brevo (transaccional) / SMTP (gomail)       |
| Almacenamiento | Cloudinary (avatares)                    |
| Autenticación | JWT (golang-jwt) + bcrypt                  |

---

## Arquitectura

Proyecto organizado en **Bounded Contexts** (`identity`, `league`) y dividido en las 3 capas hexagonales:

```
infrastructure  →  application  →  domain
```

- **`domain/`**: Entidades puras, value objects, errores de dominio, repositorios e interfaces de puertos. **Sin dependencias externas** (solo Go estándar).
- **`application/`**: Casos de uso que orquestan la lógica de negocio usando DTOs de entrada (`input/`) y salida (`output/`). Sin frameworks HTTP.
- **`infrastructure/`**: Handlers HTTP (Gin), middlewares JWT, adaptadores (bcrypt, JWT, Cloudinary), persistencia PostgreSQL (repositorios + mappers) y Redis.

### Estructura del proyecto

```
├── cmd/api/main.go                     # Punto de entrada del servidor
├── internal/
│   ├── identity/                       # Contexto: usuarios y autenticación
│   │   ├── domain/                     #   Entidades, repos, puertos, errores
│   │   ├── application/                #   Casos de uso (register, login, ...)
│   │   └── infrastructure/
│   │       ├── http/                   #   Router, handlers, middleware JWT
│   │       ├── adapter/                #   bcrypt, JWT, OTP (Redis), Cloudinary
│   │       └── persistence/            #   Repositorio GORM + mapper
│   ├── league/                         # Contexto: ligas, partidos, mercados, apuestas
│   │   ├── domain/
│   │   ├── application/
│   │   └── infrastructure/
│   │       ├── http/                   #   Router, handlers, DTOs, mapper
│   │       ├── persistence/postgres/   #   Repos GORM
│   │       └── redis/                  #   Publicador de eventos de mercado
│   └── platform/                       # Utilidades transversales
│       ├── config/                     #   Carga de variables de entorno
│       ├── database/                   #   Conexión y migraciones PostgreSQL
│       ├── redis/                      #   Cliente Redis
│       ├── email/                      #   Servicios de email (Brevo / gomail)
│       ├── http/middleware/            #   CORS
│       └── websocket/                  #   Hub + cliente WebSocket
├── Dockerfile
├── go.mod / go.sum
└── .env.example
```

---

## Requisitos previos

- **Go** 1.22+ (el módulo apunta a 1.25)
- **PostgreSQL** disponible
- **Redis** disponible (necesario para OTP, Pub/Sub y caché)
- (Opcional) Cuentas de **Brevo**, **Cloudinary** y **SMTP** para el envío de emails y almacenamiento de avatares

## ⚙️ Configuración

1. Clona el repositorio y entra al directorio:

   ```bash
   git clone https://github.com/jmlc643/twitbet-backend.git
   cd twitbet-backend
   ```

2. Crea tu archivo de entorno:

   ```bash
   cp .env.example .env
   ```

3. Completa las variables en `.env`:

   | Variable          | Descripción                                        | Ejemplo                                        |
   | ----------------- | -------------------------------------------------- | ---------------------------------------------- |
   | `PORT`            | Puerto del servidor HTTP                           | `8080`                                         |
   | `DATABASE_URL`    | URI de conexión a PostgreSQL                       | `postgres://usuario:password@host/bd_name`     |
   | `REDIS_URL`       | URI de conexión a Redis                            | `rediss://default:password@host:6379`          |
   | `JWT_SECRET`      | Secreto para firmar tokens JWT                     | `tu_clave_secreta_super_segura`                |
   | `FRONTEND_URL`    | Origen permitido para CORS                         | `https://frontend-desplegado.com`              |
   | `CLOUDINARY_URL`  | URL de Cloudinary (avatares)                       | `cloudinary://API_KEY:API_SECRET@CLOUD_NAME`   |
   | `SMTP_HOST`       | Host del servidor SMTP                             | `smtp.gmail.com`                               |
   | `SMTP_PORT`       | Puerto SMTP                                        | `587`                                          |
   | `SMTP_USER`       | Usuario SMTP                                       | `youremail@example.com`                        |
   | `SMTP_PASS`       | Contraseña / app password SMTP                     | `apppass`                                      |
   | `SMTP_SENDER`     | Remitente de los correos                           | `noreply@twitbet.com`                          |
   | `BREVO_API_KEY`   | API key de Brevo para emails transaccionales       | `your_api_key`                                 |

> **Nota:** Las migraciones (tablas) se ejecutan automáticamente al iniciar el servidor (`AutoMigrate`).

## Ejecución local

```bash
go run ./cmd/api
```

El servidor arranca en `http://localhost:8080` (por defecto).

### Healthcheck

```bash
curl http://localhost:8080/healthcheck
```

Devuelve el estado de los servicios (PostgreSQL y Redis).

### Con Docker

```bash
docker build -t twitbet-backend .
docker run -p 8080:8080 --env-file .env twitbet-backend
```

---

## API

Todas las rutas usan el prefijo `/api/v1`. Excepto las de autenticación, las rutas requieren el header `Authorization: Bearer <token>`.

### Autenticación

| Método | Ruta                          | Descripción                          |
| ------ | ----------------------------- | ------------------------------------ |
| POST   | `/api/v1/auth/register`       | Registrar un nuevo usuario           |
| POST   | `/api/v1/auth/login`          | Iniciar sesión                       |
| POST   | `/api/v1/auth/verify-account` | Verificar cuenta con OTP             |
| POST   | `/api/v1/auth/forgot-password`| Solicitar reset de contraseña (OTP)  |
| POST   | `/api/v1/auth/verify-reset-otp`| Validar OTP de reset                |
| POST   | `/api/v1/auth/reset-password` | Restablecer contraseña               |

### Usuarios (requieren JWT)

| Método | Ruta                            | Descripción                    |
| ------ | ------------------------------- | ------------------------------ |
| GET    | `/api/v1/users/me`              | Obtener perfil                 |
| PUT    | `/api/v1/users/me`              | Actualizar perfil              |
| POST   | `/api/v1/users/me/avatar`       | Subir avatar                   |
| POST   | `/api/v1/users/me/change-password` | Cambiar contraseña          |

### Ligas (requieren JWT)

| Método | Ruta                                 | Descripción                          |
| ------ | ------------------------------------ | ------------------------------------ |
| GET    | `/api/v1/leagues`                    | Ligas del usuario                    |
| POST   | `/api/v1/leagues`                    | Crear liga                           |
| POST   | `/api/v1/leagues/join`               | Unirse con código de invitación      |
| GET    | `/api/v1/leagues/:id`                | Detalles de liga                     |
| PUT    | `/api/v1/leagues/:id`                | Actualizar liga                      |
| PATCH  | `/api/v1/leagues/:id/status`         | Cambiar estado de liga               |
| DELETE | `/api/v1/leagues/:id`                | Eliminar liga                        |
| POST   | `/api/v1/leagues/:id/admins`         | Asignar admin                        |
| DELETE | `/api/v1/leagues/:id/admins/:participant_id` | Quitar admin                 |
| GET    | `/api/v1/leagues/:id/leaderboard`    | Clasificación de la liga             |
| GET    | `/api/v1/leagues/:id/me`             | Datos del participante en la liga    |
| POST   | `/api/v1/leagues/:id/recharge`       | Recargar saldo                       |
| POST   | `/api/v1/leagues/:id/bonuses`        | Otorgar bono                         |
| GET    | `/api/v1/leagues/:id/bonuses/me`     | Bonos pendientes del usuario         |
| GET    | `/api/v1/leagues/:id/bets`           | Apuestas del usuario en la liga      |

### Partidos y mercados (requieren JWT)

| Método | Ruta                                    | Descripción                          |
| ------ | --------------------------------------- | ------------------------------------ |
| POST   | `/api/v1/leagues/:id/matches`           | Crear partido en liga                |
| GET    | `/api/v1/leagues/:id/matches`           | Partidos de la liga                  |
| POST   | `/api/v1/leagues/:id/markets`           | Crear mercado en liga                |
| GET    | `/api/v1/leagues/:id/markets`           | Mercados de la liga                  |
| GET    | `/api/v1/matches/:id`                   | Detalles del partido                 |
| PATCH  | `/api/v1/matches/:id/status`            | Cambiar estado del partido           |
| POST   | `/api/v1/matches/:id/markets`           | Crear mercado del partido            |
| GET    | `/api/v1/matches/:id/markets`           | Mercados del partido                 |
| PATCH  | `/api/v1/markets/:id/status`            | Cambiar estado del mercado           |
| PATCH  | `/api/v1/markets/:id/odds`              | Actualizar cuotas del mercado        |
| POST   | `/api/v1/markets/:id/resolve`           | Resolver mercado (liquidar apuestas) |
| POST   | `/api/v1/markets/:id/cancel`            | Cancelar mercado (anular apuestas)   |

### Apuestas (requieren JWT)

| Método | Ruta                    | Descripción              |
| ------ | ----------------------- | ------------------------ |
| POST   | `/api/v1/bets`          | Colocar una apuesta      |
| POST   | `/api/v1/bets/:id/cashout` | Hacer cashout de la apuesta |

---

## WebSocket

Conecta a `ws://<host>/ws` con el token JWT en el query string (o según el flujo que maneje tu cliente). El servidor difunde eventos de mercados en tiempo real a través de **Redis Pub/Sub** (canal `market_events`).

---

## Documento de Requisitos

> **Pendiente de agregar por el equipo:**

## Diagrama de Base de Datos

> **Pendiente de agregar por el equipo:**

## Licencia

Este proyecto es de uso personal/privado. Consulta al propietario para más detalles.
