#include <epoxy/gl.h>
#include <GLFW/glfw3.h>

#include <glm/glm.hpp>
#include <glm/gtc/matrix_transform.hpp>
#include <glm/gtc/type_ptr.hpp>

#include <vector>
#include <chrono>

#include "shader.hpp"

class Cube
{
    std::vector<GLfloat> vertices = {
        -1, -1, -1,
        -1, -1,  1,
        -1,  1, -1,
        -1,  1,  1,
         1, -1, -1,
         1, -1,  1,
         1,  1, -1,
         1,  1,  1
    };

    std::vector<GLshort> indices = {
        0, 1, 3, 2,
        0, 1, 5, 4,
        0, 2, 6, 4,
        4, 5, 7, 6,
        1, 3, 7, 5,
        2, 3, 7, 6
    };

    std::vector<GLfloat> color = {
        0, 0, 0,
        0, 0, 1,
        0, 1, 0,
        0, 1, 1,
        1, 0, 0,
        1, 0, 1,
        1, 1, 0,
        1, 1, 1
    };

    GLuint vertexBuffer;
    GLuint indexBuffer;
    GLuint colorBuffer;
    GLuint vertexArray;

  public:
    Cube() 
    {
        glGenVertexArrays(1, &vertexArray);
        glBindVertexArray(vertexArray);

        glGenBuffers(1, &vertexBuffer);
        glBindBuffer(GL_ARRAY_BUFFER, vertexBuffer);
        glBufferData(GL_ARRAY_BUFFER, vertices.size() * sizeof(GLfloat), &vertices[0], GL_STATIC_DRAW);

        glGenBuffers(1, &indexBuffer);
        glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, indexBuffer);
        glBufferData(GL_ELEMENT_ARRAY_BUFFER, indices.size() * sizeof(GLshort), &indices[0], GL_STATIC_DRAW);

        glGenBuffers(1, &colorBuffer);
        glBindBuffer(GL_ARRAY_BUFFER, colorBuffer);
        glBufferData(GL_ARRAY_BUFFER, color.size() * sizeof(GLfloat), &color[0], GL_STATIC_DRAW);
    }

    void draw(GLint program)
    {
        glBindVertexArray(vertexArray);
        
        glBindBuffer(GL_ARRAY_BUFFER, vertexBuffer);
        GLint posPos = glGetAttribLocation(program, "inPos");
        glVertexAttribPointer(posPos, 3, GL_FLOAT, false, 0, nullptr);
        glEnableVertexAttribArray(posPos);

        glBindBuffer(GL_ARRAY_BUFFER, colorBuffer);
        GLint posColor = glGetAttribLocation(program, "inColor");
        glVertexAttribPointer(posColor, 3, GL_FLOAT, false, 0, nullptr);
        glEnableVertexAttribArray(posColor);

        float angle = std::chrono::duration_cast<std::chrono::milliseconds>(std::chrono::steady_clock::now().time_since_epoch()).count() / 1000.0;

        glm::mat4 rotate = glm::rotate(glm::mat4(), angle, glm::vec3(1.0f, 1.0f, 1.0f));

        glUniformMatrix4fv(glGetUniformLocation(program, "rotate"), 1, GL_FALSE, glm::value_ptr(rotate));

        glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, indexBuffer);
        glDrawElements(GL_QUADS, indices.size(), GL_UNSIGNED_SHORT, nullptr);
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

    GLuint program = loadProgram();
    glm::mat4 camera = glm::perspective(45.0f, 4.0f / 3.0f, 0.1f, 100.0f);
    camera = glm::translate(camera, glm::vec3(0.0f, 0.0f, -5.0f));
    Cube cube;

    glUseProgram(program);

    GLint posCamera = glGetUniformLocation(program, "camera");
    glUniformMatrix4fv(posCamera, 1, GL_FALSE, glm::value_ptr(camera));
    
    /* Loop until the user closes the window */
    while (!glfwWindowShouldClose(window))
    {
        /* Render here */
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
        glEnable(GL_DEPTH_TEST);

        cube.draw(program);

        /* Swap front and back buffers */
        glfwSwapBuffers(window);

        /* Poll for and process events */
        glfwPollEvents();
    }

    glfwTerminate();
    return 0;
}