# yang-to-ntriples

Convert gNMI JSON telemetry data into N-Triples format using OpenConfig YANG models.

This project provides a Go-based tool to parse gNMI (gRPC Network Management Interface) JSON output, map it to OpenConfig structs using `ygot`, and generate N-Triples (a line-based RDF format) for network interface counters. It’s designed to process network telemetry data in a structured, semantic way, making it suitable for integration with RDF-based systems or knowledge graphs.

## Features
- Parses gNMI JSON responses with interface state counters.
- Maps data to OpenConfig interface models.
- Outputs N-Triples for fields like `in-octets` and `in-broadcast-pkts`.
- Modular and extensible design following Go best practices.

## Project Structure
```
.
├── cmd/
│   └── triples/          # Command-line application
│       └── main.go       # Entry point for the triples binary
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── oc/                   # Package for OpenConfig structs
│   └── oc.go             # Generated or custom OpenConfig code
└── yang/                 # YANG model files
    ├── ietf-inet-types.yang
    ├── ietf-interfaces.yang
    ├── ietf-yang-types.yang
    ├── openconfig-extensions.yang
    ├── openconfig-interfaces.yang
    ├── openconfig-platform-types.yang
    ├── openconfig-transport-types.yang
    ├── openconfig-types.yang
    └── openconfig-yang-types.yang
```


## Prerequisites
- **Go**: Version 1.13 or higher (for error wrapping with `%w`).
- **Dependencies**:
  - `github.com/openconfig/ygot` (for YANG struct generation and utilities).
  - YANG files in `yang/` (included in this repo).

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/gwoodwa1/yang-to-ntriples.git
   cd yang-to-ntriples
