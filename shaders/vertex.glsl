#version 130

uniform mat4 camera;
uniform mat4 rotate;

in vec3 inColor;
in vec3 inPos;

out vec3 color;

void main()
{
    gl_Position = camera * rotate * vec4(inPos, 1.0);
    color = inColor;
}