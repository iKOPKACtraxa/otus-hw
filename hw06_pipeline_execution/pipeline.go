package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// chanCloser sticks done-chan with in-chan,
// so if done-chan is closed the in chan is closing.
func chanCloser(in In, done In) In {
	if done != nil {
		newIn := make(Bi)
		go func() {
			defer close(newIn)
			var val interface{}
			for {
				select {
				case <-done:
					return
				default:
					select {
					case <-done:
						return
					case val = <-in:
						select {
						case <-done:
							return
						case newIn <- val:
						}
					}
				}
			}
		}()
		return newIn
	}
	return in
}

// ExecutePipeline sends in-chan of function as in-chan
// for i stage, then it gets out-chan from i stage and
// sends it as in-chan for i+1 stage. So one by one
// it connects all stages and last stage returns
// last out-chan (as in).
// If done-chan is not nil and it closed, it is a signal
// to stop a pipeline, so all stages are getting stopped.
// If in-chan of function is nill it panics.
// If there is no stages at input in-chan returns as out-chan.
func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		panic("Input channel is nil")
	}
	for _, stage := range stages {
		in = stage(chanCloser(in, done))
	}
	return in
}
