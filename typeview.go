package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gizak/termui/v3/widgets"
)

type TypeView struct {
	View      *widgets.Paragraph
	Rendered  bool
	Type      []byte
	BigEndian bool
}

func (view *TypeView) Render() error {
	var endian binary.ByteOrder
	if view.Rendered {
		return nil
	}
	view.Rendered = true
	width := 20
	s := "Endian: "
	if view.BigEndian {
		s = s + " big - F12 to change"
		endian = binary.BigEndian
	} else {
		s = s + " little - F12 to change"
		endian = binary.LittleEndian
	}

	s = s + "\n" + fmt.Sprintf("%*s %*s %*s\n", width, "type", width, "signed", width, "unsigned")

	var i int8
	var ui uint8
	s = s + fmt.Sprintf("%*s ", width, "int8")
	signed_err := binary.Read(bytes.NewBuffer(view.Type), endian, &i)
	if signed_err == nil {
		s = s + fmt.Sprintf("%*d ", width, i)
	} else {
		s = s + fmt.Sprintf("%*s ", width, signed_err)
	}
	unsigned_err := binary.Read(bytes.NewBuffer(view.Type), endian, &ui)

	if unsigned_err == nil {
		s = s + fmt.Sprintf("%*d ", width, ui)
	} else {
		s = s + fmt.Sprintf("%*s ", width, unsigned_err)
	}
	s = s + "\n"

	var i16 int16
	var ui16 uint16
	s = s + fmt.Sprintf("%*s ", width, "int16")
	signed_err = binary.Read(bytes.NewBuffer(view.Type), endian, &i16)
	if signed_err == nil {
		s = s + fmt.Sprintf("%*d ", width, i16)
	} else {
		s = s + fmt.Sprintf("%*s ", width, signed_err)
	}
	unsigned_err = binary.Read(bytes.NewBuffer(view.Type), endian, &ui16)

	if unsigned_err == nil {
		s = s + fmt.Sprintf("%*d ", width, ui16)
	} else {
		s = s + fmt.Sprintf("%*s ", width, unsigned_err)
	}
	s = s + "\n"

	var i32 int32
	var ui32 uint32
	s = s + fmt.Sprintf("%*s ", width, "int32")
	signed_err = binary.Read(bytes.NewBuffer(view.Type), binary.LittleEndian, &i32)
	if signed_err == nil {
		s = s + fmt.Sprintf("%*d ", width, i32)
	} else {
		s = s + fmt.Sprintf("%*s ", width, signed_err)
	}
	unsigned_err = binary.Read(bytes.NewBuffer(view.Type), binary.LittleEndian, &ui32)

	if unsigned_err == nil {
		s = s + fmt.Sprintf("%*d ", width, ui32)
	} else {
		s = s + fmt.Sprintf("%*s ", width, unsigned_err)
	}
	s = s + "\n"

	var f32 float32
	s = s + fmt.Sprintf("%*s ", width, "float32")
	signed_err = binary.Read(bytes.NewBuffer(view.Type), binary.LittleEndian, &f32)
	if signed_err == nil {
		s = s + fmt.Sprintf("%f", f32)
	} else {
		s = s + fmt.Sprintf("%s", signed_err)
	}
	s = s + "\n"

	var f64 float64
	s = s + fmt.Sprintf("%*s ", width, "float64")
	signed_err = binary.Read(bytes.NewBuffer(view.Type), binary.LittleEndian, &f64)
	if signed_err == nil {
		s = s + fmt.Sprintf("%f", f64)
	} else {
		s = s + fmt.Sprintf("%s", signed_err)
	}
	s = s + "\n"
	view.View.Text = s

	return nil
}

func (view *TypeView) SetType(b []byte) {
	view.Type = b
	view.Rendered = false
}

func (view *TypeView) Init() {
	view.View = widgets.NewParagraph()
	view.View.Title = "Types"
}
