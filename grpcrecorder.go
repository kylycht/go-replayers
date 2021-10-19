package grpcreplay

import (
	"bufio"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/golang/protobuf/proto"
	pb "github.com/kylycht/go-replayers/grpcreplay/proto/grpcreplay"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// A Recorder records RPCs for later playback.
type Recorder struct {
	opts *RecorderOptions
	mu   sync.Mutex
	w    *bufio.Writer
	f    *os.File
	next int
	err  error
}

// RecorderOptions are options for a Recorder.
type RecorderOptions struct {
	// The initial state, stored in the output file for retrieval during replay.
	Initial []byte

	// A function that can inspect and modify requests and responses written to the
	// replay file. It does not modify messages sent to the service.
	//
	// The function is called with the method name and the message for each unary
	// request and response, before the messages are written to the replay file. If
	// the function returns an error, the error will be returned to the client.
	// Streaming RPCs are not supported.
	BeforeWrite func(method string, msg proto.Message) error
}

// NewRecorder creates a recorder that writes to filename.
// You must call Close on the Recorder to ensure that all data is written.
func NewRecorder(filename string, opts *RecorderOptions) (*Recorder, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	rec, err := NewRecorderWriter(f, opts)
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	rec.f = f
	return rec, nil
}

// NewRecorderWriter creates a recorder that writes to w.
// You must call Close on the Recorder to ensure that all data is written.
func NewRecorderWriter(w io.Writer, opts *RecorderOptions) (*Recorder, error) {
	if opts == nil {
		opts = &RecorderOptions{}
	}
	bw := bufio.NewWriter(w)
	if err := writeHeader(bw, opts.Initial); err != nil {
		return nil, err
	}
	return &Recorder{w: bw, opts: opts, next: 1}, nil
}

// DialOptions returns the options that must be passed to grpc.Dial
// to enable recording.
func (r *Recorder) DialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithUnaryInterceptor(r.interceptUnary),
		grpc.WithStreamInterceptor(r.interceptStream),
	}
}

// Close saves any unwritten information.
func (r *Recorder) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return r.err
	}
	err := r.w.Flush()
	if r.f != nil {
		if err2 := r.f.Close(); err == nil {
			err = err2
		}
	}
	return err
}

// Intercepts all unary (non-stream) RPCs.
func (r *Recorder) interceptUnary(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ereq := &entry{
		kind:   pb.Entry_REQUEST,
		method: method,
		msg:    message{msg: proto.Clone(req.(proto.Message))},
	}
	data, err := proto.Marshal(ereq.msg.msg)
	if err != nil {
		logrus.Error(err)
	} else {
		ereq.msg.checksum = fmt.Sprintf("%x", md5.Sum(data))
	}

	if r.opts.BeforeWrite != nil {
		if err := r.opts.BeforeWrite(method, ereq.msg.msg); err != nil {
			return err
		}
	}
	checksum, err := r.writeEntry(ereq)
	if err != nil {
		return err
	}
	ierr := invoker(ctx, method, req, res, cc, opts...)
	eres := &entry{
		kind:     pb.Entry_RESPONSE,
		checksum: checksum,
	}
	// If the error is not a gRPC status, then something more
	// serious is wrong. More significantly, we have no way
	// of serializing an arbitrary error. So just return it
	// without recording the response.
	if _, ok := status.FromError(ierr); !ok {
		r.mu.Lock()
		r.err = fmt.Errorf("saw non-status error in %s response: %v (%T)", method, ierr, ierr)
		r.mu.Unlock()
		return ierr
	}
	eres.msg.set(proto.Clone(res.(proto.Message)), ierr)
	if r.opts.BeforeWrite != nil {
		if err := r.opts.BeforeWrite(method, eres.msg.msg); err != nil {
			return err
		}
	}
	if _, err := r.writeEntry(eres); err != nil {
		return err
	}
	return ierr
}

func (r *Recorder) writeEntry(e *entry) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		return "", r.err
	}
	err := writeEntry(r.w, e)
	if err != nil {
		r.err = err
		return "", err
	}

	return e.checksum, nil
}

func (r *Recorder) interceptStream(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	cstream, serr := streamer(ctx, desc, cc, method, opts...)
	e := &entry{
		kind:   pb.Entry_CREATE_STREAM,
		method: method,
	}
	e.msg.set(nil, serr)
	checksum, err := r.writeEntry(e)
	if err != nil {
		return nil, err
	}
	return &recClientStream{
		ctx:      ctx,
		rec:      r,
		cstream:  cstream,
		checksum: checksum,
	}, serr
}

// A recClientStream implements the gprc.ClientStream interface.
// It behaves exactly like the default ClientStream, but also
// records all messages sent and received.
type recClientStream struct {
	ctx      context.Context
	rec      *Recorder
	cstream  grpc.ClientStream
	checksum string
}

func (rcs *recClientStream) Context() context.Context { return rcs.ctx }

func (rcs *recClientStream) SendMsg(m interface{}) error {
	serr := rcs.cstream.SendMsg(m)
	e := &entry{
		kind:     pb.Entry_SEND,
		checksum: rcs.checksum,
	}
	e.msg.set(m, serr)
	if _, err := rcs.rec.writeEntry(e); err != nil {
		return err
	}
	return serr
}

func (rcs *recClientStream) RecvMsg(m interface{}) error {
	serr := rcs.cstream.RecvMsg(m)
	e := &entry{
		kind:     pb.Entry_RECV,
		checksum: rcs.checksum,
	}
	e.msg.set(m, serr)
	if _, err := rcs.rec.writeEntry(e); err != nil {
		return err
	}
	return serr
}

func (rcs *recClientStream) Header() (metadata.MD, error) {
	// TODO(jba): record.
	return rcs.cstream.Header()
}

func (rcs *recClientStream) Trailer() metadata.MD {
	// TODO(jba): record.
	return rcs.cstream.Trailer()
}

func (rcs *recClientStream) CloseSend() error {
	// TODO(jba): record.
	return rcs.cstream.CloseSend()
}
