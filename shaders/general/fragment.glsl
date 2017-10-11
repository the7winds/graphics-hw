#version 130

uniform vec3 eye;
uniform vec4 torch;
uniform vec4 color;
uniform sampler2DShadow shadow;

in vec3 pos;

void main()
{
    vec3 n = normalize(cross(dFdx(pos), dFdy(pos)));
    vec3 l = normalize(torch.xyz - pos);

    float diffuse = clamp(dot(n, l), 0, 1);

    vec3 rl = reflect(-l, n);
    vec3 er = normalize(eye - pos);
    float spectacular = clamp(dot(rl, er), 0, 1);

    float power = 20 / pow(length(torch.xyz - pos), 2);

    gl_FragColor = power * (diffuse + spectacular) * color;
}