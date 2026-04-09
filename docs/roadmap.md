# BUC Roadmap

## Milestone 1 — Browser rendering foundation
- [x] Public repository created and cleaned up
- [x] Renderer-driven browser UI working
- [x] Semantic device profiles in config
- [x] Desktop browser support
- [x] Mobile browser support
- [x] Polling-based live updates
- [x] Refresh status indicator

## Milestone 2 — Multi-screen browser experiences
- [X] Support multiple screens per device
- [X] Add a device-level screen set / playlist model
- [ ] Add subtle browser navigation
- [ ] Support optional automatic screen rotation
- [ ] Add play/pause behavior for presentation mode

## Milestone 3 — Passive operational dashboards
- [ ] Support grouped overview screens for sensors and status
- [ ] Support browser-friendly climate overviews
- [ ] Support read-only IoT/device status overviews
- [ ] Prefer grouping by human use context rather than raw metric type

## Milestone 4 — Interactive control surfaces
- [x] Introduce control-oriented screen types
- [x] Support direct user interaction
- [ ] Separate passive dashboard assumptions from control-screen architecture
- [x] Support interaction-aware navigation behavior
- [x] Support panel scene activation through Home Assistant
- [x] Support config-driven panel control mappings

## Milestone 5 — Embedded player support
- [ ] Re-evaluate architecture for embedded players after browser V2 is proven
- [ ] Define player-oriented rendering strategy
- [ ] Decide client-side vs server-side rendering boundaries
- [ ] Support kiosk/touch display behavior

## Milestone 6 — Scaled deployment patterns
- [x] Support multi-device environments cleanly
- [x] Support config-driven panel identity resolution
- [x] Support config-driven panel profile resolution
- [x] Support config-driven panel command resolution
- [ ] Revisit proxy/cache strategy for shared upstream data sources
- [ ] Reduce unnecessary upstream provider load
- [ ] Improve deployability for real-world installations
- [ ] Add explicit config reload behavior without service restart

## Documentation milestones
- [ ] Keep public docs concise and useful
- [x] Reintroduce private technical docs for the active stack
- [x] Document browser-based live log and logfile access
- [ ] Publish installation and maintenance guidance when the framework has stabilized

## Guiding principles
- [ ] Keep config semantic, not pixel-based
- [ ] Group screens around human use, not raw data structure
- [ ] Prefer simple solutions unless complexity clearly pays off
- [ ] Let real-world testing drive the next design round
