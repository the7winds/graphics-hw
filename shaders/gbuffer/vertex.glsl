#version 330 core

uniform mat4 PV;
uniform mat4 M;

in vec3 Vertex;
in vec3 Norm;

out vec3 pos;
out vec3 norm;

void main() {
    gl_Position = PV * M * vec4(Vertex, 1);
    pos = (M * vec4(Vertex, 1)).xyz;
    norm = Norm;
}
