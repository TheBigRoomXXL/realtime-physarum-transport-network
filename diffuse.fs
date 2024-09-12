#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;

// Input uniform values
uniform sampler2D texture0;

// Output fragment color
out vec4 finalColor;

const float renderSize = 1024;
const float texelSize = 1 / renderSize;

float weight[2] = float[](0.85, 0.03);

void main()
{
    // Texel color fetching from texture sampler
    vec3 texelColor = texture(texture0, fragTexCoord).rgb*weight[0];

    texelColor += texture(texture0, fragTexCoord + vec2(texelSize, 0.0)).rgb*weight[1];
    texelColor += texture(texture0, fragTexCoord + vec2(texelSize, 0.0)).rgb*weight[1];
    texelColor += texture(texture0, fragTexCoord + vec2(0.0, texelSize)).rgb*weight[1];
    texelColor += texture(texture0, fragTexCoord - vec2(0.0, texelSize)).rgb*weight[1];

    finalColor = vec4(texelColor, 0.5);
}
