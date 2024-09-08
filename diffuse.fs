#version 330

// Input vertex attributes (from vertex shader)
in vec2 fragTexCoord;
in vec4 fragColor;

// Input uniform values
uniform sampler2D texture0;
uniform vec4 colDiffuse;

// Output fragment color
out vec4 finalColor;

// NOTE: Add here your custom variables

// NOTE: Render size values must be passed from code
const float renderWidth = 1024;
const float renderHeight = 1024;

float offset[3] = float[](0.0, 1.3846153846, 3.2307692308);
float weight[3] = float[](0.2270270270, 0.3162162162, 0.0702702703);

void main()
{
    // Texel color fetching from texture sampler
    vec3 texelColor = texture(texture0, fragTexCoord).rgb*weight[0];

    texelColor += texture(texture0, fragTexCoord + vec2(offset[1])/renderWidth, 0.0).rgb*weight[1];
    texelColor += texture(texture0, fragTexCoord - vec2(offset[1])/renderWidth, 0.0).rgb*weight[1];

    texelColor += texture(texture0, fragTexCoord + vec2(offset[2])/renderWidth, 0.0).rgb*weight[2];
    texelColor += texture(texture0, fragTexCoord - vec2(offset[2])/renderWidth, 0.0).rgb*weight[2];


    finalColor = vec4(texelColor, 0.5);
}
