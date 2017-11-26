#version 330 core

uniform mat4 PV;
uniform mat4 M;

in vec3 Vertex;

out vec3 npos;
out vec3 mpos;

void main() {
    vec4 mp = M * vec4(Vertex, 1);
    vec4 pp = PV * mp;
    gl_Position = pp;
    
    npos = pp.xyz / pp.w;
    mpos = mp.xyz;
}
