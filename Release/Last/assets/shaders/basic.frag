#version 330 core
out vec4 FragColor;

in VS_OUT {
    vec3 FragPos;
    vec3 Normal;
    vec2 TexCoord;
    vec4 FragPosLightSpace;
} fs_in;

uniform vec3 uLightPos;
uniform vec3 uLightColor;
uniform vec3 uLightDir;
uniform vec3 uViewPos;
uniform sampler2D uTex;
uniform sampler2D uShadowMap;

float ShadowCalculation(vec4 fragPosLightSpace, vec3 normal, vec3 lightDir) {
    vec3 projCoords = fragPosLightSpace.xyz / fragPosLightSpace.w;
    projCoords = projCoords * 0.5 + 0.5;
    if (projCoords.z > 1.0)
        return 0.0;
    float bias = max(0.0025 * (1.0 - dot(normal, lightDir)), 0.0005);
    float shadow = 0.0;
    vec2 texelSize = 1.0 / textureSize(uShadowMap, 0);
    for (int x = -1; x <= 1; ++x) {
        for (int y = -1; y <= 1; ++y) {
            float pcfDepth = texture(uShadowMap, projCoords.xy + vec2(x, y) * texelSize).r;
            float currentDepth = projCoords.z - bias;
            if (currentDepth > pcfDepth)
                shadow += 1.0;
        }
    }
    shadow /= 9.0;
    return shadow;
}

void main() {
    vec3 color = texture(uTex, fs_in.TexCoord).rgb;
    if (length(color) == 0.0) {
        color = vec3(1.0);
    }
    vec3 normal = normalize(fs_in.Normal);
    vec3 lightDir = normalize(-uLightDir);
    float diff = max(dot(normal, lightDir), 0.0);
    vec3 viewDir = normalize(uViewPos - fs_in.FragPos);
    vec3 halfwayDir = normalize(lightDir + viewDir);
    float spec = pow(max(dot(normal, halfwayDir), 0.0), 32.0);

    float shadow = ShadowCalculation(fs_in.FragPosLightSpace, normal, lightDir);
    vec3 ambient = 0.25 * color;
    vec3 lighting = (ambient + (1.0 - shadow) * (diff * color + spec * uLightColor)) * uLightColor;
    FragColor = vec4(lighting, 1.0);
}
