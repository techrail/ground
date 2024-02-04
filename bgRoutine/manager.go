package bgRoutine

type Manager struct {
	routineMap map[string]*Typ
}

func NewManager() Manager {
	return Manager{
		routineMap: make(map[string]*Typ),
	}
}

func (m *Manager) ShutdownAllRoutines() {
	for _, v := range m.routineMap {
		v.Stop()
	}
}
