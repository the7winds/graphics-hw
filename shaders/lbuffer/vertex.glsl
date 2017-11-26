#version 330 core

uniform mat4 PV;
uniform mat4 M;

in vec3 Vertex;

out vec3 mpos;
out vec4 ppos;

void main() {
    vec4 mp = M * vec4(Vertex, 1);
    
    ppos = PV * mp;
    mpos = mp.xyz;

    gl_Position = ppos;
}
