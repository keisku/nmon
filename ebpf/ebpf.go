package ebpf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"golang.org/x/exp/slog"
)

// $BPF_CLANG and $BPF_CFLAGS are set by the Makefile.
//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cflags $BPF_CFLAGS -no-global-types -type event bpf ./c/bpf_prog.c -- -I./c

var objs bpfObjects

func Start() (func(), error) {
	if err := loadBpfObjects(&objs, nil); err != nil {
		return nil, fmt.Errorf("can't load bpf: %w", err)
	}
	tcpConnect, err := link.AttachTracing(link.TracingOptions{
		Program: objs.TcpConnect,
	})
	if err != nil {
		return nil, fmt.Errorf("can't attach tracing: %w", err)
	}
	rd, err := ringbuf.NewReader(objs.bpfMaps.Events)
	if err != nil {
		return nil, fmt.Errorf("can't create ringbuf reader: %w", err)
	}
	go func() {
		// bpfEvent is generated by bpf2go.
		var event bpfEvent
		for {
			record, err := rd.Read()
			if err != nil {
				if errors.Is(err, ringbuf.ErrClosed) {
					slog.Info("exiting ringbuf reader...")
					return
				}
				slog.Warn("reading from reader", slog.Any("error", err))
				continue
			}

			// Parse the ringbuf event entry into a bpfEvent structure.
			if err := binary.Read(bytes.NewBuffer(record.RawSample), binary.BigEndian, &event); err != nil {
				slog.Warn("parsing ringbuf event", slog.Any("error", err))
				continue
			}

			slog.Info("%-16s %-15s %-6d -> %-15s %-6d",
				event.Comm,
				intToIP(event.Saddr),
				event.Sport,
				intToIP(event.Daddr),
				event.Dport,
			)
		}
	}()
	linkTracingOptions := []link.TracingOptions{
		{Program: objs.TcpClose},
		{Program: objs.TcpCloseExit},
		{Program: objs.TcpRecvmsgExit},
		{Program: objs.TcpRetransmitSkb},
		{Program: objs.TcpRetransmitSkbExit},
	}
	links := make([]link.Link, len(linkTracingOptions))
	for i, opt := range linkTracingOptions {
		links[i], err = link.AttachTracing(opt)
		if err != nil {
			return nil, fmt.Errorf("can't attach tracing: %w", err)
		}
	}
	return func() {
		if err := objs.Close(); err != nil {
			slog.Warn("can't close bpf objects", slog.Any("error", err))
		}
		if err := tcpConnect.Close(); err != nil {
			slog.Warn("can't close tracing", slog.Any("error", err))
		}
		if err := rd.Close(); err != nil {
			slog.Warn("can't close ringbuf reader", slog.Any("error", err))
		}
		for i := range links {
			if err := links[i].Close(); err != nil {
				slog.Warn("can't close tracing", slog.Any("error", err))
			}
		}
	}, nil
}

// intToIP converts IPv4 number to net.IP
func intToIP(ipNum uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, ipNum)
	return ip
}
