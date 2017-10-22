#version 330 core

uniform sampler2D shadow;
in vec2 pos;

void main()
{
    float c = texture(shadow, pos).r;
    gl_FragColor = vec4(c, c, c, 1);
}