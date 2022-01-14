package static

import "fmt"

func (v *visitor) handleSinkDetection() bool {
	outputResult = make(map[string]bool)
	log.Debugf("sink arg map len: %d", len(v.sinkArgs))
	for callInstr,m := range v.sinkArgs {
		for arg,_ := range m{
			log.Debugf("sink arg: %s=%s", arg.Name(),arg.String())
			if tags, ok := v.lattice[arg]; ok {
				//todo: report detection
				output := fmt.Sprintf("sink here %s with tag:%s ", prog.Fset.Position(callInstr.Pos()),tags.String())
				outputResult[output] = true
				//return true
			}
		}

	}

	for o,_ := range outputResult {
		log.Warning(o)
		//out("sink here", prog.Fset.Position())
		//os.Stdout.WriteString(o)
		//os.Stdout.WriteString("\n")
	}

	return false
}

