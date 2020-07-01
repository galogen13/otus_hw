package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	I   = interface{}
	In  = <-chan I
	Out = In
	Bi  = chan I
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in
	for _, stage := range stages {
		out = stage(out)
	}
	return takeAll(done, out)
}

func takeAll(done In, valueStream Out) Out {
	takeStream := make(Bi)
	go func() {
		defer close(takeStream)
		for {
			select {
			case <-done:
				return
			case value, ok := <-valueStream:
				select {
				case <-done:
					return
				default:
					if !ok {
						return
					}
					takeStream <- value
				}
			}
		}
	}()
	return takeStream
}
