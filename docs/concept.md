# BUC Concept

## What BUC is

Brockian Ultra-Cricket (BUC) is a presentation and rendering framework for dashboards, panels, and control-oriented interfaces.

It is designed for situations where the data source already exists, but the way that information is presented needs to be more deliberate, more readable, and more device-aware than a standard dashboard usually provides.

BUC is not centered around one single UI.  
It is centered around a model for describing screens, devices, components, themes, and navigation in a way that stays understandable.

## What BUC is not

BUC is not:

- a Home Assistant replacement
- a data collection platform
- an automation engine
- a card pack
- a theme pack
- a one-off dashboard for a single installation

## Security boundary

BUC should be understood as a presentation layer, not as a security boundary.

It is designed to render and organize information, not to act as a hardened public-facing access gateway. Authentication, remote access protection, and exposure to untrusted networks should be handled by the surrounding infrastructure.
It may be used with Home Assistant, but it is intended to sit above the data layer, not replace it.

## Why BUC exists

Many systems already know a lot.

They know:
- weather
- temperatures
- humidity
- energy state
- device state
- alarms
- control values

But knowing something and presenting it well are not the same thing.

BUC exists because there is often a gap between:
- what the system knows
- and what a human can comfortably understand at a glance

That gap becomes even more obvious when the same information needs to work across:
- desktop browsers
- mobile browsers
- wall panels
- embedded displays
- operator-style screens

## Core idea

BUC is built on a simple separation:

**data belongs to the source layer, presentation belongs to the UI layer**

That means BUC focuses on:
- shaping information
- grouping it meaningfully
- rendering it appropriately for the device and use case
- keeping that model understandable in configuration

In practice, this is now also visible in the control path:

- device identity is resolved separately from device behavior
- panel profiles are resolved from config at runtime
- panel control intents are resolved from config at runtime
- the server code acts as the motor, while installation-specific behavior is moving into config

## Core concepts

### Devices

A device is the presentation target.

A device describes things like:
- browser vs player vs embedded target
- orientation
- resolution
- refresh behavior
- theme
- screen set

The device is not just a screen size.  
It is a presentation context.

In the active panel stack, it helps to distinguish three different things:

- identity
  the concrete device instance, such as a specific wall panel
- device type
  the client family, such as `panel_4_3`, `panel_7`, or a future mobile app
- profile
  the room or use-case role, such as a compact room panel, a status-only panel, or a larger operator view

That separation keeps device lookup, layout family, and content role from collapsing into one hardcoded concept.

### Screens

A screen is a composed UI view.

A screen groups components into regions and layouts.  
Screens should be defined around human use, not around raw data structure.

A good screen answers a practical question such as:
- what is the weather doing?
- how is the building climate doing?
- what is the current system status?
- what can I safely control here?

### Components

Components are reusable presentation units.

Examples:
- outside summary
- wind strip
- climate overview
- status grid
- map panel
- control section

A component should express intent, not low-level layout tricks.

### Themes

Themes define visual language.

They should provide consistency across:
- large displays
- mobile screens
- browser UIs
- embedded panels

Themes exist so that devices and screens can vary without becoming visually unrelated.

### Navigation

Navigation is a first-class concept in BUC.

A useful presentation framework should not assume:
- one device = one screen forever

Some use cases need:
- manual selection
- auto-rotation
- passive information cycling
- interaction-aware pause/resume behavior

This matters especially when the same framework is used for:
- browsers
- kiosks
- touch panels
- embedded players

## Passive and interactive screens

Not all screens behave the same way.

### Passive screens
Passive screens are for observation:
- weather
- climate
- status
- energy summaries
- instrumentation-style overviews

These benefit from:
- calm refresh behavior
- long-running unattended use
- optional automatic rotation

### Interactive screens
Interactive screens are for control:
- switching
- dimming
- setpoints
- overrides
- direct operational changes

These need different assumptions:
- immediate feedback
- explicit interaction
- no unwanted automatic screen changes while in use
- For constrained embedded panels, API payloads need to be optimized for direct rendering, minimum transport cost, and predictable parsing. 
  Generality is secondary to responsiveness and robustness.

BUC should support both, but it should not confuse one for the other.

The current panel work already proves this split in practice:

- `GET /api/panel/weather`
  read-only compact state payload
- `GET /api/panel/config`
  runtime panel profile payload
- `POST /api/panel/control`
  explicit control intent path

## Browser-first, embedded-aware

BUC currently develops browser-first.

That is a deliberate choice:
- iteration is faster
- testing is easier
- rendering is easier to inspect
- concepts can be proven before embedded constraints dominate the design

Embedded support matters, but it should follow proven concepts rather than define them too early.

At the same time, the current stack is no longer only a browser thought experiment.

Real embedded panels are now part of the active design loop:

- multiple 4.3" panels resolve different runtime profiles
- panel controls and scenes work end to end through Home Assistant
- paging and compact payload design have been validated in practice
- server-side live logging makes panel behavior visible without shell access

## Example application areas

BUC is intended to be useful for more than one type of dashboard.

Examples include:
- weather dashboards
- building climate overviews
- IoT status boards
- control room style screens
- energy dashboards
- equipment monitoring
- pool plant monitoring
- greenhouse environments
- apiary or agricultural monitoring
- industrial or experimental systems

In short:
if a system has state, context, and screens that matter, BUC may be relevant.

## Design principles

BUC aims to follow a few simple principles:

- config should be semantic, not pixel-driven
- screens should be grouped by human use, not raw metric structure
- devices should define presentation context cleanly
- rendering should be explicit and controllable
- runtime configuration should be preferred over rebuild-driven behavior where possible
- complexity should only be introduced where it clearly pays off
- real testing should drive design changes

## Development philosophy

BUC is being built through real use.

That means:
- ideas are tested in practice
- browser screens are allowed to mature before embedded optimization
- once something is proven, it should move from hardcoded behavior into config
- the framework should stay understandable even while it grows
- documentation should explain the model, not bury users in implementation detail too early

BUC is not trying to hide complexity that exists in the real world.

It is trying to make that complexity presentable.
