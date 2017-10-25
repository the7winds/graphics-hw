#version 330 core

uniform mat4 PV;
uniform mat4 M;

in vec3 Vertex;

out vec3 gpos;
out vec2 tpos;

void main() {
    gl_Position = PV * M * vec4(Vertex, 1);
    gpos = gl_Position.xyz;
    tpos = 0.5 * (gl_Position.xy + vec2(1));
}
