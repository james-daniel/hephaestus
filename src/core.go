package main

import (
  "fmt"
  "flag"
  "errors"
  "strings"
  "bytes"
)

// Parses the human's plaintext input.
// It also validates whether plaintext is valid.
func parseStdIn(humArgPtr *string)(int, error) {
  splitArr := strings.Split(*humArgPtr, " ")

  acceptedVerbs := map[string]bool {
    "permit": true,
    "deny": true,
  }

  acceptedFlow :=  map[string]bool {
    "inbound": true,
    "outbound": true,
  }

  acceptedProto := map[string]bool {
    "ssh": true,
    "icmp": true,
    "http": true,
    "https": true,
  }

  if acceptedVerbs[strings.ToLower(splitArr[0])] {
    fmt.Println("[DEBUG] verb validity check: PASSED")
  } else {
    fmt.Println("[DEBUG] verb validity check: FAILED")
    return -1, errors.New("statement contains invalid verb")
  }

  if acceptedFlow[strings.ToLower(splitArr[1])] {
    fmt.Println("[DEBUG] flow validity check: PASSED")
  } else {
    fmt.Println("[DEBUG] flow validity check: FAILED")
    return -1, errors.New("statement contains invalid flow")
  }

  if acceptedProto[strings.ToLower(splitArr[2])] {
    fmt.Println("[DEBUG] protocol validity check: PASSED")
  } else {
    fmt.Println("[DEBUG] protocol validity check: FAILED")
    return -1, errors.New("statement contains invalid protocol")
  }

  fmt.Println("[DEBUG] isvalidstatement: PASSED\n-------------------------------------")
  return 0, nil
}

// Build the actual iptables-compatible ruleset.
// At this point, the human input should have been validated.
func rulesetStringBuilder(humArgPtr *string)(string){
  splitArr := strings.Split(*humArgPtr, " ")
  var buf bytes.Buffer

  protoPorts := map[string]string {
    "ssh": "22",
    "http": "80",
    "https": "443",
  }

  flowMapping := map[string]string {
    "inbound": "INPUT",
    "outbound": "OUTPUT",
  }

  if strings.Compare(strings.ToLower(splitArr[2]), "icmp") != 0 {
    buf.WriteString(fmt.Sprintf("-A %s -p tcp -m --dport %s -j %s",
      flowMapping[strings.ToLower(splitArr[1])],
      protoPorts[strings.ToLower(splitArr[2])],
      strings.ToUpper(splitArr[0])))
  } else {
    // Handling ICMP
    buf.WriteString(fmt.Sprintf("-A %s -p icmp -m icmp --icmp-type 0 -j %s",
      flowMapping[strings.ToLower(splitArr[1])],
      strings.ToUpper(splitArr[0])))
    buf.WriteString(fmt.Sprintf("\n-A %s -p icmp -m icmp --icmp-type 3 -j %s",
      flowMapping[strings.ToLower(splitArr[1])],
      strings.ToUpper(splitArr[0])))
    buf.WriteString(fmt.Sprintf("\n-A %s -p icmp -m icmp --icmp-type 11 -j %s",
      flowMapping[strings.ToLower(splitArr[1])],
      strings.ToUpper(splitArr[0])))
  }
  return buf.String()
}

func main() {
  // Defines options for the human
  var humArgPtr = flag.String("i", "null", "Accepts human definition of f/w rule.")
  flag.Parse()

  parseStdIn(humArgPtr)
  fmt.Println(rulesetStringBuilder(humArgPtr))
}