#version 330 core

in vec2 tpos;

uniform int Mode;

uniform sampler2D TexRender;
uniform sampler2D TexColor;
uniform sampler2D TexNorma;
uniform sampler2D TexHeight;

out vec4 fragColor;

void main() {
    if (Mode == 0) {
        fragColor = texture(TexRender, tpos);
    } else if (Mode == 1) {
        fragColor = texture(TexColor, tpos);
    } else if (Mode == 2) {
        fragColor = texture(TexNorma, tpos);
    } else if (Mode == 3) {
        fragColor = texture(TexHeight, tpos);
    }
}