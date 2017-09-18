#version 330

out vec3 out_color;

uniform int iterations;
uniform float reZ;
uniform float imZ;
uniform float zoom;
uniform float shiftX;
uniform float shiftY;

vec3 hsv2rgb(vec3 c)
{
    vec4 K = vec4(1.0, 2.0 / 3.0, 1.0 / 3.0, 3.0);
    vec3 p = abs(fract(c.xxx + K.xyz) * 6.0 - K.www);
    return c.z * mix(K.xxx, clamp(p - K.xxx, 0.0, 1.0), c.y);
}

void main()
{
    float newRe = (gl_FragCoord.x - shiftX - 400) / (zoom * 800);
    float newIm = (gl_FragCoord.y + shiftY - 400) / (zoom * 800);

    int i;
    for (i = 0; i < iterations; i++) {
        float oldRe = newRe;
        float oldIm = newIm;

        newRe = oldRe * oldRe - oldIm * oldIm - reZ;
        newIm = 2 * oldRe * oldIm + imZ;

        if (newRe * newRe + newIm * newIm > 4) {
            break;
        }
    }

    float c = float(i) / float(iterations);

    out_color = hsv2rgb(vec3(c, 1, 1));
}