/*
Copyright 2016 The Smudge Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package smudge

import (
	"net"
	"reflect"
	"testing"
)

// Identical but distinct instance from node1b
var node1a = Node{
	ip:          net.IP([]byte{127, 0, 0, 1}),
	port:        1234,
	timestamp:   87878787,
	status:      StatusAlive,
	emitCounter: 42,
	pingMillis:  PingNoData}

// Identical but distinct instance from node1a
var node1b = Node{
	ip:          net.IP([]byte{127, 0, 0, 1}),
	port:        1234,
	timestamp:   87878787,
	status:      StatusAlive,
	emitCounter: 42,
	pingMillis:  PingNoData}

// Different from node1a and node1b
var node2 = Node{
	ip:          net.IP([]byte{127, 0, 0, 1}),
	port:        10001,
	timestamp:   GetNowInMillis(),
	status:      StatusAlive,
	emitCounter: 42,
	pingMillis:  PingNoData}

var message1a = message{
	sender:          &node1a,
	senderHeartbeat: 255,
	verb:            verbPing}

var message1b = message{
	sender:          &node1b,
	senderHeartbeat: 255,
	verb:            verbPing}

var message2 = message{
	sender:          &node2,
	senderHeartbeat: 23,
	verb:            verbAck}

// Does deep equality of two different but identical messages return true?
func TestDeepEqualityTrue(t *testing.T) {
	if !reflect.DeepEqual(message1a, message1b) {
		t.Fail()
	}
}

// Does deep equality of two different messages return false?
func TestDeepEqualityFalse(t *testing.T) {
	if reflect.DeepEqual(message1a, message2) {
		t.Fail()
	}
}

// Endode and decode a simple message without any members, and see if
// the input/output match.
func TestEncodeDecodeBasic(t *testing.T) {
	timestamp := uint32(87878787)

	sender := Node{
		ip:         net.IP([]byte{127, 0, 0, 1}),
		port:       1234,
		timestamp:  timestamp,
		pingMillis: PingNoData,
	}

	message := message{
		sender:          &sender,
		senderHeartbeat: 255,
		verb:            verbPing}

	ip := net.IP([]byte{127, 0, 0, 1})
	bytes := message.encode()

	decoded, err := decodeMessage(ip, bytes)
	decoded.sender.timestamp = timestamp

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(message, decoded) {
		t.Error("Messages do not match:")

		t.Log(" Input:", message)
		t.Log("Output:", decoded)
		t.Log(" Input node:", message.sender)
		t.Log("Output node:", decoded.sender)
	}
}

// Endode and decode a simple IPv6 message without any members, and see if
// the input/output match.
func TestEncodeDecodeBasicIPv6(t *testing.T) {
	timestamp := uint32(87878787)

	sender := Node{
		ip:         net.IP{255, 254, 253, 252, 251, 250, 240, 230, 220, 210, 200, 10, 20, 30, 40, 50},
		port:       1234,
		timestamp:  timestamp,
		pingMillis: PingNoData,
	}

	message := message{
		sender:          &sender,
		senderHeartbeat: 255,
		verb:            verbPing}

	ipLen = net.IPv6len // encode IPv6
	ip := net.IP{255, 254, 253, 252, 251, 250, 240, 230, 220, 210, 200, 10, 20, 30, 40, 50}
	bytes := message.encode()
	decoded, err := decodeMessage(ip, bytes)
	decoded.sender.timestamp = timestamp

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(message, decoded) {
		t.Error("Messages do not match:")

		t.Log(" Input:", message)
		t.Log("Output:", decoded)
		t.Log(" Input node:", message.sender)
		t.Log("Output node:", decoded.sender)
	}

	ipLen = net.IPv4len // reset to IPv4
}

// Endode and decode a simple message with one member, and see if
// the input/output match.
func TestEncodeDecode1Member(t *testing.T) {
	timestamp := uint32(87878787)

	sender := Node{
		ip:         net.IP([]byte{127, 0, 0, 1}),
		port:       1234,
		timestamp:  timestamp,
		pingMillis: PingNoData,
	}

	member := Node{
		ip:         net.IP([]byte{127, 0, 0, 2}).To16(),
		port:       9000,
		timestamp:  timestamp,
		pingMillis: PingNoData,
	}

	message := message{
		sender:          &sender,
		senderHeartbeat: 255,
		verb:            verbPing}
	message.addMember(&member, StatusDead, 38, &member)

	if len(message.members) != 1 {
		t.Error("No member in the input members list!")
	}

	ip := net.IP([]byte{127, 0, 0, 1})
	bytes := message.encode()
	if len(bytes) != 28 {
		t.Error("Encoded message length is invalid.")
		t.Log("Should be 28 but found: ", len(bytes))
	}

	decoded, err := decodeMessage(ip, bytes)
	t.Log("bytes: ", bytes)
	decoded.sender.timestamp = timestamp
	decoded.members[0].node.timestamp = timestamp
	decoded.members[0].source.timestamp = timestamp

	if err != nil {
		t.Error(err)
	}

	if len(decoded.members) != 1 {
		t.Error("No member in the output members list!")
	}

	if !reflect.DeepEqual(message, decoded) {
		t.Error("Messages do not match")

		t.Log(" Input:", message.members[0])
		t.Log("Output:", decoded.members[0])

		t.Log(" Input node:", message.members[0].node)
		t.Log("Output node:", decoded.members[0].node)

		t.Log(" Input source:", message.members[0].source)
		t.Log("Output source:", decoded.members[0].source)
	}
}

// Endode and decode a simple message with one ipv6 member, and see if
// the input/output match.
func TestEncodeDecode1MemberIPv6(t *testing.T) {
	timestamp := uint32(87878787)

	sender := Node{
		ip:         net.IP{255, 254, 253, 252, 251, 250, 240, 230, 220, 210, 200, 10, 20, 30, 40, 50},
		port:       1234,
		timestamp:  timestamp,
		pingMillis: PingNoData,
	}

	member := Node{
		ip:         net.IP{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160},
		port:       9000,
		timestamp:  timestamp,
		pingMillis: PingNoData,
	}

	message := message{
		sender:          &sender,
		senderHeartbeat: 255,
		verb:            verbPing}
	message.addMember(&member, StatusDead, 38, &member)

	if len(message.members) != 1 {
		t.Error("No member in the input members list!")
	}

	ipLen = net.IPv6len // encode for IPv6
	ip := net.IP{255, 254, 253, 252, 251, 250, 240, 230, 220, 210, 200, 10, 20, 30, 40, 50}
	bytes := message.encode()
	if len(bytes) != 52 {
		t.Error("Encoded message length is invalid.")
		t.Log("Should be 52 but found: ", len(bytes))
	}

	decoded, err := decodeMessage(ip, bytes)
	decoded.sender.timestamp = timestamp
	decoded.members[0].node.timestamp = timestamp
	decoded.members[0].source.timestamp = timestamp

	if err != nil {
		t.Error(err)
	}

	if len(decoded.members) != 1 {
		t.Error("No member in the output members list!")
	}

	if !reflect.DeepEqual(message, decoded) {
		t.Error("Messages do not match")

		t.Log(" Input:", message.members[0])
		t.Log("Output:", decoded.members[0])
		t.Log(" Input node:", message.members[0].node)
		t.Log("Output node:", decoded.members[0].node)
		t.Log(" Input source:", message.members[0].source)
		t.Log("Output source:", decoded.members[0].source)
	}

	ipLen = net.IPv4len // reset to IPv4 for next test
}

// Endode and decode a simple message with one member and message, and see if
// the input/output match.
func TestEncodeDecode1MemberBroadcast(t *testing.T) {
	timestamp := uint32(87878787)

	sender := Node{
		ip:         net.IP([]byte{127, 0, 0, 1}),
		port:       1234,
		timestamp:  timestamp,
		pingMillis: PingNoData}

	member := Node{
		ip:         net.IP([]byte{127, 0, 0, 2}),
		port:       9000,
		timestamp:  timestamp,
		pingMillis: PingNoData}

	message := message{
		sender:          &sender,
		senderHeartbeat: 255,
		verb:            verbPing}
	message.addMember(&member, StatusDead, 38, &member)

	broadcast := Broadcast{
		bytes:  []byte("This is a message"), //len=17
		origin: &sender,
		index:  42}
	message.addBroadcast(&broadcast)

	if message.broadcast == nil {
		t.Error("Broadcast not set properly")
	}

	ip := net.IP([]byte{127, 0, 0, 1})
	bytes := message.encode()
	if len(bytes) != 57 {
		t.Error("Encoded message length is invalid.")
		t.Log("Should be 57 but found: ", len(bytes))
	}

	decoded, err := decodeMessage(ip, bytes)
	decoded.sender.timestamp = timestamp
	decoded.members[0].node.timestamp = timestamp
	decoded.members[0].source.timestamp = timestamp

	if err != nil {
		t.Error(err)
	}

	if decoded.broadcast == nil {
		t.Error("Broadcast not decoded")
	}

	message.broadcast.origin = nil
	decoded.broadcast.origin = nil

	if !reflect.DeepEqual(message.broadcast, decoded.broadcast) {
		t.Error("Broadcasts do not match:")
		t.Error(" Input bcast:", message.broadcast)
		t.Error("Output bcast:", decoded.broadcast)
	}
}

// Endode and decode a simple message with one ipv6 member and message, and see if
// the input/output match.
func TestEncodeDecode1MemberBroadcastIPv6(t *testing.T) {
	timestamp := uint32(87878787)

	sender := Node{
		ip:         net.IP{255, 254, 253, 252, 251, 250, 240, 230, 220, 210, 200, 10, 20, 30, 40, 50},
		port:       1234,
		timestamp:  timestamp,
		pingMillis: PingNoData}

	member := Node{
		ip:         net.IP{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160},
		port:       9000,
		timestamp:  timestamp,
		pingMillis: PingNoData}

	message := message{
		sender:          &sender,
		senderHeartbeat: 255,
		verb:            verbPing}
	message.addMember(&member, StatusDead, 38, &member)

	broadcast := Broadcast{
		bytes:  []byte("This is a message"),
		origin: &sender,
		index:  42}
	message.addBroadcast(&broadcast)

	if message.broadcast == nil {
		t.Error("Broadcast not set properly")
	}

	ipLen = net.IPv6len // encode for IPv6
	ip := net.IP{255, 254, 253, 252, 251, 250, 240, 230, 220, 210, 200, 10, 20, 30, 40, 50}
	bytes := message.encode()
	if len(bytes) != 93 {
		t.Error("Encoded message length is invalid.")
		t.Log("Should be 93 but found: ", len(bytes))
	}

	decoded, err := decodeMessage(ip, bytes)
	decoded.sender.timestamp = timestamp
	decoded.members[0].node.timestamp = timestamp
	decoded.members[0].source.timestamp = timestamp

	if err != nil {
		t.Error(err)
	}

	if decoded.broadcast == nil {
		t.Error("Broadcast not decoded")
	}

	message.broadcast.origin = nil
	decoded.broadcast.origin = nil

	if !reflect.DeepEqual(message.broadcast, decoded.broadcast) {
		t.Error("Broadcasts do not match:")
		t.Error(" Input bcast:", message.broadcast)
		t.Error("Output bcast:", decoded.broadcast)
	}

	ipLen = net.IPv4len
}
