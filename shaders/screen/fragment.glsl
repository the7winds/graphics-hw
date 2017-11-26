#version 330 core

in vec2 tpos;

uniform int Mode;

uniform sampler2D TexColor;
uniform sampler2D TexNorma;
uniform sampler2D TexDepth;
uniform sampler2D TexLight;
uniform sampler2D TexPos;

out vec4 fragColor;

void main() {
    if (Mode == 0) {
        fragColor = texture(TexLight, tpos);
    } else if (Mode == 1) {
        fragColor = texture(TexColor, tpos);
    } else if (Mode == 2) {
        fragColor = texture(TexNorma, tpos);
    } else if (Mode == 3) {
        float depth = texture(TexDepth, tpos).x;
        vec3 depthc = vec3(depth);
        fragColor = vec4(depthc, 1);
    }
}