#include <epoxy/gl.h>
#include <GLFW/glfw3.h>

#include <glm/glm.hpp>
#include <glm/gtc/matrix_transform.hpp>
#include <glm/gtc/type_ptr.hpp>

#include <vector>
#include <chrono>
#include <string>
#include <sstream>

#include "shader.hpp"

class Object
{
    std::vector<float> vertices;
    std::vector<short> indices;
    glm::mat4 M = glm::mat4(1);

    GLuint vao;
    GLuint vertexBuffer;
    GLuint indexBuffer;

    GLfloat speed;
    glm::vec4 colorD;
    glm::vec4 colorS;
    GLfloat powerS;

  public:
    void draw(GLuint program)
    {
        glBindVertexArray(vao);

        glBindBuffer(GL_ARRAY_BUFFER, vertexBuffer);
        GLint locPos = glGetAttribLocation(program, "inPos");
        glVertexAttribPointer(locPos, 3, GL_FLOAT, false, 0, nullptr);
        glEnableVertexAttribArray(locPos);

        M = glm::rotate(M, speed, glm::vec3(0, 1, 0));

        glUniformMatrix4fv(glGetUniformLocation(program, "M"), 1, GL_FALSE, glm::value_ptr(M));
        glUniform4fv(glGetUniformLocation(program, "colorD"), 1, glm::value_ptr(colorD));
        glUniform4fv(glGetUniformLocation(program, "colorS"), 1, glm::value_ptr(colorS));
        glUniform1f(glGetUniformLocation(program, "powerS"), powerS);

        glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, indexBuffer);
        glDrawElements(GL_TRIANGLES, indices.size(), GL_UNSIGNED_SHORT, nullptr);
    }

    void scale(const glm::vec3 &v)
    {
        M = glm::scale(M, v);
    }

    void shift(const glm::vec3 &v)
    {
        M = glm::translate(M, v);
    }

    Object(float speed, const glm::vec4 &colorD, const glm::vec4 &colorS, float powerS, const std::string &filename)
        : speed(speed),
          colorD(colorD),
          colorS(colorS),
          powerS(powerS)
    {
        parse(filename);
        init_buffers();
    }

    void parse(const std::string &filename)
    {
        std::ifstream in(filename);
        std::string line;
        while (getline(in, line))
        {
            if (line[0] == '#')
            {
                continue;
            }

            std::istringstream iline(line);

            char t;
            iline >> t;

            if (t == 'v')
            {
                float x;
                for (int i = 0; i < 3; ++i)
                {
                    iline >> x;
                    vertices.push_back(x);
                }
            }
            else if (t == 'f')
            {
                for (int i = 0; i < 3; ++i)
                {
                    short a;
                    iline >> a;
                    indices.push_back(a - 1);
                }
            }
            else
            {
                std::cerr << "unknown type \"" << t << "\"\n";
            }
        }
    }

    void init_buffers()
    {
        glGenVertexArrays(1, &vao);
        glBindVertexArray(vao);

        glGenBuffers(1, &vertexBuffer);
        glBindBuffer(GL_ARRAY_BUFFER, vertexBuffer);
        glBufferData(GL_ARRAY_BUFFER, vertices.size() * sizeof(GLfloat), &vertices[0], GL_STATIC_DRAW);

        glGenBuffers(1, &indexBuffer);
        glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, indexBuffer);
        glBufferData(GL_ELEMENT_ARRAY_BUFFER, indices.size() * sizeof(GLshort), &indices[0], GL_STATIC_DRAW);
    }
};

int main()
{
    /* Initialize the library */
    if (!glfwInit())
        return -1;

    glfwWindowHint(GLFW_SAMPLES, 4);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 3);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 3);
    glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, GL_TRUE); // To make MacOS happy; should not be needed
    glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);

    /* Create a windowed mode window and its OpenGL context */
    GLFWwindow *window = glfwCreateWindow(640, 480, "Scene", NULL, NULL);
    if (!window)
    {
        glfwTerminate();
        return -1;
    }

    /* Make the window's context current */
    glfwMakeContextCurrent(window);

    glClearColor(0.0, 0.1, 0.2, 1.0);

    /* load objects */

    Object bunny1(0.01, glm::vec4(1, 0, 0, 1), glm::vec4(0.5, 0.5, 0.5, 1), 1, "objects/cube.obj");//stanford_bunny.obj");
    bunny1.scale(glm::vec3(1, 1, 1));

    Object bunny2(0.01, glm::vec4(0, 0, 1, 1), glm::vec4(1, 1, 1, 1), 100, "objects/stanford_bunny.obj");
    bunny2.scale(glm::vec3(10, 10, 10));
    bunny2.shift(glm::vec3(0.2, 0, 0));

    Object plane(0, glm::vec4(0, 1, 0, 1), glm::vec4(0.1, 0.1, 0.1, 1), 1, "objects/plane.obj");

    /* render shadow texture */

    GLuint frameBuffer;
    glGenFramebuffers(1, &frameBuffer);
    glBindFramebuffer(GL_FRAMEBUFFER, frameBuffer);

    constexpr int S_Q = 2;

    GLuint shadowTex;
    glGenTextures(1, &shadowTex);
    glBindTexture(GL_TEXTURE_2D, shadowTex);
    glTexImage2D(GL_TEXTURE_2D, 0, GL_DEPTH_COMPONENT, S_Q * 640, S_Q * 480, 0, GL_DEPTH_COMPONENT, GL_FLOAT, 0);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE);
    glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE);

    glFramebufferTexture(GL_FRAMEBUFFER, GL_DEPTH_ATTACHMENT, shadowTex, 0);

    glDrawBuffer(GL_NONE);

    if (glCheckFramebufferStatus(GL_FRAMEBUFFER) != GL_FRAMEBUFFER_COMPLETE)
    {
        std::cerr << "can't create framebuffer";
    }

    // Compute the MVP matrix from the light's point of view
    glm::mat4 shadowP = glm::ortho<float>(-10, 10, -10, 10, -10, 30);
    glm::vec3 shadowEye = glm::vec3(5, 5, 5);
    glm::vec3 shadowDir = glm::vec3(0, 0, 0) - shadowEye;
    glm::mat4 shadowV = glm::lookAt(shadowEye, shadowEye + shadowDir, glm::vec3(0, 1, 0));
    glm::mat4 shadowCamera = shadowP * shadowV;

    GLuint shadowProgram = loadProgram("shadow/vertex.glsl", "shadow/fragment.glsl");
    glUseProgram(shadowProgram);
    glUniformMatrix4fv(glGetUniformLocation(shadowProgram, "Camera"), 1, GL_FALSE, glm::value_ptr(shadowCamera));

    /* render picture */
    GLuint program = loadProgram("general/vertex.glsl", "general/fragment.glsl");

    glUseProgram(program);
    glm::mat4 P = glm::perspective(45.0f, 4.0f / 3.0f, 0.1f, 100.0f);
    glm::vec3 eye = glm::vec3(0, 5, -5);
    glm::mat4 V = glm::lookAt(eye, glm::vec3(0, 0, 0), glm::vec3(0, 1, 0));
    glm::mat4 camera = P * V;

    GLint locCamera = glGetUniformLocation(program, "Camera");
    glUniformMatrix4fv(locCamera, 1, GL_FALSE, glm::value_ptr(camera));

    // torch

    GLint locTorchPos = glGetUniformLocation(program, "torchPos");
    glUniform3fv(locTorchPos, 1, glm::value_ptr(glm::vec3(5, 5, 5)));

    GLint locTorchPower = glGetUniformLocation(program, "torchPower");
    glUniform1f(locTorchPower, 20);

    // lamp

    GLint locLampPos = glGetUniformLocation(program, "lampPos");
    glUniform3fv(locLampPos, 1, glm::value_ptr(shadowEye));

    GLint locLampDir = glGetUniformLocation(program, "lampDir");
    glUniform3fv(locLampDir, 1, glm::value_ptr(glm::normalize(shadowDir)));

    GLint locLampAngle = glGetUniformLocation(program, "lampAngle");
    glUniform1f(locLampAngle, 0.8);

    GLint locLampPower = glGetUniformLocation(program, "lampPower");
    glUniform1f(locLampPower, 50);

    GLint locEye = glGetUniformLocation(program, "eye");
    glUniform3fv(locEye, 1, glm::value_ptr(eye));

    GLint locShadowCamera = glGetUniformLocation(program, "ShadowCamera");
    glUniformMatrix4fv(locShadowCamera, 1, GL_FALSE, glm::value_ptr(shadowCamera));

    glActiveTexture(GL_TEXTURE0);
    glBindTexture(GL_TEXTURE_2D, shadowTex);
    glUniform1i(glGetUniformLocation(program, "shadow"), 0);
    
    glEnable(GL_DEPTH_TEST);

    /* Loop until the user closes the window */
    while (!glfwWindowShouldClose(window))
    {
        glBindFramebuffer(GL_FRAMEBUFFER, frameBuffer);
        glViewport(0, 0, S_Q * 640, S_Q * 480);
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
        glUseProgram(shadowProgram);

        plane.draw(shadowProgram);
        bunny1.draw(shadowProgram);
        bunny2.draw(shadowProgram);

        ////////////
        glBindFramebuffer(GL_FRAMEBUFFER, 0);
        glViewport(0, 0, 640, 480);
        glUseProgram(program);
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);

        glUniform3fv(locTorchPos, 1, glm::value_ptr(glm::vec3(5, 5, 5 * cos((float)std::chrono::steady_clock::now().time_since_epoch().count() / 100000000.0))));

        plane.draw(program);
        bunny1.draw(program);
        bunny2.draw(program);

        /* Swap front and back buffers */
        glfwSwapBuffers(window);

        /* Poll for and process events */
        glfwPollEvents();
    }

    glfwTerminate();
    return 0;
}