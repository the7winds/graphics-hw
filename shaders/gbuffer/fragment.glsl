#version 330 core

uniform mat4 PV;
uniform vec4 Color;

in vec3 pos;

out vec3 texColor;
out vec3 texNorma;

void main() {
    texColor = Color.rgb;
    texNorma = (normalize(cross(dFdx(pos), dFdy(pos))) + 1) / 2;
}