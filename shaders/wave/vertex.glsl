#version 330 core

uniform mat4 M;

in vec3 Vertex;

out vec2 pos;

void main() {
    gl_Position = M * vec4(Vertex, 1);
    pos = Vertex.xy;
}
