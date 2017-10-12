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

    GLuint vertexBuffer;
    GLuint indexBuffer;
    GLuint colorBuffer;

  public:
    void draw(GLuint program)
    {
        glBindBuffer(GL_ARRAY_BUFFER, vertexBuffer);
        GLint posPos = glGetAttribLocation(program, "inPos");
        glVertexAttribPointer(posPos, 3, GL_FLOAT, false, 0, nullptr);
        glEnableVertexAttribArray(posPos);

        glBindBuffer(GL_ARRAY_BUFFER, colorBuffer);
        GLint posColor = glGetAttribLocation(program, "inColor");
        glVertexAttribPointer(posColor, 3, GL_FLOAT, false, 0, nullptr);
        glEnableVertexAttribArray(posColor);

        glm::mat4 rotate = glm::mat4(1.0f);
        float angle = std::chrono::steady_clock::now().time_since_epoch().count() / 1000000000.0f;
        rotate = glm::rotate(rotate, angle, glm::vec3(0, 1, 0));

        glUniformMatrix4fv(glGetUniformLocation(program, "rotate"), 1, GL_FALSE, glm::value_ptr(rotate));

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

    glm::mat4 camera = glm::perspective(45.0f, 4.0f / 3.0f, 0.1f, 100.0f);
    glm::vec3 eye = glm::vec3(0.0f, -0.1f, -0.5f);
    camera = glm::translate(camera, eye);

    Object object("objects/stanford_bunny.obj");
    object.setColor(glm::vec4(1.0f, 1.0f, 1.0f, 1.0f));

    Object plane("objects/plane.obj");
    plane.setColor(glm::vec4(1.0f, 1.0f, 1.0f, 1.0f));

    glUseProgram(program);

    GLint posCamera = glGetUniformLocation(program, "camera");
    glUniformMatrix4fv(posCamera, 1, GL_FALSE, glm::value_ptr(camera));

    GLint posTorch = glGetUniformLocation(program, "torch");
    glUniform4fv(posTorch, 1, glm::value_ptr(glm::vec3(-5, 1, 0)));

    GLint posEye = glGetUniformLocation(program, "eye");
    glUniform4fv(posEye, 1, glm::value_ptr(eye));

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