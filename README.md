# rebar

Rebar is the structural element for services at MasonHub (after we replace Buffalo everywhere)

### Major versions

- v0: first major version
- [v2](./v2): improve package structure and replace mux with gin
  - simpler subfolder and package structure
  - replaced Gorilla Mux with Gin
  - better Logger middleware with zap
  - improved graceful shutdown with cancelable context
