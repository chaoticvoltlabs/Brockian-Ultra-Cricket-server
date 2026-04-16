# Quickstart

This quickstart is intentionally simple.

It is not meant to explain every part of BUC in detail.  
It is meant to help you get from “what is this?” to “I have something running.”

## What you need

At a minimum, you need:

- a system that can run the BUC server
- a browser to view the UI
- one or more data sources that BUC can present
- a small amount of configuration

Right now, BUC is most comfortable in a browser-first workflow.

That means the easiest starting point is:
- run the server
- define one browser device
- define one screen
- open that screen in a browser

## Basic flow

The practical flow looks like this:

1. configure your data sources
2. define components
3. define a screen
4. define a device that uses that screen
5. start the BUC server
6. open the device UI in a browser

## Start small

The best first test is not a complete installation.

The best first test is something like:
- one weather page
- one climate page
- one status page

Pick one useful screen and get that working first.

## Suggested first target

If you are new to BUC, start with:
- a browser device
- one screen
- one or two components

This keeps the number of moving parts low and makes it easier to understand what the framework is doing.

## Configuration model

BUC is configured around a few core ideas:

- **sources** provide data
- **components** describe reusable presentation units
- **screens** compose components into layouts
- **devices** define where and how a screen is shown
- **themes** define visual language

You do not need to master all of this at once.

You only need enough to describe one useful screen.

## Device profiles

The main entry point for presentation setup is `config/devices.json`.

Start by copying a device profile that is already close to your target, then change the semantic parts first:

- `screen`
- `theme`
- `orientation`
- `refresh_seconds`

The repository already contains browser-style profiles for desktop, mobile, and 800x480 Waveshare panel experiments, which makes it easier to start from a known-good layout instead of inventing one from scratch.

## Panel data note

If you are building a compact weather panel, BUC also exposes `/api/panel/weather`.

That endpoint returns a reduced weather payload intended for panel use, so panel-specific frontends do not need to assemble those values client-side.

Depending on the upstream entities and configured panel commands, the same endpoint can also expose:

- `night_mode`
- `page3_target_states`
- `ambient_brightness_pct`
- `ambient_rgb`

Current panel-oriented config is split into:

- `config/panel_devices.json`
  concrete embedded panel identities such as MAC-address based records
- `config/panel_profiles.json`
  room/content profiles such as `room_alpha` or `room_beta`
- `config/panel_commands.json`
  explicit panel command to HA service-call mappings

## Live debugging note

When BUC is running, you can also inspect logs directly in a browser:

- `/log/live`
  real-time server log view
- `/log/files`
  recent hourly logfile index

This is useful while testing panel button presses, Home Assistant entity changes, and panel profile resolution without opening a second shell session to the server.

## What this section will grow into

Later, the installation docs will expand into more detailed guides such as:

- server setup
- browser device setup
- theme setup
- component configuration
- screen composition
- multi-screen devices
- embedded player setup

For now, this document exists to make one thing clear:

**BUC can start small, and it should.**

## Deployment note

The simplest and safest BUC deployment model is a trusted local network.

If remote access is required, BUC should be placed behind an appropriate access control layer such as a VPN or an authenticated reverse proxy. Direct public exposure without additional protection is not recommended.

If BUC runs behind Apache or another reverse proxy, remember that `/log/` should be proxied alongside `/api/` if you want the browser log viewer to work.
