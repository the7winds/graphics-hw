#version 330 core

uniform vec4 Color;
uniform vec3 Center;
uniform mat4 PV;

uniform sampler2D TexColor;
uniform sampler2D TexNorma;
uniform sampler2D TexDepth;

in vec3 mpos;
in vec4 ppos;

out vec4 fragColor;

void main() {
    vec2 tpos = gl_FragCoord.xy / 800;
    float z = ppos.w * (2 * texture(TexDepth, tpos).x - 1);
    vec3 rpos = (inverse(PV) * vec4(ppos.xy, z, ppos.w)).xyz;
    vec3 n = 2 * texture(TexNorma, tpos).xyz - 1;
    
    vec3 l = normalize(Center - rpos);
    float c = clamp(dot(n, l), 0, 1);

    float lr = length(rpos - Center);
    float ra = length(mpos - Center);

    float intence = float(0.01 < ra - lr) * c / pow(lr, 2);
    fragColor = Color * intence * texture(TexColor, tpos);
}