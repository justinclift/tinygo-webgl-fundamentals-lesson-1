package main

// TinyGo version of the 1st WebGL Fundamentals lesson
// https://webglfundamentals.org/webgl/lessons/webgl-fundamentals.html

import "syscall/js"

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
	gl js.Value

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
	canvasEl := doc.Call("getElementById", "mycanvas")
	width := canvasEl.Get("clientWidth").Int()
	height := canvasEl.Get("clientHeight").Int()
	canvasEl.Call("setAttribute", "width", width)
	canvasEl.Call("setAttribute", "height", height)
	gl = canvasEl.Call("getContext", "webgl")
	if gl == js.Undefined() {
		// No support
		println("Could not create WebGL context.  Seems unsupported.")
		return
	}

	// * WebGL initialisation code *

	// Create GLSL shaders, upload the GLSL source, compile the shaders
	vertexShader := createShader(gl, VERTEX_SHADER, vertCode)
	fragmentShader := createShader(gl, FRAGMENT_SHADER, fragCode)

	// Link the two shaders into a program
	program := createProgram(gl, vertexShader, fragmentShader)

	// Look up where the vertex data needs to go
	positionAttributeLocation := gl.Call("getAttribLocation", program, "a_position")

	// Create a buffer and put three 2d clip space points in it
	positionBuffer := gl.Call("createBuffer", ARRAY_BUFFER)

	// Bind it to ARRAY_BUFFER (think of it as ARRAY_BUFFER = positionBuffer)
	gl.Call("bindBuffer", ARRAY_BUFFER, positionBuffer)

	// Three 2d points
	positionsNative := []float32{
		0, 0,
		0, 0.5,
		0.7, 0,
	}
	positions := js.TypedArrayOf(positionsNative)
	gl.Call("bufferData", ARRAY_BUFFER, positions, STATIC_DRAW)

	// * WebGL rendering code *

	// Tell WebGL how to convert from clip space to pixels
	gl.Call("viewport", 0, 0, width, height)

	// Clear the canvas
	gl.Call("clearColor", 0, 0, 0, 0)
	gl.Call("clear", COLOR_BUFFER_BIT)

	// Tell it to use our program (pair of shaders)
	gl.Call("useProgram", program)

	// Turn on the attribute
	gl.Call("enableVertexAttribArray", positionAttributeLocation)

	// Bind the position buffer
	gl.Call("bindBuffer", ARRAY_BUFFER, positionBuffer)

	// Tell the attribute how to get data out of positionBuffer (ARRAY_BUFFER)
	pbSize := 2          // 2 components per iteration
	pbType := FLOAT      // the data is 32bit floats
	pbNormalize := false // don't normalize the data
	pbStride := 0        // 0 = move forward size * sizeof(pbType) each iteration to get the next position
	pbOffset := 0        // start at the beginning of the buffer
	gl.Call("vertexAttribPointer", positionAttributeLocation, pbSize, pbType, pbNormalize, pbStride, pbOffset)

	// Draw
	primType := TRIANGLES
	primOffset := 0
	primCount := 3
	gl.Call("drawArrays", primType, primOffset, primCount)
}

func createShader(gl js.Value, shaderType uint, source string) js.Value {
	shader := gl.Call("createShader", shaderType)
	gl.Call("shaderSource", shader, source)
	gl.Call("compileShader", shader)
	success := gl.Call("getShaderParameter", shader, COMPILE_STATUS).Bool()
	if success {
		return shader
	}
	println(gl.Call("getShaderInfoLog", shader).String())
	gl.Call("deleteShader", shader)
	return js.Value{}
}

func createProgram(gl js.Value, vertexShader js.Value, fragmentShader js.Value) js.Value {
	program := gl.Call("createProgram")
	gl.Call("attachShader", program, vertexShader)
	gl.Call("attachShader", program, fragmentShader)
	gl.Call("linkProgram", program)
	success := gl.Call("getProgramParameter", program, LINK_STATUS).Bool()
	if success {
		return program
	}
	println(gl.Call("getProgramInfoLog", program).String())
	gl.Call("deleteProgram", program)
	return js.Value{}
}
