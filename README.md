This is the 1st lesson of the WebGL Fundamentals code, ported from Javascript
to TinyGo:

&nbsp; &nbsp; https://webglfundamentals.org/

Running demo:

&nbsp; &nbsp; https://justinclift.github.io/tinygo-webgl-fundamentals-lessons1/

To compile the WebAssembly file:

    $ tinygo build -target wasm -no-debug -panic trap -o docs/wasm.wasm main.go

To strip the custom name section from the end (reducing file size further):

    $ wasm2wat docs/wasm.wasm -o docs/wasm.wat
    $ wat2wasm docs/wasm.wat -o docs/wasm.wasm
    $ rm -f docs/wasm.wat

