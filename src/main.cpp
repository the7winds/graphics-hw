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
    glm::vec4 color;

    GLuint vao;
    GLuint vertexBuffer;
    GLuint indexBuffer;

  public:
    void draw(GLuint program)
    {
        glBindVertexArray(vao);

        glBindBuffer(GL_ARRAY_BUFFER, vertexBuffer);
        GLint locPos = glGetAttribLocation(program, "inPos");
        glVertexAttribPointer(locPos, 3, GL_FLOAT, false, 0, nullptr);
        glEnableVertexAttribArray(locPos);

        glm::mat4 M = glm::mat4(1.0f);
        float angle = std::chrono::steady_clock::now().time_since_epoch().count() / 1000000000.0f;
        M = glm::rotate(M, angle, glm::vec3(0, 1, 0));

        glUniformMatrix4fv(glGetUniformLocation(program, "M"), 1, GL_FALSE, glm::value_ptr(M));
        glUniform4fv(glGetUniformLocation(program, "color"), 1, glm::value_ptr(color));

        glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, indexBuffer);
        glDrawElements(GL_TRIANGLES, indices.size(), GL_UNSIGNED_SHORT, nullptr);
    }

    Object(const std::string &filename)
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

    void setColor(const glm::vec4 &c)
    {
        color = c;
    }
};

int main()
{
    GLFWwindow *window;

    /* Initialize the library */
    if (!glfwInit())
        return -1;

    glfwWindowHint(GLFW_SAMPLES, 4);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MAJOR, 4);
    glfwWindowHint(GLFW_CONTEXT_VERSION_MINOR, 5);
    glfwWindowHint(GLFW_OPENGL_FORWARD_COMPAT, GL_TRUE); // To make MacOS happy; should not be needed
    glfwWindowHint(GLFW_OPENGL_PROFILE, GLFW_OPENGL_CORE_PROFILE);

    /* Create a windowed mode window and its OpenGL context */
    window = glfwCreateWindow(640, 480, "Hello World", NULL, NULL);
    if (!window)
    {
        glfwTerminate();
        return -1;
    }

    /* Make the window's context current */
    glfwMakeContextCurrent(window);

    glClearColor(0.0, 0.1, 0.2, 1.0);

    GLuint program = loadProgram("general/vertex.glsl", "general/fragment.glsl");

    glm::mat4 P = glm::perspective(45.0f, 4.0f / 3.0f, 0.1f, 100.0f);
    glm::vec3 eye = glm::vec3(-0.3, 0.3, -0.3);
    glm::mat4 V = glm::lookAt(eye, glm::vec3(0, 0, 0), glm::vec3(0, 1, 0));

    Object object("objects/stanford_bunny.obj");
    object.setColor(glm::vec4(1.0f, 1.0f, 1.0f, 1.0f));

    Object plane("objects/plane.obj");
    plane.setColor(glm::vec4(1.0f, 1.0f, 1.0f, 1.0f));

    glUseProgram(program);

    GLint locCamera = glGetUniformLocation(program, "Camera");
    glUniformMatrix4fv(locCamera, 1, GL_FALSE, glm::value_ptr(P * V));

    GLint locTorch = glGetUniformLocation(program, "torch");
    glUniform4fv(locTorch, 1, glm::value_ptr(glm::vec3(-5, 1, 0)));

    GLint locEye = glGetUniformLocation(program, "eye");
    glUniform4fv(locEye, 1, glm::value_ptr(eye));

    /* Loop until the user closes the window */
    while (!glfwWindowShouldClose(window))
    {
        /* Render here */
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
        glEnable(GL_DEPTH_TEST);

        plane.draw(program);
        object.draw(program);

        /* Swap front and back buffers */
        glfwSwapBuffers(window);

        /* Poll for and process events */
        glfwPollEvents();
    }

    glfwTerminate();
    return 0;
}