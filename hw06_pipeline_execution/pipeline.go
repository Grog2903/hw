package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	current := in

	for _, stage := range stages {
		current = runStage(current, done, stage)
	}

	return current
}

func runStage(in In, done In, stage Stage) Out {
	out := make(chan interface{})

	go func() {
		defer close(out)

		stageOut := stage(in)

		for {
			select {
			case <-done:
				return
			case v, ok := <-stageOut:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()

	return out
}
