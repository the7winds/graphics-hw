#version 330 core

in vec3 pos;
out float depth;

void main()
{
    depth = pos.z;
}