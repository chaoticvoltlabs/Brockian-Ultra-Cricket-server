# BUC Architecture

## Purpose

BUC sits between Home Assistant and renderer/control clients.

Its current responsibilities are:

- read normalized entities from Home Assistant
- adapt those entities into stable source models
- expose compact API responses for panel clients
- resolve config-driven panel control intents into HA service calls
- support broader screen/device composition work for browser and renderer clients

## Current end-to-end flow

### Read path

1. Home Assistant defines canonical entities such as:
   - `sensor.panel_weather_current`
   - `sensor.panel_weather_48h`
   - `sensor.panel_weather_daily`
   - `sensor.panel_indoor_payload`
   - `sensor.panel_overview_payload`
2. BUC reads those entities through the HA API client.
3. `internal/ha` normalizes HA state and attributes into source data.
4. `internal/httpapi` exposes:
   - source inspection endpoints
   - screen/device model endpoints
   - compact panel weather output
5. The panel firmware polls `GET /api/panel/weather`.

### Control path

1. The panel sends a compact command to `POST /api/panel/control`.
2. BUC validates the command against an explicit whitelist.
3. BUC resolves the command through the configured command map.
4. Home Assistant executes the switch toggle or scene activation.
5. BUC returns `ok: true/false`.

### Panel config path

1. The panel sends `GET /api/panel/config`.
2. BUC reads panel identity from:
   - `X-Panel-MAC`
   - `X-Panel-IP`
   - optional debug query fallback `panel_mac` / `panel_ip`
3. BUC resolves the request through three config layers:
   - `identity` -> concrete panel device record
   - `device_type` -> client family such as `panel_4_3`
   - `profile` -> room-specific content such as `room_alpha` or `room_beta`
4. BUC returns the page/profile JSON that the panel uses to label and wire page 3.

## Key modules

### Server boot

- [`cmd/buc-server/main.go`](../cmd/buc-server/main.go)

Loads config, constructs the HA client, builds the app container, and starts the router.

### Config

- [`config/sources.json`](../config/sources.json)
- [`config/screens.json`](../config/screens.json)
- [`config/devices.json`](../config/devices.json)
- [`config/device_types.json`](../config/device_types.json)
- [`config/panel_devices.json`](../config/panel_devices.json)
- [`config/panel_profiles.json`](../config/panel_profiles.json)
- [`config/panel_commands.json`](../config/panel_commands.json)
- [`config/components.json`](../config/components.json)
- [`config/themes.json`](../config/themes.json)

The split is now:

- `devices.json`
  browser and screen-composition device definitions
- `device_types.json`
  client-family definitions such as `panel_4_3`
- `panel_devices.json`
  concrete embedded panel instances, keyed by identity such as MAC address
- `panel_profiles.json`
  page-level room/content profiles returned by `/api/panel/config`
- `panel_commands.json`
  explicit `target:action -> HA service call` mappings used by `/api/panel/control`

### HA integration

- [`internal/ha/client.go`](../internal/ha/client.go)
- [`internal/ha/source_weather_current.go`](../internal/ha/source_weather_current.go)
- [`internal/ha/source_indoor_payload.go`](../internal/ha/source_indoor_payload.go)
- [`internal/ha/source_overview_payload.go`](../internal/ha/source_overview_payload.go)

### HTTP API

- [`internal/httpapi/router.go`](../internal/httpapi/router.go)
- [`internal/httpapi/handlers_panel.go`](../internal/httpapi/handlers_panel.go)
- [`internal/httpapi/handlers_panel_control.go`](../internal/httpapi/handlers_panel_control.go)
- [`internal/httpapi/handlers_panel_config.go`](../internal/httpapi/handlers_panel_config.go)
- [`internal/httpapi/handlers_log.go`](../internal/httpapi/handlers_log.go)
- [`internal/httpapi/panel_debug.go`](../internal/httpapi/panel_debug.go)

### Logging

- [`cmd/buc-server/main.go`](../cmd/buc-server/main.go)
- [`internal/logview/manager.go`](../internal/logview/manager.go)

BUC now has one central log output path.

The same human-readable log lines are written to:

- stdout or systemd journal
- hourly logfile rotation under `BUC_LOG_DIR`
- the live browser stream used by `/log/live`

## Panel contracts

### Weather contract

The panel currently depends on:

- `GET /api/panel/weather`

Current response fields include:

- `outside_temp_c`
- `feels_like_c`
- `wind_bft`
- `wind_kmh`
- `gust_bft`
- `gust_kmh`
- `wind_dir_deg`
- `humidity_pct`
- `pressure_hpa`
- `pressure_trend_24h`
- `indoor_zones`

### Control contract

The panel currently depends on:

- `POST /api/panel/control`

Request shape:

```json
{
  "target": "string",
  "action": "string"
}
```

Response shape:

```json
{
  "ok": true,
  "target": "light_a",
  "action": "toggle"
}
```

### Config contract

The panel now also depends on:

- `GET /api/panel/config`

Current response shape:

```json
{
  "profile": "room_alpha",
  "page3": {
    "scenes": [
      {"label": "Work", "target": "scene_work", "action": "activate"}
    ],
    "targets": [
      {"label": "Light A", "target": "light_a", "action": "toggle"}
    ],
    "long_press": {
      "label": "Night",
      "target": "scene_night",
      "action": "activate"
    }
  }
}
```

### Log viewing contract

BUC now also exposes a small log viewing surface intended for operators and developers:

- `GET /log/live`
- `GET /api/log/stream`
- `GET /log/files`
- `GET /log/files/{name}`

The intended usage is:

- use `/log/live` while testing panels, HA integrations, or control flows
- use `/log/files` to inspect recent hourly logs without logging into the server shell

## Current example command set

The current command whitelist is config-driven through:

- [`config/panel_commands.json`](../config/panel_commands.json)

Example direct control commands:

- `light_a:toggle`
- `light_b:toggle`
- `media_power:toggle`

Example scene commands:

- `scene_work:activate`
- `scene_evening:activate`
- `scene_movie:activate`
- `scene_night:activate`

## Operational notes

- `GET /api/panel/weather` must remain read-only.
- `POST /api/panel/control` is the only current panel writeback path.
- `GET /api/panel/config` is the current panel profile path.
- `BUC_LOG_DIR` controls where BUC stores its hourly logfile set.
- BUC is intentionally explicit today:
  - whitelisted commands
  - config-driven panel identity, profile, and command resolution
  - no generic smart-home protocol layer yet

## Live log usage

The live log feature is intentionally simple:

- `/log/live`
  small browser page that shows the current server log in real time
- `/log/files`
  browser index of recent hourly logfile segments
- `/log/files/{name}`
  plain text view of one logfile

Typical deployment notes:

- if BUC runs behind Apache or another reverse proxy, `/log/` must be proxied in addition to `/api/`
- the BUC service user must be able to create and write the configured `BUC_LOG_DIR`

This feature is meant for practical local debugging, not for hardened public log exposure.

## Panel resolve logging

BUC writes these lines to the server process log when `/api/panel/config` is called.

Depending on how BUC is run, you will typically see them:

- in `journalctl` when BUC runs as a systemd service
- in stdout/stderr when BUC runs directly in a terminal

BUC currently logs two different things for panel profile requests.

Generic request log:

```text
panel request method=GET path=/api/panel/config panel_mac=aa:bb:cc:dd:ee:ff panel_ip=192.168.1.50 remote_ip=127.0.0.1
```

This means:

- the request reached BUC
- which panel identity headers or query fallback were seen
- which remote IP opened the HTTP connection

Profile resolution log written by the BUC server:

```text
panel device resolved panel_mac=aa:bb:cc:dd:ee:ff device=panel_alpha device_type=panel_4_3 profile=room_alpha
```

This means:

- `panel_mac`
  the concrete incoming device identity
- `device`
  the matched record name in `panel_devices.json`
- `device_type`
  the client family resolved from that device record
- `profile`
  the content profile resolved from that device record and returned to the panel

In practice, this tells you not just that a panel called BUC, but exactly which configured device record matched and which profile BUC selected.
