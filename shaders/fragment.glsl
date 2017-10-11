#version 130

uniform vec4 torch;

in vec3 color;
in vec3 pos;

void main()
{
    vec3 n = normalize(cross(dFdx(pos), dFdy(pos)));
    vec3 r =  normalize(torch.xyz - pos);
    float k = dot(n, r);
    float s = float(k > 0);
    gl_FragColor = s * k * vec4(color, 1.0);
}