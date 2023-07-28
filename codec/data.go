package codec

type CType string

type Data interface {
	Type() CType
}

type Payload struct {
	T CType
	D Data
}

func (p *Payload) Type() CType {
	return p.T
}

type Func func() Data
type funcID struct {
	T CType
}

type unmarshal struct {
	r map[CType]Func
}

func (receiver unmarshal) Add(funcs ...Func) {
	for _, fn := range funcs {
		ev := fn()
		receiver.r[ev.Type()] = fn
	}
}
func (receiver unmarshal) Lookup(t CType) Func {
	return receiver.r[t]
}

func NewUnmarshal(funcs ...Func) unmarshal {
	m := unmarshal{r: map[CType]Func{}}
	m.Add(func() Data { return new(ErrorPayload) })
	m.Add(funcs...)
	return m
}
