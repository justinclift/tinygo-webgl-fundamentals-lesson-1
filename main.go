package main

// TinyGo version of the 1st WebGL Fundamentals lesson
// https://webglfundamentals.org/webgl/lessons/webgl-fundamentals.html

import "syscall/js"

var (
	gl      js.Value
	glTypes GLTypes
)

type GLTypes struct {
	arrayBuffer    js.Value
	colorBufferBit js.Value
	compileStatus  js.Value
	float          js.Value
	fragmentShader js.Value
	linkStatus     js.Value
	staticDraw     js.Value
	triangles      js.Value
	vertexShader   js.Value
}

func (types *GLTypes) New() {
	types.arrayBuffer = gl.Get("ARRAY_BUFFER")
	types.colorBufferBit = gl.Get("COLOR_BUFFER_BIT")
	types.compileStatus = gl.Get("COMPILE_STATUS")
	types.float = gl.Get("FLOAT")
	types.fragmentShader = gl.Get("FRAGMENT_SHADER")
	types.linkStatus = gl.Get("LINK_STATUS")
	types.staticDraw = gl.Get("STATIC_DRAW")
	types.triangles = gl.Get("TRIANGLES")
	types.vertexShader = gl.Get("VERTEX_SHADER")
}

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

	// Initialise the GL types
	glTypes.New()

	// * WebGL initialisation code *

	// Create GLSL shaders, upload the GLSL source, compile the shaders
	vertexShader := createShader(gl, glTypes.vertexShader, vertCode)
	fragmentShader := createShader(gl, glTypes.fragmentShader, fragCode)

	// Link the two shaders into a program
	program := createProgram(gl, vertexShader, fragmentShader)

	// Look up where the vertex data needs to go
	positionAttributeLocation := gl.Call("getAttribLocation", program, "a_position")

	// Create a buffer and put three 2d clip space points in it
	positionBuffer := gl.Call("createBuffer", glTypes.arrayBuffer)

	// Bind it to ARRAY_BUFFER (think of it as ARRAY_BUFFER = positionBuffer)
	gl.Call("bindBuffer", glTypes.arrayBuffer, positionBuffer)

	// Three 2d points
	positionsNative := []float32{
		0, 0,
		0, 0.5,
		0.7, 0,
	}
	positions := js.TypedArrayOf(positionsNative)
	gl.Call("bufferData", glTypes.arrayBuffer, positions, glTypes.staticDraw)

	// * WebGL rendering code *

	// Tell WebGL how to convert from clip space to pixels
	gl.Call("viewport", 0, 0, width, height)

	// Clear the canvas
	gl.Call("clearColor", 0, 0, 0, 0)
	gl.Call("clear", glTypes.colorBufferBit)

	// Tell it to use our program (pair of shaders)
	gl.Call("useProgram", program)

	// Turn on the attribute
	gl.Call("enableVertexAttribArray", positionAttributeLocation)

	// Bind the position buffer
	gl.Call("bindBuffer", glTypes.arrayBuffer, positionBuffer)

	// Tell the attribute how to get data out of positionBuffer (ARRAY_BUFFER)
	pbSize := 2             // 2 components per iteration
	pbType := glTypes.float // the data is 32bit floats
	pbNormalize := false    // don't normalize the data
	pbStride := 0           // 0 = move forward size * sizeof(pbType) each iteration to get the next position
	pbOffset := 0           // start at the beginning of the buffer
	gl.Call("vertexAttribPointer", positionAttributeLocation, pbSize, pbType, pbNormalize, pbStride, pbOffset)

	// Draw
	primType := glTypes.triangles
	primOffset := 0
	primCount := 3
	gl.Call("drawArrays", primType, primOffset, primCount)
}

func createShader(gl js.Value, shaderType js.Value, source string) js.Value {
	shader := gl.Call("createShader", shaderType)
	gl.Call("shaderSource", shader, source)
	gl.Call("compileShader", shader)
	success := gl.Call("getShaderParameter", shader, glTypes.compileStatus).Bool()
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
	success := gl.Call("getProgramParameter", program, glTypes.linkStatus).Bool()
	if success {
		return program
	}
	println(gl.Call("getProgramInfoLog", program).String())
	gl.Call("deleteProgram", program)
	return js.Value{}
}
