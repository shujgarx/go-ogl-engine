#version 330 core
in vec2 vTex;
out vec4 FragColor;

uniform sampler2D uTex;

void main() {
    FragColor = texture(uTex, vTex);
}
