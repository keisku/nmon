// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64 || amd64p32 || arm || arm64 || loong64 || mips64le || mips64p32le || mipsle || ppc64le || riscv64

package ebpf

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

// loadBpf returns the embedded CollectionSpec for bpf.
func loadBpf() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_BpfBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load bpf: %w", err)
	}

	return spec, err
}

// loadBpfObjects loads bpf and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*bpfObjects
//	*bpfPrograms
//	*bpfMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadBpfObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadBpf()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// bpfSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpfSpecs struct {
	bpfProgramSpecs
	bpfMapSpecs
}

// bpfSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpfProgramSpecs struct {
	TcpClose             *ebpf.ProgramSpec `ebpf:"tcp_close"`
	TcpCloseExit         *ebpf.ProgramSpec `ebpf:"tcp_close_exit"`
	TcpRecvmsgExit       *ebpf.ProgramSpec `ebpf:"tcp_recvmsg_exit"`
	TcpRetransmitSkb     *ebpf.ProgramSpec `ebpf:"tcp_retransmit_skb"`
	TcpRetransmitSkbExit *ebpf.ProgramSpec `ebpf:"tcp_retransmit_skb_exit"`
}

// bpfMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type bpfMapSpecs struct {
	ConnCloseBatch          *ebpf.MapSpec `ebpf:"conn_close_batch"`
	ConnCloseEvent          *ebpf.MapSpec `ebpf:"conn_close_event"`
	ConnStats               *ebpf.MapSpec `ebpf:"conn_stats"`
	PendingTcpRetransmitSkb *ebpf.MapSpec `ebpf:"pending_tcp_retransmit_skb"`
	PortBindings            *ebpf.MapSpec `ebpf:"port_bindings"`
	TcpOngoingConnectPid    *ebpf.MapSpec `ebpf:"tcp_ongoing_connect_pid"`
	TcpRetransmits          *ebpf.MapSpec `ebpf:"tcp_retransmits"`
	TcpStats                *ebpf.MapSpec `ebpf:"tcp_stats"`
	UdpPortBindings         *ebpf.MapSpec `ebpf:"udp_port_bindings"`
}

// bpfObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpfObjects struct {
	bpfPrograms
	bpfMaps
}

func (o *bpfObjects) Close() error {
	return _BpfClose(
		&o.bpfPrograms,
		&o.bpfMaps,
	)
}

// bpfMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpfMaps struct {
	ConnCloseBatch          *ebpf.Map `ebpf:"conn_close_batch"`
	ConnCloseEvent          *ebpf.Map `ebpf:"conn_close_event"`
	ConnStats               *ebpf.Map `ebpf:"conn_stats"`
	PendingTcpRetransmitSkb *ebpf.Map `ebpf:"pending_tcp_retransmit_skb"`
	PortBindings            *ebpf.Map `ebpf:"port_bindings"`
	TcpOngoingConnectPid    *ebpf.Map `ebpf:"tcp_ongoing_connect_pid"`
	TcpRetransmits          *ebpf.Map `ebpf:"tcp_retransmits"`
	TcpStats                *ebpf.Map `ebpf:"tcp_stats"`
	UdpPortBindings         *ebpf.Map `ebpf:"udp_port_bindings"`
}

func (m *bpfMaps) Close() error {
	return _BpfClose(
		m.ConnCloseBatch,
		m.ConnCloseEvent,
		m.ConnStats,
		m.PendingTcpRetransmitSkb,
		m.PortBindings,
		m.TcpOngoingConnectPid,
		m.TcpRetransmits,
		m.TcpStats,
		m.UdpPortBindings,
	)
}

// bpfPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadBpfObjects or ebpf.CollectionSpec.LoadAndAssign.
type bpfPrograms struct {
	TcpClose             *ebpf.Program `ebpf:"tcp_close"`
	TcpCloseExit         *ebpf.Program `ebpf:"tcp_close_exit"`
	TcpRecvmsgExit       *ebpf.Program `ebpf:"tcp_recvmsg_exit"`
	TcpRetransmitSkb     *ebpf.Program `ebpf:"tcp_retransmit_skb"`
	TcpRetransmitSkbExit *ebpf.Program `ebpf:"tcp_retransmit_skb_exit"`
}

func (p *bpfPrograms) Close() error {
	return _BpfClose(
		p.TcpClose,
		p.TcpCloseExit,
		p.TcpRecvmsgExit,
		p.TcpRetransmitSkb,
		p.TcpRetransmitSkbExit,
	)
}

func _BpfClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed bpf_bpfel.o
var _BpfBytes []byte
