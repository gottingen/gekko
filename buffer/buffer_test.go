package buffer

import (
	"bytes"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestBufferReadFrom(t *testing.T) {
	prefix := "foobar"
	expectedS := "asadfsdafsadfasdfisdsdfa"
	prefixLen := int64(len(prefix))
	expectedN := int64(len(expectedS))

	var bb Buffer
	bb.WriteString(prefix)

	rf := (io.ReaderFrom)(&bb)
	for i := 0; i < 20; i++ {
		r := bytes.NewBufferString(expectedS)
		n, err := rf.ReadFrom(r)
		if n != expectedN {
			t.Fatalf("unexpected n=%d. Expecting %d. iteration %d", n, expectedN, i)
		}
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		bbLen := int64(bb.Len())
		expectedLen := prefixLen + int64(i+1)*expectedN
		if bbLen != expectedLen {
			t.Fatalf("unexpected byteBuffer length: %d. Expecting %d", bbLen, expectedLen)
		}
		for j := 0; j < i; j++ {
			start := prefixLen + int64(j)*expectedN
			b := bb.B[start : start+expectedN]
			if string(b) != expectedS {
				t.Fatalf("unexpected byteBuffer contents: %q. Expecting %q", b, expectedS)
			}
		}
	}
}

func TestBufferWriteTo(t *testing.T) {
	expectedS := "foobarbaz"
	var bb Buffer
	bb.WriteString(expectedS[:3])
	bb.WriteString(expectedS[3:])

	wt := (io.WriterTo)(&bb)
	var w bytes.Buffer
	for i := 0; i < 10; i++ {
		n, err := wt.WriteTo(&w)
		if n != int64(len(expectedS)) {
			t.Fatalf("unexpected n returned from WriteTo: %d. Expecting %d", n, len(expectedS))
		}
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		s := string(w.Bytes())
		if s != expectedS {
			t.Fatalf("unexpected string written %q. Expecting %q", s, expectedS)
		}
		w.Reset()
	}
}

func TestBufferGetPutSerial(t *testing.T) {
	testBufferGetPut(t)
}

func TestBufferGetPutConcurrent(t *testing.T) {
	concurrency := 10
	ch := make(chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			testBufferGetPut(t)
			ch <- struct{}{}
		}()
	}

	for i := 0; i < concurrency; i++ {
		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Fatalf("timeout!")
		}
	}
}

func testBufferGetPut(t *testing.T) {
	for i := 0; i < 10; i++ {
		expectedS := fmt.Sprintf("num %d", i)
		b := Get()
		b.B = append(b.B, "num "...)
		b.B = append(b.B, fmt.Sprintf("%d", i)...)
		if string(b.B) != expectedS {
			t.Fatalf("unexpected result: %q. Expecting %q", b.B, expectedS)
		}
		Put(b)
	}
}

func testBufferGetString(t *testing.T) {
	for i := 0; i < 10; i++ {
		expectedS := fmt.Sprintf("num %d", i)
		b := Get()
		b.SetString(expectedS)
		if b.String() != expectedS {
			t.Fatalf("unexpected result: %q. Expecting %q", b.B, expectedS)
		}
		Put(b)
	}
}

func TestBufferGetStringSerial(t *testing.T) {
	testBufferGetString(t)
}

func TestBufferGetStringConcurrent(t *testing.T) {
	concurrency := 10
	ch := make(chan struct{}, concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			testBufferGetString(t)
			ch <- struct{}{}
		}()
	}

	for i := 0; i < concurrency; i++ {
		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Fatalf("timeout!")
		}
	}
}
