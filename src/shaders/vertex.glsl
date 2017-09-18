#version 330 core

layout (location = 0) in vec2 in_position;

void main()
{
    gl_Position.xy = in_position.xy;
    gl_Position.z = 1.0;
    gl_Position.w = 1.0;
}