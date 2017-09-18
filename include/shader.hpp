#pragma once

#include <iostream>
#include <fstream>

namespace my
{
GLuint compileShader(GLenum type, const char *shaderSourceName)
{
    GLuint shaderId = glCreateShader(type);

    std::string source;
    std::string line;
    std::ifstream input(shaderSourceName);

    while (getline(input, line))
    {
        source += line + "\n";
    }

    const char *cstr_source = source.c_str();

    glShaderSource(shaderId, 1, &cstr_source, nullptr);

    GLint success;
    constexpr size_t logsize = 512;
    char log[logsize];

    glCompileShader(shaderId);
    glGetShaderiv(shaderId, GL_COMPILE_STATUS, &success);

    if (!success)
    {
        glGetShaderInfoLog(shaderId, logsize, nullptr, log);
        std::cerr << "[" << shaderSourceName << "]: " << log << '\n';
        exit(1);
    }

    return shaderId;
}

GLuint loadProgram()
{
    GLuint programId = glCreateProgram();
    GLuint vs = compileShader(GL_VERTEX_SHADER, VERTEX_SHADER_PATH);
    GLuint fs = compileShader(GL_FRAGMENT_SHADER, FRAGMENT_SHADER_PATH);

    glAttachShader(programId, vs);
    glAttachShader(programId, fs);
    glLinkProgram(programId);
    glDetachShader(programId, vs);
    glDetachShader(programId, fs);
    glDeleteShader(vs);
    glDeleteShader(fs);

    GLint success;
    constexpr size_t logsize = 512;
    char log[logsize];

    glGetProgramiv(programId, GL_LINK_STATUS, &success);

    if (!success)
    {
        glGetProgramInfoLog(programId, logsize, nullptr, log);
        std::cerr << "[SHADER PROGRAM LINK]: " << log << '\n';
        exit(1);
    }

    return programId;
}
}