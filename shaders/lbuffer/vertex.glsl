#version 330 core

uniform mat4 PV;
uniform mat4 M;

in vec3 Vertex;

void main() {
    gl_Position = PV * M * vec4(Vertex, 1);
}
