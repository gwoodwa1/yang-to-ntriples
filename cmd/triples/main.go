package main

import (
    "encoding/json"
    "fmt"
    "strings"

    "github.com/openconfig/ygot/ygot"
    oc "github.com/gwoodwa1/yang-to-ntriples/oc"
)

// GnmiUpdate represents a single update in a gNMI JSON response.
type GnmiUpdate struct {
    Path   string                 `json:"Path"`
    Values map[string]interface{} `json:"values"`
}

// GnmiResponse represents a top-level gNMI JSON response.
type GnmiResponse struct {
    Source  string       `json:"source"`
    Time    string       `json:"time"`
    Updates []GnmiUpdate `json:"updates"`
}

const (
    baseURIFormat = "http://example.net/interfaces/%s"
    rdfType       = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"
    rdfInterface  = "http://openconfig.net/rdf/Interface"
    rdfInOctets   = "http://openconfig.net/rdf/inOctets"
    rdfInBcast    = "http://openconfig.net/rdf/inBroadcastPkts"
)

func main() {
    gnmiJSON := `
[
  {
    "source": "192.168.151.7:6030",
    "time": "1970-01-01T01:00:00+01:00",
    "updates": [
      {
        "Path": "interfaces/interface[name=Ethernet8]/state/counters",
        "values": {
          "interfaces/interface/state/counters": {
            "openconfig-interfaces:in-broadcast-pkts": "2367884",
            "openconfig-interfaces:in-discards": "0",
            "openconfig-interfaces:in-errors": "0",
            "openconfig-interfaces:in-fcs-errors": "0",
            "openconfig-interfaces:in-multicast-pkts": "178101",
            "openconfig-interfaces:in-octets": "25833637",
            "openconfig-interfaces:in-unicast-pkts": "8182",
            "openconfig-interfaces:out-broadcast-pkts": "2",
            "openconfig-interfaces:out-discards": "0",
            "openconfig-interfaces:out-errors": "0",
            "openconfig-interfaces:out-multicast-pkts": "152504",
            "openconfig-interfaces:out-octets": "2451633129",
            "openconfig-interfaces:out-unicast-pkts": "8177"
          }
        }
      }
    ]
  }
]`

    responses, err := parseGnmiJSON(gnmiJSON)
    if err != nil {
        fmt.Printf("Failed to parse gNMI JSON: %v\n", err)
        return
    }

    for _, resp := range responses {
        if err := processResponse(resp); err != nil {
            fmt.Printf("Failed to process response from %s: %v\n", resp.Source, err)
        }
    }
}

// parseGnmiJSON unmarshals a gNMI JSON string into a slice of GnmiResponse.
func parseGnmiJSON(jsonStr string) ([]GnmiResponse, error) {
    var responses []GnmiResponse
    if err := json.Unmarshal([]byte(jsonStr), &responses); err != nil {
        return nil, fmt.Errorf("unmarshaling JSON: %w", err)
    }
    return responses, nil
}

// processResponse processes a single gNMI response and prints its N-Triples.
func processResponse(resp GnmiResponse) error {
    for _, update := range resp.Updates {
        if !strings.Contains(update.Path, "/state/counters") {
            continue
        }

        counters, err := extractCounters(update)
        if err != nil {
            return fmt.Errorf("extracting counters: %w", err)
        }

        ifaceName := parseInterfaceName(update.Path)
        if ifaceName == "" {
            return fmt.Errorf("invalid interface name in path: %s", update.Path)
        }

        iface := newInterface(ifaceName, counters)
        triples, err := toNTriples(iface)
        if err != nil {
            return fmt.Errorf("converting to N-Triples: %w", err)
        }

        for _, triple := range triples {
            fmt.Println(triple)
        }
    }
    return nil
}

// extractCounters extracts and unmarshals the counters from a gNMI update.
func extractCounters(update GnmiUpdate) (*oc.OpenconfigInterfaces_Interfaces_Interface_State_Counters, error) {
    val, ok := update.Values["interfaces/interface/state/counters"]
    if !ok {
        return nil, fmt.Errorf("no counters found in update")
    }

    subJSON, err := json.Marshal(val)
    if err != nil {
        return nil, fmt.Errorf("marshaling counters: %w", err)
    }

    var counters oc.OpenconfigInterfaces_Interfaces_Interface_State_Counters
    if err := oc.Unmarshal(subJSON, &counters); err != nil {
        return nil, fmt.Errorf("unmarshaling counters: %w", err)
    }
    return &counters, nil
}

// parseInterfaceName extracts the interface name from a gNMI path.
func parseInterfaceName(path string) string {
    const namePrefix = "[name="
    start := strings.Index(path, namePrefix)
    if start == -1 {
        return ""
    }
    start += len(namePrefix)
    end := strings.Index(path[start:], "]")
    if end == -1 {
        return ""
    }
    return path[start : start+end]
}

// newInterface creates a new OpenConfig interface with the given name and counters.
func newInterface(name string, counters *oc.OpenconfigInterfaces_Interfaces_Interface_State_Counters) *oc.OpenconfigInterfaces_Interfaces_Interface {
    return &oc.OpenconfigInterfaces_Interfaces_Interface{
        Name: ygot.String(name),
        State: &oc.OpenconfigInterfaces_Interfaces_Interface_State{
            Counters: counters,
        },
    }
}

// toNTriples converts an OpenConfig interface to N-Triples representation.
func toNTriples(iface *oc.OpenconfigInterfaces_Interfaces_Interface) ([]string, error) {
    if iface.Name == nil {
        return nil, fmt.Errorf("interface name is nil")
    }

    baseURI := fmt.Sprintf(baseURIFormat, *iface.Name)
    triples := []string{
        fmt.Sprintf("<%s> <%s> <%s> .", baseURI, rdfType, rdfInterface),
    }

    if iface.State != nil && iface.State.Counters != nil {
        counters := iface.State.Counters
        if counters.InOctets != nil {
            triples = append(triples, fmt.Sprintf("<%s> <%s> \"%d\"^^<http://www.w3.org/2001/XMLSchema#integer> .", baseURI, rdfInOctets, *counters.InOctets))
        }
        if counters.InBroadcastPkts != nil {
            triples = append(triples, fmt.Sprintf("<%s> <%s> \"%d\"^^<http://www.w3.org/2001/XMLSchema#integer> .", baseURI, rdfInBcast, *counters.InBroadcastPkts))
        }
        // Add more counters as needed...
    }

    return triples, nil
}
