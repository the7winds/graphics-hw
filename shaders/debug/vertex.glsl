#version 330 core

in vec3 inPos;
out vec2 pos;

void main()
{
    gl_Position = vec4(inPos, 1);
    pos = inPos.xy;
    pos = (pos + vec2(1, 1)) / 0.5;
}