package codec

import (
	"fmt"
	"io"

	"x/buffers"
	"x/json"
)

type Codec struct {
	Unmarshal unmarshal
}

func NewCodec(u unmarshal) Codec {
	return Codec{Unmarshal: u}
}

type codecPayload struct {
	Payload
	Data json.Raw
}

type ErrorPayload struct {
	err error
}

func (p *ErrorPayload) Error() string {
	return p.err.Error()
}

func (p *ErrorPayload) Unwrap() error {
	return p.err
}

func (*ErrorPayload) Type() CType {
	return "__codec.Error"
}

func newCodecErr(err error, wrap string) Payload {
	if wrap != "" {
		err = fmt.Errorf("%s: %w", wrap, err)
	}
	ev := &ErrorPayload{err: err}

	return Payload{
		T: ev.Type(),
		D: ev,
	}
}

func (c Codec) DecodeFrom(r any) Payload {
	var op codecPayload
	buffer := buffers.GetInstance().GetBytes()
	op.Data = buffer
	defer buffers.GetInstance().PutBytes(buffer)

	switch x := r.(type) {
	case io.Reader:
		if err := json.DecodeStream(x, &op); err != nil {
			return newCodecErr(err, "cannot read JSON stream")
		}
	case []byte:
		if err := json.Unmarshal(x, &op); err != nil {
			return newCodecErr(err, "cannot read JSON stream")
		}
	default:
		return newCodecErr(fmt.Errorf("unknown config"), "")
	}

	fn := c.Unmarshal.Lookup(op.T)
	if fn == nil {
		return newCodecErr(fmt.Errorf("configuration not defined for %s", op.T), "")
	}
	op.Payload.D = fn()
	if err := op.Data.UnmarshalTo(op.Payload.D); err != nil {
		return newCodecErr(err, "cannot read JSON data")
	}
	return op.Payload
}

func EncodeTo(w any, payload Data) (err error) {
	op := new(codecPayload)
	op.T = payload.Type()
	marshal, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_ = op.Data.UnmarshalJSON(marshal)
	switch x := w.(type) {
	case []byte:
		x, err = json.Marshal(op)
	case io.Writer:
		err = json.EncodeStream(x, op)
	default:
		return fmt.Errorf("%T not supported", w)
	}
	return err
}
