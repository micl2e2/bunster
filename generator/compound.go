package generator

import (
	"github.com/yassinebenaid/bunster/ast"
	"github.com/yassinebenaid/bunster/ir"
)

func (g *generator) handleGroup(buf *InstructionBuffer, group ast.Group, pc *pipeContext) {
	var cmdbuf InstructionBuffer

	g.handleRedirections(&cmdbuf, "group", group.Redirections, pc, true)

	if pc == nil {
		for _, cmd := range group.Body {
			g.generate(&cmdbuf, cmd, nil)
		}
	} else {
		cmdbuf.add(ir.Literal("var done = make(chan struct{},1)"))
		cmdbuf.add(ir.Literal(`
			pipelineWaitgroup = append(pipelineWaitgroup, runtime.PiplineWaitgroupItem{
				Wait: func()error{
					<-done
					return nil
				},
			})
		`))

		var go_routing InstructionBuffer
		go_routing.add(ir.Literal("defer streamManager.Destroy()\n"))
		for _, cmd := range group.Body {
			g.generate(&go_routing, cmd, nil)
		}
		go_routing.add(ir.Literal("done<-struct{}{}\n"))
		cmdbuf.add(ir.Closure{
			Async: true,
			Body:  go_routing,
		})
	}

	*buf = append(*buf, ir.Closure{
		Body: cmdbuf,
	})
}