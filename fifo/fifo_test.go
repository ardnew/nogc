package seq

import (
	"bytes"
	"io"
	"testing"
)

func TestList_Configure(t *testing.T) {
	type fields struct {
		buf buf
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantOk bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &List{
				buf: tt.fields.buf,
			}
			if gotOk := l.Configure(tt.args.p); gotOk != tt.wantOk {
				t.Errorf("List.Configure() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestRing_Configure(t *testing.T) {
	type fields struct {
		buf buf
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantOk bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Ring{
				buf: tt.fields.buf,
			}
			if gotOk := r.Configure(tt.args.p); gotOk != tt.wantOk {
				t.Errorf("Ring.Configure() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_buf_init(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	type args struct {
		p        []byte
		capacity uint32
		mode     mode
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantOk bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			if gotOk := b.init(tt.args.p, tt.args.capacity, tt.args.mode); gotOk != tt.wantOk {
				t.Errorf("buf.init() = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func Test_buf_Len(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			if got := b.Len(); got != tt.want {
				t.Errorf("buf.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buf_Cap(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			if got := b.Cap(); got != tt.want {
				t.Errorf("buf.Cap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buf_Reset(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			b.Reset()
		})
	}
}

func Test_buf_Read(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			gotN, err := b.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("buf.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("buf.Read() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_buf_Write(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			gotN, err := b.Write(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("buf.Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("buf.Write() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_buf_readFrom(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	type args struct {
		r  io.Reader
		lo int
		hi int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			gotN, err := b.readFrom(tt.args.r, tt.args.lo, tt.args.hi)
			if (err != nil) != tt.wantErr {
				t.Errorf("buf.readFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("buf.readFrom() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_buf_ReadFrom(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			gotN, err := b.ReadFrom(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("buf.ReadFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("buf.ReadFrom() = %v, want %v", gotN, tt.wantN)
			}
		})
	}
}

func Test_buf_writeTo(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	type args struct {
		lo int
		hi int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantN   int
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			w := &bytes.Buffer{}
			gotN, err := b.writeTo(w, tt.args.lo, tt.args.hi)
			if (err != nil) != tt.wantErr {
				t.Errorf("buf.writeTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("buf.writeTo() = %v, want %v", gotN, tt.wantN)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("buf.writeTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_buf_WriteTo(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantN   int64
		wantW   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			w := &bytes.Buffer{}
			gotN, err := b.WriteTo(w)
			if (err != nil) != tt.wantErr {
				t.Errorf("buf.WriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotN != tt.wantN {
				t.Errorf("buf.WriteTo() = %v, want %v", gotN, tt.wantN)
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("buf.WriteTo() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func Test_buf_ReadByte(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantC   byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			gotC, err := b.ReadByte()
			if (err != nil) != tt.wantErr {
				t.Errorf("buf.ReadByte() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotC != tt.wantC {
				t.Errorf("buf.ReadByte() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func Test_buf_WriteByte(t *testing.T) {
	type fields struct {
		Byte  []byte
		capt  uint32
		head  uint32
		tail  uint32
		mode  mode
		valid bool
	}
	type args struct {
		c byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &buf{
				Byte:  tt.fields.Byte,
				capt:  tt.fields.capt,
				head:  tt.fields.head,
				tail:  tt.fields.tail,
				mode:  tt.fields.mode,
				valid: tt.fields.valid,
			}
			if err := b.WriteByte(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("buf.WriteByte() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
