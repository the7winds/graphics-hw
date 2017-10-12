//#version 330 core
#version 130

uniform mat4 Camera;
uniform mat4 M;

in vec3 inPos;
out vec3 pos;

void main()
{
    gl_Position = Camera * M * vec4(inPos, 1.0);
    pos = (M * vec4(inPos, 1.0)).xyz;
}