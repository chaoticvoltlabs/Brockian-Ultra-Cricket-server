# First Steps

BUC can look bigger than it really is.

That is partly because it is meant to grow into something useful, and partly because dashboards, panels, rendering, and device-specific presentation all sound more complicated than they need to be at the start.

The good news is this:

**you do not need to understand all of BUC to get something useful on the screen.**

## What you actually do

At a very simple level, BUC works like this:

1. define a device
2. define a screen
3. define the components used on that screen
4. start the server
5. open the page

That is the whole basic idea.

## A very small example

A device says where and how something is shown.

```json
{
  "devices": {
    "office_browser": {
      "mode": "web",
      "screen": "weather_dashboard_main",
      "theme": "dark_default",
      "orientation": "landscape",
      "resolution": {
        "width": 1920,
        "height": 1080
      },
      "refresh_seconds": 300
    }
  }
}
```

A screen says which components appear and where.

```json
{
  "screens": {
    "weather_dashboard_main": {
      "layout": "dashboard_two_column_footer",
      "title": "Weather Dashboard Main",
      "regions": {
        "left_top": ["outside_summary_main"],
        "left_middle": ["wind_strip_main"],
        "right_top": ["windy_embed_wind"],
        "right_middle": ["windy_embed_rain"],
        "footer": ["daily_forecast_main"]
      }
    }
  }
}
```

A component says what kind of thing it is and where its data comes from.

```json
{
  "components": {
    "outside_summary_main": {
      "type": "outside_summary",
      "source": "weather_current",
      "options": {
        "show_precip": true,
        "show_pressure": true,
        "show_humidity": true,
        "show_wind_compass": true
      }
    }
  }
}
```

That is already enough to describe something meaningful.

## What this means in practice

You are not hand-placing every pixel.

You are saying things like:

- this device is a browser
- this screen is a weather dashboard
- this region should contain a summary
- this component should use current weather data

That is the point.

BUC tries to keep configuration semantic, not pixel-driven.

## Device profiles stay central

The main place to express presentation intent is still `config/devices.json`.

That is where you define the device profile a browser, wall panel, or other target will use. In practice, the easiest way to move quickly is to copy an existing device profile and only change the screen, theme, orientation, refresh cadence, and target resolution that actually matter.

The current repository already includes examples for desktop, mobile, and 800x480 Waveshare-oriented browser profiles. Those are meant to be starting points, not rigid presets.

For the active embedded panel stack, there is now also a second config path:

- `config/panel_devices.json`
  concrete embedded panel identities
- `config/panel_profiles.json`
  room/content profiles
- `config/panel_commands.json`
  control intent to HA call mappings

That means BUC currently has both:

- browser-oriented device profiles in `config/devices.json`
- panel-oriented runtime profiles in the `panel_*.json` files

## You do not need to start big

A good first use of BUC is not:

- a full building control system
- an embedded player fleet
- a perfect touchscreen UI
- a universal automation front-end

A good first use is something like:

- one browser page
- one useful screen
- one clear purpose

For example:

- a weather overview
- a climate overview
- an IoT status page
- a monitoring page for one system

If that works, the rest can grow from there.

## Think in screens, not in features

One of the easiest mistakes is to think:

I need a huge system.

Usually you do not.

Usually you need:

- one screen that answers one practical question
- then another one
- then maybe a simple way to switch between them

That is a much easier way to build something useful.

## A good mental model

BUC becomes easier to understand when you think of it like this:

- data sources know things
- components present things
- screens compose things
- devices decide where and how those screens live

That is the main model.

## Start with the browser

If you are new to BUC, start with a browser device.

That gives you:

- fast iteration
- easy debugging
- clear visual feedback
- less embedded-specific complexity

Embedded devices, touch panels, kiosk players, and more advanced rendering strategies can come later.

If you do work with panels, the same "start small" rule still applies:

- get one panel identity resolving correctly
- confirm `/api/panel/config` returns the expected profile
- confirm `/api/panel/weather` renders correctly
- confirm one control path works end to end

That is a much better first milestone than trying to design a complete fleet up front.

## Do not panic

You do not need to solve everything at once.

If you can define:

- one device
- one screen
- one or two components

and get that on screen, then you already understand the most important part.

Everything after that is refinement.

## What next

After this document, the most useful next reads are:

- README.md — for the project overview
- docs/concept.md — for the ideas behind BUC
- docs/architecture.md — for the current working server and panel runtime model
- docs/roadmap.md — for where the framework is heading

If you are working on panel integration specifically, also read:

- docs/quickstart.md
- local_docs/api.md

Later, the documentation will grow into more formal installation and operations guides.

For now, the important part is simply this:

BUC is allowed to start small.
