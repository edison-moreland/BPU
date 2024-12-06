package netlist

import (
	"encoding/json"
	"errors"
)

// Netlist import from yosys

func MarshalNetlist(netlist *Netlist) ([]byte, error) {
	return json.MarshalIndent(netlist, "", "  ")
}

func UnmarshalNetlist(data []byte) (*Netlist, error) {
	var netlist Netlist
	err := json.Unmarshal(data, &netlist)
	if err != nil {
		return nil, err
	}

	return &netlist, nil
}

type Netlist struct {
	Creator string            `json:"creator"`
	Modules map[string]Module `json:"modules"`
}

type Module struct {
	Attributes             map[string]string `json:"attributes"`
	ParameterDefaultValues map[string]string `json:"parameter_default_values,omitempty"`
	Ports                  map[string]Port   `json:"ports"`
	Cells                  map[string]Cell   `json:"cells"`
	Memories               map[string]Memory `json:"memories,omitempty"`
	Netnames               map[string]Net    `json:"netnames"`
}

type Direction string

const (
	Direction_Input  Direction = "input"
	Direction_Output Direction = "output"
	Direction_Inout  Direction = "inout"
)

type Port struct {
	Direction Direction `json:"direction"`
	Bits      []Bit     `json:"bits"`
	Offset    int       `json:"offset,omitempty"`
	UpTo      IntBool   `json:"upto,omitempty"`
	Signed    IntBool   `json:"signed,omitempty"`
}

type Cell struct {
	HideName       IntBool              `json:"hide_name"`
	Type           string               `json:"type"`
	Model          string               `json:"model,omitempty"`
	Parameters     map[string]string    `json:"parameters"`
	Attributes     map[string]string    `json:"attributes"`
	PortDirections map[string]Direction `json:"port_directions"`
	Connections    map[string][]Bit     `json:"connections"`
}

type Constant string

const (
	Constant_0 Constant = "0"
	Constant_1 Constant = "1"
	Constant_x Constant = "x"
	Constatn_z Constant = "z"
)

// Bit represents a connection to a wire
type Bit struct {
	Index    int // Index of wire, or -1 for a Constant
	Constant Constant
}

func (b *Bit) IsConstant() bool {
	return b.Index == -1
}

func (b *Bit) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &b.Index)
	if err != nil {
		b.Index = -1
		err2 := json.Unmarshal(data, &b.Constant)
		if err2 != nil {
			return errors.Join(err, err2)
		}
	}

	return nil
}

func (b *Bit) MarshalJSON() ([]byte, error) {
	if b.IsConstant() {
		return json.Marshal(b.Constant)
	}
	return json.Marshal(b.Index)
}

type IntBool bool

func (b *IntBool) UnmarshalJSON(data []byte) error {
	var i int
	err := json.Unmarshal(data, &i)
	if err != nil {
		return err
	}

	*b = i != 0

	return nil
}

func (b IntBool) MarshalJSON() ([]byte, error) {
	if b {
		return json.Marshal(1)

	} else {
		return json.Marshal(0)

	}
}

type Memory struct {
	HideName    IntBool           `json:"hide_name"`
	Attributes  map[string]string `json:"attributes"`
	Width       int               `json:"width"`
	StartOffset int               `json:"start_offset"`
	Size        int               `json:"size"`
}

type Net struct {
	HideName   IntBool           `json:"hide_name"`
	Bits       []Bit             `json:"bits"`
	Attributes map[string]string `json:"attributes"`
	Offset     int               `json:"offset,omitempty"`
	UpTo       IntBool           `json:"upto,omitempty"`
	Signed     IntBool           `json:"signed,omitempty"`
}
