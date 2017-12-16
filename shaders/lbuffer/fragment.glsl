#version 330 core

uniform vec4 Color;
uniform vec3 Center;
uniform mat4 PV;
uniform float R;

uniform sampler2D TexColor;
uniform sampler2D TexNorma;
uniform sampler2D TexDepth;

out vec4 fragColor;
out vec4 volumesColor;

void main() {
    vec2 tpos = gl_FragCoord.xy / 800;
    float z = (2 * texture(TexDepth, tpos).x - 1);
    vec4 rpos = inverse(PV) * vec4(2 * tpos.xy - 1, z, 1);
    rpos /= rpos.w;
    vec3 n = 2 * texture(TexNorma, tpos).xyz - 1;

    vec3 l = normalize(Center - rpos.xyz);
    float c = clamp(dot(n, l), 0, 1);

    float lr = length(rpos.xyz - Center);

    float intence = 0.05 * float(0.01 < R - lr) * c / pow(lr, 2);
    fragColor = 0.5 * Color * intence * texture(TexColor, tpos);
    volumesColor = 0.5 * float(0.01 < R - lr) * c / pow(lr, 2) + 0.5 * Color;
}