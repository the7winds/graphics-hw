#include <epoxy/gl.h>
#include <GLFW/glfw3.h>

#include <glm/glm.hpp>
#include <glm/gtc/matrix_transform.hpp>
#include <glm/gtc/type_ptr.hpp>

#include <vector>
#include <chrono>
#include <string>
#include <sstream>
#include <memory>
#include <functional>

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

    GLfloat animated = 1;

  public:
    void draw(GLuint program)
    {
        glBindVertexArray(vao);

        glBindBuffer(GL_ARRAY_BUFFER, vertexBuffer);
        GLint locPos = glGetAttribLocation(program, "inPos");
        glVertexAttribPointer(locPos, 3, GL_FLOAT, false, 0, nullptr);
        glEnableVertexAttribArray(locPos);

        M = glm::rotate(M, animated * speed, glm::vec3(0, 1, 0));

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

    void setAnimated(bool animated)
    {
        this->animated = animated ? 1 : 0;
    }
};

class ShadowLamp
{
    static const int S_Q = 2;

    GLuint frameBuffer;
    GLuint shadowTex;

    glm::vec3 pos;
    glm::vec3 dir;
    GLfloat power;
    GLfloat angle;

    GLfloat verticalAngle = -1;
    GLfloat horizontalAngle = 2.2;

  public:
    ShadowLamp(const glm::vec3 &pos,GLfloat power, GLfloat angle) : pos(pos), power(power), angle(angle)
    {
        glGenFramebuffers(1, &frameBuffer);
        glBindFramebuffer(GL_FRAMEBUFFER, frameBuffer);

        update();

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
    }

    void setup(GLint shadowProgram, GLint program)
    {
        glm::mat4 P = glm::ortho<float>(-10, 10, -10, 10, -10, 30);
        glm::mat4 V = glm::lookAt(pos, pos + dir, glm::vec3(0, 1, 0));
        glm::mat4 shadowCamera = P * V;

        glUseProgram(shadowProgram);
        glUniformMatrix4fv(glGetUniformLocation(shadowProgram, "Camera"), 1, GL_FALSE, glm::value_ptr(shadowCamera));

        glUseProgram(program);

        GLint locShadowCamera = glGetUniformLocation(program, "ShadowCamera");
        glUniformMatrix4fv(locShadowCamera, 1, GL_FALSE, glm::value_ptr(shadowCamera));

        GLint locLampPos = glGetUniformLocation(program, "lampPos");
        glUniform3fv(locLampPos, 1, glm::value_ptr(pos));

        GLint locLampDir = glGetUniformLocation(program, "lampDir");
        glUniform3fv(locLampDir, 1, glm::value_ptr(glm::normalize(dir)));

        GLint locLampAngle = glGetUniformLocation(program, "lampAngle");
        glUniform1f(locLampAngle, angle);

        GLint locLampPower = glGetUniformLocation(program, "lampPower");
        glUniform1f(locLampPower, 50);

        glActiveTexture(GL_TEXTURE0);
        glBindTexture(GL_TEXTURE_2D, shadowTex);
        glUniform1i(glGetUniformLocation(program, "shadow"), 0);
    }

    void shadow(GLint shadowProgram)
    {
        glBindFramebuffer(GL_FRAMEBUFFER, frameBuffer);
        glViewport(0, 0, S_Q * 640, S_Q * 480);
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
        glUseProgram(shadowProgram);

    }
    void update()
    {
        dir = glm::vec3(cos(verticalAngle) * sin(horizontalAngle), sin(verticalAngle), cos(verticalAngle) * cos(horizontalAngle));
    }

    void angleUp()
    {
        verticalAngle += 0.01;
        update();
    }

    void angleDown()
    {
        verticalAngle -= 0.01;
        update();
    }

    void angleRight()
    {
        horizontalAngle += 0.01;
        update();
    }

    void angleLeft()
    {
        horizontalAngle -= 0.01;
        update();
    }
};

class AnimatedTorch
{
    glm::vec3 pos;
    GLfloat radius;
    GLfloat power;
    float animated;
    float dz = 0;
    float dx = 1;

  public:
    AnimatedTorch(const glm::vec3 &pos, GLfloat radius, GLfloat power) : pos(pos), radius(radius), power(power){};

    void draw(GLint program)
    {
        using namespace std::chrono;
        int time = duration_cast<milliseconds>(steady_clock::now().time_since_epoch()).count();
        if (animated)
        {
            dx = cos((float)time / 1e3);
            dz = sin((float)time / 1e3);
        }

        GLint locTorchPower = glGetUniformLocation(program, "torchPower");
        glUniform1f(locTorchPower, 20);

        GLint locTorchPos = glGetUniformLocation(program, "torchPos");
        glUniform3fv(locTorchPos, 1, glm::value_ptr(pos + radius * glm::vec3(dx, 0, dz)));
    }

    void setAnimated(bool animated)
    {
        this->animated = animated ? 1 : 0;
    }
};

class Camera
{
    glm::vec3 eye;
    glm::vec3 dir;
    glm::mat4 P;
    glm::mat4 camera;
    float verticalAngle = 0;
    float horizontalAngle = 0;
    bool animationModel = true;

  public:
    Camera()
    {
        P = glm::perspective(45.0f, 4.0f / 3.0f, 0.1f, 100.0f);
        eye = glm::vec3(0, 5, -5);
        update();
    }

    void draw(GLuint program)
    {
        GLint locCamera = glGetUniformLocation(program, "Camera");
        glUniformMatrix4fv(locCamera, 1, GL_FALSE, glm::value_ptr(camera));

        GLint locEye = glGetUniformLocation(program, "eye");
        glUniform3fv(locEye, 1, glm::value_ptr(eye));
    }

    static constexpr float dK = 0.1;

    void dForward()
    {
        eye = eye + dK * dir;
        update();
    }

    void dBack()
    {
        eye = eye - dK * dir;
        update();
    }

    void dLeft()
    {
        glm::vec3 odir = glm::cross(dir, glm::vec3(0, 1, 0));
        eye = eye - dK * odir;
        update();
    }

    void dRight()
    {
        glm::vec3 odir = glm::cross(dir, glm::vec3(0, 1, 0));
        eye = eye + dK * odir;
        update();
    }

    void update()
    {
        dir = glm::vec3(cos(verticalAngle) * sin(horizontalAngle), sin(verticalAngle), cos(verticalAngle) * cos(horizontalAngle));
        camera = P * glm::lookAt(eye, eye + dir, glm::vec3(0, 1, 0));
    }

    void rotate(float dx, float dy)
    {
        horizontalAngle += dx / 100.0;
        verticalAngle += dy / 100.0;
        update();
    }
};

class UI
{
    bool pressed = false;
    bool released = true;
    float oldxpos;
    float oldypos;

  public:
    bool animationModel = false;
    bool animationTorch = false;
    Camera camera;
    ShadowLamp *shadowLamp;

    void onMouseButton(int button, int action)
    {
        if (button == GLFW_MOUSE_BUTTON_LEFT)
        {
            pressed = action == GLFW_PRESS;
            released = action != GLFW_PRESS;
        }
    }

    void onCursorPos(float xpos, float ypos)
    {
        if (pressed)
        {
            pressed = false;
            oldxpos = xpos;
            oldypos = ypos;
        }

        if (!released)
        {
            rotate(xpos, ypos);
        }
    }

    void rotate(float xpos, float ypos)
    {
        float dx = xpos - oldxpos;
        float dy = ypos - oldypos;
        oldxpos = xpos;
        oldypos = ypos;

        camera.rotate(dx, dy);
    }

    void onKey(int key, int action)
    {

        if (key == GLFW_KEY_UP)
        {
            camera.dForward();
        }
        else if (key == GLFW_KEY_DOWN)
        {
            camera.dBack();
        }
        else if (key == GLFW_KEY_LEFT)
        {
            camera.dLeft();
        }
        else if (key == GLFW_KEY_RIGHT)
        {
            camera.dRight();
        }
        else if (key == GLFW_KEY_H)
        {
            shadowLamp->angleLeft();
        }
        else if (key == GLFW_KEY_J)
        {
            shadowLamp->angleRight();
        }
        else if (key == GLFW_KEY_K)
        {
            shadowLamp->angleUp();
        }
        else if (key == GLFW_KEY_L)
        {
            shadowLamp->angleDown();
        }
        
        if (action == GLFW_PRESS)
        {
            if (key == GLFW_KEY_1)
            {
                animationModel = !animationModel;
            }
            else if (key == GLFW_KEY_2)
            {
                animationTorch = !animationTorch;
            }
        }
    }
} ui;

void key_callback(GLFWwindow *, int key, int, int action, int)
{
    ui.onKey(key, action);
}

void cursor_pos_callback(GLFWwindow *, double xpos, double ypos)
{
    ui.onCursorPos(xpos, ypos);
}

void mouse_button_callback(GLFWwindow *, int button, int action, int)
{
    ui.onMouseButton(button, action);
}

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

    Object bunny1(0.01, glm::vec4(1, 0, 0, 1), glm::vec4(0.5, 0.5, 0.5, 1), 1, "objects/stanford_bunny.obj");
    bunny1.scale(glm::vec3(10, 10, 10));

    Object bunny2(0.01, glm::vec4(0, 0, 1, 1), glm::vec4(1, 1, 1, 1), 100, "objects/stanford_bunny.obj");
    bunny2.scale(glm::vec3(10, 10, 10));
    bunny2.shift(glm::vec3(0.2, 0, 0));

    Object plane(0, glm::vec4(0, 1, 0, 1), glm::vec4(0.1, 0.1, 0.1, 1), 1, "objects/plane.obj");

    GLuint shadowProgram = loadProgram("shadow/vertex.glsl", "shadow/fragment.glsl");
    GLuint program = loadProgram("general/vertex.glsl", "general/fragment.glsl");

    // setup lamp's uniforms
    ShadowLamp lamp(glm::vec3(-5, 5, 5), 50, 0.8);
    lamp.setup(shadowProgram, program);
    ui.shadowLamp = &lamp;

    // create a torch
    AnimatedTorch torch(glm::vec3(0, 5, 0), 5, 50);

    // configure camera
    glUseProgram(program);

    glfwSetKeyCallback(window, key_callback);
    glfwSetCursorPosCallback(window, cursor_pos_callback);
    glfwSetMouseButtonCallback(window, mouse_button_callback);

    glEnable(GL_DEPTH_TEST);

    /* Loop until the user closes the window */
    while (!glfwWindowShouldClose(window))
    {
        // evaluate shadow map
        lamp.setup(shadowProgram, program);
        lamp.shadow(shadowProgram);
        plane.draw(shadowProgram);
        bunny1.draw(shadowProgram);
        bunny2.draw(shadowProgram);

        // render
        glBindFramebuffer(GL_FRAMEBUFFER, 0);
        glViewport(0, 0, 640, 480);
        glUseProgram(program);
        glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);

        ui.camera.draw(program);

        bunny1.setAnimated(ui.animationModel);
        bunny2.setAnimated(ui.animationModel);
        torch.setAnimated(ui.animationTorch);

        plane.draw(program);
        bunny1.draw(program);
        bunny2.draw(program);
        torch.draw(program);

        /* Swap front and back buffers */
        glfwSwapBuffers(window);

        /* Poll for and process events */
        glfwPollEvents();
    }

    glfwTerminate();
    return 0;
}