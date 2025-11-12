#version 330 core
in vec3 vFragPos;
in vec3 vNormal;
in vec2 vTex;

out vec4 FragColor;

uniform sampler2D uTex;
uniform vec3 uLightPos;
uniform vec3 uLightColor;
uniform vec3 uViewPos;

void main() {
    vec3 texColor = texture(uTex, vTex).rgb;

    vec3 ambient = 0.1 * uLightColor;

    vec3 norm = normalize(vNormal);
    vec3 lightDir = normalize(uLightPos - vFragPos);
    float diff = max(dot(norm, lightDir), 0.0);
    vec3 diffuse = diff * uLightColor;

    vec3 viewDir = normalize(uViewPos - vFragPos);
    vec3 reflectDir = reflect(-lightDir, norm);
    float spec = pow(max(dot(viewDir, reflectDir), 0.0), 32.0);
    vec3 specular = 0.5 * spec * uLightColor;

    vec3 lighting = (ambient + diffuse + specular) * texColor;
    FragColor = vec4(lighting, 1.0);
}
