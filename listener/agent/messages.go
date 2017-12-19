/*
* Honeytrap
* Copyright (C) 2016-2017 DutchSec (https://dutchsec.com/)
*
* This program is free software; you can redistribute it and/or modify it under
* the terms of the GNU Affero General Public License version 3 as published by the
* Free Software Foundation.
*
* This program is distributed in the hope that it will be useful, but WITHOUT
* ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
* FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License for more
* details.
*
* You should have received a copy of the GNU Affero General Public License
* version 3 along with this program in the file "LICENSE".  If not, see
* <http://www.gnu.org/licenses/agpl-3.0.txt>.
*
* See https://honeytrap.io/ for more details. All requests should be sent to
* licensing@honeytrap.io
*
* The interactive user interfaces in modified source and object code versions
* of this program must display Appropriate Legal Notices, as required under
* Section 5 of the GNU Affero General Public License version 3.
*
* In accordance with Section 7(b) of the GNU Affero General Public License version 3,
* these Appropriate Legal Notices must retain the display of the "Powered by
* Honeytrap" logo and retain the original copyright notice. If the display of the
* logo is not reasonably feasible for technical reasons, the Appropriate Legal Notices
* must display the words "Powered by Honeytrap" and retain the original copyright notice.
 */
package agent

import (
	"errors"
	"net"
)

const (
	TypeHello             int = 0x00
	TypeReadWrite         int = 0x01
	TypePing              int = 0x05
	TypeEOF               int = 0x04
	TypeHandshake         int = 0x02
	TypeHandshakeResponse int = 0x03
)

type Handshake struct {
}

func (r *Handshake) UnmarshalBinary(data []byte) error {
	d := NewDecoder(data)

	if d.ReadUint8() != TypeHandshake {
		return errors.New("Not a handshake packet")
	}

	return nil
}

func (h Handshake) MarshalBinary() ([]byte, error) {
	e := Encoder{}
	e.WriteUint8(TypeHandshake)
	return e.Bytes(), nil
}

type HandshakeResponse struct {
	Addresses []net.Addr
}

func (h *HandshakeResponse) UnmarshalBinary(data []byte) error {
	d := NewDecoder(data)

	if d.ReadUint8() != TypeHandshakeResponse {
		return errors.New("Not a handshake packet")
	}

	n := d.ReadUint8()

	h.Addresses = make([]net.Addr, n)

	for i := 0; i < n; i++ {
		h.Addresses[i] = d.ReadAddr()
	}

	return nil
}

func (h HandshakeResponse) MarshalBinary() ([]byte, error) {
	e := Encoder{}
	e.WriteUint8(TypeHandshakeResponse)
	e.WriteUint8(len(h.Addresses))

	for _, address := range h.Addresses {
		e.WriteAddr(address)
	}

	return e.Bytes(), nil
}

type Hello struct {
	Token string
	Laddr net.Addr
	Raddr net.Addr
}

func (h Hello) MarshalBinary() ([]byte, error) {
	e := Encoder{}

	e.WriteUint8(TypeHello)
	e.WriteUint8(0)

	e.WriteString(h.Token)

	e.WriteAddr(h.Laddr)
	e.WriteAddr(h.Raddr)

	return e.Bytes(), nil
}

func (h *Hello) UnmarshalBinary(data []byte) error {
	decoder := NewDecoder(data)

	if decoder.ReadUint8() != TypeHello {
		return errors.New("Not a hello packet")
	}

	_ = decoder.ReadUint8() /* protocol */

	h.Token = decoder.ReadString()
	h.Laddr = decoder.ReadAddr()
	h.Raddr = decoder.ReadAddr()
	return nil
}

type Ping struct {
	Token string
	Laddr net.Addr
	Raddr net.Addr
}

func (h Ping) MarshalBinary() ([]byte, error) {
	e := Encoder{}

	e.WriteUint8(TypePing)
	e.WriteUint8(0)

	e.WriteString(h.Token)

	e.WriteAddr(h.Laddr)
	e.WriteAddr(h.Raddr)

	return e.Bytes(), nil
}

type EOF struct {
	Laddr net.Addr
	Raddr net.Addr
}

func (r *EOF) UnmarshalBinary(data []byte) error {
	decoder := NewDecoder(data)

	if decoder.ReadUint8() != TypeEOF {
		return errors.New("Not a eof packet")
	}

	decoder.ReadUint8()

	r.Laddr = decoder.ReadAddr()
	r.Raddr = decoder.ReadAddr()

	return nil
}

func (h EOF) MarshalBinary() ([]byte, error) {
	e := Encoder{}

	e.WriteUint8(TypeEOF)
	e.WriteUint8(0)

	e.WriteAddr(h.Laddr)
	e.WriteAddr(h.Raddr)

	return e.Bytes(), nil
}

type ReadWrite struct {
	Laddr net.Addr
	Raddr net.Addr

	Payload []byte
}

func (h ReadWrite) MarshalBinary() ([]byte, error) {
	e := Encoder{}

	e.WriteUint8(TypeReadWrite)
	e.WriteUint8(0)

	e.WriteAddr(h.Laddr)
	e.WriteAddr(h.Raddr)

	e.WriteData(h.Payload)

	return e.Bytes(), nil
}

func (r *ReadWrite) UnmarshalBinary(data []byte) error {
	decoder := NewDecoder(data)

	if decoder.ReadUint8() != TypeReadWrite {
		return errors.New("Not a read packet")
	}

	decoder.ReadUint8()

	r.Laddr = decoder.ReadAddr()
	r.Raddr = decoder.ReadAddr()

	r.Payload = decoder.ReadData()

	return nil
}