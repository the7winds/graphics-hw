#version 330 core

in vec3 Vertex;
out vec2 tpos;

void main() {
    gl_Position = vec4(Vertex, 1);
    tpos = 0.5 * (gl_Position.xy + vec2(1));
}
