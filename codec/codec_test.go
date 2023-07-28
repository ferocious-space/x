package codec

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

type TestUndefinedPayload struct {
	Data string
}

func (t *TestUndefinedPayload) Type() CType {
	return "__codec.TestUndefined"
}

type TestPayload struct {
	Data string
}

func (t *TestPayload) Type() CType {
	return "__codec.Test"
}

func TestCodec_DecodeFrom(t *testing.T) {
	type fields struct {
		Unmarshal unmarshal
	}
	type args struct {
		r any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Data
	}{
		{
			name: "DecodeFrom",
			fields: fields{
				Unmarshal: NewUnmarshal(func() Data { return new(TestPayload) }),
			},
			args: args{
				r: []byte("{\"T\":\"__codec.Test\",\"D\":null,\"Data\":{\"Data\":\"test\"}}"),
			},
			want: &TestPayload{Data: "test"},
		},
		{
			name: "DecodeFromUndefined",
			fields: fields{
				Unmarshal: NewUnmarshal(func() Data { return new(TestPayload) }),
			},
			args: args{
				r: []byte("{\"T\":\"__codec.TestUndefined\",\"D\":null,\"Data\":{\"Data\":\"test\"}}"),
			},
			want: &ErrorPayload{err: fmt.Errorf("configuration not defined for __codec.TestUndefined")},
		},
		{
			name: "DecodeFromErr",
			fields: fields{
				Unmarshal: NewUnmarshal(func() Data { return new(TestPayload) }),
			},
			args: args{
				r: []byte("{\"T\":\"__codec.Test\",\"D\":null,\"Data\":{\"Data\":\"test\"}"),
			},
			want: &ErrorPayload{err: fmt.Errorf("cannot read JSON stream: expected comma after object element")},
		},
		{
			name: "DecodeFromErrDefault",
			fields: fields{
				Unmarshal: NewUnmarshal(func() Data { return new(TestPayload) }),
			},
			args: args{
				r: "{\"T\":\"__codec.Test\",\"D\":null,\"Data\":{\"Data\":\"test\"}",
			},
			want: &ErrorPayload{err: fmt.Errorf("unknown config")},
		},
		{
			name: "DecodeFromIO",
			fields: fields{
				Unmarshal: NewUnmarshal(func() Data { return new(TestPayload) }),
			},
			args: args{
				r: bytes.NewBufferString("{\"T\":\"__codec.Test\",\"D\":null,\"Data\":{\"Data\":\"test\"}}"),
			},
			want: &TestPayload{Data: "test"},
		},
		{
			name: "DecodeFromError",
			fields: fields{
				Unmarshal: NewUnmarshal(func() Data { return new(TestPayload) }),
			},
			args: args{
				r: bytes.NewBufferString("{\"T\":\"__codec.Test\",\"D\":null,\"Data\":{\"Data\":\"test\"}"),
			},
			want: &ErrorPayload{err: fmt.Errorf("cannot read JSON stream: expected comma after object element")},
		},
		{
			name: "DecodeFromInternal",
			fields: fields{
				Unmarshal: NewUnmarshal(func() Data { return new(TestPayload) }),
			},
			args: args{
				r: bytes.NewBufferString("{\"T\":\"__codec.Test\",\"D\":null,\"Data\":{.}}"),
			},
			want: &ErrorPayload{err: fmt.Errorf("cannot read JSON data: invalid character '.' looking for beginning of value")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCodec(tt.fields.Unmarshal)
			got := c.DecodeFrom(tt.args.r)
			switch tt.want.(type) {
			case *TestPayload:
				if !reflect.DeepEqual(got.D, tt.want) {
					t.Errorf("DecodeFrom() = %v, want %v", got, tt.want)
				}
			case *ErrorPayload:
				e, ok := got.D.(*ErrorPayload)
				if !ok {
					t.Errorf("DecodeFrom() = %T, want %T", got.D, tt.want)
				}
				if e.Error() != tt.want.(*ErrorPayload).Error() {
					t.Errorf("DecodeFrom() = %v, want %v", got.D.(*ErrorPayload).Error(), tt.want.(*ErrorPayload).Error())
				}
			}

		})
	}
}

type syntaxError interface {
	Is(err error) bool
}

func TestEncodeTo(t *testing.T) {
	type args struct {
		w       any
		payload Data
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "EncodeTo",
			args: args{
				w:       &bytes.Buffer{},
				payload: &TestPayload{Data: "test"},
			},
			wantErr: false,
		},
		{
			name: "EncodeToErr",
			args: args{
				w:       "",
				payload: &TestPayload{Data: "test"},
			},
			wantErr: true,
		},
		{
			name: "EncodeToByte",
			args: args{
				w:       []byte{},
				payload: &TestPayload{Data: "test"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EncodeTo(tt.args.w, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("EncodeTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestErrorPayload_Error(t *testing.T) {
	type fields struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Error",
			fields: fields{
				err: fmt.Errorf("test"),
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ErrorPayload{
				err: tt.fields.err,
			}
			if got := p.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorPayload_Type(t *testing.T) {
	type fields struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		want   CType
	}{
		{
			name: "Type",
			fields: fields{
				err: fmt.Errorf("test"),
			},
			want: "__codec.Error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			er := &ErrorPayload{
				err: tt.fields.err,
			}
			if got := er.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestErrorPayload_Unwrap(t *testing.T) {
	type fields struct {
		err error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Unwrap",
			fields: fields{
				err: fmt.Errorf("test"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ErrorPayload{
				err: tt.fields.err,
			}
			if err := p.Unwrap(); (err != nil) != tt.wantErr {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPayload_Type(t *testing.T) {
	type fields struct {
		T CType
		D Data
	}
	tests := []struct {
		name   string
		fields fields
		want   CType
	}{
		{
			name: "Type",
			fields: fields{
				T: "__codec.Test",
			},
			want: "__codec.Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Payload{
				T: tt.fields.T,
				D: tt.fields.D,
			}
			if got := p.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newCodecErr(t *testing.T) {
	type args struct {
		err  error
		wrap string
	}
	tests := []struct {
		name string
		args args
		want Payload
	}{
		{
			name: "newCodecErr",
			args: args{
				err: fmt.Errorf("test"),
			},
			want: Payload{
				T: "__codec.Error",
				D: &ErrorPayload{
					err: fmt.Errorf("test"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCodecErr(tt.args.err, tt.args.wrap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCodecErr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newCodecErr1(t *testing.T) {
	type args struct {
		err  error
		wrap string
	}
	tests := []struct {
		name string
		args args
		want Payload
	}{
		{
			name: "newCodecErr",
			args: args{
				err:  fmt.Errorf("test"),
				wrap: "wrap",
			},
			want: Payload{
				T: "__codec.Error",
				D: &ErrorPayload{
					err: fmt.Errorf("wrap: %w", fmt.Errorf("test")),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newCodecErr(tt.args.err, tt.args.wrap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCodecErr() = %v, want %v", got, tt.want)
			}
		})
	}
}
