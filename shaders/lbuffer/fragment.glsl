#version 330 core

uniform vec4 Color;
uniform vec3 Center;
uniform mat4 InvPV;

uniform sampler2D TexColor;
uniform sampler2D TexNorma;
uniform sampler2D TexDepth;

in vec3 gpos;
in vec2 tpos;

out vec4 fragColor;

vec3 getRealPos() {
    float depth = texture(TexDepth, tpos).x;
    float z = 2 * depth - 1;
    vec3 pos = vec3(gpos.xy, z);
    return (InvPV * vec4(pos, 1)).xyz;
}

vec3 getRealNormal() {
    vec3 n = texture(TexNorma, tpos).xyz;
    return 2 * n - vec3(1);
}

void main() {
    vec3 rpos = getRealPos();

    vec3 n = getRealNormal();

    vec3 l = normalize(Center - rpos);
    float c = clamp(dot(n, l), 0, 1);

    float d = 2 * texture(TexDepth, tpos).x - 1;
    float vis = 0;
    if (d - gpos.z < 0) {
        vis = 1;
    }

    fragColor = vis * c * texture(TexColor, tpos);
}