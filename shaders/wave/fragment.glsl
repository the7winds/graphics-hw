#version 330 core

uniform float Time;

in vec2 pos;
out float waveHeight;

void main() {
    float rt = Time;
    float r  = length(pos);
    float d  = rt - r;

    float v = float(abs(d) < 0.05);

    waveHeight = v * pow(2.7, -30 * r) * sin(62.8 * d);
}
