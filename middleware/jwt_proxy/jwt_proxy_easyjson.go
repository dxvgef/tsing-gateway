// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package jwt_proxy

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

func easyjsonAbf2d450DecodeLocalMiddlewareJwtProxy(in *jlexer.Lexer, out *JWTProxy) {
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
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "source_type":
			out.SourceType = string(in.String())
		case "source_name":
			out.SourceName = string(in.String())
		case "upstream_url":
			out.UpstreamURL = string(in.String())
		case "send_type":
			out.SendType = string(in.String())
		case "send_method":
			out.SendMethod = string(in.String())
		case "send_name":
			out.SendName = string(in.String())
		case "upstream_success_body":
			out.UpstreamSuccessBody = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonAbf2d450EncodeLocalMiddlewareJwtProxy(out *jwriter.Writer, in JWTProxy) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"source_type\":"
		out.RawString(prefix[1:])
		out.String(string(in.SourceType))
	}
	{
		const prefix string = ",\"source_name\":"
		out.RawString(prefix)
		out.String(string(in.SourceName))
	}
	{
		const prefix string = ",\"upstream_url\":"
		out.RawString(prefix)
		out.String(string(in.UpstreamURL))
	}
	{
		const prefix string = ",\"send_type\":"
		out.RawString(prefix)
		out.String(string(in.SendType))
	}
	{
		const prefix string = ",\"send_method\":"
		out.RawString(prefix)
		out.String(string(in.SendMethod))
	}
	{
		const prefix string = ",\"send_name\":"
		out.RawString(prefix)
		out.String(string(in.SendName))
	}
	if in.UpstreamSuccessBody != "" {
		const prefix string = ",\"upstream_success_body\":"
		out.RawString(prefix)
		out.String(string(in.UpstreamSuccessBody))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v JWTProxy) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonAbf2d450EncodeLocalMiddlewareJwtProxy(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JWTProxy) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonAbf2d450EncodeLocalMiddlewareJwtProxy(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JWTProxy) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonAbf2d450DecodeLocalMiddlewareJwtProxy(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JWTProxy) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonAbf2d450DecodeLocalMiddlewareJwtProxy(l, v)
}
