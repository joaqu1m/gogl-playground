package engine

import "github.com/go-gl/gl/v4.1-core/gl"

// RenderPass represents a single rendering pass with a target framebuffer.
type RenderPass struct {
	Name        string
	Framebuffer uint32 // 0 = default framebuffer
	DrawFunc    func()
}

// Execute binds the framebuffer and runs the pass.
func (p *RenderPass) Execute() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, p.Framebuffer)
	p.DrawFunc()
}
