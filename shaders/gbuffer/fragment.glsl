#version 330 core

uniform mat4 PV;
uniform vec4 Color;
uniform bool NormaPassed;

in vec3 pos;
in vec3 norm;

out vec3 texColor;
out vec3 texNorma;

void main() {
    texColor = Color.rgb;
    vec3 n1 = (normalize(cross(dFdx(pos), dFdy(pos))) + 1) / 2;
    vec3 n2 = (norm + 1) / 2;
    texNorma = int(!NormaPassed) * n1 + int(NormaPassed) * n2; 
}