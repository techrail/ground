package bgRoutine

import (
	"fmt"
	"time"

	`github.com/techrail/ground/channels`
	`github.com/techrail/ground/constants`
	`github.com/techrail/ground/logger`
	`github.com/techrail/ground/typs/appError`
)

const (
	StateUninitialized = "uninitialized"
	StateInitialized   = "initialized"
	StatePaused        = "paused"
	StateRunning       = "running"
	StateTerminated    = "terminated"
)

const (
	EventRunResult = "runResult"
)

type Typ struct {
	Name            string                 // Name of the routine
	done            chan bool              // Send to this channel to stop the routine TODO: unexport it
	ticker          *time.Ticker           // Time ticker to call the routine Every day at 8:45 AM TODO: unexport it
	function        func() appError.Typ    // The function to run on each tick TODO: unexport it
	state           string                 // What is the state of this routine
	instanceRunning bool                   // Is the function already running (used to prevent parallel runs only)
	monitorHook     func(typ appError.Typ) // The function which is called for each run
}

var routineMap map[string]*Typ

func init() {
	routineMap = make(map[string]*Typ)
}

func addRoutine(name string, routine *Typ) appError.Typ {
	// Check if another routine already exists with the same name
	if _, ok := routineMap[name]; ok {
		// already exists
		return appError.NewError(appError.Error, "1NCFF9", "Another routine by that name already exists")
	}
	routineMap[name] = routine
	return appError.BlankError
}

func New(name string, tickerDuration time.Duration, runnerFunc func() appError.Typ) (*Typ, appError.Typ) {
	r := Typ{
		Name:            name,
		done:            make(chan bool),
		ticker:          time.NewTicker(tickerDuration),
		function:        runnerFunc,
		state:           StateInitialized,
		instanceRunning: false,
		monitorHook:     nil,
	}
	errt := addRoutine(name, &r)
	if errt.IsNotBlank() {
		return nil, errt
	}
	return &r, errt
}

func (r *Typ) AddMonitorFunc(f func(typ appError.Typ)) {
	r.monitorHook = f
}

func (r *Typ) Start(launchRightNow bool) appError.Typ {
	if r.Name == constants.EmptyString {
		return appError.NewError(appError.Error, "1NCFFT", "Cannot launch nameless routine")
	}
	if r.state == StatePaused {
		return appError.NewError(
			appError.Error, "1NCFFY", fmt.Sprintf("Routine %v is paused. You can resume it, but not start it.", r.Name))
	}

	if r.state == StateRunning {
		return appError.NewError(
			appError.Error, "1NCFG1", fmt.Sprintf("Routine %v is already running.", r.Name))
	}

	if r.state == StateTerminated {
		return appError.NewError(
			appError.Error, "1NCFG4", fmt.Sprintf("Routine %v has been terminated. It cannot be started again.", r.Name))
	}

	if r.state == StateUninitialized {
		return appError.NewError(
			appError.Error, "1NCFG7", fmt.Sprintf("Routine %v has not been initialized yet.", r.Name))
	}

	if r.state == StateInitialized {
		// We are supposed to start the routine
		r.state = StateRunning
	}

	// The actual function that launches the routine function (once)
	runRoutineOnce := func() {
		if r.instanceRunning == false {
			r.instanceRunning = true
			err := r.function()
			if err != appError.BlankError {
				e := appError.NewError(appError.Error, "1NCHIL", fmt.Sprintf("function for routine %v could not run: %v", r.Name, err))
				if r.monitorHook != nil {
					r.monitorHook(e)
				}
				channels.ErrorTypChan <- e
				r.instanceRunning = false
			} else {
				e := appError.NewError(appError.Info, "1NCHKC", fmt.Sprintf("function for routine %v finished running", r.Name))
				if r.monitorHook != nil {
					r.monitorHook(e)
				}
				channels.ErrorTypChan <- e
			}
			r.instanceRunning = false
		} else {
			e := appError.NewError(appError.Notice, "1NCHML", fmt.Sprintf("function for routine %v seems to be running already", r.Name))
			if r.monitorHook != nil {
				r.monitorHook(e)
			}
			channels.ErrorTypChan <- e
		}
	}

	// goroutine to keep running the routine function periodically and to stop doing that when done channel gets data
	// 	indicating that the routine is no longer needed
	go func() {
		for {
			select {
			case <-r.done:
				e := appError.NewError(appError.Info, "1NCFGS", fmt.Sprintf("Routine %v shutting down.", r.Name))
				if r.monitorHook != nil {
					r.monitorHook(e)
				}
				channels.ErrorTypChan <- e
				r.ticker.Stop()
				return
			case t := <-r.ticker.C:
				if r.state != StateRunning {
					logger.Println(fmt.Sprintf("I#1NCFGY - Tick for Routine %s was received at %v but the routine is %v.", r.Name, t, r.state))
				} else {
					channels.ErrorTypChan <- fmt.Sprintf("I#1NCFJF - Tick for routine %v at %v", r.Name, t)
					runRoutineOnce()
				}
			}
		}
	}()

	if launchRightNow {
		go runRoutineOnce()
	}

	return appError.BlankError
}

func (r *Typ) Pause() {
	if r.state == StateRunning || r.state == StateInitialized {
		r.state = StatePaused
	}
}

func (r *Typ) Resume() {
	if r.state == StatePaused {
		r.state = StateRunning
	}
}

func (r *Typ) Stop() {
	r.done <- true
}

func (r *Typ) IsTerminated() bool {
	if r.state == StateTerminated {
		return true
	}
	return false
}

func (r *Typ) IsPaused() bool {
	if r.state == StatePaused {
		return true
	}
	return false
}

func (r *Typ) IsRunning() bool {
	if r.state == StateRunning {
		return true
	}
	return false
}

func (r *Typ) GetCurrentState() string {
	return r.state
}

func ShutdownAllRoutines() {
	for _, v := range routineMap {
		v.Stop()
	}
}
