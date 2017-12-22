#version 330 core

uniform mat4 PV;
uniform sampler2D TexHeight;
uniform samplerCube TexEnv;
uniform vec3 CameraPos;

in vec3 pos;
in vec2 tpos;

out vec3 texRender;
out vec3 texColor;
out vec3 texNorma;

vec3 evalNorma() {
    float step = 0.001;
    float h0 = texture(TexHeight, tpos).r;
    float h1 = texture(TexHeight, vec2(tpos.x, tpos.y + step)).r;
    float h2 = texture(TexHeight, vec2(tpos.x, tpos.y - step)).r;
    float h3 = texture(TexHeight, vec2(tpos.x + step, tpos.y)).r;
    float h4 = texture(TexHeight, vec2(tpos.x - step, tpos.y)).r;

    vec3 v1 = normalize(vec3(0, h1 - h0, step));
    vec3 v2 = normalize(vec3(0, h2 - h0, -step));
    vec3 v3 = normalize(vec3(step, h3 - h0, 0));
    vec3 v4 = normalize(vec3(-step, h4 - h0, 0));
    
    vec3 n1 = cross(v1, v3);
    vec3 n2 = cross(v4, v1);
    vec3 n3 = cross(v2, v4);
    vec3 n4 = cross(v3, v2);

    vec3 n = normalize((n1 + n2 + n3 + n4) / 4);
    return n;
}

vec3 light() {
    vec3 c = vec3(0, 2, 1);
    vec3 n = evalNorma();
    vec3 l = normalize(c - pos);
    return dot(l, n) * vec3(1);
}

void main() {
    vec3 n = evalNorma();
    vec3 l = reflect(pos - CameraPos, n);
    texRender = texture(TexEnv, l).rgb;
    texNorma  = n;
    texColor  = light();
}