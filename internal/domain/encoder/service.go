package encoder

type Encoder struct {
	key string
}

func (e *Encoder) Encode(input string) (string, error)
func (e *Encoder) Decode(input string) (string, error)
