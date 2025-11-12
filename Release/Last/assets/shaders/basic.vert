#version 330 core
layout(location = 0) in vec3 aPos;
layout(location = 1) in vec3 aNormal;
layout(location = 2) in vec2 aTex;

uniform mat4 uModel;
uniform mat4 uView;
uniform mat4 uProjection;
uniform mat3 uNormalMatrix;
uniform mat4 uLightSpaceMatrix;

out VS_OUT {
    vec3 FragPos;
    vec3 Normal;
    vec2 TexCoord;
    vec4 FragPosLightSpace;
} vs_out;

void main() {
    vec4 worldPos = uModel * vec4(aPos, 1.0);
    vs_out.FragPos = worldPos.xyz;
    vs_out.Normal = normalize(uNormalMatrix * aNormal);
    vs_out.TexCoord = aTex;
    vs_out.FragPosLightSpace = uLightSpaceMatrix * worldPos;
    gl_Position = uProjection * uView * worldPos;
}
