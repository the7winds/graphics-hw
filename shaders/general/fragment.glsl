#version 330 core

uniform vec3 eye;

uniform vec3 lampPos;
uniform vec3 lampDir;
uniform float lampAngle;
uniform float lampPower;

uniform vec3 torchPos;
uniform float torchPower;

// object parameters
uniform vec4 colorD;
uniform vec4 colorS;
uniform float powerS;

uniform sampler2D shadow;

in vec3 pos;
in vec3 shadowPos;

void main()
{
    float d = texture(shadow, 0.5 * (1 + shadowPos.xy)).x;
    float K = 1.0;
    
    if (0.5 * (1  + shadowPos.z) - d > 0.001) {
        K = 0.5;
    }

    vec3 n = normalize(cross(dFdx(pos), dFdy(pos)));

    gl_FragColor = vec4(0);

    // lamp
    vec3 l = normalize(lampPos - pos);
    float power = max(int(dot(-l, lampDir) > lampAngle), 0.3) * lampPower / pow(length(lampPos - pos), 2);
    vec4 diffuse = colorD * clamp(dot(n, l), 0, 1) * power;

    gl_FragColor = diffuse;

    // torch
    l = normalize(torchPos - pos);
    power = torchPower / pow(length(torchPos - pos), 2);
    diffuse = colorD * clamp(dot(n, l), 0, 1) * power;

    gl_FragColor += diffuse;

    // spec
    vec3 rl = normalize(reflect(-l, n));
    vec3 er = normalize(eye - pos);
    vec4 spectacular = colorS * pow(clamp(dot(rl, er), 0, 1), powerS);

    gl_FragColor += spectacular;

    gl_FragColor *= K;
}