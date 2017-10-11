#version 130

uniform vec4 torch;
uniform mat4 camera;
uniform mat4 rotate;

in vec3 inColor;
in vec3 inPos;

out vec3 color;
out vec3 pos;

void main()
{
    gl_Position = camera * rotate * vec4(inPos, 1.0);
    color = inColor;
    pos = (rotate * vec4(inPos, 1.0)).xyz;
}