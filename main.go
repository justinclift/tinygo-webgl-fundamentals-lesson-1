package main

// TinyGo version of the 1st WebGL Fundamentals lesson
// https://webglfundamentals.org/webgl/lessons/webgl-fundamentals.html

import (
	"syscall/js"

	"github.com/justinclift/webgl"
)

const (
	// https://developer.mozilla.org/en-US/docs/Web/API/WebGL_API/Constants
	ARRAY_BUFFER     = 0x8892
	COLOR_BUFFER_BIT = 0x00004000
	COMPILE_STATUS   = 0x8B81
	FLOAT            = 0x1406
	FRAGMENT_SHADER  = 0x8B30
	LINK_STATUS      = 0x8B82
	STATIC_DRAW      = 0x88E4
	TRIANGLES        = 0x0004
	VERTEX_SHADER    = 0x8B31
)

var (
	// Vertex shader source code
	vertCode = `
	// an attribute will receive data from a buffer
	attribute vec4 a_position;

	// all shaders have a main function
	void main() {
		// gl_Position is a special variable a vertex shader
		// is responsible for setting
		gl_Position = a_position;
	}`

	// Fragment shader source code
	fragCode = `
	// fragment shaders don't have a default precision so we need
	// to pick one. mediump is a good default. It means "medium precision"
	precision mediump float;

	void main() {
		// gl_FragColor is a special variable a fragment shader
		// is responsible for setting
		gl_FragColor = vec4(1, 0, 0.5, 1); // return redish-purple
	}`
)

func main() {
	// Set up WebGL context
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "mycanvas")
	width := canvas.Get("clientWidth").Int()
	height := canvas.Get("clientHeight").Int()
	canvas.Call("setAttribute", "width", width)
	canvas.Call("setAttribute", "height", height)

	// Set up the WebGL context attributes we want to use
	// https://developer.mozilla.org/en-US/docs/Web/API/HTMLCanvasElement/getContext (search for "WebGL context
	// attributes" on the page)
	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false

	// Create the WebGL context
	gl, err := webgl.NewContext(&canvas, attrs)
	if err != nil {
		js.Global().Call("alert", "Error: "+err.Error())
		return
	}

	// * WebGL initialisation code *

	// Create GLSL shaders, upload the GLSL source, compile the shaders
	vertexShader := createShader(gl, VERTEX_SHADER, vertCode)
	fragmentShader := createShader(gl, FRAGMENT_SHADER, fragCode)

	// Link the two shaders into a program
	program := createProgram(gl, vertexShader, fragmentShader)

	// Look up where the vertex data needs to go
	positionAttributeLocation := gl.GetAttribLocation(program, "a_position")

	// Create a buffer and put three 2d clip space points in it
	positionBuffer := gl.CreateArrayBuffer()
	// positionBuffer := gl.Call("createBuffer", ARRAY_BUFFER)

	// Bind it to ARRAY_BUFFER (think of it as ARRAY_BUFFER = positionBuffer)
	gl.BindBuffer(ARRAY_BUFFER, positionBuffer)

	// Three 2d points
	positionsNative := []float32{
		0, 0,
		0, 0.5,
		0.7, 0,
	}
	positions := js.TypedArrayOf(positionsNative)
	gl.BufferData(ARRAY_BUFFER, positions, STATIC_DRAW)

	// * WebGL rendering code *

	// Tell WebGL how to convert from clip space to pixels
	gl.Viewport(0, 0, width, height)

	// Clear the canvas
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(COLOR_BUFFER_BIT)

	// Tell it to use our program (pair of shaders)
	gl.UseProgram(program)

	// Turn on the attribute
	gl.EnableVertexAttribArray(positionAttributeLocation)

	// Bind the position buffer
	gl.BindBuffer(ARRAY_BUFFER, positionBuffer)

	// Tell the attribute how to get data out of positionBuffer (ARRAY_BUFFER)
	pbSize := 2          // 2 components per iteration
	pbType := FLOAT      // the data is 32bit floats
	pbNormalize := false // don't normalize the data
	pbStride := 0        // 0 = move forward size * sizeof(pbType) each iteration to get the next position
	pbOffset := 0        // start at the beginning of the buffer
	gl.VertexAttribPointer(positionAttributeLocation, pbSize, pbType, pbNormalize, pbStride, pbOffset)

	// Draw
	primType := TRIANGLES
	primOffset := 0
	primCount := 3
	gl.DrawArrays(primType, primOffset, primCount)
}

func createShader(gl *webgl.Context, shaderType int, source string) *js.Value {
	shader := gl.CreateShader(shaderType)
	gl.ShaderSource(shader, source)
	gl.CompileShader(shader)
	success := gl.GetShaderParameter(shader, COMPILE_STATUS).Bool()
	if success {
		return shader
	}
	println(gl.GetShaderInfoLog(shader))
	gl.DeleteShader(shader)
	return &js.Value{}
}

func createProgram(gl *webgl.Context, vertexShader *js.Value, fragmentShader *js.Value) *js.Value {
	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)
	success := gl.GetProgramParameterb(program, LINK_STATUS)
	if success {
		return program
	}
	println(gl.GetProgramInfoLog(program))
	gl.DeleteProgram(program)
	return &js.Value{}
}
