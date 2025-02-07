package pu2

type EventType int

const (
	EventNone EventType = iota
	EventVanish
	EventFall
	EventHandle
	EventOjamaFall
)

type SoloRecord struct {
	Yama Yama
	Step int
	Tind int

	stepType []EventType
	hHistory []Handle // zero value if stepType[i] != EventHandle
	fHistory []Field  // psudo formula: f[i+1] = f[i] + hHistory[i] if stepType[i] == EventHandle else f[i].Next()
}

func NewSoloRecord(yama Yama) SoloRecord {
	fh := make([]Field, 1, 100)
	fh[0] = NewField()
	return SoloRecord{
		Yama:     yama,
		hHistory: make([]Handle, 0, 100),
		fHistory: fh,
		Step:     0,
	}
}

// Field returns the copied current field.
func (r *SoloRecord) Field() Field {
	return r.fHistory[r.Step]
}

func (r *SoloRecord) StepType() EventType {
	return r.stepType[r.Step]
}

// Push pushes the handle to hHistory[step] and fHistory[step+1]. no validation.
func (r *SoloRecord) Push(h Handle) {
	tsumo := r.Yama.Get(r.Tind)
	f := r.Field()
	f.AddHandle(tsumo, h)
	r.stepType = append(r.stepType[:r.Step], EventHandle)
	r.hHistory = append(r.hHistory[:r.Step], h)
	r.fHistory = append(r.fHistory[:r.Step+1], f)
	r.Step++
	r.Tind++
}

// rollback
func (r *SoloRecord) Pop() {
	r.stepType = r.stepType[:r.Step-1]
	r.hHistory = r.hHistory[:r.Step-1]
	r.fHistory = r.fHistory[:r.Step]
	r.Step--
	r.Tind--
}

// Vanish seeks states and returns EventVanish if possible.
func (r *SoloRecord) Vanish() bool {
	f := r.Field()
	if ok, _ := f.Vanish(); ok {
		r.stepType = append(r.stepType[:r.Step], EventVanish)
		r.hHistory = append(r.hHistory[:r.Step], Handle{}) // as zero
		r.fHistory = append(r.fHistory[:r.Step+1], f)
		r.Step++
		return true
	}
	return false
}

// Fall seeks states and returns EventVanish if possible.
func (r *SoloRecord) Fall() bool {
	f := r.Field()
	if f.Fall() {
		r.stepType = append(r.stepType[:r.Step], EventFall)
		r.hHistory = append(r.hHistory[:r.Step], Handle{}) // as zero
		r.fHistory = append(r.fHistory[:r.Step+1], f)
		r.Step++
		return true
	}
	return false
}

func (r *SoloRecord) StepBack() EventType {
	if r.Step == 0 {
		return EventNone
	}
	s := r.stepType[r.Step-1]
	if s == EventHandle {
		r.Tind--
	}
	r.Step--
	return s
}

func (r *SoloRecord) StepForward() EventType {
	if r.Step >= len(r.fHistory)-1 {
		return EventNone
	}
	s := r.stepType[r.Step]
	if s == EventHandle {
		r.Tind++
	}
	r.Step++
	return s
}

// Undo rollback until r.Step == 0 or r.stepType[r.Step] must be EventHandle.
// returns true if find EventHandle.
func (r *SoloRecord) Undo() bool {
	for {
		s := r.StepBack()
		if s == EventNone {
			return false
		}
		if s == EventHandle {
			return true
		}
	}
}

func (r *SoloRecord) Redo() bool {
	for {
		s := r.StepForward()
		if s == EventNone {
			return false
		}
		if s == EventHandle {
			return true
		}
	}
}
