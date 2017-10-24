#version 330 core

in vec3 pos;
uniform vec4 Color;

void main() {
    vec3 lampPos = vec3(0,0,0);
    vec3 n = normalize(cross(dFdx(pos), dFdy(pos)));
    vec3 l = normalize(lampPos - pos);
    float power = 10 / pow(length(lampPos - pos), 2);
    gl_FragColor = Color * clamp(dot(n, l), 0, 1);
}