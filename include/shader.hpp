#pragma once

#include <iostream>
#include <fstream>
#include <string>
#include <stdlib.h>

#include <boost/filesystem.hpp>

namespace fs = boost::filesystem;

GLuint compileShader(GLenum type, const std::string &shaderSourceName)
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

    fs::path shaders(SHADERS_DIR);

    GLuint vs = compileShader(GL_VERTEX_SHADER, (shaders / fs::path("vertex.glsl")).string());
    GLuint fs = compileShader(GL_FRAGMENT_SHADER, (shaders / fs::path("fragment.glsl")).string());

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