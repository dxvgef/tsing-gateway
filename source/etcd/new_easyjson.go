// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package etcd

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

func easyjson16d7ec28DecodeGithubComDxvgefTsingGatewaySourceEtcd(in *jlexer.Lexer, out *Etcd) {
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
		case "key_prefix":
			out.KeyPrefix = string(in.String())
		case "endpoints":
			if in.IsNull() {
				in.Skip()
				out.Endpoints = nil
			} else {
				in.Delim('[')
				if out.Endpoints == nil {
					if !in.IsDelim(']') {
						out.Endpoints = make([]string, 0, 4)
					} else {
						out.Endpoints = []string{}
					}
				} else {
					out.Endpoints = (out.Endpoints)[:0]
				}
				for !in.IsDelim(']') {
					var v1 string
					v1 = string(in.String())
					out.Endpoints = append(out.Endpoints, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "dial_timeout":
			out.DialTimeout = uint(in.Uint())
		case "username":
			out.Username = string(in.String())
		case "password":
			out.Password = string(in.String())
		case "auto_sync_interval":
			out.AutoSyncInterval = uint(in.Uint())
		case "dial_keep_alive_time":
			out.DialKeepAliveTime = uint(in.Uint())
		case "dial_keep_alive_timeout":
			out.DialKeepAliveTimeout = uint(in.Uint())
		case "max_call_send_msg_size":
			out.MaxCallSendMsgSize = uint(in.Uint())
		case "max_call_recv_msg_size":
			out.MaxCallRecvMsgSize = uint(in.Uint())
		case "reject_old_cluster":
			out.RejectOldCluster = bool(in.Bool())
		case "permit_without_stream":
			out.PermitWithoutStream = bool(in.Bool())
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
func easyjson16d7ec28EncodeGithubComDxvgefTsingGatewaySourceEtcd(out *jwriter.Writer, in Etcd) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"key_prefix\":"
		out.RawString(prefix[1:])
		out.String(string(in.KeyPrefix))
	}
	{
		const prefix string = ",\"endpoints\":"
		out.RawString(prefix)
		if in.Endpoints == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Endpoints {
				if v2 > 0 {
					out.RawByte(',')
				}
				out.String(string(v3))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"dial_timeout\":"
		out.RawString(prefix)
		out.Uint(uint(in.DialTimeout))
	}
	{
		const prefix string = ",\"username\":"
		out.RawString(prefix)
		out.String(string(in.Username))
	}
	{
		const prefix string = ",\"password\":"
		out.RawString(prefix)
		out.String(string(in.Password))
	}
	{
		const prefix string = ",\"auto_sync_interval\":"
		out.RawString(prefix)
		out.Uint(uint(in.AutoSyncInterval))
	}
	{
		const prefix string = ",\"dial_keep_alive_time\":"
		out.RawString(prefix)
		out.Uint(uint(in.DialKeepAliveTime))
	}
	{
		const prefix string = ",\"dial_keep_alive_timeout\":"
		out.RawString(prefix)
		out.Uint(uint(in.DialKeepAliveTimeout))
	}
	{
		const prefix string = ",\"max_call_send_msg_size\":"
		out.RawString(prefix)
		out.Uint(uint(in.MaxCallSendMsgSize))
	}
	{
		const prefix string = ",\"max_call_recv_msg_size\":"
		out.RawString(prefix)
		out.Uint(uint(in.MaxCallRecvMsgSize))
	}
	{
		const prefix string = ",\"reject_old_cluster\":"
		out.RawString(prefix)
		out.Bool(bool(in.RejectOldCluster))
	}
	{
		const prefix string = ",\"permit_without_stream\":"
		out.RawString(prefix)
		out.Bool(bool(in.PermitWithoutStream))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Etcd) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson16d7ec28EncodeGithubComDxvgefTsingGatewaySourceEtcd(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Etcd) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson16d7ec28EncodeGithubComDxvgefTsingGatewaySourceEtcd(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Etcd) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson16d7ec28DecodeGithubComDxvgefTsingGatewaySourceEtcd(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Etcd) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson16d7ec28DecodeGithubComDxvgefTsingGatewaySourceEtcd(l, v)
}