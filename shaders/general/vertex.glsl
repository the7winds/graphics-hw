#version 330 core

uniform mat4 Camera;
uniform mat4 ShadowCamera;
uniform mat4 M;

in vec3 inPos;

out vec3 pos;
out vec3 shadowPos;

void main()
{
    pos = (M * vec4(inPos, 1.0)).xyz;
    shadowPos = (ShadowCamera * M * vec4(inPos, 1.0)).xyz;
    gl_Position = Camera * M * vec4(inPos, 1.0);
}