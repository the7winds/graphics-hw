#version 130

uniform vec3 eye;
uniform vec4 torch;
uniform vec4 color;

in vec3 pos;

void main()
{
    vec3 n = normalize(cross(dFdx(pos), dFdy(pos)));
    vec3 l = normalize(torch.xyz - pos);

    float diffuse = clamp(dot(n, l), 0, 1);

    vec3 rl = reflect(-l, n);
    vec3 er = normalize(eye - pos);
    float spectacular = clamp(dot(rl, er), 0, 1);

    gl_FragColor = (diffuse + spectacular) * color;
}