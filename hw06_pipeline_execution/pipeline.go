package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func prepareChanelForStage(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer func() {
			close(out)
		}()

		for {
			select {
			case <-done:
				return
			case value, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- value:
				}
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	lastChanel := in

	for _, stage := range stages {
		lastChanel = stage(prepareChanelForStage(lastChanel, done))
	}

	return lastChanel
}
