# Sync

Sync is a synchronization mechanism for data storage.

## Overview

We need to be able to sync between different Store types and locations. Often we describe 
this as local, regional, global or cloud, edge, dev. Sync provides a way to quite literally 
sync data between different stores and provides a Key-Value abstraction with built in 
data encoding for efficiency and timestamp values.

What we're fundamentally fighting is replication versus a layered distributed architecture. 
Replication is a flat design which works at one layer but multi-layer is ultimately 
the model for a large scale distributed system.

## Design

Ideally we operate like a computer. Cache misses walk the chain, writes as well.

A computer model

- CPU Register, L1, L2, L3, Ram, Disk

Our model

- local, region, global
- memory, cache, database, blob store

Or more concretely

- memory, etcd, cockroach, s3/blob/github
- service, cache, store, blob

## Architecture

Ultimately what we want is to replicate data without the need for data replication, where 
every cache miss results in recursively walking the chain. We find that federated models 
are far superior to replication alone. Again replication operates at a single layer 
and federation layers on top of it.

Walking through a real example. Where we're using the micro runtime we have reached 
limitations in terms of APi rates for cloudflare and global data storage is 
something we've found is expensive or unsupported by other services without 
using vpn or wireguard to support replication.

By using a federated model built entirely in micro, we can allow each layer to 
operate with their respective abstraction or Store and layer on top simply building 
the primitive for synchronisation.

Order of retrieval and storage

1. Local (memory)
2. Cache (etcd)
3. Database (cockroach)
4. Blob (github)

## Source of Truth

GitHub is and will always be our source of truth, for code, for configuration, for packages and now 
potentially for blob storage. By creating a central point that is not a server but in some ways 
cold storage we have a place for long term storage all things.

What we store in GitHub

- Code
- Configuration
- Binaries
- Packages
- Events
- Blobs
- Files

We can optionally load the entirety of our source of truth into DigitalOcean for higher throughput 
at very low cost and may choose to provide APIs as a central point through DO or elsewhere.

We have attempted to use CloudFlare as a distributed source of truth but without fully immersing 
ourselves in workers this will not work. In fact workers pushes us more down the path of a complete 
micro runtime on the edge using wasm (2022).

## References

- https://en.wikipedia.org/wiki/Microsoft_Sync_Framework
- GitHub Large File Storage https://help.github.com/en/github/managing-large-files/versioning-large-files
- https://www.digitalocean.com/products/block-storage/
- BigCache https://github.com/allegro/bigcache
- GroupCache https://github.com/golang/groupcache

