#version 330 core

uniform vec4 Color;
uniform vec3 Center;
uniform mat4 InvPV;
uniform mat4 PV;

uniform sampler2D TexColor;
uniform sampler2D TexNorma;
uniform sampler2D TexDepth;

in vec3 mpos;
in vec3 npos;

out vec4 fragColor;

vec3 getRealPos() {
    vec2 tpos = (npos.xy + 1) / 2;
    float z = texture(TexDepth, tpos).x;
    vec4 pos = vec4(tpos.xy, z, 1);
    return (InvPV * pos).xyz;
}

vec3 getRealNormal() {
    vec2 tpos = (npos.xy + 1) / 2;
    vec3 n = texture(TexNorma,  tpos).xyz;
    return 2 * n - 1;
}

void main() {
    vec3 rpos = getRealPos();

    vec3 n = getRealNormal();

    vec3 l = normalize(Center - rpos);
    float c = clamp(dot(n, l), 0, 1);

    float lr = length(rpos - Center);
    float ra = length(mpos - Center);

    float vis = float(lr < ra) * c / pow(lr, 2);

    fragColor = vis * texture(TexColor, (npos.xy + 1) / 2);
    fragColor = float(lr < ra) * texture(TexColor, (npos.xy + 1) / 2);
}