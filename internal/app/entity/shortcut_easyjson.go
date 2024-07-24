// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package entity

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson935ce6eeDecodeShortenerInternalAppEntity(in *jlexer.Lexer, out *Shortcut) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "uuid":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.ID).UnmarshalText(data))
			}
		case "user_id":
			if data := in.UnsafeBytes(); in.Ok() {
				in.AddError((out.UserID).UnmarshalText(data))
			}
		case "short_url":
			out.ShortURL = string(in.String())
		case "original_url":
			out.OriginalURL = string(in.String())
		default:
			in.AddError(&jlexer.LexerError{
				Offset: in.GetPos(),
				Reason: "unknown field",
				Data:   key,
			})
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson935ce6eeEncodeShortenerInternalAppEntity(out *jwriter.Writer, in Shortcut) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"uuid\":"
		out.RawString(prefix[1:])
		out.RawText((in.ID).MarshalText())
	}
	{
		const prefix string = ",\"user_id\":"
		out.RawString(prefix)
		out.RawText((in.UserID).MarshalText())
	}
	{
		const prefix string = ",\"short_url\":"
		out.RawString(prefix)
		out.String(string(in.ShortURL))
	}
	{
		const prefix string = ",\"original_url\":"
		out.RawString(prefix)
		out.String(string(in.OriginalURL))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Shortcut) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson935ce6eeEncodeShortenerInternalAppEntity(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Shortcut) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson935ce6eeEncodeShortenerInternalAppEntity(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Shortcut) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson935ce6eeDecodeShortenerInternalAppEntity(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Shortcut) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson935ce6eeDecodeShortenerInternalAppEntity(l, v)
}
