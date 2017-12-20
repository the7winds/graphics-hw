package program

import "log"
import "os"
import "bufio"

import "github.com/go-gl/gl/v4.1-core/gl"

const logBufferSize = 512

var logBuffer [logBufferSize]byte

func loadShader(filename string, shaderType uint32) uint32 {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatalln("can't open:", filename)
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	source := ""
	for scanner.Scan() {
		source += scanner.Text() + "\n"
	}
	source += "\x00"

	shaderID := gl.CreateShader(shaderType)

	cSources, free := gl.Strs(source)
	gl.ShaderSource(shaderID, 1, cSources, nil)
	gl.CompileShader(shaderID)
	free()

	var success int32
	gl.GetShaderiv(shaderID, gl.COMPILE_STATUS, &success)

	if success != gl.TRUE {
		gl.GetShaderInfoLog(shaderID, logBufferSize, nil, &logBuffer[0])
		log.Fatalln("compile error in", filename, gl.GoStr(&logBuffer[0]))
		panic("can't compile shader")
	}

	return shaderID
}

func New(vertex string, fragment string) uint32 {
	programID := gl.CreateProgram()

	vertexID := loadShader(vertex, gl.VERTEX_SHADER)
	fragmentID := loadShader(fragment, gl.FRAGMENT_SHADER)

	gl.AttachShader(programID, vertexID)
	gl.AttachShader(programID, fragmentID)
	gl.LinkProgram(programID)

	var success int32
	gl.GetProgramiv(programID, gl.LINK_STATUS, &success)

	if success != gl.TRUE {
		gl.GetProgramInfoLog(programID, logBufferSize, nil, &logBuffer[0])
		log.Fatalln("link error:", gl.GoStr(&logBuffer[0]))
		panic("link error")
	}

	return programID
}
