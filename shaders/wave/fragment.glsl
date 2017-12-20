#version 330 core

uniform float Time;

in vec2 pos;
out float waveHeight;

void main() {
    float rt = 10 * Time;
    float r  = 10 * length(pos);

    float v = float(abs(rt - r) < 0.3);

    waveHeight = v * pow(2.7, -2 * r) * sin((rt - r));
}
