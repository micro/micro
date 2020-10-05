## Roadmap

This is a high level overview for the open source. See [m3o/dev](https://github.com/m3o/dev) for further info.

- [X] consolidate all libraries into go-micro
- [X] gRPC API for the micro proxy
- [X] quic as a default transport
- [X] nats as a default broker
- [x] [gRPC](design/framework/grpc.md) integration and interop
- [x] embedded nats as the default broker
- [x] go, java, typescript, ruby, python [clients](https://github.com/micro/clients)
- [x] standalone micro server
- [x] better documentation and end to end tutorials
- [ ] simpler kubernetes integration by default
- [ ] pluggable wrappers for go-micro itself
- [ ] dynamically defined cli, flags and env vars
- [ ] improved micro api configuration
- [ ] graphql handler for the micro api
- [ ] wireguard support in the micro network
- [ ] reusable foundation [services](https://github.com/micro/services)
- [ ] define the mucp protocol
- [ ] define the mu language spec
- [ ] implement the [wasm](https://github.com/bytecodealliance/wasm-micro-runtime) runtime
