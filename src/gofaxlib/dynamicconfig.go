// This file is part of the GOfax.IP project - https://github.com/gonicus/gofaxip
// Copyright (C) 2014 GONICUS GmbH, Germany - http://www.gonicus.de
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; version 2
// of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

package gofaxlib

// Functions to parse and write
// HylaFax configuration files
//
// Right now this only contains a very simple
// implementation of DynamicConfig.
// The only supported option is "RejectCall: true"
//
// TODO: Merge this with qfile.go
// because it's pretty similar

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type param struct {
	Tag   string
	Value string
}

type HylaConfig struct {
	params []param
}

func (h *HylaConfig) GetFirst(tag string) string {
	for _, param := range h.params {
		if param.Tag == tag {
			return param.Value
		}
	}
	return ""
}

func DynamicConfig(command string, cidnum string, cidname string, recipient string) (*HylaConfig, error) {
	if Config.Gofaxd.DynamicConfig == "" {
		return nil, errors.New("No DynamicConfig command provided")
	}

	cmd := exec.Command(Config.Gofaxd.DynamicConfig, cidnum, cidname, recipient)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	h := new(HylaConfig)

	scanner := bufio.NewScanner(bytes.NewBuffer(out))
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ":", 2)
		if len(parts) != 2 {
			continue
		}
		h.params = append(h.params, param{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])})
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return h, nil
}