# Web Assembly

[Web Assembly](https://webassembly.org) (WASM) is a new standard runtime for the internet.

## Overview

We want to integrate wasm as the final runtime for Micro to complete local, kubernetes, wasm. WebAssembly 
is a new standard runtime which has major investment from all major browsers and tech companies with the 
goal of achieving the dream of write once, run anywhere. Intel has produced a lightweight runtime 
called the [Intel Micro WASM Runtime](https://github.com/bytecodealliance/wasm-micro-runtime) which 
we want to use as the basis for our wasm runtime. The idea there is to embed the runtime directly into 
micro so that we can run WebAssembly binaries in process, anywhere.

## Goals

By adopting wasm we can gain true ubiquity at scale by enabling Micro and Micro services to be run 
absolutely anywhere and everywhere. To turn everything into a service. A Micro service.
