// Package bgroutine is for managing the background routines
package bgroutine

import (
	"fmt"
	"strconv"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/techrail/bark/appRuntime"

	"github.com/techrail/ground/channels"
	"github.com/techrail/ground/constants"
	"github.com/techrail/ground/typs/appError"
)

const (
	StateUninitialized = "uninitialized"
	StateInitialized   = "initialized"
	StatePaused        = "paused"
	StateRunning       = "running"
	StateTerminated    = "terminated"
)

const (
	CronMode   = "cronMode"   // When we are working in the cron mode
	TickerMode = "tickerMode" // When we are working in the ticker mode
)

type Typ struct {
	Name            string                 // Name of the routine
	done            chan bool              // Send to this channel to stop the routine
	operationMode   string                 // Which more are we working in?
	cronExpr        string                 // The Cron expression to run the function repeatedly
	ticker          *time.Ticker           // Time ticker to call the function in case we are in ticker mode
	schedule        cron.Schedule          // Schedule using which we call the function when we are in cron mode
	function        func() appError.Typ    // The function to run on each tick
	state           string                 // What is the state of this routine
	instanceRunning bool                   // Is the function already running (used to prevent parallel runs only)
	monitorHook     func(typ appError.Typ) // The function which is called for each run
}

func (m *Manager) AddRoutine(name string, cronExpression string, runnerFunc func() appError.Typ) appError.Typ {
	mode := CronMode
	tickr := time.NewTicker(87600 * time.Hour) // 10 years default duration for the routine
	// check if we have a valid cron expression or not
	s, err := cron.ParseStandard(cronExpression)
	if err != nil {
		// The expression is not in the standard format. Let's check if we can convert this into an integer
		tickerMills, convErr := strconv.Atoi(cronExpression)
		if convErr != nil {
			// Can't convert it to integer either. It's definitely an error
			return appError.NewError(appError.Error, "1NIV49", fmt.Sprintf("Could not parse cron expression %v for routine %v. Parser error: %v and Atoi error: %v", cronExpression, name, err, convErr))
		}
		// Looks like the expression is that of milliseconds
		mode = TickerMode
		tickr = time.NewTicker(time.Duration(tickerMills) * time.Millisecond)
	}

	r := Typ{
		Name:            name,
		done:            make(chan bool),
		operationMode:   mode,
		cronExpr:        cronExpression,
		ticker:          tickr,
		schedule:        s,
		function:        runnerFunc,
		state:           StateInitialized,
		instanceRunning: false,
		monitorHook:     nil,
	}

	if _, ok := m.routineMap[name]; ok {
		// already exists
		return appError.NewError(appError.Error, "1NCFF9", "Another routine by that name already exists")
	}
	m.routineMap[name] = &r

	return appError.BlankError
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
		if !r.instanceRunning {
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
				if r.operationMode == TickerMode {
					if r.state != StateRunning {
						e := appError.NewError(appError.Info, "1NI933", fmt.Sprintf("Tick for Routine %s was received at %v but the routine is %v.", r.Name, t, r.state))
						if r.monitorHook != nil {
							r.monitorHook(e)
						}
						channels.ErrorTypChan <- e
					} else {
						e := appError.NewError(appError.Info, "1NI9P9", fmt.Sprintf("Tick for routine %v at %v", r.Name, t))
						if r.monitorHook != nil {
							r.monitorHook(e)
						}
						channels.ErrorTypChan <- e
						runRoutineOnce()
					}
				} else {
					e := appError.NewError(appError.Info, "1NPBCQ", fmt.Sprintf("Tick for routine %v recieved when it should not have happened at %v", r.Name, t))
					if r.monitorHook != nil {
						r.monitorHook(e)
					}
					channels.ErrorTypChan <- e
				}
			case t := <-time.After(time.Second):
				if appRuntime.ShutdownRequested.Load() {
					r.done <- true
					e := appError.NewError(appError.Info, "1NPB3F", fmt.Sprintf("Shutdown was requested. Stopping routine %v at %v", r.Name, t))
					if r.monitorHook != nil {
						r.monitorHook(e)
					}
					channels.ErrorTypChan <- e
					return
				}

				if r.operationMode == CronMode && time.Now().UTC().After(r.schedule.Next(time.Now().UTC())) {
					// Time to execute
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

func (r *Typ) CurrentMode() string {
	return r.operationMode
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
