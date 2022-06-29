package seq

import (
	"io"

	"github.com/ardnew/nogc"
)

// mode defines the behavior when writing to a full buf.
type mode bool

const (
	retain  mode = iota != 0 // retain bytes in buf, never allow Write
	dequeue                  // dequeue bytes from buf, always allow Write
)

// List defines a fixed-length queue of bytes in which no bytes may be added
// when the queue is full.
type List struct{ buf }

// Ring defines a fixed-length queue of bytes in which old bytes are dequeued if
// new bytes are added when the queue is full.
type Ring struct{ buf }

// buf defines a first-in, first-out (FIFO) queue of bytes.
type buf struct {
	Byte  []byte
	capt  uint32
	head  uint32
	tail  uint32
	mode  mode
	valid bool
}

// Configure initializes l using all of p as storage.
// The initial length of l is 0; any data already in p may be overwritten.
// The capacity of l is permanently len(p).
// Callers must not modify p after initializing.
func (l *List) Configure(p []byte) (ok bool) {
	l.valid = l.init(p, uint32(len(p)), retain)
	return l.valid
}

// Configure initializes r using all of p as storage.
// The initial length of r is 0; any data already in p may be overwritten.
// The capacity of r is permanently len(p).
// Callers must not modify p after initializing.
func (r *Ring) Configure(p []byte) (ok bool) {
	r.valid = r.init(p, uint32(len(p)), dequeue)
	return r.valid
}

// init initializes the configuration.
func (b *buf) init(p []byte, capacity uint32, mode mode) (ok bool) {
	if b == nil {
		return false
	}
	b.Byte = p
	b.capt = capacity
	b.head = 0
	b.tail = 0
	b.mode = mode
	ok = p != nil
	return
}

// Len returns the number of bytes.
func (b *buf) Len() int {
	if b == nil || !b.valid {
		return 0
	}
	return int(b.tail - b.head)
}

// Cap returns the byte capacity.
func (b *buf) Cap() int {
	if b == nil || !b.valid {
		return 0
	}
	return int(b.capt)
}

// Reset sets the number of bytes to 0.
func (b *buf) Reset() {
	if b == nil || !b.valid {
		return
	}
	b.head = 0
	b.tail = 0
}

// Read copies up to len(p) unread bytes from b to p and returns the number of
// bytes copied.
func (b *buf) Read(p []byte) (n int, err error) {
	if b == nil || !b.valid {
		return 0, &nogc.ErrInvalidReceiver
	}
	if p == nil {
		return 0, &nogc.ErrInvalidArgument
	}
	var ns int
	ns, n = b.Len(), len(p)
	if ns <= n {
		n, err = ns, io.EOF
	}
	h := b.head
	for i := range p[:n] {
		p[i] = b.Byte[h%b.capt]
		h++
	}
	b.head = h
	return
}

// Write appends up to len(p) bytes from p to b and returns the number of bytes
// copied.
//
// Write will only write to the free space in b and then return ErrWriteOverflow
// if all of p could not be copied.
func (b *buf) Write(p []byte) (n int, err error) {
	if b == nil || !b.valid {
		return 0, &nogc.ErrInvalidReceiver
	}
	if p == nil {
		return 0, &nogc.ErrInvalidArgument
	}
	np := len(p)
	if np == 0 {
		// Source buffer is empty; there are no bytes from p to copy into b.
		// This isn't considered a Write error; it only means nothing was written.
		return
	}
	h, t := b.head, b.tail
	for _, c := range p {
		if t-h >= b.capt {
			break
		}
		b.Byte[t%b.capt] = c
		t++
		n++
	}
	b.tail = t
	if n < np {
		err = &nogc.ErrWriteOverflow
	}
	return
}

func (b *buf) readFrom(r io.Reader, lo, hi int) (n int, err error) {
	// The caller is responsible for coordinating calls to readFrom when the
	// elements of b will not be stored contiguously in the backing array.
	//
	// We can do a brief sanity check on the indices to prevent A/V errors, but no
	// attempt is made to normalize, split the range into slices, or verify the
	// range starts at tail and spans only the free-space region.
	if lo >= hi || lo < 0 || hi > int(b.capt) {
		// The above condition implies 0<=lo < hi<=N:
		//   If lo<hi and lo>=0, then hi>0 (i.e.: 0<=lo<hi => hi>0).
		//   If lo<hi and hi<=N, then lo<N (i.e.: lo<hi<=N => lo<N).
		return 0, &nogc.ErrOutOfRange
	}
	n, err = r.Read(b.Byte[lo:hi])
	// Catch any attempt to return io.EOF and return nil instead.
	// See documentation on io.ReaderFrom, and io.Copy.
	if err == io.EOF {
		err = nil
	}
	// Extend the length of b by the number of bytes copied.
	b.tail += uint32(n)
	return
}

// ReadFrom copies bytes from r to b until all bytes have been read or an error
// was encountered. Returns the number of bytes successfully copied.
//
// A successful ReadFrom returns err == nil and not err == io.EOF.
// ReadFrom is defined to read from r until all bytes have been read (io.EOF),
// so it does not treat io.EOF from r as an error to be reported.
//
// Bytes are copied directly without any buffering, so r and b must not overlap
// if both are implemented as buffers of physical memory.
func (b *buf) ReadFrom(r io.Reader) (n int64, err error) {
	if b == nil || !b.valid {
		return 0, &nogc.ErrInvalidReceiver
	}
	if r == nil {
		return 0, &nogc.ErrInvalidArgument
	}
	h, t := b.head, b.tail
	// Convert head and tail to physical array indices to determine if the used
	// elements span a contiguous region of memory in the backing array.
	ih, it := h%b.capt, t%b.capt
	// If the array indices are equal, with head not eqaul to tail, then the
	// backing array is filled to capacity. We have nowhere to store the bytes
	// from r. We can either overwrite the existing buffer or retain it and return
	// an error. Opting for the latter so that no bytes are lost, and it gives the
	// caller an opportunity to remedy the situation.
	if h != t && ih == it {
		return 0, &nogc.ErrReadOverflow
	}
	// When first-in (head) element is located at index 0 (i.e., the start of the
	// backing array), then the unused space spans exactly one region only;
	// namely, the range from last-in (tail) element to end of the array:
	//   (0123456789A) === Array index reference
	//   [HxxxT......]     Free-space forms contiguous span [4..A]
	if ih == 0 {
		nr, errr := b.readFrom(r, int(it), int(b.capt))
		return int64(nr), errr
	}
	// Tail grows as elements are added to the ring buffer. Thus, if tail is less
	// than head, then the tail index has wrapped around after growing beyond the
	// backing array's high index (capacity-1), but the head index has not yet
	// wrapped around.
	if it > ih {
		// Tail has not overflowed its storage prior to head, which is the normal
		// case, and thus the unused elements potentially exist in two separate
		// contiguous regions of the backing array. The first region (1) spans from
		// the last-in (tail) element to the end of the array, and the second
		// region (2) spans from the start of the array to the first-in (head)
		// element:
		//   (0123456789A) === Array index reference
		//   [...HxxxT...]     Free-space in region 1 [7..A] and region 2 [0..2]
		//   [......HxxxT]     Free-space in region 1 (A) and region 2 [0..5]
		var (
			n1, n2     int
			err1, err2 error
		)
		// (1.) Copy into tail to end of the backing array
		if n1, err1 = b.readFrom(r, int(it), int(b.capt)); err1 != nil {
			return int64(n1), err1
		}
		// (2.) Copy into start of the backing array to head (if region length > 0).
		if ih > 0 {
			n2, err2 = b.readFrom(r, 0, int(ih))
		}
		return int64(n1 + n2), err2
	}
	// Unused elements form contiguous span in backing array from last-in (tail)
	// to first-in (head), which may include the entire front of the backing array
	// when tail is exactly equal to a multiple of array capacity.
	//   (0123456789A) === Array index reference
	//   [xxT......Hx]     Free-space forms contiguous span [2..8]
	//   [T......Hxxx]     Free-space forms contiguous span [0..6]
	nr, errr := b.readFrom(r, int(it), int(ih))
	return int64(nr), errr
}

func (b *buf) writeTo(w io.Writer, lo, hi int) (n int, err error) {
	// The caller is responsible for coordinating calls to writeTo when the
	// elements of b are not stored contiguously in the backing array.
	//
	// We can do a brief sanity check on the indices to prevent A/V errors,
	// but no attempt is made to normalize or split the range into slices.
	if lo >= hi || lo < 0 || hi > int(b.capt) {
		// The above condition implies 0<=lo < hi<=N:
		//   If lo<hi and lo>=0, then hi>0 (i.e.: 0<=lo<hi => hi>0).
		//   If lo<hi and hi<=N, then lo<N (i.e.: lo<hi<=N => lo<N).
		return 0, &nogc.ErrOutOfRange
	}
	n, err = w.Write(b.Byte[lo:hi])
	// Decrease length by the number of bytes copied.
	b.head += uint32(n)
	return
}

// WriteTo copies bytes from b to w until all bytes have been written or an
// error was encountered. Returns the number of bytes successfully copied.
//
// Bytes are copied directly without any buffering, so w and b must not overlap
// if both are implemented as buffers of physical memory.
func (b *buf) WriteTo(w io.Writer) (n int64, err error) {
	if b == nil || !b.valid {
		return 0, &nogc.ErrInvalidReceiver
	}
	if w == nil {
		return 0, &nogc.ErrInvalidArgument
	}
	h, t := b.head, b.tail
	if h == t {
		// Buffer is empty, writing zero bytes to w.
		return 0, io.EOF
	}
	// Convert head and tail to physical array indices to determine if the used
	// elements span a contiguous region of memory in the backing array.
	ih, it := h%b.capt, t%b.capt
	// Tail grows as elements are added to the ring buffer. Thus, if tail is less
	// than head, then the tail index has wrapped around after growing beyond the
	// backing array's high index (capacity-1), but the head index has not yet
	// wrapped around.
	if it <= ih {
		// Tail has overflowed its storage prior to head, so the elements are not
		// contiguous in the backing array and/or the array is filled to capacity.
		// The first region (1) spans from the first-in (head) element to the end of
		// the backing array, and the second region (2) spans from the start of the
		// array to the last-in (tail) element:
		//   (0123456789A) === Array index reference
		//   [xxT......Hx]     Elements in region 1 [9..A] and region 2 [0..1]
		//   [T......Hxxx]     Elements in region 1 [7..A] only
		//   [xxxxxHxxxxx]     Elements in region 1 [5..A] and region 2 [0..4]
		//   [Hxxxxxxxxxx]     Elements in region 1 [0..A] only
		// So we will potentially need to copy the elements in two phases:
		// (1.) Copy from head to the end of the backing array.
		n1, err1 := b.writeTo(w, int(ih), int(b.capt))
		// If the number of bytes written equals the backing array's capacity, then
		// b was filled to capacity and is now empty; nothing to copy in phase 2.
		if err1 != nil || n1 == int(b.capt) {
			return int64(n1), err1
		}
		// (2.) Copy from start of the backing array to tail.
		n2, err2 := b.writeTo(w, 0, int(it))
		return int64(n1 + n2), err2
	}
	// Elements form contiguous span in backing array from first-in (head) element
	// to last-in (tail) element:
	//   (0123456789A) === Array index reference
	//   [HxxxT......]     Elements forms contiguous span [0..3]
	//   [......HxxxT]     Elements forms contiguous span [6..9]
	nw, errw := b.writeTo(w, int(ih), int(it))
	return int64(nw), errw
}

// ReadByte returns the next unread byte from b and a nil error.
// If b is empty, returns 0, io.EOF.
//
// To avoid ambiguous validity of the returned byte, ReadByte will always return
// either a valid byte and nil error, or an invalid byte and non-nil error.
// In particular, ReadByte never returns a byte read along with error == io.EOF.
func (b *buf) ReadByte() (c byte, err error) {
	if b == nil || !b.valid {
		return 0, &nogc.ErrInvalidReceiver
	}
	h, t := b.head, b.tail
	if h == t {
		// Reading zero bytes from b (empty), return io.EOF.
		return 0, io.EOF
	}
	// Reading 1 byte from b, reduce length by 1.
	b.head++
	// Return the byte from original head position.
	return b.Byte[h%b.capt], nil
}

// UnreadByte causes the next call to ReadByte to return the last byte read.
// If the last operation was not a successful call to ReadByte, UnreadByte will
// unread the last byte read.
func (b *buf) UnreadByte() error {
	if b == nil || !b.valid {
		return &nogc.ErrInvalidReceiver
	}
	if b.head > 0 {
		b.head--
	}
	return nil
}

// WriteByte appends c to b and returns nil.
// If b is full, returns ErrWriteOverflow.
func (b *buf) WriteByte(c byte) (err error) {
	if b == nil || !b.valid {
		return &nogc.ErrInvalidReceiver
	}
	h, t := b.head, b.tail
	ih, it := h%b.capt, t%b.capt
	// If the array indices are equal, with head not eqaul to tail, then the
	// backing array is filled to capacity. We have nowhere to store the byte.
	// We can either discard head or retain it and return an error. Opting for
	// the latter so that no byte is lost, and it gives the caller an opportunity
	// to remedy the situation.
	if h != t && ih == it {
		return &nogc.ErrWriteOverflow
	}
	// Write the byte into tail position and increment length by 1.
	b.Byte[it] = c
	b.tail++
	return nil
}
