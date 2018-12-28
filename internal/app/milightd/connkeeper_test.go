package milightd

import (
	"testing"
	"time"
)

type TestConnectionManagerer struct {
	lightController  *TestLightController
	allocateCalls    int
	releaseCalls     int
	isAllocatedCalls int
	terminateCalls   int
	isAllocated      bool
	isCreated        bool
}

func NewTestConnectionManagerer() *TestConnectionManagerer {
	return &TestConnectionManagerer{
		lightController: &TestLightController{},
	}
}

func (m *TestConnectionManagerer) Allocate() (LightController, error) {
	m.allocateCalls++
	m.isAllocated = true
	m.isCreated = true
	return m.lightController, nil
}

func (m *TestConnectionManagerer) Release() {
	m.releaseCalls++
	m.isAllocated = false
}

func (m *TestConnectionManagerer) GetStatus() (bool, bool) {
	m.isAllocatedCalls++
	return m.isAllocated, m.isCreated
}

func (m *TestConnectionManagerer) Terminate() {
	m.terminateCalls++
	m.isAllocated = false
	m.isCreated = false
}

func TestSingleAllocation(t *testing.T) {
	checkPeriod := 3 * time.Second
	connman := NewTestConnectionManagerer()
	keeper := NewConnectionKeeper(connman, checkPeriod)
	keeper.Allocate()
	time.Sleep(1 * time.Second)
	keeper.Release()
	keeper.Terminate()

	if connman.terminateCalls != 1 {
		t.Errorf("expected %d terminations, got %d", 1, connman.terminateCalls)
	}

	if connman.allocateCalls != 1 {
		t.Errorf("expected %d allocations, got %d", 1, connman.allocateCalls)
	}

	if connman.releaseCalls != 1 {
		t.Errorf("expected %d releases, got %d", 1, connman.releaseCalls)
	}
}

func TestManyAllocations(t *testing.T) {
	checkPeriod := 2 * time.Second
	connman := NewTestConnectionManagerer()
	keeper := NewConnectionKeeper(connman, checkPeriod)

	for i := 0; i < 6; i++ {
		keeper.Allocate()
		time.Sleep(1 * time.Second)
		keeper.Release()
	}

	keeper.Terminate()

	if connman.terminateCalls != 1 {
		t.Errorf("expected %d terminations, got %d", 1, connman.terminateCalls)
	}

	if connman.allocateCalls != 6 {
		t.Errorf("expected %d allocations, got %d", 6, connman.allocateCalls)
	}

	if connman.releaseCalls != 6 {
		t.Errorf("expected %d releases, got %d", 6, connman.releaseCalls)
	}
}

func TestTerminations(t *testing.T) {
	checkPeriod := 2 * time.Second
	connman := NewTestConnectionManagerer()
	keeper := NewConnectionKeeper(connman, checkPeriod)

	keeper.Allocate()
	time.Sleep(1 * time.Second)
	keeper.Release()

	time.Sleep(6 * time.Second)

	keeper.Terminate()

	if connman.terminateCalls != 2 {
		t.Errorf("expected %d terminations, got %d", 2, connman.terminateCalls)
	}

	if connman.allocateCalls != 1 {
		t.Errorf("expected %d allocations, got %d", 1, connman.allocateCalls)
	}

	if connman.releaseCalls != 1 {
		t.Errorf("expected %d releases, got %d", 1, connman.releaseCalls)
	}
}
