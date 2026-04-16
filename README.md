# Brockian Ultra-Cricket (BUC)

**A presentation and rendering framework for Home Assistant and beyond — for people who need more than a Lovelace dashboard.**

## Why this exists

At some point, many Home Assistant setups hit the same wall:

- the data is there
- the automations work
- the integrations are fine
- but the UI is still not quite what you want

Sometimes you want a cleaner weather screen.  
Sometimes you want a wall panel that behaves like an instrument, not a collection of cards.  
Sometimes you want the same information rendered for a browser, an embedded display, or something more specialized.

That is where **BUC** comes in.

BUC is a presentation layer: it sits on top of your data sources and focuses on **how information is shaped, composed, and rendered**.

## What it is for

BUC is meant for people who want to build things like:

- better weather dashboards
- cleaner control room style screens
- energy and system overviews
- touch panels
- embedded display UIs
- dedicated interfaces for specific use cases

It is especially useful when you want more control over layout, rendering, and visual behavior than a standard dashboard normally gives you.

<img src="docs/BUC-1024.jpg" width="320" alt="hardware mockup">

## Why it may be useful to you

BUC is built around a simple idea:

**data collection and UI rendering are different jobs**

If your system already knows things, BUC helps present them in a way that is:

- more deliberate
- more readable
- more device-aware
- more visually consistent
- easier to evolve into something purpose-built

## What makes it different

BUC is not trying to be yet another theme or card pack.

It is about treating UI as its own layer:

- composable
- renderer-aware
- device-aware
- suitable for both general dashboards and highly specific interfaces

If you have ever thought:

> “Home Assistant knows the data, but I want the presentation to be mine.”

then BUC may be useful.

## Current direction

BUC started in the weather-and-panel space, but it is intended to grow beyond that.

The long-term direction is broader:
- weather
- energy
- control systems
- instrumentation-style displays
- alternative frontends for complex installations

In other words:

BUC is for the gloriously convoluted world where data exists, screens matter, and the rules of the game are not always obvious.

## Current working stack

BUC is no longer only a presentation concept.

The current stack already proves a working end-to-end flow where Home Assistant provides normalized data, BUC exposes compact panel-facing APIs, and an embedded touch panel both renders that data locally and sends back explicit control intents.

That working loop now covers:

- read-only weather and status delivery for a dedicated panel UI
- separate writeback control endpoints
- direct device control
- scene activation from the panel

The implementation is still evolving, but the core architecture is already practical and usable.

## Status

BUC is under active development.

It already proves the core idea:
a separate presentation and rendering layer can produce cleaner, more intentional interfaces than a default dashboard flow.

Over time, BUC is also expected to support event-driven temporary screen takeovers, such as doorbell camera views or other high-priority operational interrupts.

The details will evolve.  
That is part of the point.

## Current config shape

The main presentation entry point remains [`config/devices.json`](config/devices.json).

That file defines the semantic device profiles that map a presentation target to:

- a screen
- a theme
- an orientation
- a resolution
- a refresh cadence

The repository currently includes browser-oriented profiles for larger desktop views, mobile views, and multiple 800x480 Waveshare panel variants. The intended workflow is to copy an existing device profile, keep the config semantic, and then point it at the screen composition you want.

## Panel-oriented weather API

For compact panel layouts, BUC also exposes `/api/panel/weather`.

That endpoint is designed to provide a concise weather payload for panel rendering, including:

- outside temperature
- feels-like temperature
- wind and gust values
- humidity
- pressure
- a 24-hour pressure trend when available
- an optional `night_mode` flag
- optional live page-3 target states
- optional ambient light state such as brightness and RGB

The handler reads from the configured `weather_current` source and, when present, enriches the result with trend data from `overview_payload`.

## Related repositories

BUC is intended to work as part of a broader stack.

- [Brockian-Ultra-Cricket-panel](https://github.com/chaoticvoltlabs/Brockian-Ultra-Cricket-panel)
  - embedded panel firmware client
- [Brockian-Ultra-Cricket-homeassistant](https://github.com/chaoticvoltlabs/Brockian-Ultra-Cricket-homeassistant)
  - Home Assistant package and configuration layer
- additional technical notes in the `docs/` folder

## In one sentence

**BUC exists to give serious tinkerers, builders, and system integrators more control over how their systems are actually seen.**

** Final thought and Warning **
BUC is currently not designed as a hardened public internet frontend.
It does not provide built-in authentication, authorization, or secure remote access features.
If you use BUC outside a trusted local environment, place it behind appropriate access controls such as a VPN, reverse proxy, authentication layer, or other protective infrastructure.

## Copyright & license

PolyForm Noncommercial License 1.0.0
with Commercial Use by Explicit Permission Only
See LICENSE.txt


Copyright (c) 2026 Robin Kluit / Chaoticvolt.
