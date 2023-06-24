package server

type EventKey int32

const (
	EK_STAGE_SWITCH = 1000

)

func (r *Robot) fireEvent(k EventKey) {
	for _, e := range r.executors {
		e.onEvent(k)
	}
}